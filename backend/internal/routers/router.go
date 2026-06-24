package routers

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	ginSwagger "github.com/swaggo/gin-swagger" //nolint: goimports
	"github.com/swaggo/gin-swagger/swaggerFiles"

	"github.com/insight/backend/internal/handler/v1/user"
	mw "github.com/insight/backend/internal/middleware"
	"github.com/insight/backend/pkg/app"
	"github.com/insight/backend/pkg/middleware"
)

// NewRouter loads the middlewares, routes, handlers.
func NewRouter() *gin.Engine {
	g := gin.New()

	// 全局中间件
	g.Use(gin.Recovery())
	g.Use(middleware.NoCache)
	g.Use(middleware.Options)
	g.Use(middleware.Secure)
	g.Use(middleware.Logging())
	g.Use(middleware.RequestID())
	g.Use(middleware.Metrics(app.Conf.Name))
	g.Use(middleware.Tracing(app.Conf.Name))
	g.Use(middleware.Timeout(3 * time.Second))
	g.Use(mw.Translations())

	// 跨域：方便前端 dev server 直连
	corsCfg := cors.DefaultConfig()
	corsCfg.AllowAllOrigins = true
	corsCfg.AllowHeaders = []string{"Origin", "Content-Type", "Authorization", "Accept", "X-Requested-With"}
	corsCfg.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	g.Use(cors.New(corsCfg))

	// 404 / 405
	g.NoRoute(app.RouteNotFound)
	g.NoMethod(app.RouteNotFound)

	// swagger api docs
	g.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// pprof（开发期）
	if app.Conf.EnablePprof {
		pprof.Register(g)
	}

	// 健康检查与 metrics
	g.GET("/health", app.HealthCheck)
	g.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// v1 路由
	apiV1 := g.Group("/v1")
	{
		// 公开接口
		// 公开注册入口已关闭（账号仅能由管理员从后台 /admin/users 创建）
		// 保留 user.Register 实现以供 admin 路径复用
		// apiV1.POST("/register", user.Register)
		apiV1.POST("/login", user.Login)
		apiV1.GET("/users/:id", user.Get)

		// 认证接口
		authed := apiV1.Group("")
		authed.Use(middleware.Auth())
		{
			authed.GET("/users/me", user.Me)
			authed.PUT("/users/:id", user.Update)
		}

		// 管理员接口（仅 admin 角色可访问）
		admin := apiV1.Group("/admin")
		admin.Use(middleware.Auth(), mw.AdminOnly())
		{
			admin.GET("/users", user.List)
			// 复用 user.Register 作为管理员创建账号的接口
			admin.POST("/users", user.Register)
			admin.PUT("/users/:id", user.AdminUpdate)
			admin.DELETE("/users/:id", user.AdminDelete)
		}
	}

	return g
}
