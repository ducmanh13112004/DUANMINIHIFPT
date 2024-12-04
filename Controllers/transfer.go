package controllers

import (
	"MiniHIFPT/database"

	"github.com/gofiber/fiber/v2"
)

// API chuyển sở hữu hợp đồng
func TransferOwnership(c *fiber.Ctx) error {
	// Lấy tham số từ yêu cầu JSON
	var request struct {
		OldCustomerID string `json:"oldCustomerId"`
		NewCustomerID string `json:"newCustomerId"`
	}

	// Phân tích dữ liệu JSON từ yêu cầu
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Giá trị nhập vào không hợp lệ",
		})
	}

	// Kiểm tra xem khách hàng cũ có tồn tại không
	oldCustomer, err := database.FindCustomerByID(request.OldCustomerID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Không tìm thấy khách hàng cũ",
		})
	}

	// Kiểm tra xem khách hàng mới có tồn tại không
	newCustomer, err := database.FindCustomerByID(request.NewCustomerID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Không tìm thấy khách hàng mới",
		})
	}

	// Lấy tất cả hợp đồng của khách hàng cũ dựa trên số điện thoại
	customerContracts, err := database.FindCustomerContractsByPhoneNumber(oldCustomer.SoDienThoai)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Không thể tìm thấy các hợp đồng của khách hàng cũ",
		})
	}

	// Nếu không có hợp đồng nào thuộc khách hàng cũ
	if len(customerContracts) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Không có hợp đồng nào thuộc về khách hàng cũ",
		})
	}

	// Cập nhật tất cả các hợp đồng để chuyển sang khách hàng mới
	for _, contract := range customerContracts {
		// Cập nhật cột `SoDienThoai` sang số điện thoại của khách hàng mới
		if err := database.TransferContractOwnership(&contract, newCustomer.SoDienThoai); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Không thể chuyển nhượng hợp đồng",
			})
		}
	}

	// Trả về kết quả
	return c.JSON(fiber.Map{
		"message": "Tất cả hợp đồng đã được chuyển giao thành công",
	})
}
