package database

import (
	"MiniHIFPT/models"
	"fmt"
)

// Lấy tất cả các hợp đồng
func GetContracts() ([]models.Contract, error) {
	var contracts []models.Contract
	result := DB.Find(&contracts)
	return contracts, result.Error
}

// Tạo hợp đồng mới
func CreateContract(contract *models.Contract) error {
	result := DB.Create(&contract)
	return result.Error
}

// Lấy hợp đồng theo ID
func GetContractByID(idUUID string) (models.Contract, error) {
	var contract models.Contract

	// In thông tin để kiểm tra
	fmt.Println("Querying contract with ID:", idUUID)

	// Truy vấn dữ liệu từ database
	result := DB.First(&contract, "id_uuid = ?", idUUID)

	// In lỗi nếu xảy ra
	if result.Error != nil {
		fmt.Println("Error:", result.Error)
	}

	return contract, result.Error
}

// Cập nhật hợp đồng
func UpdateContract(contract *models.Contract, updates map[string]interface{}) error {
	result := DB.Model(&contract).Updates(updates)
	return result.Error
}

// Xóa hợp đồng
func DeleteContract(contract *models.Contract) error {
	result := DB.Delete(&contract)
	return result.Error
}
