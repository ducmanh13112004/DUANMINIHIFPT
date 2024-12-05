package database

import (
	"MiniHIFPT/models"
)

func GetContractByIDs(contractID string) (*models.Contract, error) {
	var contract models.Contract
	if err := DB.Where("id_uuid = ?", contractID).First(&contract).Error; err != nil {
		return nil, err
	}
	return &contract, nil
}

func CheckAccess(accountID, contractID string) (bool, error) {
	var accountContract models.Account_Contract
	if err := DB.Where("account_id = ? AND contract_id = ?", accountID, contractID).First(&accountContract).Error; err != nil {
		return false, err
	}
	return true, nil
}
