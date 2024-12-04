package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"strings"
)

var jwtSecretKey = []byte("your_secret_key")

// wtSecretKey là khóa bí mật được sử dụng để ký và xác thực JWT.

func Authenticate(c *fiber.Ctx) error {
	// Lấy token từ header Authorization
	tokenString := c.Get("Authorization")
	tokenString = strings.Replace(tokenString, "Bearer ", "", 1)

	// Kiểm tra nếu không có token
	if tokenString == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Token không hợp lệ",
		})
	}

	// Parse token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Kiểm tra phương thức ký
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fiber.ErrUnauthorized
		}
		return jwtSecretKey, nil // Sử dụng jwtSecretKey để xác thực
	})
	//SigningMethodHMAC là một phương thức ký trong JWT (JSON Web Token) được sử dụng trong quá trình ký và xác thực token
	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Token không hợp lệ: ",
			//+ err.Error()
		})
	}

	return c.Next()
}
