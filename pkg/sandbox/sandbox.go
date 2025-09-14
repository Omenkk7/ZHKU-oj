package sandbox

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"zhku-oj/pkg/logger"
)

// SandboxClient go-judge沙箱客户端
// 提供与go-judge沙箱服务的完整接口，支持代码编译、运行、文件管理等功能
type SandboxClient struct {
	baseURL    string         // go-judge服务地址
	httpClient *http.Client   // HTTP客户端
	mutex      sync.RWMutex   // 读写锁
	isHealthy  bool           // 健康状态
	config     *SandboxConfig // 沙箱配置
}

// SandboxConfig 沙箱配置
type SandboxConfig struct {
	URL                 string        `yaml:"url"`                   // 沙箱服务地址
	Weight              int           `yaml:"weight"`                // 权重
	MaxConcurrent       int           `yaml:"max_concurrent"`        // 最大并发数
	Timeout             time.Duration `yaml:"timeout"`               // 请求超时时间
	HealthCheckInterval time.Duration `yaml:"health_check_interval"` // 健康检查间隔
	Enabled             bool          `yaml:"enabled"`               // 是否启用
	RetryTimes          int           `yaml:"retry_times"`           // 重试次数
	RetryInterval       time.Duration `yaml:"retry_interval"`        // 重试间隔
}

// RunRequest go-judge运行请求结构
// 对应go-judge的/run接口请求参数
type RunRequest struct {
	Cmd         []CommandConfig `json:"cmd"`                   // 命令配置数组
	PipeMapping []PipeMapping   `json:"pipeMapping,omitempty"` // 管道映射配置
}

// CommandConfig 单个命令配置
// 包含编译或运行一个程序所需的所有参数
type CommandConfig struct {
	Args              []string          `json:"args"`                        // 命令参数数组，如["/usr/bin/javac", "Main.java"]
	Env               []string          `json:"env,omitempty"`               // 环境变量，如["PATH=/usr/bin:/bin"]
	Files             []FileDescriptor  `json:"files"`                       // 文件描述符配置
	CPULimit          int64             `json:"cpuLimit"`                    // CPU时间限制(纳秒)
	ClockLimit        int64             `json:"clockLimit,omitempty"`        // 墙钟时间限制(纳秒)
	MemoryLimit       int64             `json:"memoryLimit"`                 // 内存限制(字节)
	StackLimit        int64             `json:"stackLimit,omitempty"`        // 栈限制(字节)
	ProcLimit         int               `json:"procLimit"`                   // 进程数限制
	CPURate           float64           `json:"cpuRate,omitempty"`           // CPU使用率限制
	StrictMemoryLimit bool              `json:"strictMemoryLimit,omitempty"` // 严格内存限制
	CopyIn            map[string]CopyIn `json:"copyIn,omitempty"`            // 输入文件映射(直接内容)
	CopyInCached      map[string]string `json:"copyInCached,omitempty"`      // 缓存文件映射(引用文件ID)
	CopyOut           []string          `json:"copyOut,omitempty"`           // 输出文件列表
	CopyOutCached     []string          `json:"copyOutCached,omitempty"`     // 缓存输出文件列表
	CopyOutMax        int64             `json:"copyOutMax,omitempty"`        // 输出文件大小限制
}

// FileDescriptor 文件描述符配置
// 定义标准输入/输出/错误的处理方式
type FileDescriptor struct {
	Content string `json:"content,omitempty"` // 文件内容(用于stdin)
	Name    string `json:"name,omitempty"`    // 文件名(如"stdout", "stderr")
	Max     int64  `json:"max,omitempty"`     // 文件大小限制
}

// CopyIn 输入文件配置
// 支持直接内容或引用缓存文件
type CopyIn struct {
	Content string `json:"content,omitempty"` // 直接文件内容
	FileID  string `json:"fileId,omitempty"`  // 缓存文件ID
}

// PipeMapping 管道映射配置
// 用于多命令间的数据传递
type PipeMapping struct {
	In  PipeEnd `json:"in"`  // 输入端
	Out PipeEnd `json:"out"` // 输出端
}

// PipeEnd 管道端点
type PipeEnd struct {
	Index int `json:"index"` // 命令索引
	FD    int `json:"fd"`    // 文件描述符
}

// RunResponse go-judge运行响应结构
// 对应go-judge的/run接口返回结果
type RunResponse struct {
	Status     string                 `json:"status"`              // 执行状态
	ExitStatus int                    `json:"exitStatus"`          // 退出状态码
	Time       int64                  `json:"time"`                // CPU时间(纳秒)
	Memory     int64                  `json:"memory"`              // 内存使用(字节)
	RunTime    int64                  `json:"runTime"`             // 运行时间(纳秒)
	Files      map[string]string      `json:"files,omitempty"`     // 输出文件内容
	FileIDs    map[string]string      `json:"fileIds,omitempty"`   // 缓存文件ID映射
	FileError  []FileError            `json:"fileError,omitempty"` // 文件错误信息
	ProcPeak   int                    `json:"procPeak,omitempty"`  // 峰值进程数
	Raw        map[string]interface{} `json:"-"`                   // 原始响应数据(用于调试)
}

// FileError 文件操作错误
type FileError struct {
	Name    string `json:"name"`    // 文件名
	Type    string `json:"type"`    // 错误类型
	Message string `json:"message"` // 错误信息
}

// CompileRequest 编译请求封装
// 基于go-judge的编译阶段参数封装
type CompileRequest struct {
	SourceCode  string   // Java源代码
	SourceFile  string   // 源文件名(如Main.java)
	CompileCmd  []string // 编译命令
	CompileEnv  []string // 编译环境变量
	CPULimit    int64    // CPU时间限制(纳秒)
	MemoryLimit int64    // 内存限制(字节)
	StackLimit  int64    // 栈限制(字节)
	ProcLimit   int      // 进程数限制
	OutputLimit int64    // 输出限制(字节)
}

// RunTestRequest 运行测试请求封装
// 基于go-judge的运行阶段参数封装
type RunTestRequest struct {
	ClassFileID string   // 编译生成的class文件ID
	Input       string   // 测试用例输入
	RunCmd      []string // 运行命令
	RunEnv      []string // 运行环境变量
	CPULimit    int64    // CPU时间限制(纳秒)
	MemoryLimit int64    // 内存限制(字节)
	StackLimit  int64    // 栈限制(字节)
	ProcLimit   int      // 进程数限制
	OutputLimit int64    // 输出限制(字节)
}

// NewSandboxClient 创建沙箱客户端
// 初始化与go-judge沙箱服务的连接
//
// 参数:
//   - config: 沙箱配置信息
//
// 返回:
//   - *SandboxClient: 沙箱客户端实例
//   - error: 错误信息
//
// 示例:
//
//	config := &SandboxConfig{
//	    URL: "http://localhost:5050",
//	    Timeout: 30 * time.Second,
//	    MaxConcurrent: 10,
//	}
//	client, err := NewSandboxClient(config)
func NewSandboxClient(config *SandboxConfig) (*SandboxClient, error) {
	if config == nil {
		return nil, fmt.Errorf("沙箱配置不能为空")
	}

	if config.URL == "" {
		return nil, fmt.Errorf("沙箱URL不能为空")
	}

	// 创建HTTP客户端
	httpClient := &http.Client{
		Timeout: config.Timeout,
		Transport: &http.Transport{
			MaxIdleConns:        100,
			IdleConnTimeout:     90 * time.Second,
			DisableCompression:  true,
			MaxIdleConnsPerHost: config.MaxConcurrent,
		},
	}

	client := &SandboxClient{
		baseURL:    config.URL,
		httpClient: httpClient,
		isHealthy:  true,
		config:     config,
	}

	// 启动健康检查
	go client.startHealthCheck()

	logger.Info("沙箱客户端创建成功", "url", config.URL)
	return client, nil
}

// IsHealthy 检查沙箱是否健康
// 返回沙箱服务的健康状态
//
// 返回:
//   - bool: 健康状态，true表示健康
//
// 示例:
//
//	if client.IsHealthy() {
//	    // 执行判题任务
//	}
func (c *SandboxClient) IsHealthy() bool {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.isHealthy
}

// CompileJava 编译Java代码
// 将Java源代码编译为class文件，对应go-judge的编译阶段
//
// 参数:
//   - ctx: 上下文
//   - req: 编译请求参数
//
// 返回:
//   - *RunResponse: 编译结果，包含class文件ID
//   - error: 错误信息
//
// 示例:
//
//	req := &CompileRequest{
//	    SourceCode: "public class Main { public static void main(String[] args) { System.out.println(\"Hello\"); } }",
//	    SourceFile: "Main.java",
//	    CompileCmd: []string{"/usr/bin/javac", "-encoding", "UTF-8", "-cp", "/w", "Main.java"},
//	    CompileEnv: []string{"PATH=/usr/bin:/bin", "JAVA_HOME=/usr/lib/jvm/java-17-openjdk-amd64"},
//	    CPULimit: 10000000000,
//	    MemoryLimit: 268435456,
//	}
//	result, err := client.CompileJava(ctx, req)
func (c *SandboxClient) CompileJava(ctx context.Context, req *CompileRequest) (*RunResponse, error) {
	// 构建go-judge编译请求
	runReq := &RunRequest{
		Cmd: []CommandConfig{
			{
				Args: req.CompileCmd,
				Env:  req.CompileEnv,
				Files: []FileDescriptor{
					{Content: ""},                          // stdin
					{Name: "stdout", Max: req.OutputLimit}, // stdout
					{Name: "stderr", Max: req.OutputLimit}, // stderr
				},
				CPULimit:    req.CPULimit,
				MemoryLimit: req.MemoryLimit,
				StackLimit:  req.StackLimit,
				ProcLimit:   req.ProcLimit,
				CopyIn: map[string]CopyIn{
					req.SourceFile: {Content: req.SourceCode},
				},
				CopyOut:       []string{"stdout", "stderr"},
				CopyOutCached: []string{"Main.class"},
				CopyOutMax:    req.OutputLimit,
			},
		},
	}

	// 执行编译请求
	responses, err := c.run(ctx, runReq)
	if err != nil {
		return nil, fmt.Errorf("编译请求失败: %w", err)
	}

	if len(responses) == 0 {
		return nil, fmt.Errorf("编译响应为空")
	}

	response := &responses[0]
	logger.Info("Java编译完成",
		"status", response.Status,
		"exit_status", response.ExitStatus,
		"time", response.Time,
		"memory", response.Memory,
	)

	return response, nil
}

// RunJava 运行Java程序
// 使用编译后的class文件运行Java程序，对应go-judge的运行阶段
//
// 参数:
//   - ctx: 上下文
//   - req: 运行请求参数
//
// 返回:
//   - *RunResponse: 运行结果，包含程序输出
//   - error: 错误信息
//
// 示例:
//
//	req := &RunTestRequest{
//	    ClassFileID: "ABC123DEF456",
//	    Input: "Hello World\n",
//	    RunCmd: []string{"/usr/bin/java", "-cp", "/w", "Main"},
//	    RunEnv: []string{"PATH=/usr/bin:/bin", "JAVA_HOME=/usr/lib/jvm/java-17-openjdk-amd64"},
//	    CPULimit: 5000000000,
//	    MemoryLimit: 134217728,
//	}
//	result, err := client.RunJava(ctx, req)
func (c *SandboxClient) RunJava(ctx context.Context, req *RunTestRequest) (*RunResponse, error) {
	// 构建go-judge运行请求
	runReq := &RunRequest{
		Cmd: []CommandConfig{
			{
				Args: req.RunCmd,
				Env:  req.RunEnv,
				Files: []FileDescriptor{
					{Content: req.Input},                   // stdin
					{Name: "stdout", Max: req.OutputLimit}, // stdout
					{Name: "stderr", Max: req.OutputLimit}, // stderr
				},
				CPULimit:    req.CPULimit,
				MemoryLimit: req.MemoryLimit,
				StackLimit:  req.StackLimit,
				ProcLimit:   req.ProcLimit,
				CopyInCached: map[string]string{
					"Main.class": req.ClassFileID, // 使用copyInCached字段引用缓存文件
				},
				CopyOut:    []string{"stdout", "stderr"},
				CopyOutMax: req.OutputLimit,
			},
		},
	}

	// 执行运行请求
	responses, err := c.run(ctx, runReq)
	if err != nil {
		return nil, fmt.Errorf("运行请求失败: %w", err)
	}

	if len(responses) == 0 {
		return nil, fmt.Errorf("运行响应为空")
	}

	response := &responses[0]
	logger.Info("Java运行完成",
		"status", response.Status,
		"exit_status", response.ExitStatus,
		"time", response.Time,
		"memory", response.Memory,
	)

	return response, nil
}

// GetVersion 获取go-judge版本信息
// 查询沙箱服务的版本和配置信息
//
// 参数:
//   - ctx: 上下文
//
// 返回:
//   - map[string]interface{}: 版本信息
//   - error: 错误信息
//
// 示例:
//
//	version, err := client.GetVersion(ctx)
//	if err == nil {
//	    fmt.Printf("go-judge版本: %v\n", version)
//	}
func (c *SandboxClient) GetVersion(ctx context.Context) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/version", c.baseURL)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("创建版本请求失败: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("获取版本信息失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("获取版本信息失败: HTTP %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("解析版本信息失败: %w", err)
	}

	return result, nil
}

// ListFiles 列出所有缓存文件
// 获取go-judge服务中缓存的所有文件列表
//
// 参数:
//   - ctx: 上下文
//
// 返回:
//   - []string: 文件ID列表
//   - error: 错误信息
//
// 示例:
//
//	files, err := client.ListFiles(ctx)
//	if err == nil {
//	    fmt.Printf("缓存文件: %v\n", files)
//	}
func (c *SandboxClient) ListFiles(ctx context.Context) ([]string, error) {
	url := fmt.Sprintf("%s/file", c.baseURL)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("创建文件列表请求失败: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("获取文件列表失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("获取文件列表失败: HTTP %d", resp.StatusCode)
	}

	var files []string
	if err := json.NewDecoder(resp.Body).Decode(&files); err != nil {
		return nil, fmt.Errorf("解析文件列表失败: %w", err)
	}

	return files, nil
}

// CompileAndRunJava 编译并运行Java代码的便捷方法
// 将编译和运行两个步骤合并为一个操作，自动处理文件缓存和清理
//
// 参数:
//   - ctx: 上下文
//   - sourceCode: Java源代码
//   - inputData: 测试用例输入数据
//   - compileConfig: 编译配置(可选，使用默认值)
//   - runConfig: 运行配置(可选，使用默认值)
//
// 返回:
//   - *JudgeResult: 统一的判题结果
//   - error: 错误信息
//
// 示例:
//
//	result, err := client.CompileAndRunJava(ctx, javaCode, "1 2\n", nil, nil)
//	if err == nil {
//	    fmt.Printf("状态: %s, 输出: %s\n", result.Status, result.Output)
//	}
func (c *SandboxClient) CompileAndRunJava(ctx context.Context, sourceCode, inputData string, compileConfig *CompileRequest, runConfig *RunTestRequest) (*JudgeResult, error) {
	// 使用默认编译配置
	if compileConfig == nil {
		compileConfig = &CompileRequest{
			SourceCode:  sourceCode,
			SourceFile:  "Main.java",
			CompileCmd:  []string{"/usr/bin/javac", "-encoding", "UTF-8", "-cp", "/w", "Main.java"},
			CompileEnv:  []string{"PATH=/usr/bin:/bin", "JAVA_HOME=/usr/lib/jvm/java-17-openjdk-amd64"},
			CPULimit:    10000000000, // 10秒
			MemoryLimit: 268435456,   // 256MB
			StackLimit:  134217728,   // 128MB
			ProcLimit:   50,
			OutputLimit: 10240, // 10KB
		}
	} else {
		compileConfig.SourceCode = sourceCode
		compileConfig.SourceFile = "Main.java"
	}

	// 步骤1: 编译Java代码
	compileResult, err := c.CompileJava(ctx, compileConfig)
	if err != nil {
		return &JudgeResult{
			Status:       "Internal Error",
			Success:      false,
			CompileError: fmt.Sprintf("编译请求失败: %v", err),
		}, err
	}

	// 检查编译状态
	if compileResult.Status != "Accepted" {
		return &JudgeResult{
			Status:       "Compile Error",
			Success:      false,
			CompileError: compileResult.Files["stderr"],
			TimeUsed:     compileResult.Time / 1000000, // 转换为毫秒
			MemoryUsed:   compileResult.Memory / 1024,  // 转换为KB
		}, nil
	}

	// 获取编译生成的class文件ID
	classFileID, exists := compileResult.FileIDs["Main.class"]
	if !exists {
		return &JudgeResult{
			Status:       "Compile Error",
			Success:      false,
			CompileError: "编译未生成Main.class文件",
		}, nil
	}

	// 使用默认运行配置
	if runConfig == nil {
		runConfig = &RunTestRequest{
			ClassFileID: classFileID,
			Input:       inputData,
			RunCmd:      []string{"/usr/bin/java", "-cp", "/w", "Main"},
			RunEnv:      []string{"PATH=/usr/bin:/bin", "JAVA_HOME=/usr/lib/jvm/java-17-openjdk-amd64"},
			CPULimit:    5000000000, // 5秒
			MemoryLimit: 134217728,  // 128MB
			StackLimit:  134217728,  // 128MB
			ProcLimit:   1,
			OutputLimit: 10240, // 10KB
		}
	} else {
		runConfig.ClassFileID = classFileID
		runConfig.Input = inputData
	}

	// 步骤2: 运行Java程序
	runResult, err := c.RunJava(ctx, runConfig)
	if err != nil {
		// 运行失败时尝试清理缓存文件
		c.DeleteFile(ctx, classFileID)
		return &JudgeResult{
			Status:       "Internal Error",
			Success:      false,
			CompileError: fmt.Sprintf("运行请求失败: %v", err),
		}, err
	}

	// 步骤3: 清理缓存文件
	go func() {
		cleanupCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := c.DeleteFile(cleanupCtx, classFileID); err != nil {
			logger.Warn("清理缓存文件失败", "file_id", classFileID, "error", err)
		}
	}()

	// 构建统一返回结果
	result := &JudgeResult{
		Status:        runResult.Status,
		Success:       runResult.Status == "Accepted",
		Output:        runResult.Files["stdout"],
		ErrorOutput:   runResult.Files["stderr"],
		TimeUsed:      (compileResult.Time + runResult.Time) / 1000000,         // 总耗时(毫秒)
		MemoryUsed:    maxInt64(compileResult.Memory, runResult.Memory) / 1024, // 峰值内存(KB)
		ExitStatus:    runResult.ExitStatus,
		CompileTime:   compileResult.Time / 1000000, // 编译时间(毫秒)
		RunTime:       runResult.RunTime / 1000000,  // 运行时间(毫秒)
		CompileMemory: compileResult.Memory / 1024,  // 编译内存(KB)
		RunMemory:     runResult.Memory / 1024,      // 运行内存(KB)
	}

	return result, nil
}

// JudgeResult 统一的判题结果结构
// 整合编译和运行阶段的所有结果信息
type JudgeResult struct {
	Status        string `json:"status"`         // 最终状态
	Success       bool   `json:"success"`        // 是否成功
	Output        string `json:"output"`         // 程序输出
	ErrorOutput   string `json:"error_output"`   // 错误输出
	CompileError  string `json:"compile_error"`  // 编译错误信息
	TimeUsed      int64  `json:"time_used"`      // 总耗时(毫秒)
	MemoryUsed    int64  `json:"memory_used"`    // 峰值内存使用(KB)
	ExitStatus    int    `json:"exit_status"`    // 程序退出码
	CompileTime   int64  `json:"compile_time"`   // 编译时间(毫秒)
	RunTime       int64  `json:"run_time"`       // 运行时间(毫秒)
	CompileMemory int64  `json:"compile_memory"` // 编译内存(KB)
	RunMemory     int64  `json:"run_memory"`     // 运行内存(KB)
}

// maxInt64 返回两个int64中的较大值
func maxInt64(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

// RunCode 通用代码执行方法
// 支持多种编程语言的编译和运行
//
// 参数:
//   - ctx: 上下文
//   - req: 通用执行请求
//
// 返回:
//   - *JudgeResult: 统一的判题结果
//   - error: 错误信息
//
// 示例:
//
//	req := &CodeExecutionRequest{
//	    Language: "java",
//	    SourceCode: "public class Main { ... }",
//	    Input: "1 2\n",
//	    TimeLimit: 1000,
//	    MemoryLimit: 128,
//	}
//	result, err := client.RunCode(ctx, req)
func (c *SandboxClient) RunCode(ctx context.Context, req *CodeExecutionRequest) (*JudgeResult, error) {
	switch req.Language {
	case "java":
		// Java代码执行
		compileConfig := &CompileRequest{
			SourceCode:  req.SourceCode,
			SourceFile:  "Main.java",
			CompileCmd:  []string{"/usr/bin/javac", "-encoding", "UTF-8", "-cp", "/w", "Main.java"},
			CompileEnv:  []string{"PATH=/usr/bin:/bin", "JAVA_HOME=/usr/lib/jvm/java-17-openjdk-amd64"},
			CPULimit:    req.CompileTimeLimit * 1000000,       // 转换为纳秒
			MemoryLimit: req.CompileMemoryLimit * 1024 * 1024, // 转换为字节
			StackLimit:  req.CompileMemoryLimit * 1024 * 1024 / 2,
			ProcLimit:   50,
			OutputLimit: 10240,
		}

		runConfig := &RunTestRequest{
			Input:       req.Input,
			RunCmd:      []string{"/usr/bin/java", "-cp", "/w", "Main"},
			RunEnv:      []string{"PATH=/usr/bin:/bin", "JAVA_HOME=/usr/lib/jvm/java-17-openjdk-amd64"},
			CPULimit:    req.TimeLimit * 1000000,       // 转换为纳秒
			MemoryLimit: req.MemoryLimit * 1024 * 1024, // 转换为字节
			StackLimit:  req.MemoryLimit * 1024 * 1024 / 2,
			ProcLimit:   1,
			OutputLimit: 10240,
		}

		return c.CompileAndRunJava(ctx, req.SourceCode, req.Input, compileConfig, runConfig)

	case "python", "python3":
		// Python代码执行(无编译阶段)
		return c.runPython(ctx, req)

	case "cpp", "c++":
		// C++代码执行
		return c.runCpp(ctx, req)

	case "c":
		// C代码执行
		return c.runC(ctx, req)

	default:
		return &JudgeResult{
			Status:       "Unsupported Language",
			Success:      false,
			CompileError: fmt.Sprintf("不支持的编程语言: %s", req.Language),
		}, fmt.Errorf("不支持的编程语言: %s", req.Language)
	}
}

// CodeExecutionRequest 通用代码执行请求
type CodeExecutionRequest struct {
	Language           string `json:"language"`             // 编程语言
	SourceCode         string `json:"source_code"`          // 源代码
	Input              string `json:"input"`                // 输入数据
	TimeLimit          int64  `json:"time_limit"`           // 运行时间限制(毫秒)
	MemoryLimit        int64  `json:"memory_limit"`         // 运行内存限制(MB)
	CompileTimeLimit   int64  `json:"compile_time_limit"`   // 编译时间限制(毫秒)
	CompileMemoryLimit int64  `json:"compile_memory_limit"` // 编译内存限制(MB)
	OutputLimit        int64  `json:"output_limit"`         // 输出限制(字节)
}

// runPython 执行Python代码
func (c *SandboxClient) runPython(ctx context.Context, req *CodeExecutionRequest) (*JudgeResult, error) {
	// Python代码直接运行，无编译阶段
	runReq := &RunRequest{
		Cmd: []CommandConfig{
			{
				Args: []string{"/usr/bin/python3", "main.py"},
				Env:  []string{"PATH=/usr/bin:/bin", "PYTHONPATH=/usr/lib/python3/dist-packages"},
				Files: []FileDescriptor{
					{Content: req.Input},                   // stdin
					{Name: "stdout", Max: req.OutputLimit}, // stdout
					{Name: "stderr", Max: req.OutputLimit}, // stderr
				},
				CPULimit:    req.TimeLimit * 1000000,       // 转换为纳秒
				MemoryLimit: req.MemoryLimit * 1024 * 1024, // 转换为字节
				ProcLimit:   1,
				CopyIn: map[string]CopyIn{
					"main.py": {Content: req.SourceCode},
				},
				CopyOut:    []string{"stdout", "stderr"},
				CopyOutMax: req.OutputLimit,
			},
		},
	}

	responses, err := c.run(ctx, runReq)
	if err != nil {
		return &JudgeResult{
			Status:       "Internal Error",
			Success:      false,
			CompileError: fmt.Sprintf("Python运行失败: %v", err),
		}, err
	}

	if len(responses) == 0 {
		return &JudgeResult{
			Status:       "Internal Error",
			Success:      false,
			CompileError: "Python运行响应为空",
		}, fmt.Errorf("Python运行响应为空")
	}

	response := &responses[0]
	return &JudgeResult{
		Status:      response.Status,
		Success:     response.Status == "Accepted",
		Output:      response.Files["stdout"],
		ErrorOutput: response.Files["stderr"],
		TimeUsed:    response.Time / 1000000, // 转换为毫秒
		MemoryUsed:  response.Memory / 1024,  // 转换为KB
		ExitStatus:  response.ExitStatus,
		RunTime:     response.RunTime / 1000000, // 转换为毫秒
		RunMemory:   response.Memory / 1024,     // 转换为KB
	}, nil
}

// runCpp 执行C++代码
func (c *SandboxClient) runCpp(ctx context.Context, req *CodeExecutionRequest) (*JudgeResult, error) {
	// C++编译和运行的一体化请求
	runReq := &RunRequest{
		Cmd: []CommandConfig{
			// 编译阶段
			{
				Args: []string{"/usr/bin/g++", "-o", "main", "main.cpp", "-std=c++17"},
				Env:  []string{"PATH=/usr/bin:/bin"},
				Files: []FileDescriptor{
					{Content: ""},                          // stdin
					{Name: "stdout", Max: req.OutputLimit}, // stdout
					{Name: "stderr", Max: req.OutputLimit}, // stderr
				},
				CPULimit:    req.CompileTimeLimit * 1000000,       // 转换为纳秒
				MemoryLimit: req.CompileMemoryLimit * 1024 * 1024, // 转换为字节
				ProcLimit:   50,
				CopyIn: map[string]CopyIn{
					"main.cpp": {Content: req.SourceCode},
				},
				CopyOut:       []string{"stdout", "stderr"},
				CopyOutCached: []string{"main"},
				CopyOutMax:    req.OutputLimit,
			},
			// 运行阶段
			{
				Args: []string{"./main"},
				Env:  []string{"PATH=/usr/bin:/bin"},
				Files: []FileDescriptor{
					{Content: req.Input},                   // stdin
					{Name: "stdout", Max: req.OutputLimit}, // stdout
					{Name: "stderr", Max: req.OutputLimit}, // stderr
				},
				CPULimit:    req.TimeLimit * 1000000,       // 转换为纳秒
				MemoryLimit: req.MemoryLimit * 1024 * 1024, // 转换为字节
				ProcLimit:   1,
				CopyOut:     []string{"stdout", "stderr"},
				CopyOutMax:  req.OutputLimit,
			},
		},
		// 管道映射: 编译阶段的输出传递给运行阶段
		PipeMapping: []PipeMapping{
			{In: PipeEnd{Index: 0, FD: 1}, Out: PipeEnd{Index: 1, FD: 0}},
		},
	}

	// 执行请求
	responses, err := c.run(ctx, runReq)
	if err != nil {
		return &JudgeResult{
			Status:       "Internal Error",
			Success:      false,
			CompileError: fmt.Sprintf("C++执行失败: %v", err),
		}, err
	}

	if len(responses) < 2 {
		return &JudgeResult{
			Status:       "Internal Error",
			Success:      false,
			CompileError: "C++执行响应不完整",
		}, fmt.Errorf("C++执行响应不完整")
	}

	compileResponse := &responses[0]
	runResponse := &responses[1]

	// 检查编译状态
	if compileResponse.Status != "Accepted" {
		return &JudgeResult{
			Status:        "Compile Error",
			Success:       false,
			CompileError:  compileResponse.Files["stderr"],
			CompileTime:   compileResponse.Time / 1000000,
			CompileMemory: compileResponse.Memory / 1024,
		}, nil
	}

	return &JudgeResult{
		Status:        runResponse.Status,
		Success:       runResponse.Status == "Accepted",
		Output:        runResponse.Files["stdout"],
		ErrorOutput:   runResponse.Files["stderr"],
		TimeUsed:      (compileResponse.Time + runResponse.Time) / 1000000,
		MemoryUsed:    maxInt64(compileResponse.Memory, runResponse.Memory) / 1024,
		ExitStatus:    runResponse.ExitStatus,
		CompileTime:   compileResponse.Time / 1000000,
		RunTime:       runResponse.RunTime / 1000000,
		CompileMemory: compileResponse.Memory / 1024,
		RunMemory:     runResponse.Memory / 1024,
	}, nil
}

// runC 执行C代码
func (c *SandboxClient) runC(ctx context.Context, req *CodeExecutionRequest) (*JudgeResult, error) {
	// C编译和运行的一体化请求
	runReq := &RunRequest{
		Cmd: []CommandConfig{
			// 编译阶段
			{
				Args: []string{"/usr/bin/gcc", "-o", "main", "main.c", "-std=c11"},
				Env:  []string{"PATH=/usr/bin:/bin"},
				Files: []FileDescriptor{
					{Content: ""},                          // stdin
					{Name: "stdout", Max: req.OutputLimit}, // stdout
					{Name: "stderr", Max: req.OutputLimit}, // stderr
				},
				CPULimit:    req.CompileTimeLimit * 1000000,       // 转换为纳秒
				MemoryLimit: req.CompileMemoryLimit * 1024 * 1024, // 转换为字节
				ProcLimit:   50,
				CopyIn: map[string]CopyIn{
					"main.c": {Content: req.SourceCode},
				},
				CopyOut:       []string{"stdout", "stderr"},
				CopyOutCached: []string{"main"},
				CopyOutMax:    req.OutputLimit,
			},
			// 运行阶段
			{
				Args: []string{"./main"},
				Env:  []string{"PATH=/usr/bin:/bin"},
				Files: []FileDescriptor{
					{Content: req.Input},                   // stdin
					{Name: "stdout", Max: req.OutputLimit}, // stdout
					{Name: "stderr", Max: req.OutputLimit}, // stderr
				},
				CPULimit:    req.TimeLimit * 1000000,       // 转换为纳秒
				MemoryLimit: req.MemoryLimit * 1024 * 1024, // 转换为字节
				ProcLimit:   1,
				CopyOut:     []string{"stdout", "stderr"},
				CopyOutMax:  req.OutputLimit,
			},
		},
		// 管道映射: 编译阶段的输出传递给运行阶段
		PipeMapping: []PipeMapping{
			{In: PipeEnd{Index: 0, FD: 1}, Out: PipeEnd{Index: 1, FD: 0}},
		},
	}

	// 执行请求
	responses, err := c.run(ctx, runReq)
	if err != nil {
		return &JudgeResult{
			Status:       "Internal Error",
			Success:      false,
			CompileError: fmt.Sprintf("C执行失败: %v", err),
		}, err
	}

	if len(responses) < 2 {
		return &JudgeResult{
			Status:       "Internal Error",
			Success:      false,
			CompileError: "C执行响应不完整",
		}, fmt.Errorf("C执行响应不完整")
	}

	compileResponse := &responses[0]
	runResponse := &responses[1]

	// 检查编译状态
	if compileResponse.Status != "Accepted" {
		return &JudgeResult{
			Status:        "Compile Error",
			Success:       false,
			CompileError:  compileResponse.Files["stderr"],
			CompileTime:   compileResponse.Time / 1000000,
			CompileMemory: compileResponse.Memory / 1024,
		}, nil
	}

	return &JudgeResult{
		Status:        runResponse.Status,
		Success:       runResponse.Status == "Accepted",
		Output:        runResponse.Files["stdout"],
		ErrorOutput:   runResponse.Files["stderr"],
		TimeUsed:      (compileResponse.Time + runResponse.Time) / 1000000,
		MemoryUsed:    maxInt64(compileResponse.Memory, runResponse.Memory) / 1024,
		ExitStatus:    runResponse.ExitStatus,
		CompileTime:   compileResponse.Time / 1000000,
		RunTime:       runResponse.RunTime / 1000000,
		CompileMemory: compileResponse.Memory / 1024,
		RunMemory:     runResponse.Memory / 1024,
	}, nil
}

// GetConfig 获取go-judge配置信息
// 查询沙箱服务的系统配置信息
//
// 参数:
//   - ctx: 上下文
//
// 返回:
//   - map[string]interface{}: 配置信息
//   - error: 错误信息
//
// 示例:
//
//	config, err := client.GetConfig(ctx)
//	if err == nil {
//	    fmt.Printf("go-judge配置: %v\n", config)
//	}
func (c *SandboxClient) GetConfig(ctx context.Context) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/config", c.baseURL)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("创建配置请求失败: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("获取配置信息失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("获取配置信息失败: HTTP %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("解析配置信息失败: %w", err)
	}

	return result, nil
}

// UploadFile 上传文件到缓存
// 将文件内容上传到go-judge服务的缓存中
//
// 参数:
//   - ctx: 上下文
//   - content: 文件内容
//
// 返回:
//   - string: 文件ID
//   - error: 错误信息
//
// 示例:
//
//	fileID, err := client.UploadFile(ctx, "file content")
//	if err == nil {
//	    fmt.Printf("文件ID: %s\n", fileID)
//	}
func (c *SandboxClient) UploadFile(ctx context.Context, content string) (string, error) {
	url := fmt.Sprintf("%s/file", c.baseURL)

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader([]byte(content)))
	if err != nil {
		return "", fmt.Errorf("创建上传请求失败: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("上传文件失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("上传文件失败: HTTP %d, %s", resp.StatusCode, string(body))
	}

	// 读取返回的文件ID
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取上传响应失败: %w", err)
	}

	return string(body), nil
}

// DownloadFile 下载缓存文件
// 从 go-judge 服务中下载指定的缓存文件
//
// 参数:
//   - ctx: 上下文
//   - fileID: 文件ID
//
// 返回:
//   - []byte: 文件内容
//   - error: 错误信息
//
// 示例:
//
//	content, err := client.DownloadFile(ctx, "ABC123DEF456")
//	if err == nil {
//	    fmt.Printf("文件内容: %s\n", string(content))
//	}
func (c *SandboxClient) DownloadFile(ctx context.Context, fileID string) ([]byte, error) {
	url := fmt.Sprintf("%s/file/%s", c.baseURL, fileID)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("创建下载请求失败: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("下载文件失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("下载文件失败: HTTP %d", resp.StatusCode)
	}

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取文件内容失败: %w", err)
	}

	return content, nil
}

// DeleteFile 删除缓存文件
// 从go-judge服务中删除指定的缓存文件，释放存储空间
//
// 参数:
//   - ctx: 上下文
//   - fileID: 文件ID
//
// 返回:
//   - error: 错误信息
//
// 示例:
//
//	err := client.DeleteFile(ctx, "ABC123DEF456")
//	if err != nil {
//	    log.Printf("删除文件失败: %v", err)
//	}
func (c *SandboxClient) DeleteFile(ctx context.Context, fileID string) error {
	url := fmt.Sprintf("%s/file/%s", c.baseURL, fileID)

	req, err := http.NewRequestWithContext(ctx, "DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("创建删除请求失败: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("删除文件失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNotFound {
		return fmt.Errorf("删除文件失败: HTTP %d", resp.StatusCode)
	}

	logger.Info("文件删除成功", "file_id", fileID)
	return nil
}

// run 执行go-judge请求的核心方法
// 发送HTTP请求到go-judge的/run接口
func (c *SandboxClient) run(ctx context.Context, req *RunRequest) ([]RunResponse, error) {
	// 序列化请求
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("序列化请求失败: %w", err)
	}

	// 创建HTTP请求
	url := fmt.Sprintf("%s/run", c.baseURL)
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("创建HTTP请求失败: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	// 执行请求
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		// 标记沙箱不健康
		c.mutex.Lock()
		c.isHealthy = false
		c.mutex.Unlock()
		return nil, fmt.Errorf("执行HTTP请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("请求失败: HTTP %d, 响应: %s", resp.StatusCode, string(body))
	}

	// 解析响应
	var responses []RunResponse
	if err := json.Unmarshal(body, &responses); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w, 响应: %s", err, string(body))
	}

	// 保存原始响应数据用于调试
	for i := range responses {
		var raw map[string]interface{}
		json.Unmarshal(body, &raw)
		responses[i].Raw = raw
	}

	return responses, nil
}

// startHealthCheck 启动健康检查协程
// 定期检查沙箱服务的健康状态
func (c *SandboxClient) startHealthCheck() {
	if c.config.HealthCheckInterval <= 0 {
		return
	}

	ticker := time.NewTicker(c.config.HealthCheckInterval)
	defer ticker.Stop()

	for range ticker.C {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		healthy := c.checkHealth(ctx)
		cancel()

		c.mutex.Lock()
		c.isHealthy = healthy
		c.mutex.Unlock()

		if !healthy {
			logger.Warn("沙箱健康检查失败", "url", c.baseURL)
		}
	}
}

// checkHealth 检查沙箱健康状态
// 通过访问/version接口来检查服务是否正常
func (c *SandboxClient) checkHealth(ctx context.Context) bool {
	_, err := c.GetVersion(ctx)
	return err == nil
}

// GetBaseURL 获取沙箱服务地址
func (c *SandboxClient) GetBaseURL() string {
	return c.baseURL
}

// GetClientConfig 获取客户端配置
func (c *SandboxClient) GetClientConfig() *SandboxConfig {
	return c.config
}

// Close 关闭沙箱客户端
func (c *SandboxClient) Close() error {
	c.httpClient.CloseIdleConnections()
	logger.Info("沙箱客户端已关闭", "url", c.baseURL)
	return nil
}
