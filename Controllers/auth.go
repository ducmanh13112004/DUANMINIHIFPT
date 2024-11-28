package controllers

import (
	"MiniHIFPT/database"
	"MiniHIFPT/models"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"strconv"
	"time"
)

// Hàm đăng ký tài khoản
func Register(c *fiber.Ctx) error {
	// Lấy thông tin đăng ký từ request body
	var newAccount models.Accounts
	if err := c.BodyParser(&newAccount); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Dữ liệu đầu vào không hợp lệ",
		})
	}

	// Kiểm tra nếu số điện thoại đã tồn tại trong hệ thống
	var existingAccount models.Accounts
	result := database.DB.Where("SoDienThoai = ?", newAccount.SoDienThoai).First(&existingAccount)

	if result.Error == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "Số điện thoại đã tồn tại trong hệ thống",
		})
	}

	// Mã hóa mật khẩu (ví dụ sử dụng bcrypt)
	// := khai báo và gán giá trị ...
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newAccount.MatKhau), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Không thể mã hóa mật khẩu",
		})
	}
	newAccount.MatKhau = string(hashedPassword)

	// Tạo tài khoản mới
	if err := database.DB.Create(&newAccount).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Không thể tạo tài khoản mới",
		})
	}

	// Trả về kết quả thành công
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Tạo tài khoản thành công. Vui lòng đăng nhập.",
	})
}

// Hàm login
func Login(c *fiber.Ctx) error {
	// Lấy thông tin số điện thoại và mật khẩu từ request body
	var loginCredentials models.Accounts
	if err := c.BodyParser(&loginCredentials); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Dữ liệu đầu vào không hợp lệ",
		})
	}

	// Kiểm tra số điện thoại có tồn tại trong hệ thống không
	var account models.Accounts
	result := database.DB.Where("SoDienThoai = ?", loginCredentials.SoDienThoai).First(&account)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Số điện thoại hoặc mật khẩu không đúng",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": result.Error.Error(),
		})
	}

	// Kiểm tra số lần nhập sai trong ngày
	var loginAttempts models.LoginAttempt
	dailyAttempts := database.DB.Where("SoDienThoai = ? AND DATE(Ngay) = CURDATE()", loginCredentials.SoDienThoai).First(&loginAttempts)

	// Nếu đã nhập sai 5 lần trong ngày
	if dailyAttempts.Error == nil && loginAttempts.SoLanSai >= 5 {
		// Kiểm tra thời gian lần nhập sai cuối cùng
		if time.Since(loginAttempts.Ngay) < time.Minute {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Bạn đã nhập sai mật khẩu quá 5 lần. Vui lòng thử lại sau 1 phút.",
			})
		}
		//time.Minute: đại diện cho 1 phút.
		// Nếu thời gian đã quá 1 phút, reset số lần sai
		loginAttempts.SoLanSai = 0
		loginAttempts.Ngay = time.Now() // Cập nhật thời gian
		database.DB.Save(&loginAttempts)
	}

	// Kiểm tra mật khẩu có đúng không
	if err := bcrypt.CompareHashAndPassword([]byte(account.MatKhau), []byte(loginCredentials.MatKhau)); err != nil {
		// Nếu mật khẩu sai, tăng số lần sai đăng nhập
		if loginAttempts.ID == "" {
			// Nếu chưa có bản ghi, tạo mới
			loginAttempts = models.LoginAttempt{
				SoDienThoai: loginCredentials.SoDienThoai,
				SoLanSai:    1,
			}
			database.DB.Create(&loginAttempts)
		} else {
			// Nếu có bản ghi, tăng số lần sai
			loginAttempts.SoLanSai++
			loginAttempts.Ngay = time.Now() // Cập nhật thời gian nhập sai
			database.DB.Save(&loginAttempts)
		}
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Số điện thoại hoặc mật khẩu không đúng",
		})
	}

	// Nếu đăng nhập thành công, reset số lần sai
	if loginAttempts.ID != "" {
		loginAttempts.SoLanSai = 0
		loginAttempts.KhoiPhuc = true
		loginAttempts.Ngay = time.Now() // Cập nhật thời gian thành công
		database.DB.Save(&loginAttempts)
	}

	// Kiểm tra nếu là thiết bị mới (cần OTP)
	var device models.Devices
	resultDevice := database.DB.Where("SoDienThoai = ?", account.SoDienThoai).First(&device)

	// Nếu chưa có thiết bị hoặc lần đăng nhập từ thiết bị mới, gửi OTP
	if resultDevice.Error != nil {
		otpCode := generateOTP()

		// Tạo OTP mới trong database
		otp := models.OTPCode{
			SoDienThoai: account.SoDienThoai,
			OTP_Code:    otpCode,
			HetHan:      time.Now().Add(5 * time.Minute),
		}
		if err := database.DB.Create(&otp).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Không thể tạo mã OTP",
			})
		}

		// Trả về yêu cầu nhập OTP
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Mã OTP đã được gửi. Vui lòng nhập mã OTP.",
		})
	}

	// Nếu đã đăng nhập trên thiết bị này trước đó, trả về kết quả thành công
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Đăng nhập thành công, không cần nhập OTP.",
	})
}

// Hàm xác thực OTP
func VerifyOTP(c *fiber.Ctx) error {
	// Lấy số điện thoại và OTP từ request body
	var otpRequest struct {
		SoDienThoai string `json:"SoDienThoai"`
		OTPCode     string `json:"OTPCode"`
	}
	if err := c.BodyParser(&otpRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Dữ liệu đầu vào không hợp lệ",
		})
	}

	// Lấy OTP từ cơ sở dữ liệu
	var otp models.OTPCode
	result := database.DB.Where("SoDienThoai = ? AND OTP_Code = ?", otpRequest.SoDienThoai, otpRequest.OTPCode).First(&otp)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Mã OTP không hợp lệ",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": result.Error.Error(),
		})
	}

	// Kiểm tra mã OTP đã hết hạn chưa
	if otp.HetHan.Before(time.Now()) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Mã OTP đã hết hạn",
		})
	}

	// Đánh dấu OTP là đã xác thực
	otp.DaXacThuc = true
	if err := database.DB.Save(&otp).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Không thể xác thực mã OTP",
		})
	}

	// Lưu thiết bị vào bảng Devices
	device := models.Devices{
		ID:          utils.UUIDv4(),
		SoDienThoai: otpRequest.SoDienThoai,
	}
	if err := database.DB.Create(&device).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Không thể lưu thiết bị",
		})
	}

	// Trả về kết quả đăng nhập thành công
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Đăng nhập thành công.",
	})
}

// Hàm tạo mã OTP ngẫu nhiên
func generateOTP() string {
	rand.Seed(time.Now().UnixNano())
	otp := rand.Intn(999999-100000) + 100000
	return strconv.Itoa(otp)
}
