# ACGGoods Go 从 0 到 1 重写实战指南

这是一套“从空目录手动重写简化版 ACGGoods 后端”的实战教材，不是单纯的源码阅读说明。

你会新建一个训练项目，从零完成下面的业务闭环：

~~~~
用户 User
  -> 店铺 Store
      -> 商品 Product
~~~~

每个阶段都包含学习目标、输入命令、需要创建的文件、预期输出、错误排查、验收标准和 Git 提交点。

训练项目只做高频后端能力：

- Go 基础语法和工程目录
- Gin HTTP 服务
- Docker Compose 启动 MySQL 和 Redis
- SQL 迁移和 GORM CRUD
- 用户注册、登录、退出
- Token Cookie 鉴权
- JWT 鉴权对比
- 店铺和商品归属权限
- 参数校验
- 统一响应和业务错误
- Zap 结构化日志
- request_id、trace_id 和 recovery
- Redis 缓存
- 单元测试、HTTP 测试、集成测试
- Delve + VS Code 调试
- 编译、排错和日常开发流程

暂不实现支付、IM、Meilisearch、订单、购物车、Cron、异步队列、后台 RBAC、国际化和 Wooacry。

> 所有练习代码都写入新项目目录。不要在 acggoods-go 中创建 Go 文件、配置文件或迁移文件。

---

## 一、最终训练项目

### 1. 最终目录

~~~~
acggoods-go-practice/
├── cmd/server/main.go
├── config/local.yaml
├── migrations/001_init.sql
├── internal/
│   ├── app/
│   ├── config/
│   ├── dto/
│   ├── handler/
│   ├── logger/
│   ├── middleware/
│   ├── model/
│   ├── repository/
│   ├── response/
│   └── service/
├── tests/integration/
├── .env.example
├── .gitignore
├── docker-compose.yml
├── go.mod
└── README.md
~~~~

第一阶段使用手动依赖注入。基础功能跑通后再理解 Wire，避免一开始被生成代码挡住。

### 2. 最终接口

统一响应：

~~~~
{
  "code": 0,
  "message": "success",
  "data": {},
  "trace_id": "01J...",
  "request_id": "01J..."
}
~~~~

认证接口：

| 方法 | 路径 | 说明 |
| --- | --- | --- |
| POST | /api/auth/register | 注册 |
| POST | /api/auth/login | 登录 |
| POST | /api/auth/logout | 退出 |
| GET | /api/me | 当前用户 |

店铺接口：

| 方法 | 路径 | 说明 |
| --- | --- | --- |
| POST | /api/stores | 创建店铺 |
| GET | /api/stores | 当前用户店铺列表 |
| GET | /api/stores/:id | 店铺详情 |
| PUT | /api/stores/:id | 修改自己的店铺 |
| DELETE | /api/stores/:id | 删除自己的店铺 |

商品接口：

| 方法 | 路径 | 说明 |
| --- | --- | --- |
| POST | /api/stores/:store_id/products | 创建商品 |
| GET | /api/stores/:store_id/products | 商品分页 |
| GET | /api/products/:id | 商品详情 |
| PUT | /api/products/:id | 修改商品 |
| DELETE | /api/products/:id | 删除商品 |

系统接口：

| 方法 | 路径 | 说明 |
| --- | --- | --- |
| GET | /health | 健康检查 |
| GET | /api/hello | 最小业务接口 |

### 3. 最终验收

完成后可以：

~~~~
docker compose up -d
go run ./cmd/server
curl http://localhost:8080/health
~~~~

然后完成：

1. 注册用户 Alice。
2. 登录 Alice。
3. 创建店铺。
4. 创建和修改商品。
5. 注册用户 Bob。
6. 验证 Bob 无法修改 Alice 的店铺和商品。
7. 退出 Alice，验证 Cookie 失效。
8. 切换到 JWT 模式，重复登录和鉴权。
9. 观察店铺详情的 Redis 缓存命中和失效。
10. 使用 Delve 在 Handler、Service、Repository 设置断点。

---

## 二、环境安装

## 第 0 阶段：Go、Docker、VS Code

### 0.1 安装 Go 1.23.1

当前项目 go.mod 要求 Go 1.23，并声明 toolchain 1.23.1。训练项目固定使用 Go 1.23.1。

官方下载页：

https://go.dev/dl/

选择对应安装包：

- macOS Apple Silicon：go1.23.1.darwin-arm64.pkg
- macOS Intel：go1.23.1.darwin-amd64.pkg
- Windows AMD64：go1.23.1.windows-amd64.msi
- Linux AMD64：go1.23.1.linux-amd64.tar.gz

安装后打开新终端：

~~~~
go version
which go
go env GOPATH GOMODCACHE
~~~~

预期：

~~~~
go version go1.23.1 darwin/arm64
~~~~

Linux 使用压缩包安装时确认 PATH：

~~~~
export PATH=/usr/local/go/bin:$PATH
go version
~~~~

不要把 GOROOT 写死在项目中。

### 0.2 安装 Docker Desktop

官方下载页：

https://www.docker.com/products/docker-desktop/

启动 Docker Desktop 后执行：

~~~~
docker version
docker compose version
docker run --rm hello-world
~~~~

如果 docker version 只有 Client 没有 Server，先启动 Docker Desktop。

### 0.3 安装 VS Code 和 Go 扩展

安装：

- Visual Studio Code
- VS Code 官方 Go 扩展
- Delve 在调试阶段安装

---

## 三、Go 基础和空项目

## 第 1 阶段：从空目录写第一个程序

### 1.1 创建项目

在 acggoods-go 同级目录执行：

~~~~
cd /Users/xiaolei/Desktop/overseas
mkdir acggoods-go-practice
cd acggoods-go-practice

go mod init example.com/acggoods-go-practice
mkdir -p cmd/hello internal/catalog
~~~~

预期：

~~~~
go: creating new go.mod: module example.com/acggoods-go-practice
~~~~

查看：

~~~~
cat go.mod
find . -maxdepth 3 -type d | sort
~~~~

### 1.2 第一个 main

创建 cmd/hello/main.go：

~~~~
package main

import "fmt"

func main() {
	fmt.Println("ACGGoods practice started")
}
~~~~

运行：

~~~~
go run ./cmd/hello
mkdir -p bin
go build -o bin/hello ./cmd/hello
./bin/hello
~~~~

预期两次输出：

~~~~
ACGGoods practice started
~~~~

格式和静态检查：

~~~~
gofmt -w cmd/hello/main.go
go vet ./cmd/hello
~~~~

需要理解：

- package main 是可执行包。
- func main 是入口。
- go run 适合开发运行。
- go build 生成二进制。
- gofmt 是固定格式化工具。
- go vet 发现部分可疑代码。

### 1.3 struct、方法和 error

创建 internal/catalog/catalog.go：

~~~~
package catalog

import (
	"errors"
	"fmt"
)

var ErrInvalidPrice = errors.New("price must be greater than zero")

type Product struct {
	Name  string
	Price int64
}

func NewProduct(name string, price int64) (*Product, error) {
	if name == "" {
		return nil, errors.New("name is required")
	}
	if price <= 0 {
		return nil, ErrInvalidPrice
	}

	return &Product{Name: name, Price: price}, nil
}

func (p Product) DisplayPrice() string {
	return fmt.Sprintf("$%d.%02d", p.Price/100, p.Price%100)
}
~~~~

创建 internal/catalog/catalog_test.go：

~~~~
package catalog

import (
	"errors"
	"testing"
)

func TestNewProduct(t *testing.T) {
	product, err := NewProduct("Acrylic Stand", 1299)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if product.DisplayPrice() != "$12.99" {
		t.Fatalf("unexpected price: %s", product.DisplayPrice())
	}
}

func TestNewProductRejectsInvalidPrice(t *testing.T) {
	_, err := NewProduct("Acrylic Stand", 0)
	if !errors.Is(err, ErrInvalidPrice) {
		t.Fatalf("expected ErrInvalidPrice, got %v", err)
	}
}
~~~~

运行：

~~~~
go test ./internal/catalog -v
go test ./...
~~~~

预期：

~~~~
PASS
ok  	example.com/acggoods-go-practice/internal/catalog
~~~~

### 1.4 本阶段验收

你能回答：

- struct、slice、map 的用途分别是什么？
- 为什么函数返回指针和 error？
- 为什么用 errors.Is 判断业务错误？
- 值接收者和指针接收者有什么区别？
- go test ./... 中的 ./... 是什么意思？

提交：

~~~~
git init
printf "bin/\n.env\nconfig/local.yaml\nlogs/\n" > .gitignore
git add .
git commit -m "chore: bootstrap go practice project"
~~~~

---

## 四、HTTP 服务骨架

## 第 2 阶段：Gin、统一响应、请求 ID

### 2.1 安装 Gin 和 UUID

~~~~
go get github.com/gin-gonic/gin@v1.9.1
go get github.com/google/uuid@v1.6.0
go mod tidy
mkdir -p internal/response internal/middleware
~~~~

### 2.2 统一响应

创建 internal/response/response.go：

~~~~
package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Envelope struct {
	Code      int    `json:"code"`
	Message   string `json:"message"`
	Data      any    `json:"data"`
	TraceID   string `json:"trace_id"`
	RequestID string `json:"request_id"`
}

func Success(c *gin.Context, data any) {
	c.JSON(http.StatusOK, Envelope{
		Code:      0,
		Message:   "success",
		Data:      data,
		TraceID:   GetTraceID(c),
		RequestID: GetRequestID(c),
	})
}

func Fail(c *gin.Context, status int, code int, message string, data any) {
	c.JSON(status, Envelope{
		Code:      code,
		Message:   message,
		Data:      data,
		TraceID:   GetTraceID(c),
		RequestID: GetRequestID(c),
	})
}

func GetRequestID(c *gin.Context) string {
	value, _ := c.Get("request_id")
	id, _ := value.(string)
	return id
}

func GetTraceID(c *gin.Context) string {
	value, _ := c.Get("trace_id")
	id, _ := value.(string)
	return id
}
~~~~

### 2.3 请求 ID 和 trace ID

创建 internal/middleware/request_id.go：

~~~~
package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.GetHeader("X-Request-ID")
		if id == "" {
			id = uuid.NewString()
		}
		c.Set("request_id", id)
		c.Header("X-Request-ID", id)
		c.Next()
	}
}

func TraceID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := uuid.NewString()
		c.Set("trace_id", id)
		c.Header("X-Trace-ID", id)
		c.Next()
	}
}
~~~~

### 2.4 HTTP 入口

创建 cmd/server/main.go：

~~~~
package main

import (
	"net/http"

	"example.com/acggoods-go-practice/internal/middleware"
	"example.com/acggoods-go-practice/internal/response"
	"github.com/gin-gonic/gin"
)

func main() {
	engine := gin.New()
	engine.Use(
		middleware.RequestID(),
		middleware.TraceID(),
		gin.Logger(),
		gin.Recovery(),
	)

	engine.GET("/health", func(c *gin.Context) {
		response.Success(c, gin.H{"status": "ok"})
	})

	engine.GET("/api/hello", func(c *gin.Context) {
		response.Success(c, gin.H{
			"message": "hello from acggoods practice",
		})
	})

	server := &http.Server{
		Addr:    ":8080",
		Handler: engine,
	}

	if err := server.ListenAndServe(); err != nil &&
		err != http.ErrServerClosed {
		panic(err)
	}
}
~~~~

运行：

~~~~
go run ./cmd/server
~~~~

另开终端：

~~~~
curl -i http://localhost:8080/health
curl -i http://localhost:8080/api/hello
~~~~

预期：

~~~~
{
  "code": 0,
  "message": "success",
  "data": {
    "status": "ok"
  },
  "trace_id": "...",
  "request_id": "..."
}
~~~~

响应头有 X-Request-ID 和 X-Trace-ID。

### 2.5 优雅关闭

把 server 启动逻辑改为：

~~~~
启动 server goroutine
  -> signal.Notify 等待 SIGINT/SIGTERM
  -> context.WithTimeout(..., 5*time.Second)
  -> server.Shutdown(ctx)
  -> 进程退出
~~~~

验收：

1. 启动服务。
2. curl 请求 /health。
3. 按 Ctrl-C。
4. 服务正常退出，没有强制终止堆栈。

### 2.6 当前项目对照

完成后阅读：

- internal/router/base.go
- internal/response/response.go
- internal/middleware/logger.go
- internal/middleware/recovery.go
- cmd/cli/server.go

不要直接复制当前项目的敏感请求 body 日志。

提交：

~~~~
go fmt ./...
go test ./...
go build ./cmd/server
git add .
git commit -m "feat: add gin http foundation"
~~~~

---

## 五、Docker MySQL 和 Redis

## 第 3 阶段：本地基础设施

### 3.1 docker-compose.yml

创建项目根目录的 docker-compose.yml：

~~~~
services:
  mysql:
    image: mysql:8.0
    container_name: acggoods-practice-mysql
    restart: unless-stopped
    environment:
      MYSQL_ROOT_PASSWORD: root_password
      MYSQL_DATABASE: acggoods_practice
      MYSQL_USER: app
      MYSQL_PASSWORD: app_password
    ports:
      - "${MYSQL_PORT:-13306}:3306"
    volumes:
      - mysql_data:/var/lib/mysql
      - ./migrations:/docker-entrypoint-initdb.d:ro
    healthcheck:
      test: ["CMD", "mysqladmin", "ping", "-h", "localhost", "-u", "root", "-proot_password"]
      interval: 5s
      timeout: 5s
      retries: 20

  redis:
    image: redis:7-alpine
    container_name: acggoods-practice-redis
    restart: unless-stopped
    ports:
      - "${REDIS_PORT:-16379}:6379"
    volumes:
      - redis_data:/data
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 5s
      timeout: 3s
      retries: 20

volumes:
  mysql_data:
  redis_data:
~~~~

先创建空目录：

~~~~
mkdir -p migrations
docker compose up -d
docker compose ps
~~~~

检查：

~~~~
docker compose exec mysql mysql -uapp -papp_password acggoods_practice \
  -e "SELECT DATABASE(), VERSION();"

redis-cli -h 127.0.0.1 -p 16379 ping
~~~~

预期：

- MySQL 返回数据库名和版本。
- Redis 返回 PONG。
- 两个服务都是 running/healthy。

端口冲突：

~~~~
lsof -nP -iTCP:13306 -sTCP:LISTEN
lsof -nP -iTCP:16379 -sTCP:LISTEN
~~~~

如果需要改端口：

~~~~
"${MYSQL_PORT:-13306}:3306"
"${REDIS_PORT:-16379}:6379"
~~~~

再同步修改应用配置。

彻底重建本训练项目的数据：

~~~~
docker compose down -v
docker compose up -d
~~~~

只对训练项目执行 down -v，它会删除本训练项目的数据卷。

提交：

~~~~
git add docker-compose.yml migrations
git commit -m "chore: add docker mysql and redis"
~~~~

---

## 六、配置、迁移和 GORM

## 第 4 阶段：Viper、MySQL、Redis

### 4.1 安装依赖

~~~~
go get github.com/spf13/viper@v1.17.0
go get gorm.io/gorm@v1.25.5
go get gorm.io/driver/mysql@v1.5.2
go get github.com/redis/go-redis/v9@v9.6.1
go get go.uber.org/zap@v1.26.0
go get golang.org/x/crypto@v0.39.0
go mod tidy
~~~~

### 4.2 本地配置

创建 config/local.yaml：

~~~~
app:
  name: acggoods-practice
  env: development
  port: 8080

database:
  dsn: "app:app_password@tcp(127.0.0.1:13306)/acggoods_practice?charset=utf8mb4&parseTime=True&loc=Local"
  max_open_conns: 20
  max_idle_conns: 5

redis:
  addr: "127.0.0.1:16379"
  password: ""
  db: 0

logger:
  level: debug
  format: json

auth:
  mode: cookie
  cookie_name: acggoods_practice_auth
  jwt_secret: local-only-secret-change-me
  jwt_expire_minutes: 60
~~~~

将它加入 .gitignore：

~~~~
printf "config/local.yaml\n.env\nbin/\nlogs/\n" >> .gitignore
~~~~

### 4.3 Viper 配置加载

创建 internal/config/config.go：

~~~~
package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	App      AppConfig      `mapstructure:"app"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	Logger   LoggerConfig   `mapstructure:"logger"`
	Auth     AuthConfig     `mapstructure:"auth"`
}

type AppConfig struct {
	Name string `mapstructure:"name"`
	Env  string `mapstructure:"env"`
	Port int    `mapstructure:"port"`
}

type DatabaseConfig struct {
	DSN          string `mapstructure:"dsn"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
}

type RedisConfig struct {
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type LoggerConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
}

type AuthConfig struct {
	Mode             string `mapstructure:"mode"`
	CookieName       string `mapstructure:"cookie_name"`
	JWTSecret        string `mapstructure:"jwt_secret"`
	JWTExpireMinutes int    `mapstructure:"jwt_expire_minutes"`
}

func Load() (*Config, error) {
	v := viper.New()
	v.SetConfigFile("config/local.yaml")
	v.SetConfigType("yaml")
	v.SetEnvPrefix("ACG")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("decode config: %w", err)
	}

	if cfg.App.Port == 0 {
		cfg.App.Port = 8080
	}
	if cfg.Auth.Mode == "" {
		cfg.Auth.Mode = "cookie"
	}

	return &cfg, nil
}
~~~~

测试配置缺失：

~~~~
mv config/local.yaml config/local.yaml.disabled
go run ./cmd/server
mv config/local.yaml.disabled config/local.yaml
~~~~

预期看到 read config 错误。

### 4.4 SQL 迁移

创建 migrations/001_init.sql：

~~~~
CREATE TABLE IF NOT EXISTS users (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    account VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    nickname VARCHAR(100) NOT NULL DEFAULT '',
    auth_token_hash CHAR(64) NOT NULL DEFAULT '',
    state TINYINT UNSIGNED NOT NULL DEFAULT 1,
    created_at DATETIME(6) NOT NULL,
    updated_at DATETIME(6) NOT NULL,
    deleted_at DATETIME(6) NULL,
    PRIMARY KEY (id),
    UNIQUE KEY uk_users_account (account),
    KEY idx_users_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS stores (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    owner_id BIGINT UNSIGNED NOT NULL,
    name VARCHAR(100) NOT NULL,
    slug VARCHAR(100) NOT NULL,
    description TEXT NOT NULL,
    status TINYINT UNSIGNED NOT NULL DEFAULT 1,
    created_at DATETIME(6) NOT NULL,
    updated_at DATETIME(6) NOT NULL,
    deleted_at DATETIME(6) NULL,
    PRIMARY KEY (id),
    UNIQUE KEY uk_stores_slug (slug),
    KEY idx_stores_owner_id (owner_id),
    CONSTRAINT fk_stores_owner_id FOREIGN KEY (owner_id) REFERENCES users(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS products (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    store_id BIGINT UNSIGNED NOT NULL,
    name VARCHAR(150) NOT NULL,
    description TEXT NOT NULL,
    price_cents BIGINT NOT NULL,
    stock INT NOT NULL DEFAULT 0,
    status TINYINT UNSIGNED NOT NULL DEFAULT 1,
    created_at DATETIME(6) NOT NULL,
    updated_at DATETIME(6) NOT NULL,
    deleted_at DATETIME(6) NULL,
    PRIMARY KEY (id),
    KEY idx_products_store_id (store_id),
    KEY idx_products_store_status (store_id, status),
    CONSTRAINT fk_products_store_id FOREIGN KEY (store_id) REFERENCES stores(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
~~~~

训练版使用随机 Token 的 SHA-256 哈希做查询条件，但不为初始空字符串建立唯一索引；否则多个尚未登录的用户会因为相同的默认空值产生冲突。

重新初始化：

~~~~
docker compose down -v
docker compose up -d
docker compose exec mysql mysql -uapp -papp_password acggoods_practice \
  -e "SHOW TABLES;"
~~~~

预期：

~~~~
products
stores
users
~~~~

### 4.5 GORM 初始化

创建 internal/app/database.go：

~~~~
package app

import (
	"fmt"

	"example.com/acggoods-go-practice/internal/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func OpenDatabase(cfg *config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(cfg.Database.DSN), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("get sql database: %w", err)
	}

	sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("ping database: %w", err)
	}

	return db, nil
}
~~~~

当前项目对应阅读：

- internal/config/config.go
- internal/database/database.go
- internal/redis/redis.go

训练项目先使用一个数据库和一个 Redis。

---

## 七、User CRUD

## 第 5 阶段：Handler → Service → Repository → GORM

这一阶段先不加鉴权，专门练习完整 CRUD 分层。

### 5.1 Model

创建 internal/model/user.go：

~~~~
package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID            uint64         `gorm:"primaryKey"`
	Account       string         `gorm:"size:255;not null;uniqueIndex"`
	PasswordHash  string         `gorm:"size:255;not null"`
	Nickname      string         `gorm:"size:100;not null"`
    AuthTokenHash string         `gorm:"size:64;not null"`
	State         uint8          `gorm:"not null;default:1"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"`
}
~~~~

### 5.2 DTO

创建 internal/dto/auth.go：

~~~~
package dto

type RegisterRequest struct {
	Account  string `json:"account" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Nickname string `json:"nickname" binding:"required,min=1,max=100"`
}

type LoginRequest struct {
	Account  string `json:"account" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type UserResponse struct {
	ID       uint64 `json:"id"`
	Account  string `json:"account"`
	Nickname string `json:"nickname"`
}
~~~~

### 5.3 Repository

创建 internal/repository/user.go：

~~~~
package repository

import (
	"context"

	"example.com/acggoods-go-practice/internal/model"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	GetByID(ctx context.Context, id uint64) (*model.User, error)
	GetByAccount(ctx context.Context, account string) (*model.User, error)
	UpdateNickname(ctx context.Context, id uint64, nickname string) error
	Delete(ctx context.Context, id uint64) error
	UpdateAuthTokenHash(ctx context.Context, id uint64, hash string) error
	ClearAuthTokenHash(ctx context.Context, id uint64) error
	GetByAuthTokenHash(ctx context.Context, hash string) (*model.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepository) GetByID(ctx context.Context, id uint64) (*model.User, error) {
	var user model.User
	if err := r.db.WithContext(ctx).First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByAccount(ctx context.Context, account string) (*model.User, error) {
	var user model.User
	if err := r.db.WithContext(ctx).
		Where("account = ?", account).
		First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) UpdateNickname(ctx context.Context, id uint64, nickname string) error {
	return r.db.WithContext(ctx).
		Model(&model.User{}).
		Where("id = ?", id).
		Update("nickname", nickname).Error
}

func (r *userRepository) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Delete(&model.User{}, id).Error
}

func (r *userRepository) UpdateAuthTokenHash(ctx context.Context, id uint64, hash string) error {
	return r.db.WithContext(ctx).
		Model(&model.User{}).
		Where("id = ?", id).
		Update("auth_token_hash", hash).Error
}

func (r *userRepository) ClearAuthTokenHash(ctx context.Context, id uint64) error {
	return r.UpdateAuthTokenHash(ctx, id, "")
}

func (r *userRepository) GetByAuthTokenHash(ctx context.Context, hash string) (*model.User, error) {
	var user model.User
	if err := r.db.WithContext(ctx).
		Where("auth_token_hash = ?", hash).
		First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
~~~~

要求：

- Repository 只处理数据访问。
- Repository 接收 context.Context。
- Repository 不返回 HTTP status。
- 查询使用参数绑定，禁止拼接 SQL。

### 5.4 Service

创建 internal/service/user.go，先实现：

- 注册前检查账号。
- 查询用户。
- 修改昵称。
- 删除用户。
- 将 gorm.ErrRecordNotFound 转换为业务错误。

接口：

~~~~
type UserService interface {
	Register(ctx context.Context, req dto.RegisterRequest) (*dto.UserResponse, error)
	Get(ctx context.Context, id uint64) (*dto.UserResponse, error)
	UpdateNickname(ctx context.Context, id uint64, nickname string) error
	Delete(ctx context.Context, id uint64) error
}
~~~~

错误：

~~~~
var (
	ErrUserNotFound = errors.New("user not found")
	ErrAccountTaken = errors.New("account already exists")
)
~~~~

注册流程：

~~~~
GetByAccount
  -> gorm.ErrRecordNotFound：可以创建
  -> 找到用户：ErrAccountTaken
  -> 其他错误：向上返回并记录日志
Create
  -> 返回脱敏 UserResponse
~~~~

### 5.5 Handler 和手动依赖注入

创建 internal/handler/user.go：

~~~~
type UserHandler struct {
	users service.UserService
}

func NewUserHandler(users service.UserService) *UserHandler {
	return &UserHandler{users: users}
}
~~~~

Handler 必须：

- 使用 ShouldBindJSON。
- 读取并校验 URL 参数。
- 调用 Service。
- 把业务错误映射为统一响应。
- 不写 GORM 查询。
- 不把数据库原始错误返回客户端。

main 中按顺序创建：

~~~~
cfg
  -> db
      -> userRepository
          -> userService
              -> userHandler
                  -> router
~~~~

### 5.6 请求和输出

注册：

~~~~
curl -i -X POST http://localhost:8080/api/auth/register \
  -H 'Content-Type: application/json' \
  -d '{"account":"alice@example.com","password":"not-used-yet","nickname":"Alice"}'
~~~~

预期：

~~~~
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "account": "alice@example.com",
    "nickname": "Alice"
  },
  "trace_id": "...",
  "request_id": "..."
}
~~~~

错误练习：

~~~~
curl -i http://localhost:8080/api/users/abc
curl -i http://localhost:8080/api/users/999999
~~~~

预期分别为 400 和 404。

---

## 八、Token Cookie 鉴权

## 第 6 阶段：bcrypt、登录、退出、中间件

### 6.1 bcrypt

安装：

~~~~
go get golang.org/x/crypto@v0.39.0
go mod tidy
~~~~

注册：

~~~~
passwordHash, err := bcrypt.GenerateFromPassword(
	[]byte(req.Password),
	bcrypt.DefaultCost,
)
~~~~

登录：

~~~~
err := bcrypt.CompareHashAndPassword(
	[]byte(user.PasswordHash),
	[]byte(password),
)
~~~~

规则：

- 只保存 password_hash。
- 响应不返回密码和哈希。
- 日志不记录密码。
- 不使用 MD5。当前项目的 MD5 是遗留实现，不作为新代码范例。

### 6.2 Token 设计

Cookie 保存原始 Token，数据库保存 Token 哈希：

~~~~
rawToken = random token in cookie
tokenHash = SHA-256(rawToken) in database
~~~~

辅助函数：

~~~~
func hashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}
~~~~

随机 Token：

~~~~
token := uuid.NewString() + uuid.NewString()
~~~~

登录流程：

~~~~
查询账号
  -> bcrypt 校验密码
  -> 生成 rawToken
  -> 保存 hashToken(rawToken)
  -> 写入 HttpOnly Cookie
  -> 返回脱敏用户
~~~~

### 6.3 Cookie 和中间件

Cookie：

~~~~
c.SetCookie(
	cfg.Auth.CookieName,
	rawToken,
	7*24*60*60,
	"/",
	"",
	false,
	true,
)
~~~~

本地 HTTP 使用 Secure=false，HttpOnly=true。

RequiredAuth 中间件必须：

~~~~
读取 Cookie
  -> 为空：401 并 Abort
  -> Token 无效：401 并 Abort
  -> 找到用户：Set current_user_id
  -> c.Next()
~~~~

Service 对 Token 做 SHA-256 后查询 auth_token_hash。

### 6.4 登录、当前用户、退出

接口：

~~~~
POST /api/auth/login
GET  /api/me
POST /api/auth/logout
~~~~

登录：

~~~~
curl -c cookies.txt -i -X POST http://localhost:8080/api/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"account":"alice@example.com","password":"correct-horse-battery-staple"}'
~~~~

预期：

- 响应头有 Set-Cookie。
- Cookie 有 HttpOnly。
- JSON 没有密码和 Token。

当前用户：

~~~~
curl -b cookies.txt -i http://localhost:8080/api/me
~~~~

退出：

~~~~
curl -b cookies.txt -i -X POST http://localhost:8080/api/auth/logout
curl -b cookies.txt -i http://localhost:8080/api/me
~~~~

最后一次必须返回 401。

### 6.5 鉴权矩阵

| 场景 | 预期 |
| --- | --- |
| 不带 Cookie 访问 /api/me | 401 |
| 随机 Cookie | 401 |
| 密码错误 | 401 |
| 重复账号 | 409 |
| 登录成功 | 200 + Set-Cookie |
| 退出后旧 Cookie | 401 |
| 密码出现在响应 | 测试失败 |
| Token 出现在日志 | 测试失败 |

对照当前项目：

- internal/router/web_auth.go
- internal/controller/web/auth.go
- internal/middleware/web.go
- internal/cookie/cookie.go

---

## 九、店铺 CRUD 和归属权限

## 第 7 阶段：防止 IDOR 越权

### 7.1 Store Model 和 DTO

创建 internal/model/store.go：

~~~~
package model

import (
	"time"

	"gorm.io/gorm"
)

type Store struct {
	ID          uint64         `gorm:"primaryKey"`
	OwnerID     uint64         `gorm:"not null;index"`
	Name        string         `gorm:"size:100;not null"`
	Slug        string         `gorm:"size:100;not null;uniqueIndex"`
	Description string         `gorm:"type:text;not null"`
	Status      uint8          `gorm:"not null;default:1"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}
~~~~

DTO：

~~~~
type CreateStoreRequest struct {
	Name        string `json:"name" binding:"required,min=1,max=100"`
	Slug        string `json:"slug" binding:"required,min=1,max=100"`
	Description string `json:"description" binding:"max=2000"`
}

type UpdateStoreRequest struct {
	Name        string `json:"name" binding:"required,min=1,max=100"`
	Description string `json:"description" binding:"max=2000"`
}
~~~~

### 7.2 Service 接口和 owner 校验

接口：

~~~~
type StoreService interface {
	Create(ctx context.Context, ownerID uint64, req dto.CreateStoreRequest) (*dto.StoreResponse, error)
	ListByOwner(ctx context.Context, ownerID uint64) ([]dto.StoreResponse, error)
	Get(ctx context.Context, currentUserID uint64, storeID uint64) (*dto.StoreResponse, error)
	Update(ctx context.Context, currentUserID uint64, storeID uint64, req dto.UpdateStoreRequest) (*dto.StoreResponse, error)
	Delete(ctx context.Context, currentUserID uint64, storeID uint64) error
}
~~~~

更新和删除：

~~~~
查询资源
  -> 不存在：404
  -> resource.OwnerID != currentUserID：403
  -> 通过校验后才允许修改
~~~~

不要接受请求 body 中的 owner_id。owner 必须来自认证上下文。

### 7.3 请求

创建：

~~~~
curl -b cookies.txt -i -X POST http://localhost:8080/api/stores \
  -H 'Content-Type: application/json' \
  -d '{"name":"Alice Studio","slug":"alice-studio","description":"Acrylic goods"}'
~~~~

预期：

~~~~
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "owner_id": 1,
    "name": "Alice Studio",
    "slug": "alice-studio",
    "description": "Acrylic goods",
    "status": 1
  }
}
~~~~

用户 B 修改 Alice 的店铺：

~~~~
curl -b bob-cookies.txt -i -X PUT http://localhost:8080/api/stores/1 \
  -H 'Content-Type: application/json' \
  -d '{"name":"Hacked","description":"changed"}'
~~~~

预期 HTTP 403。

### 7.4 IDOR 练习

错误实现：

~~~~
UPDATE stores SET name = ? WHERE id = ?
~~~~

安全实现的前置条件：

~~~~
store.OwnerID == currentUserID
~~~~

Service 必须接收当前用户 ID 和资源 ID：

~~~~
Update(ctx, currentUserID, storeID, req)
~~~~

不能只接收 storeID。

### 7.5 必测场景

- owner 可以修改自己的店铺。
- 非 owner 无法修改和删除。
- 不存在返回 404。
- 未登录返回 401。
- 重复 slug 返回 409。
- 伪造 owner_id 不改变真实 owner。

---

## 十、商品 CRUD、分页和状态

## 第 8 阶段：嵌套资源和查询参数

### 8.1 Product Model 和 DTO

创建 internal/model/product.go：

~~~~
package model

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	ID          uint64         `gorm:"primaryKey"`
	StoreID     uint64         `gorm:"not null;index"`
	Name        string         `gorm:"size:150;not null"`
	Description string         `gorm:"type:text;not null"`
	PriceCents  int64          `gorm:"not null"`
	Stock       int            `gorm:"not null;default:0"`
	Status      uint8          `gorm:"not null;default:1"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}
~~~~

DTO：

~~~~
type CreateProductRequest struct {
	Name        string `json:"name" binding:"required,min=1,max=150"`
	Description string `json:"description" binding:"max=5000"`
	PriceCents  int64  `json:"price_cents" binding:"gte=0"`
	Stock       int    `json:"stock" binding:"gte=0"`
}

type UpdateProductRequest struct {
	Name        string `json:"name" binding:"required,min=1,max=150"`
	Description string `json:"description" binding:"max=5000"`
	PriceCents  int64  `json:"price_cents" binding:"gte=0"`
	Stock       int    `json:"stock" binding:"gte=0"`
	Status      uint8  `json:"status" binding:"oneof=1 2"`
}

type ProductListQuery struct {
	Page     int    `form:"page,default=1" binding:"gte=1"`
	PageSize int    `form:"page_size,default=20" binding:"gte=1,lte=100"`
	Status   *uint8 `form:"status"`
}
~~~~

价格使用 price_cents 整数，避免浮点金额误差。

### 8.2 Repository 和分页

分页查询至少包含：

~~~~
Where("store_id = ?", storeID).
Order("id DESC").
Offset(offset).
Limit(pageSize).
Find(&products)
~~~~

同时查询 total：

~~~~
Count(&total)
Find(&products)
~~~~

必须满足：

- page_size 最大 100。
- 所有商品查询限定 store_id 或通过 owner 校验。
- 删除使用 GORM soft delete。
- 禁止拼接客户端传入的 SQL 和排序字段。

### 8.3 Service 流程

创建商品：

~~~~
检查店铺存在
  -> 检查当前用户是店铺 owner
  -> 校验 price 和 stock
  -> 创建商品
~~~~

修改和删除：

~~~~
查询商品
  -> 查询商品所属店铺
  -> 校验店铺 owner
  -> 更新或软删除
~~~~

### 8.4 请求

创建：

~~~~
curl -b cookies.txt -i -X POST \
  http://localhost:8080/api/stores/1/products \
  -H 'Content-Type: application/json' \
  -d '{
    "name":"Blue Acrylic Stand",
    "description":"A blue stand",
    "price_cents":1299,
    "stock":10
  }'
~~~~

列表：

~~~~
curl -b cookies.txt -i \
  "http://localhost:8080/api/stores/1/products?page=1&page_size=20"
~~~~

预期：

~~~~
{
  "code": 0,
  "message": "success",
  "data": {
    "list": [
      {
        "id": 1,
        "store_id": 1,
        "name": "Blue Acrylic Stand",
        "description": "A blue stand",
        "price_cents": 1299,
        "stock": 10,
        "status": 1
      }
    ],
    "total": 1,
    "page": 1,
    "page_size": 20,
    "total_pages": 1
  }
}
~~~~

### 8.5 故障练习

- price_cents 为负数返回 400。
- stock 为负数返回 400。
- page_size 为 1000 返回 400。
- 不存在的店铺返回 404。
- B 修改 A 的商品返回 403。
- URL 店铺 ID 与商品真实店铺不一致时拒绝。
- 删除商品后列表不再显示。

重点：URL 中的 /stores/1 不是权限证明，权限必须由 Service 查询数据库确认。

---

## 十一、错误、日志和 recovery

## 第 9 阶段：统一处理高频异常

### 9.1 错误分类

~~~~
400 validation_error
401 unauthorized
403 forbidden
404 not_found
409 conflict
500 internal_error
~~~~

错误响应：

~~~~
{
  "code": 404,
  "message": "store not found",
  "data": null,
  "trace_id": "...",
  "request_id": "..."
}
~~~~

Handler 使用 errors.Is 将 Service 错误映射为 HTTP 响应。未知错误统一返回 500，详细错误只写日志。

### 9.2 Zap 日志

创建 internal/logger/logger.go：

~~~~
package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New(level string) (*zap.Logger, error) {
	parsed, err := zapcore.ParseLevel(level)
	if err != nil {
		parsed = zapcore.InfoLevel
	}

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		zapcore.AddSync(os.Stdout),
		parsed,
	)

	return zap.New(core), nil
}
~~~~

使用结构化字段：

~~~~
logger.Info("store created",
	zap.Uint64("store_id", store.ID),
	zap.Uint64("owner_id", store.OwnerID),
)

logger.Error("failed to create store",
	zap.Error(err),
	zap.Uint64("owner_id", ownerID),
)
~~~~

禁止：

~~~~
zap.String("password", password)
zap.String("token", token)
zap.String("cookie", cookieHeader)
zap.String("dsn", dsn)
~~~~

当前项目部分遗留日志存在记录敏感 Token 或密码的问题，训练项目不能复制。

### 9.3 请求日志

请求日志字段：

~~~~
method
path
status
latency_ms
client_ip
request_id
trace_id
user_id
~~~~

不要记录完整请求 body。需要观察请求时使用断点或脱敏字段。

### 9.4 Recovery

Recovery 必须：

1. 捕获 panic。
2. 写 error 日志。
3. 包含 path、method、request_id、trace_id 和 stack。
4. 返回统一 500。
5. 防止进程退出。

故意制造 nil pointer：

~~~~
var user *model.User
fmt.Println(user.ID)
~~~~

预期：

- 客户端收到 500。
- 服务仍然运行。
- 日志包含 panic 和调用栈。
- 日志不包含认证信息。

### 9.5 错误请求

~~~~
curl -i -X POST http://localhost:8080/api/auth/register \
  -H 'Content-Type: application/json' \
  -d '{}'

curl -i http://localhost:8080/api/users/abc
curl -i http://localhost:8080/api/users/999999
~~~~

分别观察 400、400、404 和 trace_id。

---

## 十二、Redis 缓存

## 第 10 阶段：缓存店铺详情并处理失效

### 10.1 Redis 连接

启动时创建 Redis Client 并 Ping：

~~~~
client := redis.NewClient(&redis.Options{
	Addr:     cfg.Redis.Addr,
	Password: cfg.Redis.Password,
	DB:       cfg.Redis.DB,
})

if err := client.Ping(ctx).Err(); err != nil {
	return nil, fmt.Errorf("ping redis: %w", err)
}
~~~~

### 10.2 缓存规则

key：

~~~~
store:detail:{store_id}
~~~~

查询流程：

~~~~
GET /api/stores/1
  -> Redis GET
      -> 命中：反序列化并返回
      -> 未命中：查询 MySQL
          -> 序列化
          -> Redis SETEX 5m
          -> 返回
~~~~

更新和删除：

~~~~
PUT /api/stores/1
  -> 更新 MySQL
  -> DEL store:detail:1

DELETE /api/stores/1
  -> 删除 MySQL
  -> DEL store:detail:1
~~~~

### 10.3 观察缓存

~~~~
curl -b cookies.txt http://localhost:8080/api/stores/1
redis-cli GET store:detail:1
redis-cli TTL store:detail:1
~~~~

第二次请求日志应显示 cache hit。

修改店铺：

~~~~
curl -b cookies.txt -X PUT http://localhost:8080/api/stores/1 \
  -H 'Content-Type: application/json' \
  -d '{"name":"Alice Studio Updated","description":"updated"}'

redis-cli EXISTS store:detail:1
~~~~

预期 EXISTS 为 0。

### 10.4 故障练习

1. 临时删除更新时的 DEL，观察旧数据。
2. 恢复 DEL 并增加测试。
3. 停止 Redis，观察错误。
4. 讨论缓存失败时应该报错还是降级，不能无条件吞错。

---

## 十三、JWT 对比实现

## 第 11 阶段：同一组接口支持两种鉴权

### 11.1 安装 JWT

~~~~
go get github.com/golang-jwt/jwt/v5@v5.2.1
go mod tidy
~~~~

### 11.2 模式配置

~~~~
auth:
  mode: jwt
  cookie_name: acggoods_practice_auth
  jwt_secret: local-only-secret-change-me
  jwt_expire_minutes: 60
~~~~

支持：

~~~~
auth.mode=cookie
auth.mode=jwt
~~~~

两种模式共用账号查询和 bcrypt 校验，不共用凭证生成和验证代码。

### 11.3 JWT Claims

~~~~
type Claims struct {
	UserID uint64 `json:"user_id"`
	jwt.RegisteredClaims
}
~~~~

签发时包含 user_id、iat、exp。

验证时必须检查：

- 签名算法。
- 签名是否正确。
- exp 是否过期。
- user_id 是否存在。
- 用户是否仍然有效。

### 11.4 JWT 请求

登录后把响应中的 token 保存为变量：

~~~~
TOKEN=YOUR_TOKEN
~~~~

访问：

~~~~
curl -i http://localhost:8080/api/me \
  -H "Authorization: Bearer $TOKEN"
~~~~

JWT 模式登录响应：

~~~~
{
  "code": 0,
  "message": "success",
  "data": {
    "token": "eyJ...",
    "user": {
      "id": 1,
      "account": "alice@example.com",
      "nickname": "Alice"
    }
  }
}
~~~~

Cookie 模式不在 JSON 中返回 Token，只写 HttpOnly Cookie。

### 11.5 方案比较

| 项目 | Cookie Token | JWT |
| --- | --- | --- |
| 请求传递 | Cookie 自动发送 | Authorization Header |
| 服务端状态 | 数据库保存 Token 哈希 | 默认无服务端会话 |
| 退出 | 清 Cookie 并清数据库 Token | 客户端删除或增加黑名单 |
| 浏览器场景 | 需处理 SameSite/CORS | 前端需手动加 Header |
| 训练重点 | Cookie、Session、中间件 | Claims、签名、过期 |

---

## 十四、测试和调试

## 第 12 阶段：三层测试

### 12.1 Service 单元测试

Repository 使用接口，可以写 fake：

~~~~
type fakeUserRepository struct {
	users map[uint64]*model.User
}

func (f *fakeUserRepository) GetByID(
	ctx context.Context,
	id uint64,
) (*model.User, error) {
	user, ok := f.users[id]
	if !ok {
		return nil, gorm.ErrRecordNotFound
	}
	return user, nil
}
~~~~

至少测试：

- 注册成功。
- 重复账号。
- 错误密码。
- 用户不存在。
- 店铺 owner 校验。
- 商品 owner 校验。
- 无效价格和库存。
- 业务错误类型。

运行：

~~~~
go test ./internal/service/... -v
go test -race ./internal/service/...
~~~~

### 12.2 Handler HTTP 测试

使用 httptest：

~~~~
request := httptest.NewRequest(http.MethodGet, "/api/me", nil)
recorder := httptest.NewRecorder()

engine.ServeHTTP(recorder, request)

if recorder.Code != http.StatusUnauthorized {
	t.Fatalf("expected 401, got %d", recorder.Code)
}
~~~~

验证：

- HTTP status。
- JSON code。
- message。
- trace_id。
- 未登录不进入业务 Handler。
- 登录成功设置 Cookie。
- 错误响应不暴露数据库细节。

### 12.3 Docker 集成测试

~~~~
docker compose up -d
go test ./tests/integration -v
~~~~

流程：

1. 连接本地 MySQL。
2. 使用独立测试数据。
3. 创建用户、店铺、商品。
4. 验证真实 CRUD。
5. 验证唯一约束。
6. 验证软删除。
7. 验证 Redis 缓存失效。

不要让测试读取当前项目的 QA、Webtest 或生产配置。

### 12.4 全量质量检查

~~~~
go fmt ./...
go test ./...
go test -race ./...
go vet ./...
go build ./...
git diff --check
git status --short
~~~~

---

## 十五、Delve 和 VS Code 调试

## 第 13 阶段：调试真实请求链路

### 13.1 安装 Delve

~~~~
go install github.com/go-delve/delve/cmd/dlv@v1.24.0
export PATH="$(go env GOPATH)/bin:$PATH"
dlv version
~~~~

### 13.2 launch.json

创建 .vscode/launch.json：

~~~~
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Debug ACGGoods Practice",
      "type": "go",
      "request": "launch",
      "mode": "auto",
      "program": "./cmd/server",
      "env": {
        "ACG_AUTH_MODE": "cookie"
      },
      "args": []
    }
  ]
}
~~~~

### 13.3 断点位置

设置断点：

- handler/auth.go：JSON 绑定和登录入口。
- service/auth.go：bcrypt 和 Token 生成。
- repository/user.go：用户查询。
- middleware/auth.go：Cookie/JWT 解析。
- service/store.go：owner 校验。
- repository/product.go：分页查询。

使用 Alice 登录：

~~~~
curl -c cookies.txt -X POST http://localhost:8080/api/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"account":"alice@example.com","password":"correct-horse-battery-staple"}'
~~~~

观察：

1. 请求 DTO。
2. bcrypt 校验结果。
3. raw Token。
4. Token 哈希。
5. 当前用户 ID。
6. SQL 查询参数。
7. owner 校验结果。
8. Set-Cookie 参数。

### 13.4 故障调试练习

忘记注册路由：

~~~~
rg -n "api/hello|api/stores|api/products" cmd internal
~~~~

中间件没有调用 c.Next：

~~~~
c.Next()
~~~~

owner 查询条件遗漏：

~~~~
Where("id = ?", storeID)
~~~~

恢复为先查询资源、再比较 OwnerID 和 currentUserID 的安全实现，并增加测试。

缓存没有失效：

- 删除更新后的 Redis DEL。
- 观察旧数据。
- 恢复 DEL。
- 增加回归测试。

---

## 十六、当前项目对照阅读

训练项目完成对应阶段后，再回到当前项目：

| 训练能力 | 当前项目参考 |
| --- | --- |
| HTTP 入口 | cmd/server/main.go、cmd/cli/server.go |
| 配置 | internal/config/config.go |
| 响应 | internal/response/response.go |
| 错误 | internal/errors/errors.go |
| 用户模型 | internal/model/user.go |
| 店铺模型 | internal/model/store.go |
| 商品模型 | internal/model/product.go |
| 用户登录路由 | internal/router/web_auth.go |
| 用户鉴权 | internal/middleware/web.go |
| 店铺 CRUD | internal/router/web_store.go |
| 商品 CRUD | internal/controller/web/store_product.go |
| 数据库 | internal/database/database.go |
| Redis | internal/redis/redis.go |
| 日志 | internal/logger/logger.go、internal/middleware/logger.go |

阅读顺序：

1. 比较两边的启动链路。
2. 比较用户注册和 Cookie 鉴权。
3. 比较店铺 owner 校验。
4. 比较商品分页和状态。
5. 比较统一响应与错误。
6. 再阅读当前项目的订单、支付和搜索复杂依赖。

不要在基础项目还没有跑通时先读 wire_gen.go、订单 Service 或支付 Webhook。

---

## 十七、日常开发工作流

### 1. 先看状态和测试

~~~~
git status --short
go test ./...
~~~~

### 2. 先写失败测试

明确输入、期望输出、期望错误和被保护的业务规则。

### 3. 按分层开发

~~~~
DTO
  -> Handler
  -> Service
  -> Repository
  -> Model / SQL
~~~~

### 4. 使用日志和调试器

优先使用 request_id、trace_id、结构化日志、go test -run TestName -v、Delve 断点和 go test -race。

不要用大量 fmt.Println 代替日志和调试器。

### 5. 提交前检查

~~~~
go fmt ./...
go test ./...
go test -race ./...
go vet ./...
go build ./...
git diff --check
git status --short
~~~~

每完成一个独立能力提交一次：

~~~~
git add .
git commit -m "feat: add store ownership checks"
~~~~

---

## 十八、最终验收清单

### 环境

- [ ] Go 1.23.1 安装成功。
- [ ] Docker Desktop 启动成功。
- [ ] MySQL 和 Redis 由 Docker Compose 启动。
- [ ] Redis 返回 PONG。
- [ ] MySQL 可以进入 acggoods_practice。

### Go

- [ ] 能写 package、struct、method、interface。
- [ ] 能使用 slice、map、pointer。
- [ ] 能返回和判断 error。
- [ ] 能使用 context。
- [ ] 能运行 go run、go build、go fmt、go test、go vet。

### HTTP

- [ ] /health 可访问。
- [ ] /api/hello 可访问。
- [ ] 响应使用统一 envelope。
- [ ] 响应有 request_id 和 trace_id。
- [ ] Ctrl-C 可以优雅退出。

### 数据库和 CRUD

- [ ] SQL 迁移可以执行。
- [ ] GORM 可以连接 MySQL。
- [ ] User CRUD 完成。
- [ ] Store CRUD 完成。
- [ ] Product CRUD 完成。
- [ ] 能理解软删除、唯一索引、外键和分页。

### 鉴权

- [ ] 注册使用 bcrypt。
- [ ] Cookie Token 模式可登录和退出。
- [ ] 未登录返回 401。
- [ ] JWT 模式可登录。
- [ ] JWT 过期和签名错误返回 401。
- [ ] 密码、Cookie、Token 不出现在响应和日志中。

### 权限和错误

- [ ] 店铺写操作校验 owner。
- [ ] 商品写操作校验所属店铺 owner。
- [ ] 越权返回 403。
- [ ] 不存在返回 404。
- [ ] 重复账号和 slug 返回 409。
- [ ] 参数错误返回 400。
- [ ] panic 返回统一 500。

### Redis

- [ ] 店铺详情可以缓存。
- [ ] 更新和删除会清理缓存。
- [ ] 能用 Redis CLI 观察 key 和 TTL。
- [ ] 有缓存命中和失效测试。

### 测试和调试

- [ ] 有 Service 单元测试。
- [ ] 有 Handler HTTP 测试。
- [ ] 有 Docker MySQL 集成测试。
- [ ] 能使用 Delve 在三层之间断点。
- [ ] 能通过 go test ./...、go test -race ./...、go vet ./...、go build ./...。

完成后，再进入当前项目的库存、订单、支付、搜索和异步任务模块。
