package controllers

import (
	"MiniHIFPT/database"
	"MiniHIFPT/models"
	"github.com/gofiber/fiber/v2"
)

// Lấy thông tin các hợp đồng
func GetContracts(c *fiber.Ctx) error {
	var contracts []models.Contract //khởi tạo danh sách

	// Lấy danh sách hợp đồng
	result := database.DB.Find(&contracts)
	if result.Error != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Không thể lấy thông tin hợp đồng",
		})
	}

	// Trả về danh sách hợp đồng
	return c.JSON(contracts)
}

// Tạo hợp đồng mới (thêm)
func CreateContract(c *fiber.Ctx) error {
	var contract models.Contract

	// Phân tích dữ liệu JSON từ yêu cầu POST
	if err := c.BodyParser(&contract); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Dữ liệu đầu vào không hợp lệ",
			"details": err.Error(), // In chi tiết lỗi ra
		})
	}
	//StatusBadRequest lỗi cú pháp và dữ liệu không hợp lệ hoặc thiếu thông tin.

	// Kiểm tra dữ liệu hợp đồng, ví dụ kiểm tra trường hợp hợp đồng đã có
	if contract.TenKhachHang == "" || contract.DiaChi == "" || contract.MaTinh == "" || contract.MaQuanHuyen == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Thiếu thông tin cần thiết",
		})
	}

	// Thêm hợp đồng vào cơ sở dữ liệu
	result := database.DB.Create(&contract)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Không thể tạo hợp đồng",
			"details": result.Error.Error(), // Trả về chi tiết lỗi từ kết quả tạo
		})
	}

	// Trả về hợp đồng đã được thêm vào
	return c.Status(fiber.StatusCreated).JSON(contract)
}

// Sửa thông tin hợp đồng
func UpdateContract(c *fiber.Ctx) error {
	// Lấy ID hợp đồng từ URL (được truyền dưới dạng tham số)
	id := c.Params("id")

	// Tìm hợp đồng theo ID
	var contract models.Contract
	if err := database.DB.First(&contract, "id = ?", id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Hợp đồng không tồn tại",
		})
	}

	// Phân tích dữ liệu từ yêu cầu PUT/PATCH
	var updatedData models.Contract
	if err := c.BodyParser(&updatedData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Dữ liệu đầu vào không hợp lệ",
			"details": err.Error(),
		})
	}

	// Cập nhật thông tin hợp đồng
	if updatedData.TenKhachHang != "" {
		contract.TenKhachHang = updatedData.TenKhachHang
	}
	if updatedData.DiaChi != "" {
		contract.DiaChi = updatedData.DiaChi
	}
	if updatedData.MaTinh != "" {
		contract.MaTinh = updatedData.MaTinh
	}
	if updatedData.MaQuanHuyen != "" {
		contract.MaQuanHuyen = updatedData.MaQuanHuyen
	}
	if updatedData.MaPhuongXa != "" {
		contract.MaPhuongXa = updatedData.MaPhuongXa
	}
	if updatedData.MaDuong != "" {
		contract.MaDuong = updatedData.MaDuong
	}
	if updatedData.SoNha != "" {
		contract.SoNha = updatedData.SoNha
	}

	// Lưu các thay đổi vào cơ sở dữ liệu
	if err := database.DB.Save(&contract).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Không thể cập nhật hợp đồng",
			"details": err.Error(),
		})
	} //

	// Trả về hợp đồng đã được cập nhật
	return c.JSON(contract)
}

// Xóa
func DeleteContract(c *fiber.Ctx) error {
	// Lấy ID hợp đồng từ tham số URL
	id := c.Params("id")

	// Tìm hợp đồng trong cơ sở dữ liệu
	var contract models.Contract
	if err := database.DB.First(&contract, "id = ?", id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Hợp đồng không tồn tại",
		})
	}

	// Xóa hợp đồng
	if err := database.DB.Delete(&contract).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Không thể xóa hợp đồng",
		})
	}

	// Trả về phản hồi thành công
	return c.JSON(fiber.Map{
		"message": "Xóa hợp đồng thành công",
	})
}