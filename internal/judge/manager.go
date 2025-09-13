package judge

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"zhku-oj/internal/config"
	"zhku-oj/internal/model"
	"zhku-oj/internal/pkg/logger"
	"zhku-oj/internal/repository/interfaces"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Manager 判题管理器
type Manager struct {
	cfg            config.JudgeConfig
	submissionRepo interfaces.SubmissionRepository
	problemRepo    interfaces.ProblemRepository
	balancer       *Balancer
	fileManager    *FileManager
	processor      *ResultProcessor
	wg             sync.WaitGroup
	shutdown       chan struct{}
}

// NewManager 创建判题管理器
func NewManager(
	cfg config.JudgeConfig,
	submissionRepo interfaces.SubmissionRepository,
	problemRepo interfaces.ProblemRepository,
) (*Manager, error) {
	// 创建沙箱负载均衡器
	balancer, err := NewBalancer(cfg.Sandboxes)
	if err != nil {
		return nil, fmt.Errorf("创建负载均衡器失败: %w", err)
	}

	// 创建文件管理器
	fileManager := NewFileManager(cfg.FileManagement)

	// 创建结果处理器
	processor := NewResultProcessor()

	return &Manager{
		cfg:            cfg,
		submissionRepo: submissionRepo,
		problemRepo:    problemRepo,
		balancer:       balancer,
		fileManager:    fileManager,
		processor:      processor,
		shutdown:       make(chan struct{}),
	}, nil
}

// ProcessTask 处理判题任务
// 接收代码提交任务，执行Java代码编译和运行
func (m *Manager) ProcessTask(ctx context.Context, task *JudgeTask) error {
	logger.Info("开始处理判题任务", "submission_id", task.SubmissionID.Hex())

	// 更新提交状态为判题中
	if err := m.updateSubmissionStatus(ctx, task.SubmissionID, model.StatusJudging); err != nil {
		logger.Error("更新提交状态失败", "error", err)
		return err
	}

	// 获取题目信息
	problem, err := m.problemRepo.GetByID(ctx, task.ProblemID)
	if err != nil {
		logger.Error("获取题目信息失败", "error", err)
		return err
	}

	// 选择可用的沙箱实例
	sandbox := m.balancer.SelectSandbox()
	if sandbox == nil {
		logger.Error("没有可用的沙箱实例")
		return fmt.Errorf("没有可用的沙箱实例")
	}

	// 创建Java判题器
	javaJudge := NewJavaJudge(sandbox, m.cfg.Compile.Java, m.cfg.Runtime.Java)

	// 执行判题
	result, err := m.executeJudge(ctx, javaJudge, task, problem)
	if err != nil {
		logger.Error("执行判题失败", "error", err)
		// 更新为系统错误
		m.updateSubmissionWithError(ctx, task.SubmissionID, model.StatusSystemError, err.Error())
		return err
	}

	// 更新提交结果
	if err := m.updateSubmissionResult(ctx, task.SubmissionID, result); err != nil {
		logger.Error("更新提交结果失败", "error", err)
		return err
	}

	logger.Info("判题任务完成", "submission_id", task.SubmissionID.Hex(), "status", result.Status)
	return nil
}

// executeJudge 执行判题逻辑
func (m *Manager) executeJudge(ctx context.Context, judge *JavaJudge, task *JudgeTask, problem *model.Problem) (*JudgeResult, error) {
	// 1. 编译Java代码
	compileResult, err := judge.Compile(ctx, task.Code)
	if err != nil {
		return nil, fmt.Errorf("编译失败: %w", err)
	}

	// 检查编译是否成功
	if compileResult.Status != "Accepted" {
		return &JudgeResult{
			Status:     model.StatusCompileError,
			Score:      0,
			TimeUsed:   int(compileResult.Time / 1000000), // 纳秒转毫秒
			MemoryUsed: int(compileResult.Memory / 1024),  // 字节转KB
			CompileInfo: model.CompileInfo{
				Status:     "FAILED",
				TimeUsed:   compileResult.Time,
				MemoryUsed: compileResult.Memory,
				Message:    compileResult.ErrorMessage,
			},
			TestResults: []model.TestResult{},
			JudgedAt:    time.Now(),
		}, nil
	}

	// 编译成功，记录编译信息
	compileInfo := model.CompileInfo{
		Status:     "SUCCESS",
		TimeUsed:   compileResult.Time,
		MemoryUsed: compileResult.Memory,
		Message:    "",
	}

	// 2. 运行测试用例
	testResults := make([]model.TestResult, 0, len(problem.TestCases))
	totalScore := 0
	maxTime := 0
	maxMemory := 0

	for _, testCase := range problem.TestCases {
		runResult, err := judge.Run(ctx, compileResult.ClassFileID, testCase.Input)
		if err != nil {
			logger.Error("运行测试用例失败", "test_case_id", testCase.ID, "error", err)
			continue
		}

		// 比对输出结果
		status := model.StatusWrongAnswer
		score := 0
		if strings.TrimSpace(runResult.Output) == strings.TrimSpace(testCase.Output) {
			status = model.StatusAccepted
			score = testCase.Score
			totalScore += score
		}

		// 检查时间和内存限制
		timeUsed := int(runResult.Time / 1000000)  // 纳秒转毫秒
		memoryUsed := int(runResult.Memory / 1024) // 字节转KB

		if runResult.Status == "Time Limit Exceeded" {
			status = model.StatusTimeLimitExceeded
			score = 0
		} else if runResult.Status == "Memory Limit Exceeded" {
			status = model.StatusMemoryLimitExceeded
			score = 0
		} else if runResult.Status != "Accepted" {
			status = model.StatusRuntimeError
			score = 0
		}

		if timeUsed > maxTime {
			maxTime = timeUsed
		}
		if memoryUsed > maxMemory {
			maxMemory = memoryUsed
		}

		testResult := model.TestResult{
			TestCaseID:     testCase.ID,
			Status:         status,
			TimeUsed:       timeUsed,
			MemoryUsed:     memoryUsed,
			Score:          score,
			Input:          testCase.Input,
			ExpectedOutput: testCase.Output,
			ActualOutput:   runResult.Output,
			JudgeDetails: model.JudgeDetail{
				GoJudgeStatus: runResult.Status,
				ExitStatus:    runResult.ExitStatus,
				RuntimeNS:     runResult.Time,
			},
		}
		testResults = append(testResults, testResult)
	}

	// 3. 清理缓存文件
	if err := judge.CleanupFile(ctx, compileResult.ClassFileID); err != nil {
		logger.Error("清理缓存文件失败", "file_id", compileResult.ClassFileID, "error", err)
	}

	// 4. 计算最终状态
	finalStatus := model.StatusAccepted
	if totalScore == 0 {
		finalStatus = model.StatusWrongAnswer
		// 检查是否有特殊错误状态
		for _, result := range testResults {
			if result.Status == model.StatusTimeLimitExceeded ||
				result.Status == model.StatusMemoryLimitExceeded ||
				result.Status == model.StatusRuntimeError {
				finalStatus = result.Status
				break
			}
		}
	}

	return &JudgeResult{
		Status:      finalStatus,
		Score:       totalScore,
		TimeUsed:    maxTime,
		MemoryUsed:  maxMemory,
		CompileInfo: compileInfo,
		TestResults: testResults,
		JudgedAt:    time.Now(),
	}, nil
}

// ProcessResult 处理判题结果
func (m *Manager) ProcessResult(ctx context.Context, result *JudgeResult) error {
	// 这里可以添加结果后处理逻辑
	// 比如更新统计信息、发送通知等
	return nil
}

// updateSubmissionStatus 更新提交状态
func (m *Manager) updateSubmissionStatus(ctx context.Context, submissionID primitive.ObjectID, status string) error {
	return m.submissionRepo.UpdateStatus(ctx, submissionID, status)
}

// updateSubmissionResult 更新提交结果
func (m *Manager) updateSubmissionResult(ctx context.Context, submissionID primitive.ObjectID, result *JudgeResult) error {
	submission := &model.Submission{
		ID:          submissionID,
		Status:      result.Status,
		Score:       result.Score,
		TimeUsed:    result.TimeUsed,
		MemoryUsed:  result.MemoryUsed,
		CompileInfo: result.CompileInfo,
		TestResults: result.TestResults,
		JudgedAt:    &result.JudgedAt,
	}
	return m.submissionRepo.UpdateResult(ctx, submission)
}

// updateSubmissionWithError 更新提交错误状态
func (m *Manager) updateSubmissionWithError(ctx context.Context, submissionID primitive.ObjectID, status, errorMsg string) error {
	compileInfo := model.CompileInfo{
		Status:  "FAILED",
		Message: errorMsg,
	}
	submission := &model.Submission{
		ID:          submissionID,
		Status:      status,
		Score:       0,
		CompileInfo: compileInfo,
		JudgedAt:    &time.Time{},
	}
	*submission.JudgedAt = time.Now()
	return m.submissionRepo.UpdateResult(ctx, submission)
}

// Shutdown 优雅关闭
func (m *Manager) Shutdown() {
	close(m.shutdown)
	m.wg.Wait()
	logger.Info("判题管理器已关闭")
}
