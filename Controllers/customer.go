package controllers

import (
	"MiniHIFPT/database"
	"MiniHIFPT/models"
	"github.com/gofiber/fiber/v2"
)

// Lấy thông tin các hợp đồng
func GetCustomers(c *fiber.Ctx) error {
	var customers []models.Customer //khởi tạo danh sách

	// Lấy danh sách hợp đồng
	result := database.DB.Find(&customers)
	if result.Error != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Không thể lấy thông tin khách hàng",
		})
	}

	// Trả về danh sách hợp đồng
	return c.JSON(customers)
}

// Tạo hợp đồng mới (thêm)
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

	// Kiểm tra dữ liệu hợp đồng, ví dụ kiểm tra trường hợp hợp đồng đã có
	if customer.SoDienThoai == "" || customer.TenKhachHang == "" || customer.GioiTinh == "" || customer.NgaySinh == nil || customer.Email == "" || customer.LoaiKhachHang == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Thiếu thông tin cần thiết",
		})
	}

	// Thêm hợp đồng vào cơ sở dữ liệu
	result := database.DB.Create(&customer)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Không thể tạo khách hàng",
			"details": result.Error.Error(), // Trả về chi tiết lỗi từ kết quả tạo
		})
	}

	// Trả về hợp đồng đã được thêm vào
	return c.Status(fiber.StatusCreated).JSON(customer)
}
