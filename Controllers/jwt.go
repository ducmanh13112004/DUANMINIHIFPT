package controllers

import (
	"MiniHIFPT/models"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

var jwtSecretKey = []byte("your_secret_key") // Khai báo khóa bí mật cho việc ký và xác thực token

// Hàm tạo JWT token
func generateJWT(account *models.Accounts) (string, error) {
	claims := jwt.MapClaims{
		"id":          account.ID,
		"soDienThoai": account.SoDienThoai,
		"exp":         time.Now().Add(time.Hour).Unix(), // Token hết hạn trong 1 giờ

	}

	// Tạo token bằng cách ký với secret key
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// sử dụng HMAC SHA-256 để ký token ( thuật toán mã hóa sử dụng khóa bí mật để ký và kiểm tra tính toàn vẹn).
	// Ký token bằng jwtSecretKey
	signedToken, err := token.SignedString(jwtSecretKey)
	if err != nil {
		return "", err // Nếu có lỗi trong quá trình ký token
	}

	return signedToken, nil // Trả về token đã ký
}
