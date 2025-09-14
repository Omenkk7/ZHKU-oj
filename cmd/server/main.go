package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	database2 "zhku-oj/pkg/database"
	"zhku-oj/pkg/logger"

	"zhku-oj/internal/config"
	"zhku-oj/internal/handler/admin"
	"zhku-oj/internal/handler/auth"
	"zhku-oj/internal/handler/problem"
	"zhku-oj/internal/handler/submission"
	"zhku-oj/internal/handler/user"
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
	mongoClient, err := database2.NewMongoDB(cfg.MongoDB)
	if err != nil {
		log.Fatalf("连接MongoDB失败: %v", err)
	}
	defer mongoClient.Disconnect(context.Background())

	redisClient, err := database2.NewRedis(cfg.Redis)
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

	// 创建路由器
	router := gin.New()

	// 创建路由管理器并设置所有路由
	routerManager := router.NewRouterManager(
		authHandler,
		userHandler,
		problemHandler,
		submissionHandler,
		adminHandler,
	)
	routerManager.SetupRoutes(router)

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
