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
	"zhku-oj/internal/judge"
	"zhku-oj/internal/queue"
	"zhku-oj/internal/repository/mongodb"
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
	submissionRepo := mongodb.NewSubmissionRepository(mongoClient, cfg.MongoDB.Database)
	problemRepo := mongodb.NewProblemRepository(mongoClient, cfg.MongoDB.Database)

	// 初始化判题管理器
	judgeManager, err := judge.NewManager(cfg.Judge, submissionRepo, problemRepo)
	if err != nil {
		log.Fatalf("初始化判题管理器失败: %v", err)
	}

	// 初始化消息队列消费者
	consumer, err := queue.NewConsumer(cfg.RabbitMQ)
	if err != nil {
		log.Fatalf("初始化消息队列失败: %v", err)
	}
	defer consumer.Close()

	// 启动判题任务消费者
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 启动判题服务
	go func() {
		logger.Info("判题服务已启动")
		if err := consumer.ConsumeJudgeTasks(ctx, judgeManager.ProcessTask); err != nil {
			logger.Error("判题任务消费失败: %v", err)
		}
	}()

	// 启动结果处理服务
	go func() {
		logger.Info("结果处理服务已启动")
		if err := consumer.ConsumeJudgeResults(ctx, judgeManager.ProcessResult); err != nil {
			logger.Error("结果处理失败: %v", err)
		}
	}()

	// 等待中断信号以优雅关闭服务
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("正在关闭判题服务...")

	// 停止消费任务
	cancel()

	// 等待判题任务完成
	judgeManager.Shutdown()

	logger.Info("判题服务已关闭")
}
