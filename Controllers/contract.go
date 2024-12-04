package controllers

import (
	"MiniHIFPT/database"
	"MiniHIFPT/models"
	"errors"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// Lấy thông tin các hợp đồng
func GetContracts(c *fiber.Ctx) error {
	contracts, err := database.GetContracts()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Không thể lấy thông tin hợp đồng",
		})
	}
	return c.JSON(contracts)
}

// Lấy thông tin chi tiết một hợp đồng
// func GetContractByID(c *fiber.Ctx) error {
// 	// Lấy ID từ URL
// 	idUUID := c.Params("id_uuid")
// 	// Gọi hàm truy vấn từ package database
// 	contract, err := database.GetContractByID(idUUID)
// 	if err != nil {
// 		// Kiểm tra lỗi không tìm thấy
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
// 				"error": "Không tìm thấy hợp đồng với ID được cung cấp",
// 			})
// 		}
// 		// Xử lý lỗi khác
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 			"error": "Đã xảy ra lỗi khi truy vấn hợp đồng",
// 		})
// 	}

// 	// Trả về thông tin hợp đồng dưới dạng JSON
// 	// return c.JSON(fiber.Map{
// 	// 	"message": "Hợp đồng bạn yêu cầu là: ",
// 	// })
// 	return c.JSON(contract)
// }

// Thêm quyền truy cập cho tài khoản đối với hợp đồng
func AddContractAccess(c *fiber.Ctx) error {
	var data struct {
		AccountID  string `json:"accountID"`
		ContractID string `json:"contractID"`
	}

	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Dữ liệu đầu vào không hợp lệ",
		})
	}

	// Thêm quyền truy cập vào bảng trung gian
	accountContract := models.Account_Contract{
		AccountID:  data.AccountID,
		ContractID: data.ContractID,
	}

	if err := database.DB.Create(&accountContract).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Không thể thêm quyền truy cập",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Thêm quyền truy cập thành công",
	})
}

// Lấy thông tin chi tiết một hợp đồng với kiểm tra quyền truy cập
func GetContractByID(c *fiber.Ctx) error {
	// Lấy ID từ URL
	idUUID := c.Params("id_uuid")

	// Kiểm tra quyền truy cập (nếu cần)
	accountID := c.Locals("accountID").(string) // Nếu sử dụng middleware để lấy thông tin tài khoản
	var count int64
	if err := database.DB.Model(&models.Account_Contract{}).
		Where("AccountID = ? AND ContractID = ?", accountID, idUUID).
		Count(&count).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Đã xảy ra lỗi khi kiểm tra quyền truy cập",
		})
	}
	if count == 0 {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Bạn không có quyền truy cập hợp đồng này",
		})
	}

	// Gọi hàm truy vấn từ package database
	contract := models.Contract{}
	err := database.DB.Where("id_uuid = ?", idUUID).First(&contract).Error

	// Kiểm tra lỗi khi truy vấn dữ liệu
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Không tìm thấy hợp đồng với ID được cung cấp",
			})
		}
		// Lỗi khác (có thể là lỗi kết nối DB hoặc lỗi khác)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Đã xảy ra lỗi khi truy vấn hợp đồng: " + err.Error(),
		})
	}

	// Kiểm tra xem đối tượng hợp đồng có dữ liệu hợp lệ không (ID là string)
	if contract.ID == "" { // Kiểm tra với chuỗi trống nếu ID là string
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Hợp đồng không tồn tại",
		})
	}

	// Trả về dữ liệu hợp đồng
	return c.JSON(fiber.Map{
		"data": contract,
	})
}

// Tạo hợp đồng mới (thêm)
func CreateContract(c *fiber.Ctx) error {
	var contract models.Contract
	if err := c.BodyParser(&contract); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Dữ liệu đầu vào không hợp lệ",
		})
	}

	// Kiểm tra dữ liệu hợp đồng
	if contract.TenKhachHang == "" || contract.DiaChi == "" || contract.MaTinh == "" || contract.MaQuanHuyen == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Thiếu thông tin cần thiết",
		})
	}

	if err := database.CreateContract(&contract); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Không thể tạo hợp đồng",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Tạo hợp đồng thành công",
	})
}

// Sửa thông tin hợp đồng
func UpdateContract(c *fiber.Ctx) error {
	id := c.Params("id")
	contract, err := database.GetContractByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Hợp đồng không tồn tại",
		})
	}

	var updatedData models.Contract
	if err := c.BodyParser(&updatedData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Dữ liệu đầu vào không hợp lệ",
		})
	}

	// Cập nhật các trường hợp hợp đồng
	updates := map[string]interface{}{}
	if updatedData.TenKhachHang != "" {
		updates["TenKhachHang"] = updatedData.TenKhachHang
	}
	if updatedData.DiaChi != "" {
		updates["DiaChi"] = updatedData.DiaChi
	}
	if updatedData.MaTinh != "" {
		updates["MaTinh"] = updatedData.MaTinh
	}
	if updatedData.MaQuanHuyen != "" {
		updates["MaQuanHuyen"] = updatedData.MaQuanHuyen
	}
	if updatedData.MaPhuongXa != "" {
		updates["MaPhuongXa"] = updatedData.MaPhuongXa
	}
	if updatedData.MaDuong != "" {
		updates["MaDuong"] = updatedData.MaDuong
	}
	if updatedData.SoNha != "" {
		updates["SoNha"] = updatedData.SoNha
	}

	if len(updates) > 0 {
		if err := database.UpdateContract(&contract, updates); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Không thể cập nhật hợp đồng",
			})
		}
	}

	return c.JSON(fiber.Map{
		"message": "Sửa hợp đồng thành công",
	})
}

// Xóa hợp đồng
func DeleteContract(c *fiber.Ctx) error {
	id := c.Params("id")
	contract, err := database.GetContractByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Hợp đồng không tồn tại",
		})
	}

	if err := database.DeleteContract(&contract); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Không thể xóa hợp đồng",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Xóa hợp đồng thành công",
	})
}
