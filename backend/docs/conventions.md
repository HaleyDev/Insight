# 后端开发规范

本文档描述 Insight 后端的 API 交互格式、新接口开发流程以及数据库 Model 编写规范。

---

## 一、API 交互 JSON 格式

### 1.1 统一响应结构

所有接口（无论成功或失败）均返回如下 JSON 结构：

```json
{
  "code": 0,
  "message": "Ok",
  "data": {},
  "details": []
}
```

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `code` | `int` | 业务状态码，`0` 表示成功，非 `0` 表示错误 |
| `message` | `string` | 状态描述信息 |
| `data` | `object/array/null` | 业务数据，成功时返回具体数据，失败时为 `{}` |
| `details` | `string[]` | 错误详情，仅错误时可能携带，成功时省略 |

> 对应代码：`backend/pkg/app/response.go` 中的 `Response` 结构体。

### 1.2 成功响应

```json
{
  "code": 0,
  "message": "Ok",
  "data": {
    "id": 1,
    "username": "admin",
    "email": "admin@example.com",
    "avatar": "/images/avatars/avatar-1.png",
    "role": "admin"
  }
}
```

- HTTP 状态码：`200`
- `code` 为 `0`，`data` 为实际业务数据
- 如果接口无需返回数据，`data` 为 `{}`

### 1.3 错误响应

```json
{
  "code": 10001,
  "message": "Invalid params",
  "data": {},
  "details": ["email is required"]
}
```

- HTTP 状态码根据错误码映射（见下方错误码表），常见如 `400`、`401`、`404`、`500`
- `code` 为非零业务错误码
- `details` 可选，携带额外错误信息

### 1.4 业务错误码规划

错误码定义在两个位置：

**通用错误码**（`backend/pkg/errcode/code.go`）：

| 错误码 | 说明 | HTTP 状态码 |
| --- | --- | --- |
| `0` | 成功 | 200 |
| `10000` | 内部服务器错误 | 500 |
| `10001` | 参数无效 | 400 |
| `10002` | 未授权 | 401 |
| `10003` | 资源不存在 | 404 |
| `10006` | 访问被拒绝 | 403 |
| `10014` | Token 生成错误 | 401 |
| `10015` | 无效 Token | 401 |
| `10016` | Token 过期 | 401 |

**业务错误码**（`backend/internal/ecode/`，按模块分文件）：

以用户模块为例（`ecode/user.go`）：

| 错误码 | 说明 |
| --- | --- |
| `20101` | 用户不存在 |
| `20102` | 账号或密码错误 |
| `20109` | 邮箱或密码错误 |
| `20110` | 两次密码输入不一致 |
| `20111` | 注册失败 |

> 错误码命名规则：通用码 `1xxxx`，业务码按模块分段（如用户模块 `201xx`），新增模块依次递增（如 `202xx`、`203xx`）。

### 1.5 请求参数绑定

请求参数通过 Gin 的 `ShouldBindJSON`（JSON Body）或 `c.Param`/`c.Query`（路径/查询参数）获取，并使用 struct tag 声明：

```go
type LoginRequest struct {
    Email    string `json:"email" form:"email" binding:"required"`
    Password string `json:"password" form:"password" binding:"required"`
}
```

- `json` tag：JSON 请求体的字段名
- `form` tag：表单/查询参数的字段名
- `binding` tag：校验规则（如 `required`），使用 `go-playground/validator`

---

## 二、新接口开发流程

以「创建一个文章详情接口 `GET /v1/articles/:id`」为例，完整流程如下：

### 第 1 步：定义 Model

在 `backend/internal/model/` 下新建 `article.go`：

```go
package model

import "time"

// ArticleModel 文章表
type ArticleModel struct {
    ID        uint64    `gorm:"primary_key;AUTO_INCREMENT;column:id" json:"id"`
    Title     string    `gorm:"column:title;not null" json:"title" binding:"required"`
    Content   string    `gorm:"column:content;type:text" json:"content"`
    AuthorID  uint64    `gorm:"column:author_id" json:"author_id"`
    CreatedAt time.Time `gorm:"column:created_at" json:"-"`
    UpdatedAt time.Time `gorm:"column:updated_at" json:"-"`
}

// TableName 表名
func (a *ArticleModel) TableName() string {
    return "article"
}

// ArticleInfo 对外暴露的文章结构
type ArticleInfo struct {
    ID       uint64 `json:"id"`
    Title    string `json:"title"`
    Content  string `json:"content"`
    AuthorID uint64 `json:"author_id"`
}

// ToInfo 转换为对外结构
func (a *ArticleModel) ToInfo() *ArticleInfo {
    if a == nil {
        return &ArticleInfo{}
    }
    return &ArticleInfo{
        ID:       a.ID,
        Title:    a.Title,
        Content:  a.Content,
        AuthorID: a.AuthorID,
    }
}
```

### 第 2 步：扩展 Repository 接口与实现

**接口定义**（`backend/internal/repository/repository.go`）：

```go
type Repository interface {
    // ... 已有方法
    GetArticle(ctx context.Context, id uint64) (*model.ArticleModel, error)
}
```

**实现**（新建 `backend/internal/repository/article_repo.go`）：

```go
package repository

import (
    "context"

    "github.com/pkg/errors"
    "gorm.io/gorm"

    "github.com/insight/backend/internal/model"
)

// GetArticle 根据 id 获取文章
func (d *repository) GetArticle(ctx context.Context, id uint64) (*model.ArticleModel, error) {
    var article model.ArticleModel
    err := d.orm.WithContext(ctx).First(&article, id).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, ErrNotFound
        }
        return nil, errors.Wrap(err, "[repo.article] get article err")
    }
    return &article, nil
}
```

### 第 3 步：扩展 Service 接口与实现

**接口定义**（`backend/internal/service/service.go`）：

```go
type Service interface {
    Users() UserService
    Articles() ArticleService  // 新增
}
```

**接口声明与实现**（新建 `backend/internal/service/article_service.go`）：

```go
package service

import (
    "context"

    "github.com/insight/backend/internal/model"
    "github.com/insight/backend/internal/repository"
)

// ArticleService 文章服务接口
type ArticleService interface {
    GetArticle(ctx context.Context, id uint64) (*model.ArticleInfo, error)
}

type articleService struct {
    repo repository.Repository
}

func newArticles(svc *service) *articleService {
    return &articleService{repo: svc.repo}
}

// GetArticle 获取文章详情
func (s *articleService) GetArticle(ctx context.Context, id uint64) (*model.ArticleInfo, error) {
    article, err := s.repo.GetArticle(ctx, id)
    if err != nil {
        return nil, err
    }
    return article.ToInfo(), nil
}
```

**注册到 Service**（`backend/internal/service/service.go`）：

```go
func (s *service) Articles() ArticleService {
    return newArticles(s)
}
```

### 第 4 步：定义业务错误码

在 `backend/internal/ecode/` 下新建 `article.go`：

```go
package ecode

import "github.com/insight/backend/pkg/errcode"

var (
    ErrArticleNotFound = errcode.NewError(20201, "文章不存在")
)
```

### 第 5 步：编写 Handler

新建 `backend/internal/handler/v1/article/article.go`：

```go
package article

import (
    "github.com/insight/backend/pkg/app"
)

var response = app.NewResponse()
```

新建 `backend/internal/handler/v1/article/get.go`：

```go
package article

import (
    "errors"

    "github.com/gin-gonic/gin"
    "github.com/spf13/cast"

    "github.com/insight/backend/internal/ecode"
    "github.com/insight/backend/internal/repository"
    "github.com/insight/backend/internal/service"
    "github.com/insight/backend/pkg/errcode"
    "github.com/insight/backend/pkg/log"
)

// Get 获取文章详情
// @Summary 获取文章详情
// @Tags 文章
// @Produce json
// @Param id path int true "文章 id"
// @Success 200 {object} model.ArticleInfo
// @Router /articles/{id} [get]
func Get(c *gin.Context) {
    articleID := cast.ToUint64(c.Param("id"))
    if articleID == 0 {
        response.Error(c, errcode.ErrInvalidParam)
        return
    }

    info, err := service.Svc.Articles().GetArticle(c.Request.Context(), articleID)
    if errors.Is(err, repository.ErrNotFound) {
        response.Error(c, ecode.ErrArticleNotFound)
        return
    }
    if err != nil {
        log.Errorf("get article err: %+v", err)
        response.Error(c, errcode.ErrInternalServer.WithDetails(err.Error()))
        return
    }

    response.Success(c, info)
}
```

### 第 6 步：注册路由

在 `backend/internal/routers/router.go` 中添加：

```go
import (
    "github.com/insight/backend/internal/handler/v1/article"
)

// 在 v1 路由组中添加
apiV1 := g.Group("/v1")
{
    apiV1.GET("/articles/:id", article.Get)
}
```

### 流程总结

```
Model（数据结构） → Repository（数据访问） → Service（业务逻辑） → Handler（请求处理） → Router（路由注册）
```

每一层的职责：

| 层 | 目录 | 职责 |
| --- | --- | --- |
| Model | `internal/model/` | 定义数据库表结构和对外 DTO |
| Repository | `internal/repository/` | 数据库 CRUD，封装缓存逻辑 |
| Service | `internal/service/` | 业务逻辑编排，调用 Repository |
| Handler | `internal/handler/v1/` | 接收请求、参数校验、调用 Service、返回响应 |
| Router | `internal/routers/` | 路由注册、中间件挂载 |

---

## 三、数据库 Model 编写规范

### 3.1 基本结构

每个数据表对应一个 Model 文件，放在 `backend/internal/model/` 目录下，文件名与表名一致（如 `user_base.go` 对应 `user_base` 表）。

```go
package model

import "time"

// UserBaseModel 用户基础表
type UserBaseModel struct {
    ID        uint64    `gorm:"primary_key;AUTO_INCREMENT;column:id" json:"id"`
    Username  string    `gorm:"column:username;not null" json:"username" binding:"required" validate:"min=1,max=32"`
    Password  string    `gorm:"column:password;not null" json:"-" binding:"required"`
    Email     string    `gorm:"column:email;not null" json:"email"`
    Avatar    string    `gorm:"column:avatar" json:"avatar"`
    Role      string    `gorm:"column:role;not null;default:user" json:"role"`
    CreatedAt time.Time `gorm:"column:created_at" json:"-"`
    UpdatedAt time.Time `gorm:"column:updated_at" json:"-"`
}
```

### 3.2 Tag 规范

**GORM Tag**（数据库映射）：

| Tag | 说明 | 示例 |
| --- | --- | --- |
| `primary_key` | 主键 | `gorm:"primary_key"` |
| `AUTO_INCREMENT` | 自增 | `gorm:"AUTO_INCREMENT"` |
| `column` | 列名 | `gorm:"column:username"` |
| `not null` | 非空约束 | `gorm:"not null"` |
| `default` | 默认值 | `gorm:"default:user"` |
| `type` | 指定类型 | `gorm:"type:text"` |

**JSON Tag**（序列化控制）：

| Tag | 说明 | 示例 |
| --- | --- | --- |
| 字段名 | 正常输出 | `json:"username"` |
| `-` | 不输出 | `json:"-"` |

> 敏感字段（如密码）必须使用 `json:"-"`，确保不会通过 API 泄露。

**Binding Tag**（请求参数校验）：

| Tag | 说明 | 示例 |
| --- | --- | --- |
| `required` | 必填 | `binding:"required"` |
| `min/max` | 长度限制 | `validate:"min=1,max=32"` |

### 3.3 必须实现的方法

每个 Model 必须实现 `TableName()` 方法：

```go
// TableName 表名
func (u *UserBaseModel) TableName() string {
    return "user_base"
}
```

### 3.4 对外 DTO 结构

数据库 Model 不直接返回给前端，需定义对外 DTO 结构（以 `Info` 结尾），并提供转换方法：

```go
// UserInfo 对外暴露的用户结构
type UserInfo struct {
    ID       uint64 `json:"id"`
    Username string `json:"username"`
    Email    string `json:"email"`
    Avatar   string `json:"avatar"`
    Role     string `json:"role"`
}

// ToUserInfo 转换为对外结构
func (u *UserBaseModel) ToUserInfo() *UserInfo {
    if u == nil {
        return &UserInfo{}
    }
    return &UserInfo{
        ID:       u.ID,
        Username: u.Username,
        Email:    u.Email,
        Avatar:   u.Avatar,
        Role:     u.Role,
    }
}
```

> DTO 结构只包含前端需要的字段，屏蔽密码、时间戳等内部字段。

### 3.5 字段命名

- 数据库列名：`snake_case`（如 `created_at`、`author_id`）
- Go 字段名：`CamelCase`（如 `CreatedAt`、`AuthorID`）
- JSON 字段名：`snake_case`（如 `created_at`、`author_id`）

### 3.6 常量定义

与 Model 相关的常量（如默认值、角色枚举）定义在同一个文件中：

```go
const (
    DefaultAvatar = "/images/avatars/avatar-1.png"
)

const (
    RoleAdmin = "admin"
    RoleUser  = "user"
)
```
