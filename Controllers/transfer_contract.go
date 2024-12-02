package controllers

import (
	"MiniHIFPT/database"
	"MiniHIFPT/models"
	"errors"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// Hàm chuyển sở hữu hợp đồng từ một khách hàng sang một khách hàng khác
func ChuyenSoHuu(db *gorm.DB, soDienThoaiNguon string, soDienThoaiDich string) error {
	var countNguon, countDich int64

	// Kiểm tra khách hàng nguồn
	if err := db.Model(&models.Customer{}).Where("SoDienThoai = ?", soDienThoaiNguon).Count(&countNguon).Error; err != nil {
		// Ghi log lỗi vào database
		db.Create(&models.Log{
			Action:  "Chuyển sở hữu hợp đồng",
			Details: "Lỗi khi kiểm tra khách hàng nguồn: " + err.Error(),
		})
		return errors.New("lỗi khi kiểm tra khách hàng nguồn")
	}
	if countNguon == 0 {
		// Ghi log nếu không tìm thấy khách hàng nguồn
		db.Create(&models.Log{
			Action:  "Chuyển sở hữu hợp đồng",
			Details: "Khách hàng nguồn không tồn tại: " + soDienThoaiNguon,
		})
		return errors.New("khách hàng nguồn không tồn tại")
	}

	// Kiểm tra khách hàng đích
	if err := db.Model(&models.Customer{}).Where("SoDienThoai = ?", soDienThoaiDich).Count(&countDich).Error; err != nil {
		// Ghi log lỗi nếu không thể kiểm tra khách hàng đích
		db.Create(&models.Log{
			Action:  "Chuyển sở hữu hợp đồng",
			Details: "Lỗi khi kiểm tra khách hàng đích: " + err.Error(),
		})
		return errors.New("lỗi khi kiểm tra khách hàng đích")
	}
	if countDich == 0 {
		// Ghi log nếu không tìm thấy khách hàng đích
		db.Create(&models.Log{
			Action:  "Chuyển sở hữu hợp đồng",
			Details: "Khách hàng đích không tồn tại: " + soDienThoaiDich,
		})
		return errors.New("khách hàng đích không tồn tại")
	}

	// Tiến hành chuyển sở hữu hợp đồng
	if err := db.Model(&models.Customer_Contractt{}).
		Where("SoDienThoai = ?", soDienThoaiNguon).
		Update("SoDienThoai", soDienThoaiDich).Error; err != nil {
		// Ghi log lỗi khi chuyển sở hữu hợp đồng
		db.Create(&models.Log{
			Action:  "Chuyển sở hữu hợp đồng",
			Details: "Lỗi khi chuyển sở hữu hợp đồng: " + err.Error(),
		})
		return errors.New("lỗi khi chuyển sở hữu hợp đồng")
	}

	// Ghi nhật ký thành công
	db.Create(&models.Log{
		Action:  "Chuyển sở hữu hợp đồng",
		Details: "Chuyển sở hữu từ " + soDienThoaiNguon + " sang " + soDienThoaiDich + " thành công",
	})

	return nil
}

// Handler cho Fiber để gọi hàm ChuyenSoHuu
func ChuyenSoHuuHandler(c *fiber.Ctx) error {
	// Lấy tham số từ body JSON
	var request struct {
		SoDienThoaiNguon string `json:"soDienThoaiNguon"`
		SoDienThoaiDich  string `json:"soDienThoaiDich"`
	}

	// Kiểm tra lỗi khi phân tích dữ liệu JSON
	if err := c.BodyParser(&request); err != nil {
		// Trả về lỗi nếu body không hợp lệ
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Dữ liệu không hợp lệ",
			"details": err.Error(), // In chi tiết lỗi ra
		})
	}

	// Kiểm tra dữ liệu đầu vào
	if request.SoDienThoaiNguon == "" || request.SoDienThoaiDich == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Thiếu thông tin số điện thoại nguồn hoặc đích",
			"details": "Số điện thoại nguồn hoặc đích không hợp lệ",
		})
	}

	// Gọi hàm ChuyenSoHuu từ Controller
	if err := ChuyenSoHuu(database.DB, request.SoDienThoaiNguon, request.SoDienThoaiDich); err != nil {
		// Trả về lỗi nếu gặp vấn đề khi chuyển sở hữu
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Lỗi khi chuyển sở hữu",
			"details": err.Error(),
		})
	}

	// Trả về phản hồi thành công
	return c.JSON(fiber.Map{
		"message": "Chuyển sở hữu thành công",
	})
}
