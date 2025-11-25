# Insight Backend API

一个基于 Gin 框架构建的高性能 Go 后端服务，支持用户管理、权限控制、日志记录等功能。

## 项目特性

- 🚀 基于 Gin 框架的高性能 Web 服务
- 🔐 完整的 JWT 认证和权限管理系统
- 📝 结构化日志记录（支持文件切割和压缩）
- 🗄️ GORM ORM 数据库操作
- ⚙️ 灵活的配置管理
- 🕐 定时任务支持
- 🧪 完整的测试覆盖
- 📋 数据验证和错误处理
- 🔧 命令行工具集

## 系统要求

- Go 1.24.4 或更高版本
- MySQL 5.7 或更高版本

## 快速开始

### 1. 克隆项目

```bash
git clone <repository-url>
cd backend
```

### 2. 安装依赖

```bash
go mod tidy
```

### 3. 配置文件

复制配置模板并修改：

```bash
cp config/config_template.yaml config/config.yaml
```

编辑 `config/config.yaml` 文件，配置数据库连接信息：

```yaml
mysql:
  host: "localhost"
  port: 3306
  username: "root"
  password: "your_password"
  database: "insight"
  # ... 其他配置
```

### 4. 初始化数据库

```bash
go run main.go migrate
```

### 5. 创建管理员用户

```bash
go run main.go admin create --username=admin --password=123456
```

### 6. 启动服务

```bash
go run main.go server
```

服务将在 `http://localhost:8099` 启动。

## 命令行工具

项目提供了丰富的命令行工具来管理应用：

### 主命令

```bash
# 显示帮助信息
go run main.go -h

# 显示版本信息
go run main.go version
```

### 服务器管理

```bash
# 启动 HTTP 服务器
go run main.go server

# 指定配置文件启动
go run main.go server -c config.yaml
```

### 数据库迁移

```bash
# 执行数据库迁移
go run main.go migrate
```

这将自动创建所需的数据库表结构。

### 管理员用户管理

#### 创建管理员用户

```bash
# 创建管理员用户（必需参数）
go run main.go admin create --username=admin --password=123456

# 创建用户并指定详细信息
go run main.go admin create \
  --username=admin \
  --password=123456 \
  --email=admin@example.com \
  --mobile=13800138000 \
  --nickname="系统管理员"

# 创建普通用户（非管理员）
go run main.go admin create \
  --username=user1 \
  --password=123456 \
  --admin=false
```

#### 查看用户列表

```bash
# 列出所有管理员用户
go run main.go admin list
```

#### 删除用户

```bash
# 删除指定用户
go run main.go admin delete --username=user1
```

#### 重置密码

```bash
# 重置用户密码
go run main.go admin reset-password --username=admin --password=newpassword
```

### 定时任务

```bash
# 启动定时任务服务
go run main.go cron
```

### 命令模式

```bash
# 启动命令模式服务
go run main.go command

# 运行演示命令
go run main.go command demo
```

## API 接口文档

### 认证接口

#### 登录
```
POST /api/auth/login
```

请求体：
```json
{
  "username": "admin",
  "password": "123456"
}
```

响应：
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "username": "admin",
      "nickname": "系统管理员"
    }
  }
}
```

### 用户管理接口

#### 获取用户列表
```
GET /api/admin/users
Authorization: Bearer <token>
```

#### 创建用户
```
POST /api/admin/users
Authorization: Bearer <token>
```

#### 更新用户
```
PUT /api/admin/users/{id}
Authorization: Bearer <token>
```

#### 删除用户
```
DELETE /api/admin/users/{id}
Authorization: Bearer <token>
```

### 示例接口

#### Hello 接口
```
GET /api/hello
```

响应：
```json
{
  "code": 200,
  "message": "success",
  "data": "Hello, Gin!"
}
```

## 开发指南

### 项目结构

```
.
├── cmd/                    # 命令行工具
│   ├── admin/             # 管理员命令
│   ├── command/           # 自定义命令
│   ├── cron/              # 定时任务
│   ├── migrate/           # 数据库迁移
│   ├── server/            # 服务器启动
│   └── version/           # 版本信息
├── config/                # 配置文件
├── data/                  # 数据层
├── internal/              # 内部包
│   ├── controller/        # 控制器
│   ├── service/           # 业务逻辑
│   ├── model/             # 数据模型
│   ├── middleware/        # 中间件
│   ├── routers/           # 路由
│   ├── validator/         # 数据验证
│   └── pkg/               # 工具包
├── logs/                  # 日志文件
└── main.go               # 程序入口
```

### 如何使用日志系统

项目使用 `zap` 日志库，支持结构化日志记录：

```go
package yourpackage

import (
    log "insight/internal/pkg/logger"
    "go.uber.org/zap"
)

func ExampleFunction() {
    // 信息日志
    log.Logger.Info("用户登录成功", 
        zap.String("username", "admin"),
        zap.String("ip", "192.168.1.1"),
    )
    
    // 错误日志
    log.Logger.Error("数据库连接失败", 
        zap.Error(err),
        zap.String("database", "insight"),
    )
    
    // 警告日志
    log.Logger.Warn("用户名已存在", 
        zap.String("username", "admin"),
    )
    
    // 调试日志（仅在开发模式下显示）
    log.Logger.Debug("调试信息", 
        zap.Any("data", someData),
    )
}
```

#### 日志配置

在 `config.yaml` 中配置日志参数：

```yaml
logger:
  file_name: "app.log"              # 日志文件名
  default_division: "size"          # 切割方式：size/time
  division_time:                    # 按时间切割
    max_age: 7                      # 保留天数
    rotation_time: 24               # 切割间隔（小时）
  division_size:                    # 按大小切割
    max_size: 100                   # 单个文件最大大小（MB）
    max_backups: 10                 # 保留文件数量
    max_age: 7                      # 保留天数
    compress: true                  # 是否压缩
```

### 如何创建新的 API 接口

#### 1. 创建数据模型（可选）

在 `internal/model/` 目录下创建模型文件：

```go
// internal/model/product.go
package model

import "gorm.io/gorm"

type Product struct {
    BaseModel
    Name        string  `json:"name" gorm:"size:100;not null"`
    Price       float64 `json:"price" gorm:"type:decimal(10,2)"`
    Description string  `json:"description" gorm:"type:text"`
    Status      int     `json:"status" gorm:"default:1"`
}

func NewProduct() *Product {
    return &Product{}
}

func (p *Product) GetList() []Product {
    var products []Product
    // 实现获取产品列表的逻辑
    return products
}
```

#### 2. 创建验证器（可选）

在 `internal/validator/form/` 目录下创建验证文件：

```go
// internal/validator/form/product.go
package form

type ProductCreateForm struct {
    Name        string  `json:"name" binding:"required,min=2,max=100"`
    Price       float64 `json:"price" binding:"required,min=0"`
    Description string  `json:"description" binding:"max=500"`
}

type ProductUpdateForm struct {
    Name        string  `json:"name" binding:"omitempty,min=2,max=100"`
    Price       float64 `json:"price" binding:"omitempty,min=0"`
    Description string  `json:"description" binding:"omitempty,max=500"`
    Status      int     `json:"status" binding:"omitempty,oneof=0 1"`
}
```

#### 3. 创建服务层

在 `internal/service/` 目录下创建服务文件：

```go
// internal/service/product.go
package service

import (
    "insight/internal/model"
    "insight/internal/validator/form"
)

type ProductService interface {
    GetList() ([]model.Product, error)
    Create(form *form.ProductCreateForm) (*model.Product, error)
    Update(id uint, form *form.ProductUpdateForm) (*model.Product, error)
    Delete(id uint) error
}

type productServiceImpl struct{}

func NewProductService() ProductService {
    return &productServiceImpl{}
}

func (s *productServiceImpl) GetList() ([]model.Product, error) {
    productModel := model.NewProduct()
    products := productModel.GetList()
    return products, nil
}

func (s *productServiceImpl) Create(form *form.ProductCreateForm) (*model.Product, error) {
    product := &model.Product{
        Name:        form.Name,
        Price:       form.Price,
        Description: form.Description,
        Status:      1,
    }
    
    // 保存到数据库的逻辑
    // ...
    
    return product, nil
}

func (s *productServiceImpl) Update(id uint, form *form.ProductUpdateForm) (*model.Product, error) {
    // 更新逻辑
    return nil, nil
}

func (s *productServiceImpl) Delete(id uint) error {
    // 删除逻辑
    return nil
}
```

#### 4. 创建控制器

在 `internal/controller/` 目录下创建控制器文件：

```go
// internal/controller/product/product.go
package product

import (
    "net/http"
    "strconv"
    
    "insight/internal/controller"
    log "insight/internal/pkg/logger"
    "insight/internal/service"
    "insight/internal/validator/form"
    
    "github.com/gin-gonic/gin"
    "go.uber.org/zap"
)

type ProductController struct {
    controller.Api
}

func NewProductController() *ProductController {
    return &ProductController{}
}

func (api *ProductController) GetList(c *gin.Context) {
    log.Logger.Info("获取产品列表", zap.String("path", c.Request.URL.Path))
    
    products, err := service.NewProductService().GetList()
    if err != nil {
        log.Logger.Error("获取产品列表失败", zap.Error(err))
        api.Err(c, err)
        return
    }
    
    api.Success(c, products)
}

func (api *ProductController) Create(c *gin.Context) {
    var form form.ProductCreateForm
    
    if err := c.ShouldBindJSON(&form); err != nil {
        log.Logger.Error("参数验证失败", zap.Error(err))
        api.ValidatorError(c, err)
        return
    }
    
    product, err := service.NewProductService().Create(&form)
    if err != nil {
        log.Logger.Error("创建产品失败", zap.Error(err))
        api.Err(c, err)
        return
    }
    
    log.Logger.Info("产品创建成功", zap.String("name", form.Name))
    api.Success(c, product)
}

func (api *ProductController) Update(c *gin.Context) {
    idStr := c.Param("id")
    id, err := strconv.ParseUint(idStr, 10, 32)
    if err != nil {
        api.ParameterError(c, "无效的产品ID")
        return
    }
    
    var form form.ProductUpdateForm
    if err := c.ShouldBindJSON(&form); err != nil {
        api.ValidatorError(c, err)
        return
    }
    
    product, err := service.NewProductService().Update(uint(id), &form)
    if err != nil {
        api.Err(c, err)
        return
    }
    
    api.Success(c, product)
}

func (api *ProductController) Delete(c *gin.Context) {
    idStr := c.Param("id")
    id, err := strconv.ParseUint(idStr, 10, 32)
    if err != nil {
        api.ParameterError(c, "无效的产品ID")
        return
    }
    
    err = service.NewProductService().Delete(uint(id))
    if err != nil {
        api.Err(c, err)
        return
    }
    
    api.Success(c, "删除成功")
}
```

#### 5. 注册控制器

在 `internal/routers/setup/controllers.go` 中添加新控制器：

```go
func NewControllers() *Controllers {
    return &Controllers{
        HelloController:   hello.NewHelloController(),
        DemoController:    demo.NewDemoController(),
        AuthController:    admin.NewAuthController(),
        AdminUserController: admin.NewAdminUserController(),
        PermissionController: admin.NewPermissionController(),
        ProductController: product.NewProductController(), // 添加新控制器
    }
}
```

#### 6. 创建路由组

在 `internal/routers/groups/` 目录下创建路由文件：

```go
// internal/routers/groups/product.go
package groups

import (
    "insight/internal/routers/setup"
    
    "github.com/gin-gonic/gin"
)

func ProductRouters(router *gin.RouterGroup, controller setup.Controllers) {
    productGroup := router.Group("/products")
    {
        productGroup.GET("", controller.ProductController.GetList)
        productGroup.POST("", controller.ProductController.Create)
        productGroup.PUT("/:id", controller.ProductController.Update)
        productGroup.DELETE("/:id", controller.ProductController.Delete)
    }
}
```

#### 7. 注册路由

在 `internal/routers/router.go` 中添加路由组：

```go
func SetupRouter(router *gin.Engine) {
    Controllers := setup.NewControllers()
    api := router.Group("/api")
    groups.HelloRouters(api, *Controllers)
    groups.DemoRouters(api, *Controllers)
    groups.AdminRouters(api, *Controllers)
    groups.ProductRouters(api, *Controllers) // 添加新路由组
}
```

### 响应格式

所有 API 响应都遵循统一格式：

```json
{
  "code": 200,           // 状态码
  "message": "success",  // 消息
  "data": {}            // 数据内容
}
```

#### 成功响应
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "id": 1,
    "name": "产品名称"
  }
}
```

#### 错误响应
```json
{
  "code": 400,
  "message": "参数验证失败",
  "data": null
}
```

### 中间件使用

#### 添加认证中间件

```go
// 需要认证的路由组
authGroup := api.Group("/admin")
authGroup.Use(middleware.AdminAuthMiddleware())
{
    authGroup.GET("/users", controller.AdminUserController.GetList)
    authGroup.POST("/users", controller.AdminUserController.Create)
}
```

### 数据库操作

使用 GORM ORM 进行数据库操作：

```go
import "insight/data"

// 查询
var user model.AdminUser
data.MysqlDB.Where("username = ?", "admin").First(&user)

// 创建
data.MysqlDB.Create(&user)

// 更新
data.MysqlDB.Model(&user).Updates(map[string]interface{}{
    "nickname": "新昵称",
    "email": "new@example.com",
})

// 软删除
data.MysqlDB.Delete(&user)
```

## 配置说明

### 数据库配置

```yaml
mysql:
  host: "localhost"        # 数据库地址
  port: 3306              # 端口
  username: "root"        # 用户名
  password: "password"    # 密码
  database: "insight"     # 数据库名
  print_sql: false        # 是否打印 SQL
  log_level: "info"       # 日志级别
  table_prefix: ""        # 表前缀
  max_idle_conns: 10      # 最大空闲连接数
  max_open_conns: 100     # 最大开放连接数
  max_life_time: 3600     # 连接最大生命周期（秒）
  enable: true            # 是否启用
```

### 系统配置

```yaml
system:
  host: "0.0.0.0"         # 监听地址
  port: 8099              # 监听端口
  language: "zh_CN"       # 语言
  debug: false            # 调试模式
```

### JWT 配置

```yaml
jwt:
  secret: "insight"       # JWT 密钥
  header_prefix: "Bearer" # 请求头前缀
  expiration: 7200        # 过期时间（秒）
  refresh_time: 86400     # 刷新时间
  ttl: 7200s             # 生存时间
```

## 部署

### 构建

```bash
# 构建二进制文件
go build -o insight main.go

# 交叉编译（Linux）
GOOS=linux GOARCH=amd64 go build -o insight-linux main.go

# 交叉编译（Windows）
GOOS=windows GOARCH=amd64 go build -o insight-windows.exe main.go
```

### Docker 部署

创建 `Dockerfile`：

```dockerfile
FROM golang:1.24.4-alpine AS builder

WORKDIR /app
COPY . .
RUN go mod tidy && go build -o insight main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/
COPY --from=builder /app/insight .
COPY --from=builder /app/config ./config

EXPOSE 8099
CMD ["./insight", "server"]
```

构建和运行：

```bash
# 构建镜像
docker build -t insight-backend .

# 运行容器
docker run -p 8099:8099 -v $(pwd)/config:/root/config insight-backend
```

### 生产环境建议

1. 使用环境变量覆盖敏感配置
2. 启用 HTTPS
3. 配置反向代理（Nginx）
4. 设置日志级别为 `info` 或 `warn`
5. 定期备份数据库
6. 监控服务状态和性能

## 故障排除

### 常见问题

#### 1. 数据库连接失败
检查配置文件中的数据库配置，确保数据库服务正常运行。

#### 2. 端口占用
```bash
# 查看端口占用
lsof -i :8099

# 杀死进程
kill -9 <PID>
```

#### 3. 日志文件权限问题
确保应用有写入日志目录的权限：
```bash
chmod 755 logs/
```
