# Simple Golang Clean Architecture Template

Template Golang đơn giản cho người mới bắt đầu, sử dụng Clean Architecture với 3 tầng: **Handler → Usecase → Repository**.

## 📋 Mục Lục

- [Giới Thiệu](#giới-thiệu)
- [Kiến Trúc](#kiến-trúc)
- [Công Nghệ Sử Dụng](#công-nghệ-sử-dụng)
- [Cấu Trúc Thư Mục](#cấu-trúc-thư-mục)
- [Yêu Cầu](#yêu-cầu)
- [Cài Đặt](#cài-đặt)
- [Sử Dụng](#sử-dụng)
- [API Documentation](#api-documentation)
- [Database Migration](#database-migration)

## 🎯 Giới Thiệu

Đây là một template Golang đơn giản, dễ hiểu, được thiết kế đặc biệt cho người mới bắt đầu.

**Đặc điểm:**
- ✅ Clean Architecture với 3 tầng rõ ràng
- ✅ Không sử dụng interface để giảm độ phức tạp
- ✅ Fiber framework (tương tự Express.js)
- ✅ MySQL database với thư viện goqu (nhẹ, dễ hiểu)
- ✅ Docker & Docker Compose ready
- ✅ Code comments bằng tiếng Việt
- ✅ Ví dụ CRUD hoàn chỉnh

## 🏗️ Kiến Trúc

Template này sử dụng Clean Architecture với 3 tầng:

```
┌─────────────────────────────────────────────────┐
│                   Handler Layer                  │
│  - Nhận HTTP requests                           │
│  - Parse request body/params                    │
│  - Gọi Usecase                                  │
│  - Trả về HTTP responses                        │
└─────────────────┬───────────────────────────────┘
                  │ Gọi trực tiếp (no interface)
┌─────────────────▼───────────────────────────────┐
│                  Usecase Layer                   │
│  - Business logic                               │
│  - Validation                                   │
│  - Data transformation                          │
│  - Gọi Repository                               │
└─────────────────┬───────────────────────────────┘
                  │ Gọi trực tiếp (no interface)
┌─────────────────▼───────────────────────────────┐
│                 Repository Layer                 │
│  - Database operations (CRUD)                   │
│  - Query building với goqu                      │
│  - Trả về data                                  │
└─────────────────────────────────────────────────┘
```

### Luồng Xử Lý Request

1. **Client** gửi HTTP request → **Handler**
2. **Handler** parse request → gọi **Usecase**
3. **Usecase** xử lý logic → gọi **Repository**
4. **Repository** query database → trả về data
5. **Usecase** xử lý data → trả về cho **Handler**
6. **Handler** format response → trả về **Client**

## 🛠️ Công Nghệ Sử Dụng

- **Language**: Go 1.21+
- **Web Framework**: [Fiber v2](https://gofiber.io/) - tương tự Express.js
- **Database**: MySQL 8.0
- **Query Builder**: [goqu](https://github.com/doug-martin/goqu) - nhẹ và dễ hiểu
- **Config**: [godotenv](https://github.com/joho/godotenv) - load .env file
- **Containerization**: Docker & Docker Compose

## 📁 Cấu Trúc Thư Mục

```
simple-golang-code/
├── cmd/
│   └── api/
│       └── main.go                 # Entry point của ứng dụng
├── internal/                       # Code nội bộ (không export)
│   ├── config/
│   │   └── config.go              # Quản lý cấu hình
│   ├── database/
│   │   └── mysql.go               # Kết nối database
│   ├── handler/
│   │   └── user_handler.go        # HTTP handlers
│   ├── usecase/
│   │   └── user_usecase.go        # Business logic
│   ├── repository/
│   │   └── user_repository.go     # Database operations
│   ├── model/
│   │   └── user.go                # Data models
│   └── middleware/
│       ├── logger.go              # Request logging
│       └── error.go               # Error handling
├── pkg/                            # Code có thể export
│   └── response/
│       └── response.go            # Standard API responses
├── migrations/
│   └── 001_create_users_table.sql # Database migrations
├── .env.example                    # Template cho environment variables
├── .gitignore
├── docker-compose.yml
├── Dockerfile
├── Makefile                        # Common commands
├── go.mod
├── go.sum
└── README.md
```

## ✅ Yêu Cầu

- Go 1.21 hoặc mới hơn
- Docker & Docker Compose (nếu chạy với Docker)
- MySQL 8.0 (nếu chạy local không dùng Docker)
- Make (optional, để sử dụng Makefile)

## 🚀 Cài Đặt

### 1. Clone hoặc sử dụng template này

```bash
cd simple-template
```

### 2. Khởi tạo project

```bash
# Sử dụng Makefile
make init

# Hoặc thủ công
cp .env.example .env
go mod download
```

### 3. Cấu hình environment variables

Chỉnh sửa file `.env` theo môi trường của bạn:

```env
# Server Configuration
SERVER_HOST=0.0.0.0
SERVER_PORT=8080
APP_ENV=development

# Database Configuration
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=secret
DB_NAME=simple_golang_db
```

## 🎮 Sử Dụng

### Chạy với Docker (Khuyến nghị cho người mới)

```bash
# Khởi động tất cả services (MySQL + API)
make docker-up

# Hoặc
docker-compose up -d

# Xem logs
make docker-logs

# Dừng services
make docker-down
```

API sẽ chạy tại: `http://localhost:8080`

### Chạy Local (Không dùng Docker)

**Bước 1**: Cài đặt và khởi động MySQL

**Bước 2**: Chạy migration để tạo database và tables

```bash
make migrate
```

**Bước 3**: Chạy ứng dụng

```bash
# Sử dụng Makefile
make run

# Hoặc trực tiếp
go run cmd/api/main.go
```

### Build Binary

```bash
# Build binary
make build

# Chạy binary
./bin/simple-golang-api
```

## 📚 API Documentation

### Base URL

```
http://localhost:8080
```

### Health Check

```http
GET /health
```

**Response:**
```json
{
  "success": true,
  "message": "Service is healthy",
  "data": {
    "status": "ok",
    "database": "connected"
  }
}
```

### User Endpoints

#### 1. Tạo User Mới

```http
POST /api/v1/users
Content-Type: application/json

{
  "name": "John Doe",
  "email": "john@example.com"
}
```

**Response (201 Created):**
```json
{
  "success": true,
  "message": "User created successfully",
  "data": {
    "id": 1,
    "name": "John Doe",
    "email": "john@example.com",
    "created_at": "2024-01-01T10:00:00Z",
    "updated_at": "2024-01-01T10:00:00Z"
  }
}
```

#### 2. Lấy Danh Sách Users

```http
GET /api/v1/users
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Users retrieved successfully",
  "data": [
    {
      "id": 1,
      "name": "John Doe",
      "email": "john@example.com",
      "created_at": "2024-01-01T10:00:00Z",
      "updated_at": "2024-01-01T10:00:00Z"
    }
  ]
}
```

#### 3. Lấy User Theo ID

```http
GET /api/v1/users/:id
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "User retrieved successfully",
  "data": {
    "id": 1,
    "name": "John Doe",
    "email": "john@example.com",
    "created_at": "2024-01-01T10:00:00Z",
    "updated_at": "2024-01-01T10:00:00Z"
  }
}
```

#### 4. Cập Nhật User

```http
PUT /api/v1/users/:id
Content-Type: application/json

{
  "name": "John Updated",
  "email": "john.updated@example.com"
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "User updated successfully",
  "data": {
    "id": 1,
    "name": "John Updated",
    "email": "john.updated@example.com",
    "created_at": "2024-01-01T10:00:00Z",
    "updated_at": "2024-01-01T10:30:00Z"
  }
}
```

#### 5. Xóa User

```http
DELETE /api/v1/users/:id
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "User deleted successfully"
}
```

### Error Response Format

```json
{
  "success": false,
  "message": "Error message",
  "error": "Detailed error description"
}
```

## 🗄️ Database Migration

### Chạy Migration

```bash
# Sử dụng Makefile
make migrate

# Hoặc thủ công với MySQL client
mysql -h localhost -u root -p simple_golang_db < migrations/001_create_users_table.sql
```

### Tạo Migration Mới

Tạo file SQL mới trong thư mục `migrations/`:

```sql
-- migrations/002_add_new_table.sql
CREATE TABLE IF NOT EXISTS new_table (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    -- columns here
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

## 📖 Hướng Dẫn Cho Người Mới

### Thêm Feature Mới (Ví dụ: Product)

**Bước 1**: Tạo model

```go
// internal/model/product.go
type Product struct {
    ID    int64  `db:"id" json:"id"`
    Name  string `db:"name" json:"name"`
    Price int64  `db:"price" json:"price"`
}
```

**Bước 2**: Tạo repository

```go
// internal/repository/product_repository.go
type ProductRepository struct {
    db *database.DB
}

func NewProductRepository(db *database.DB) *ProductRepository {
    return &ProductRepository{db: db}
}

func (r *ProductRepository) Create(ctx context.Context, product *model.Product) error {
    // Implementation
}
```

**Bước 3**: Tạo usecase

```go
// internal/usecase/product_usecase.go
type ProductUsecase struct {
    productRepo *repository.ProductRepository
}

func NewProductUsecase(productRepo *repository.ProductRepository) *ProductUsecase {
    return &ProductUsecase{productRepo: productRepo}
}

func (u *ProductUsecase) CreateProduct(ctx context.Context, req *model.CreateProductRequest) (*model.Product, error) {
    // Business logic
    return u.productRepo.Create(ctx, product)
}
```

**Bước 4**: Tạo handler

```go
// internal/handler/product_handler.go
type ProductHandler struct {
    productUsecase *usecase.ProductUsecase
}

func NewProductHandler(productUsecase *usecase.ProductUsecase) *ProductHandler {
    return &ProductHandler{productUsecase: productUsecase}
}

func (h *ProductHandler) CreateProduct(c *fiber.Ctx) error {
    // Parse request và gọi usecase
}
```

**Bước 5**: Đăng ký routes trong `main.go`

```go
productRepo := repository.NewProductRepository(db)
productUsecase := usecase.NewProductUsecase(productRepo)
productHandler := handler.NewProductHandler(productUsecase)

products := api.Group("/products")
products.Post("/", productHandler.CreateProduct)
```

## 🤝 Contributing

Mọi đóng góp đều được chào đón! Hãy tạo issue hoặc pull request.

## 📝 License

MIT License - Tự do sử dụng cho mục đích học tập và thương mại.

## 🙋 Hỗ Trợ

Nếu bạn gặp vấn đề hoặc có câu hỏi, hãy tạo issue trên GitHub.

---

**Happy Coding! 🚀**
