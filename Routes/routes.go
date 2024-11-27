package routes

import (
	"MiniHIFPT/controllers"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	// Route lấy thông tin khách hàng
	app.Get("/customers", controllers.GetCustomers)

	// Route thêm khách hàng mới
	app.Post("/customers", controllers.CreateCustomers)

	// Route lấy thông tin hợp đồng

	app.Get("/contracts", controllers.GetContracts)
	// Định nghĩa route để tạo hợp đồng mới
	app.Post("/contracts", controllers.CreateContract)
	// Sửa hợp đồng
	app.Put("/contracts/:id", controllers.UpdateContract)
	// Xóa hợp đồng
	app.Delete("/contracts/:id", controllers.DeleteContract)

	//ctm_tract
	app.Get("/ctmtract", controllers.Getctm_tracts)
	app.Post("/Createctmtract", controllers.Createctm_tracts)
	// Route đăng ký tài khoản
	app.Post("/register", controllers.Register) // Đảm bảo là POST để tạo tài khoản mới

	// Route đăng nhập
	app.Post("/login", controllers.Login) // Đảm bảo là POST để đăng nhập

	// Route xác thực OTP
	app.Post("/otp", controllers.VerifyOTP) // Đảm bảo là POST để xác thực OTP
}
