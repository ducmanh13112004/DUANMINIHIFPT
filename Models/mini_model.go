package models

import (
	"time"
)

// ----------------------- Khách hàng -----------------------
type Customer struct {
	ID           string     `gorm:"primaryKey;type:uuid;default:uuid_generate_v4();column:id_uuid"` // UUID tự động
	SoDienThoai  string     `json:"soDienThoai" gorm:"column:SoDienThoai"`                          // Số điện thoại khách hàng
	TenKhachHang string     `json:"tenKhachHang" gorm:"column:TenKhachHang"`                        // Tên khách hàng
	GioiTinh     string     `json:"gioiTinh" gorm:"column:GioiTinh"`                                // Giới tính khách hàng
	NgaySinh     *time.Time `gorm:"type:date;column:NgaySinh" json:"ngaySinh"`

	Email         string     `json:"email" gorm:"column:Email"`                                          // Email khách hàng
	LoaiKhachHang string     `gorm:"type:char(1);default:'T';column:LoaiKhachHang" json:"loaiKhachHang"` // Loại khách hàng: Tiềm năng (T) hoặc Sử dụng dịch vụ (S)
	Contracts     []Contract `gorm:"many2many:customer_contracts;" json:"contracts"`                     // Quan hệ nhiều-nhiều với Contract
}

// Chỉ định tên bảng
func (Customer) TableName() string {
	return "Customer" // Tên bảng thực tế trong cơ sở dữ liệu
}

// ----------------------- Hợp đồng -----------------------
type Contract struct {
	ID           string `gorm:"primaryKey;type:uuid;default:uuid_generate_v4();column:id_uuid"` // UUID tự động
	TenKhachHang string `gorm:"column:TenKhachHang;not null"`
	DiaChi       string `gorm:"column:DiaChi;not null"`
	MaTinh       string `gorm:"column:MaTinh;not null"`
	MaQuanHuyen  string `gorm:"column:MaQuanHuyen;not null"`
	MaPhuongXa   string `gorm:"column:MaPhuongXa;not null"`
	MaDuong      string `gorm:"column:MaDuong;null"`
	SoNha        string `gorm:"column:SoNha;null"`
}

// Chỉ định tên bảng trong MySQL
func (Contract) TableName() string {
	return "contractt" // Tên bảng trong cơ sở dữ liệu
}

// ----------------------- Bảng trung gian -----------------------
type Customer_Contractt struct {
	ID          string `gorm:"primaryKey;type:uuid;default:uuid_generate_v4();column:id_uuid"` // UUID tự động
	SoDienThoai string `json:"soDienThoai" gorm:"unique;not null;column:SoDienThoai"`          // Số điện thoại (duy nhất)
	HopDongID   string `gorm:"index;not null;column:HopDongID"`                                // ID hợp đồng
}

// Chỉ định tên bảng trong MySQL
func (Customer_Contractt) TableName() string {
	return "customer_contractt" // Tên bảng trong cơ sở dữ liệu
}

// ----------------------- Tài khoản người dùng -----------------------
type Accounts struct {
	ID                 string     `gorm:"primaryKey;type:uuid;default:uuid_generate_v4();column:id_uuid"`     // UUID tự động
	SoDienThoai        string     `json:"soDienThoai" gorm:"unique;not null;column:SoDienThoai"`              // Số điện thoại (duy nhất)
	MatKhau            string     `json:"matKhau" gorm:"not null;column:MatKhau"`                             // Mật khẩu (bắt buộc)
	NgayTao            time.Time  `json:"ngayTao" gorm:"autoCreateTime;column:NgayTao"`                       // Ngày tạo tài khoản
	NgayCapNhat        time.Time  `json:"ngayCapNhat" gorm:"autoUpdateTime;column:NgayCapNhat"`               // Ngày cập nhật tài khoản
	LanDangNhapGanNhat *time.Time `json:"lanDangNhapGanNhat" gorm:"autoUpdateTime;column:LanDangNhapGanNhat"` // Lần đăng nhập gần nhất
}

// bảng LoginAttempts
type LoginAttempt struct {
	ID          string    `gorm:"primaryKey;type:uuid;default:uuid_generate_v4();column:id_uuid"`
	SoDienThoai string    `gorm:"column:SoDienThoai;not null"`
	SoLanSai    int       `gorm:"column:SoLanSai;default:0"`
	Ngay        time.Time `gorm:"column:Ngay;default:current_timestamp"`
	KhoiPhuc    bool      `gorm:"column:KhoiPhuc;default:false"`
}

func (LoginAttempt) TableName() string {
	return "login_attempts"
}

type OTPCode struct {
	ID          string    `gorm:"primaryKey;type:uuid;default:uuid_generate_v4();column:id_uuid"` // UUID tự động
	SoDienThoai string    `json:"soDienThoai" gorm:"not null;index;column:SoDienThoai"`           // Số điện thoại (có index)
	OTP_Code    string    `json:"otpCode" gorm:"not null;column:OTP_Code"`                        // Mã OTP (bắt buộc)
	NgayTao     time.Time `json:"ngayTao" gorm:"autoCreateTime;column:NgayTao"`                   // Thời gian tạo OTP
	HetHan      time.Time `json:"hetHan" gorm:"not null;column:HetHan"`                           // Thời gian hết hạn OTP
	DaXacThuc   bool      `json:"daXacThuc" gorm:"default:false;column:DaXacThuc"`                // Trạng thái xác thực OTP (mặc định false)
}

type Devices struct {
	ID             string    `gorm:"primaryKey;type:uuid;default:uuid_generate_v4();column:id_uuid"` // UUID tự động
	SoDienThoai    string    `json:"soDienThoai" gorm:"not null;index;column:SoDienThoai"`           // Số điện thoại (có index)
	LanDungGanNhat time.Time `json:"lanDungGanNhat" gorm:"autoUpdateTime;column:LanDungGanNhat"`     // Lần sử dụng thiết bị gần nhất
}

type Log struct {
	ID        string    `gorm:"primaryKey;type:uuid;default:uuid_generate_v4();column:id_uuid"`
	Action    string    `gorm:"not null;column:Action"`          // Hành động
	Details   string    `gorm:"not null;column:Details"`         // Chi tiết
	CreatedAt time.Time `gorm:"autoCreateTime;column:CreatedAt"` // Thời gian thực hiện
}

func (Log) TableName() string {
	return "logs"
}
