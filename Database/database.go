package database

// import (
// 	"MiniHIFPT/models"
// 	"github.com/jinzhu/gorm"
// )

// // Kiểm tra số điện thoại đã tồn tại trong hệ thống
// func CheckPhoneExist(soDienThoai string) (*models.Accounts, error) {
// 	var account models.Accounts
// 	result := DB.Where("SoDienThoai = ?", soDienThoai).First(&account)
// 	if result.Error != nil {
// 		return nil, result.Error
// 	}
// 	return &account, nil
// }

// // Tạo tài khoản mới
// func CreateAccount(account *models.Accounts) error {
// 	if err := DB.Create(account).Error; err != nil {
// 		return err
// 	}
// 	return nil
// }

// // Kiểm tra số điện thoại đã tồn tại trong hệ thống cho Login
// func FindAccountByPhone(soDienThoai string) (*models.Accounts, error) {
// 	var account models.Accounts
// 	result := DB.Where("SoDienThoai = ?", soDienThoai).First(&account)
// 	if result.Error != nil {
// 		return nil, result.Error
// 	}
// 	return &account, nil
// }

// // Kiểm tra số lần đăng nhập sai trong ngày
// func CheckLoginAttempts(soDienThoai string) (*models.LoginAttempt, error) {
// 	var loginAttempt models.LoginAttempt
// 	result := DB.Where("SoDienThoai = ? AND DATE(Ngay) = CURDATE()", soDienThoai).First(&loginAttempt)
// 	if result.Error != nil {
// 		return nil, result.Error
// 	}
// 	return &loginAttempt, nil
// }

// // Cập nhật số lần đăng nhập sai
// func UpdateLoginAttempt(attempt *models.LoginAttempt) error {
// 	if err := DB.Save(attempt).Error; err != nil {
// 		return err
// 	}
// 	return nil
// }

// // Lưu thiết bị vào bảng Devices
// func SaveDevice(device *models.Devices) error {
// 	if err := DB.Create(device).Error; err != nil {
// 		return err
// 	}
// 	return nil
// }

// // Kiểm tra OTP
// func CheckOTP(soDienThoai, otpCode string) (*models.OTPCode, error) {
// 	var otp models.OTPCode
// 	result := DB.Where("SoDienThoai = ? AND OTP_Code = ?", soDienThoai, otpCode).First(&otp)
// 	if result.Error != nil {
// 		return nil, result.Error
// 	}
// 	return &otp, nil
// }

// // Cập nhật OTP
// func UpdateOTP(otp *models.OTPCode) error {
// 	if err := DB.Save(otp).Error; err != nil {
// 		return err
// 	}
// 	return nil
// }

// // Tạo mới OTP
// func CreateOTP(otp *models.OTPCode) error {
// 	if err := DB.Create(otp).Error; err != nil {
// 		return err
// 	}
// 	return nil
// }
