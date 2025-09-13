package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"zhku-oj/internal/config"
	"zhku-oj/internal/handler/admin"
	"zhku-oj/internal/handler/auth"
	"zhku-oj/internal/handler/problem"
	"zhku-oj/internal/handler/submission"
	"zhku-oj/internal/handler/user"
	"zhku-oj/internal/middleware"
	"zhku-oj/internal/pkg/database"
	"zhku-oj/internal/pkg/logger"
	"zhku-oj/internal/repository/mongodb"
	"zhku-oj/internal/service/impl"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化配置
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 初始化日志
	logger.Init(cfg.Logging)

	// 初始化数据库连接
	mongoClient, err := database.NewMongoDB(cfg.MongoDB)
	if err != nil {
		log.Fatalf("连接MongoDB失败: %v", err)
	}
	defer mongoClient.Disconnect(context.Background())

	redisClient, err := database.NewRedis(cfg.Redis)
	if err != nil {
		log.Fatalf("连接Redis失败: %v", err)
	}
	defer redisClient.Close()

	// 初始化Repository层
	userRepo := mongodb.NewUserRepository(mongoClient, cfg.MongoDB.Database)
	problemRepo := mongodb.NewProblemRepository(mongoClient, cfg.MongoDB.Database)
	submissionRepo := mongodb.NewSubmissionRepository(mongoClient, cfg.MongoDB.Database)

	// 初始化Service层
	authService := impl.NewAuthService(userRepo, redisClient, cfg)
	userService := impl.NewUserService(userRepo, redisClient)
	problemService := impl.NewProblemService(problemRepo, redisClient)
	submissionService := impl.NewSubmissionService(submissionRepo, problemRepo, redisClient, cfg)

	// 初始化Handler层
	authHandler := auth.NewAuthHandler(authService)
	userHandler := user.NewUserHandler(userService)
	problemHandler := problem.NewProblemHandler(problemService)
	submissionHandler := submission.NewSubmissionHandler(submissionService)
	adminHandler := admin.NewAdminHandler(userService, problemService, submissionService)

	// 设置Gin模式
	gin.SetMode(cfg.Server.Mode)

	// 创建路由
	router := gin.New()

	// 添加中间件
	router.Use(middleware.Logger())
	router.Use(middleware.Recovery())
	router.Use(middleware.CORS())

	// 健康检查接口
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "healthy",
			"timestamp": time.Now().Unix(),
			"version":   "v1.0.0",
		})
	})

	// API路由组
	v1 := router.Group("/api/v1")
	{
		// 认证相关路由
		authGroup := v1.Group("/auth")
		{
			authGroup.POST("/register", authHandler.Register)
			authGroup.POST("/login", authHandler.Login)
			authGroup.POST("/logout", middleware.AuthRequired(), authHandler.Logout)
		}

		// 用户相关路由
		userGroup := v1.Group("/users")
		userGroup.Use(middleware.AuthRequired())
		{
			// 当前用户相关接口
			userGroup.GET("/profile", userHandler.GetProfile)
			userGroup.PUT("/profile", userHandler.UpdateProfile)
			userGroup.PUT("/password", userHandler.ChangePassword)

			// 用户信息查询接口
			userGroup.GET("/:id", userHandler.GetUser)
			userGroup.GET("/:id/stats", userHandler.GetUserStats)
		}

		// 题目相关路由
		problemGroup := v1.Group("/problems")
		problemGroup.Use(middleware.AuthRequired())
		{
			problemGroup.GET("", problemHandler.ListProblems)
			problemGroup.GET("/:id", problemHandler.GetProblem)
			problemGroup.POST("", middleware.RoleRequired("teacher", "admin"), problemHandler.CreateProblem)
			problemGroup.PUT("/:id", middleware.RoleRequired("teacher", "admin"), problemHandler.UpdateProblem)
			problemGroup.DELETE("/:id", middleware.RoleRequired("admin"), problemHandler.DeleteProblem)
		}

		// 提交相关路由
		submissionGroup := v1.Group("/submissions")
		submissionGroup.Use(middleware.AuthRequired())
		{
			submissionGroup.POST("", submissionHandler.Submit)
			submissionGroup.GET("/:id", submissionHandler.GetSubmission)
			submissionGroup.GET("", submissionHandler.ListSubmissions)
		}

		// 管理员路由
		adminGroup := v1.Group("/admin")
		adminGroup.Use(middleware.AuthRequired(), middleware.RoleRequired("admin"))
		{
			// 系统管理
			adminGroup.GET("/dashboard", adminHandler.Dashboard)
			adminGroup.GET("/system/status", adminHandler.SystemStatus)

			// 用户管理 CRUD 接口
			adminGroup.POST("/users", userHandler.CreateUser)                   // 创建用户
			adminGroup.GET("/users", userHandler.ListUsers)                     // 获取用户列表
			adminGroup.PUT("/users/:id", userHandler.UpdateUser)                // 更新用户
			adminGroup.DELETE("/users/:id", userHandler.DeleteUser)             // 删除用户
			adminGroup.PUT("/users/:id/activate", userHandler.ActivateUser)     // 激活用户
			adminGroup.PUT("/users/:id/deactivate", userHandler.DeactivateUser) // 停用用户
		}
	}

	// 创建HTTP服务器
	server := &http.Server{
		Addr:         cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	// 启动服务器
	go func() {
		logger.Info("服务器启动在端口: %s", cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("服务器启动失败: %v", err)
		}
	}()

	// 等待中断信号以优雅关闭服务器
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("正在关闭服务器...")

	// 优雅关闭服务器，等待5秒
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Fatal("服务器强制关闭: %v", err)
	}

	logger.Info("服务器已关闭")
}
