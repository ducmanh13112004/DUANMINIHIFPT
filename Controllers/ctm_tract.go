package controllers

import (
	"MiniHIFPT/database"
	"MiniHIFPT/models"
	"github.com/gofiber/fiber/v2"
)

// Lấy thông tin các hợp đồng
func Getctm_tracts(c *fiber.Ctx) error {
	var ctm_tracts []models.Customer_Contractt //khởi tạo danh sách

	// Lấy danh sách hợp đồng
	result := database.DB.Find(&ctm_tracts)
	if result.Error != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Không thể lấy thông tin ",
		})
	}

	// Trả về danh sách hợp đồng
	return c.JSON(ctm_tracts)
}

// (thêm liên kết )
func Createctm_tracts(c *fiber.Ctx) error {
	var ctm_tract models.Customer_Contractt

	// Phân tích dữ liệu JSON từ yêu cầu POST
	if err := c.BodyParser(&ctm_tract); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Dữ liệu đầu vào không hợp lệ",
			"details": err.Error(), // In chi tiết lỗi ra
		})
	}
	//StatusBadRequest lỗi cú pháp và dữ liệu không hợp lệ hoặc thiếu thông tin.

	// Kiểm tra dữ liệu hợp đồng, ví dụ kiểm tra trường hợp hợp đồng đã có
	if ctm_tract.SoDienThoai == "" || ctm_tract.HopDongID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Thiếu thông tin cần thiết",
		})
	}

	// Thêm hợp đồng vào cơ sở dữ liệu
	result := database.DB.Create(&ctm_tract)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Không thể tạo liên kết  hợp đồng",
			"details": result.Error.Error(), // Trả về chi tiết lỗi từ kết quả tạo
		})
	}

	// Trả về hợp đồng đã được thêm vào
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Tạo liên kết số điện thoại và hợp đồng thành công",
	})

}
