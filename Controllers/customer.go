package controllers

import (
	"MiniHIFPT/database"
	"MiniHIFPT/models"
	"github.com/gofiber/fiber/v2"
)

// Lấy thông tin các kh
func GetCustomers(c *fiber.Ctx) error {
	var customers []models.Customer //khởi tạo danh sách

	// Lấy danh sách kh
	result := database.DB.Find(&customers)
	if result.Error != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Không thể lấy thông tin khách hàng",
		})
	}

	// Trả về danh sách kh
	return c.JSON(customers)
}

// Tạo khách hàng mới (thêm)
func CreateCustomers(c *fiber.Ctx) error {
	var customer models.Customer

	// Phân tích dữ liệu JSON từ yêu cầu POST
	if err := c.BodyParser(&customer); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Dữ liệu đầu vào không hợp lệ",
			"details": err.Error(), // In chi tiết lỗi ra
		})
	}
	//StatusBadRequest lỗi cú pháp và dữ liệu không hợp lệ hoặc thiếu thông tin.

	// Kiểm tra dữ liệu khách hàng
	if customer.SoDienThoai == "" || customer.TenKhachHang == "" || customer.GioiTinh == "" || customer.NgaySinh == nil || customer.Email == "" || customer.LoaiKhachHang == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Thiếu thông tin cần thiết",
		})
	}

	// Thêm khách hàng vào cơ sở dữ liệu
	result := database.DB.Create(&customer)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Không thể tạo khách hàng",
			"details": result.Error.Error(), // Trả về chi tiết lỗi từ kết quả tạo
		})
	}

	// Trả về khách hàng đã được thêm vào (customer)
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Thêm khách hàng thành công",
	})

}
