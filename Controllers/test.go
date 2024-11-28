package controllers

import (
	"MiniHIFPT/database"
	"MiniHIFPT/models"
	"github.com/gofiber/fiber/v2"
)

// Struct chứa dữ liệu yêu cầu chuyển nhượng
type TransferRequest struct {
	NewPhone string `json:"newPhone"`
	NewName  string `json:"newName"`
}

// Chuyển sở hữu hợp đồng từ một khách hàng sang khách hàng khác (bao gồm cả số điện thoại và tên)
func TransferOwnership(c *fiber.Ctx) error {
	var contract models.Contract
	var newCustomer models.Customer
	var customerContract models.Customer_Contractt

	// Lấy ID hợp đồng từ tham số URL
	contractID := c.Params("contractID")
	if contractID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID hợp đồng không hợp lệ",
		})
	}

	// Lấy dữ liệu JSON từ body yêu cầu
	var req TransferRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Dữ liệu không hợp lệ hoặc thiếu thông tin",
		})
	}

	// Kiểm tra nếu số điện thoại hoặc tên khách hàng mới trống
	if req.NewPhone == "" || req.NewName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Số điện thoại và tên khách hàng mới không được để trống",
		})
	}

	// Tìm hợp đồng cần chuyển nhượng
	if err := database.DB.First(&contract, "id_uuid = ?", contractID).Error; err != nil {
		// Nếu không tìm thấy hợp đồng, trả về lỗi
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Hợp đồng không tồn tại",
		})
	}

	// Tìm khách hàng mới bằng số điện thoại
	if err := database.DB.First(&newCustomer, "SoDienThoai = ?", req.NewPhone).Error; err != nil {
		// Nếu không tìm thấy khách hàng với số điện thoại, trả về lỗi
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Khách hàng không tồn tại với số điện thoại đã cung cấp",
		})
	}

	// Cập nhật tên khách hàng mới cho hợp đồng
	contract.TenKhachHang = req.NewName
	// Cập nhật thông tin khác nếu cần

	// Lưu lại hợp đồng đã cập nhật
	if err := database.DB.Save(&contract).Error; err != nil {
		// Nếu không thể lưu hợp đồng, trả về lỗi
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Không thể chuyển nhượng hợp đồng",
		})
	}

	// Tạo hoặc cập nhật Customer_Contractt để liên kết hợp đồng với khách hàng mới
	customerContract.HopDongID = contractID
	customerContract.SoDienThoai = req.NewPhone

	// Lưu thông tin liên kết hợp đồng và khách hàng
	if err := database.DB.Save(&customerContract).Error; err != nil {
		// Nếu không thể lưu liên kết, trả về lỗi
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Không thể liên kết hợp đồng với khách hàng mới",
		})
	}

	// Trả về thông tin hợp đồng sau khi chuyển nhượng
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":  "Chuyển nhượng hợp đồng thành công",
		"contract": contract,
	})
}
