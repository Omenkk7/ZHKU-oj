package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	database2 "zhku-oj/pkg/database"
	"zhku-oj/pkg/logger"

	"zhku-oj/internal/config"
	"zhku-oj/internal/queue"
	"zhku-oj/internal/repository/mongodb"
	"zhku-oj/internal/service/impl"
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
	statsService := impl.NewStatsService(userRepo, problemRepo, submissionRepo, redisClient)

	// 初始化消息队列消费者
	consumer, err := queue.NewConsumer(cfg.RabbitMQ)
	if err != nil {
		log.Fatalf("初始化消息队列失败: %v", err)
	}
	defer consumer.Close()

	// 启动统计更新服务
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		logger.Info("统计更新服务已启动")
		if err := consumer.ConsumeStatsUpdates(ctx, statsService.UpdateStats); err != nil {
			logger.Error("统计更新服务失败: %v", err)
		}
	}()

	// 启动通知服务
	go func() {
		logger.Info("通知服务已启动")
		if err := consumer.ConsumeNotifications(ctx, statsService.ProcessNotification); err != nil {
			logger.Error("通知服务失败: %v", err)
		}
	}()

	// 等待中断信号以优雅关闭服务
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("正在关闭工作服务...")

	// 停止消费任务
	cancel()

	logger.Info("工作服务已关闭")
}
