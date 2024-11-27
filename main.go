package main

import (
    "log"
    "github.com/gofiber/fiber/v2" // Đảm bảo import Fiber
    "MiniHIFPT/database"
    "MiniHIFPT/routes"
)

func main() {
    // Kết nối cơ sở dữ liệu
    database.ConnectDB()
    log.Println("Kết nối cơ sở dữ liệu thành công!")

    // Khởi tạo ứng dụng Fiber
    app := fiber.New()

    // Định nghĩa các Routes
    routes.Setup(app)

    // Chạy ứng dụng
    log.Fatal(app.Listen(":3000"))
}
