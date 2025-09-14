package sandbox

import (
	"fmt"
	"strings"
)

// StatusCodeMapping go-judge状态码映射
// 将go-judge的状态码映射为OJ系统的标准状态
var StatusCodeMapping = map[string]string{
	"Accepted":              "ACCEPTED",              // 正常完成
	"Memory Limit Exceeded": "MEMORY_LIMIT_EXCEEDED", // 内存超限
	"Time Limit Exceeded":   "TIME_LIMIT_EXCEEDED",   // 时间超限
	"Output Limit Exceeded": "OUTPUT_LIMIT_EXCEEDED", // 输出超限
	"File Error":            "SYSTEM_ERROR",          // 文件错误
	"Nonzero Exit Status":   "RUNTIME_ERROR",         // 非零退出(可能是编译错误或运行时错误)
	"Signalled":             "RUNTIME_ERROR",         // 信号终止
	"Dangerous Syscall":     "SYSTEM_ERROR",          // 危险系统调用
	"Internal Error":        "SYSTEM_ERROR",          // 内部错误
}

// MapStatus 映射go-judge状态码到OJ标准状态
// 将go-judge返回的状态码转换为OJ系统使用的标准状态码
//
// 参数:
//   - goJudgeStatus: go-judge返回的状态
//   - exitStatus: 程序退出状态码
//   - isCompileStage: 是否为编译阶段
//
// 返回:
//   - string: OJ标准状态码
//
// 示例:
//
//	status := MapStatus("Nonzero Exit Status", 1, true)
//	// 返回: "COMPILE_ERROR"
func MapStatus(goJudgeStatus string, exitStatus int, isCompileStage bool) string {
	// 特殊处理编译阶段的非零退出
	if goJudgeStatus == "Nonzero Exit Status" && isCompileStage {
		return "COMPILE_ERROR"
	}

	// 使用映射表
	if mapped, exists := StatusCodeMapping[goJudgeStatus]; exists {
		return mapped
	}

	// 默认返回系统错误
	return "SYSTEM_ERROR"
}

// ConvertTimeToMS 将纳秒转换为毫秒
// go-judge返回纳秒级时间，转换为OJ系统使用的毫秒
//
// 参数:
//   - nanoseconds: 纳秒时间
//
// 返回:
//   - int: 毫秒时间
//
// 示例:
//
//	ms := ConvertTimeToMS(123456789)
//	// 返回: 123
func ConvertTimeToMS(nanoseconds int64) int {
	return int(nanoseconds / 1000000)
}

// ConvertMemoryToKB 将字节转换为KB
// go-judge返回字节级内存，转换为OJ系统使用的KB
//
// 参数:
//   - bytes: 字节数
//
// 返回:
//   - int: KB数
//
// 示例:
//
//	kb := ConvertMemoryToKB(1048576)
//	// 返回: 1024
func ConvertMemoryToKB(bytes int64) int {
	return int(bytes / 1024)
}

// ConvertMSToNS 将毫秒转换为纳秒
// OJ系统的毫秒配置转换为go-judge需要的纳秒
//
// 参数:
//   - milliseconds: 毫秒时间
//
// 返回:
//   - int64: 纳秒时间
//
// 示例:
//
//	ns := ConvertMSToNS(1000)
//	// 返回: 1000000000
func ConvertMSToNS(milliseconds int) int64 {
	return int64(milliseconds) * 1000000
}

// ConvertMBToBytes 将MB转换为字节
// OJ系统的MB配置转换为go-judge需要的字节
//
// 参数:
//   - megabytes: MB数
//
// 返回:
//   - int64: 字节数
//
// 示例:
//
//	bytes := ConvertMBToBytes(128)
//	// 返回: 134217728
func ConvertMBToBytes(megabytes int) int64 {
	return int64(megabytes) * 1024 * 1024
}

// ConvertKBToBytes 将KB转换为字节
// OJ系统的KB配置转换为go-judge需要的字节
//
// 参数:
//   - kilobytes: KB数
//
// 返回:
//   - int64: 字节数
//
// 示例:
//
//	bytes := ConvertKBToBytes(10)
//	// 返回: 10240
func ConvertKBToBytes(kilobytes int) int64 {
	return int64(kilobytes) * 1024
}

// IsCompileError 判断是否为编译错误
// 根据go-judge返回的状态和退出码判断是否为编译错误
//
// 参数:
//   - status: go-judge状态
//   - exitStatus: 退出状态码
//   - stderr: 错误输出
//
// 返回:
//   - bool: 是否为编译错误
//
// 示例:
//
//	isCompileErr := IsCompileError("Nonzero Exit Status", 1, "syntax error")
//	// 返回: true
func IsCompileError(status string, exitStatus int, stderr string) bool {
	if status == "Nonzero Exit Status" && exitStatus != 0 {
		// 检查错误输出是否包含编译相关关键词
		compileKeywords := []string{"error:", "错误", "syntax", "cannot find symbol", "class", "package"}
		stderrLower := strings.ToLower(stderr)
		for _, keyword := range compileKeywords {
			if strings.Contains(stderrLower, keyword) {
				return true
			}
		}
	}
	return false
}

// IsRuntimeError 判断是否为运行时错误
// 根据go-judge返回的状态判断是否为运行时错误
//
// 参数:
//   - status: go-judge状态
//   - exitStatus: 退出状态码
//
// 返回:
//   - bool: 是否为运行时错误
//
// 示例:
//
//	isRuntimeErr := IsRuntimeError("Signalled", 11)
//	// 返回: true
func IsRuntimeError(status string, exitStatus int) bool {
	return status == "Signalled" ||
		(status == "Nonzero Exit Status" && exitStatus != 0)
}

// FormatErrorMessage 格式化错误信息
// 将go-judge的错误信息格式化为用户友好的提示
//
// 参数:
//   - status: go-judge状态
//   - stderr: 错误输出
//   - exitStatus: 退出状态码
//
// 返回:
//   - string: 格式化的错误信息
//
// 示例:
//
//	msg := FormatErrorMessage("Memory Limit Exceeded", "", 0)
//	// 返回: "程序内存使用超出限制，请优化算法减少内存占用"
func FormatErrorMessage(status string, stderr string, exitStatus int) string {
	switch status {
	case "Memory Limit Exceeded":
		return "程序内存使用超出限制，请优化算法减少内存占用"
	case "Time Limit Exceeded":
		return "程序运行时间超出限制，请优化算法提高执行效率"
	case "Output Limit Exceeded":
		return "程序输出内容过多，请检查是否存在无限循环输出"
	case "Nonzero Exit Status":
		if stderr != "" {
			return fmt.Sprintf("程序执行出错：%s", stderr)
		}
		return fmt.Sprintf("程序异常退出，退出码：%d", exitStatus)
	case "Signalled":
		return "程序运行时发生错误，可能存在数组越界、空指针等问题"
	case "Dangerous Syscall":
		return "程序尝试执行危险的系统调用，已被安全机制阻止"
	case "File Error":
		return "文件操作错误，请检查输入输出处理逻辑"
	case "Internal Error":
		return "系统内部错误，请稍后重试或联系管理员"
	default:
		return "未知错误，请检查代码并重试"
	}
}

// BuildJavaCompileRequest 构建Java编译请求
// 根据配置参数构建标准的Java编译请求
//
// 参数:
//   - sourceCode: Java源代码
//   - compileConfig: 编译配置
//
// 返回:
//   - *CompileRequest: 编译请求
//
// 示例:
//
//	config := JavaCompileConfig{
//	    Command: []string{"/usr/bin/javac", "-encoding", "UTF-8", "-cp", "/w", "Main.java"},
//	    Env: []string{"PATH=/usr/bin:/bin", "JAVA_HOME=/usr/lib/jvm/java-17-openjdk-amd64"},
//	    CPULimit: 10000000000,
//	    MemoryLimit: 268435456,
//	}
//	req := BuildJavaCompileRequest(code, config)
func BuildJavaCompileRequest(sourceCode string, compileConfig JavaCompileConfig) *CompileRequest {
	return &CompileRequest{
		SourceCode:  sourceCode,
		SourceFile:  "Main.java",
		CompileCmd:  compileConfig.Command,
		CompileEnv:  compileConfig.Env,
		CPULimit:    compileConfig.CPULimit,
		MemoryLimit: compileConfig.MemoryLimit,
		StackLimit:  compileConfig.StackLimit,
		ProcLimit:   compileConfig.ProcLimit,
		OutputLimit: ConvertKBToBytes(compileConfig.OutputLimit),
	}
}

// BuildJavaRunRequest 构建Java运行请求
// 根据配置参数构建标准的Java运行请求
//
// 参数:
//   - classFileID: 编译后的class文件ID
//   - input: 测试用例输入
//   - runtimeConfig: 运行时配置
//
// 返回:
//   - *RunTestRequest: 运行请求
//
// 示例:
//
//	config := JavaRuntimeConfig{
//	    Command: []string{"/usr/bin/java", "-cp", "/w", "Main"},
//	    Env: []string{"PATH=/usr/bin:/bin", "JAVA_HOME=/usr/lib/jvm/java-17-openjdk-amd64"},
//	    CPULimit: 5000000000,
//	    MemoryLimit: 134217728,
//	}
//	req := BuildJavaRunRequest("file123", "test input", config)
func BuildJavaRunRequest(classFileID, input string, runtimeConfig JavaRuntimeConfig) *RunTestRequest {
	return &RunTestRequest{
		ClassFileID: classFileID,
		Input:       input,
		RunCmd:      runtimeConfig.Command,
		RunEnv:      runtimeConfig.Env,
		CPULimit:    runtimeConfig.CPULimit,
		MemoryLimit: runtimeConfig.MemoryLimit,
		StackLimit:  runtimeConfig.StackLimit,
		ProcLimit:   runtimeConfig.ProcLimit,
		OutputLimit: ConvertKBToBytes(runtimeConfig.OutputLimit),
	}
}

// JavaCompileConfig Java编译配置
type JavaCompileConfig struct {
	Command     []string `yaml:"command"`      // 编译命令
	Env         []string `yaml:"env"`          // 环境变量
	CPULimit    int64    `yaml:"cpu_limit"`    // CPU限制(纳秒)
	MemoryLimit int64    `yaml:"memory_limit"` // 内存限制(字节)
	StackLimit  int64    `yaml:"stack_limit"`  // 栈限制(字节)
	ProcLimit   int      `yaml:"proc_limit"`   // 进程限制
	OutputLimit int      `yaml:"output_limit"` // 输出限制(KB)
}

// JavaRuntimeConfig Java运行时配置
type JavaRuntimeConfig struct {
	Command     []string `yaml:"command"`      // 运行命令
	Env         []string `yaml:"env"`          // 环境变量
	CPULimit    int64    `yaml:"cpu_limit"`    // CPU限制(纳秒)
	MemoryLimit int64    `yaml:"memory_limit"` // 内存限制(字节)
	StackLimit  int64    `yaml:"stack_limit"`  // 栈限制(字节)
	ProcLimit   int      `yaml:"proc_limit"`   // 进程限制
	OutputLimit int      `yaml:"output_limit"` // 输出限制(KB)
}
