package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Problem 题目模型 - 存储编程题目的完整信息
// 对应MongoDB集合: problems
type Problem struct {
	ID                 primitive.ObjectID `bson:"_id,omitempty" json:"id"`                          // 题目唯一标识ID
	Title              string             `bson:"title" json:"title"`                               // 题目标题
	Description        string             `bson:"description" json:"description"`                   // 题目描述（支持Markdown格式）
	InputFormat        string             `bson:"input_format" json:"input_format"`                 // 输入格式说明
	OutputFormat       string             `bson:"output_format" json:"output_format"`               // 输出格式说明
	SampleInput        string             `bson:"sample_input" json:"sample_input"`                 // 示例输入
	SampleOutput       string             `bson:"sample_output" json:"sample_output"`               // 示例输出
	TimeLimit          int                `bson:"time_limit" json:"time_limit"`                     // 默认时间限制（毫秒），对应go-judge的cpuLimit
	MemoryLimit        int                `bson:"memory_limit" json:"memory_limit"`                 // 默认内存限制（MB），对应go-judge的memoryLimit
	StackLimit         int                `bson:"stack_limit" json:"stack_limit"`                   // 栈空间限制（MB），对应go-judge的stackLimit
	OutputLimit        int                `bson:"output_limit" json:"output_limit"`                 // 输出大小限制(KB)，对应go-judge的copyOutMax
	CompileTimeLimit   int                `bson:"compile_time_limit" json:"compile_time_limit"`     // 编译时间限制（毫秒），默认10秒
	CompileMemoryLimit int                `bson:"compile_memory_limit" json:"compile_memory_limit"` // 编译内存限制（MB），默认128MB
	Difficulty         string             `bson:"difficulty" json:"difficulty"`                     // 难度级别: easy(简单), medium(中等), hard(困难)
	JudgeMode          string             `bson:"judge_mode" json:"judge_mode"`                     // 判题模式: interface(标准输入输出), acm(函数调用)
	Languages          []LanguageConfig   `bson:"languages" json:"languages"`                       // 支持的编程语言配置
	SpecialJudge       bool               `bson:"special_judge" json:"special_judge"`               // 是否启用特殊判题（SPJ）
	Tags               []string           `bson:"tags" json:"tags"`                                 // 题目标签（如：算法、数据结构、动态规划等）
	TestCases          []TestCase         `bson:"test_cases" json:"test_cases"`                     // 测试用例集合
	Stats              ProblemStats       `bson:"stats" json:"stats"`                               // 题目统计信息
	IsPublic           bool               `bson:"is_public" json:"is_public"`                       // 是否公开可见
	CreatedBy          primitive.ObjectID `bson:"created_by" json:"created_by"`                     // 创建者用户ID
	CreatedAt          time.Time          `bson:"created_at" json:"created_at"`                     // 创建时间
	UpdatedAt          time.Time          `bson:"updated_at" json:"updated_at"`                     // 最后更新时间
}

// LanguageConfig 编程语言配置 - 针对go-judge沙箱的语言特定设置
type LanguageConfig struct {
	Language      string   `bson:"language" json:"language"`           // 语言名称（java/cpp/python/go等）
	CompileCmd    []string `bson:"compile_cmd" json:"compile_cmd"`     // 编译命令参数，如: ["/usr/bin/javac", "-encoding", "UTF-8", "-cp", "/w", "Main.java"]
	RunCmd        []string `bson:"run_cmd" json:"run_cmd"`             // 运行命令参数，如: ["/usr/bin/java", "-cp", "/w", "Main"]
	CompileEnv    []string `bson:"compile_env" json:"compile_env"`     // 编译环境变量，如: ["PATH=/usr/bin:/bin"]
	RunEnv        []string `bson:"run_env" json:"run_env"`             // 运行环境变量
	SourceFile    string   `bson:"source_file" json:"source_file"`     // 源文件名（如Main.java）
	Executable    string   `bson:"executable" json:"executable"`       // 可执行文件名（如Main.class）
	CompileTime   int64    `bson:"compile_time" json:"compile_time"`   // 编译时间限制（纳秒）
	CompileMemory int64    `bson:"compile_memory" json:"compile_memory"` // 编译内存限制（字节）
	IsEnabled     bool     `bson:"is_enabled" json:"is_enabled"`       // 是否启用该语言
}

// TestCase 测试用例模型 - 存储测试用例的完整信息
// 支持两种模式：Interface模式(标准输入输出)和ACM模式(函数调用)
type TestCase struct {
	ID          string            `bson:"id" json:"id"`                                   // 测试用例唯一标识
	Input       string            `bson:"input" json:"input"`                             // 测试用例输入数据
	Output      string            `bson:"output" json:"output"`                           // 期望输出结果
	Score       int               `bson:"score" json:"score"`                             // 该测试用例的分值
	IsPublic    bool              `bson:"is_public" json:"is_public"`                     // 是否公开给学生查看
	IsPretest   bool              `bson:"is_pretest" json:"is_pretest"`                   // 是否为预测试用例（编译后立即执行）
	Group       string            `bson:"group" json:"group"`                             // 测试用例分组（用于子任务）
	TimeLimit   int               `bson:"time_limit" json:"time_limit"`                   // 该测试用例的时间限制(毫秒)，可覆盖题目默认值
	MemoryLimit int               `bson:"memory_limit" json:"memory_limit"`               // 该测试用例的内存限制(MB)，可覆盖题目默认值
	FilePath    string            `bson:"file_path,omitempty" json:"file_path,omitempty"` // 大数据文件存放路径（可选）
	Metadata    map[string]string `bson:"metadata,omitempty" json:"metadata,omitempty"`   // 额外元数据
}

// ProblemStats 题目统计信息 - 记录题目的提交和通过数据
// 由后台异步任务定时更新，用于题目难度分析和推荐
type ProblemStats struct {
	TotalSubmissions int                     `bson:"total_submissions" json:"total_submissions"` // 总提交次数
	AcceptedCount    int                     `bson:"accepted_count" json:"accepted_count"`       // 通过次数
	AcceptanceRate   float64                 `bson:"acceptance_rate" json:"acceptance_rate"`     // 通过率（百分比）
	LanguageStats    map[string]LanguageStat `bson:"language_stats" json:"language_stats"`       // 各语言统计信息
	AvgTimeUsed      int                     `bson:"avg_time_used" json:"avg_time_used"`         // 平均执行时间（毫秒）
	AvgMemoryUsed    int                     `bson:"avg_memory_used" json:"avg_memory_used"`     // 平均内存使用（KB）
}

// LanguageStat 单语言统计信息
type LanguageStat struct {
	Submissions    int     `bson:"submissions" json:"submissions"`       // 该语言提交次数
	Accepted       int     `bson:"accepted" json:"accepted"`             // 该语言通过次数
	AcceptanceRate float64 `bson:"acceptance_rate" json:"acceptance_rate"` // 该语言通过率
	AvgTime        int     `bson:"avg_time" json:"avg_time"`             // 该语言平均执行时间
	AvgMemory      int     `bson:"avg_memory" json:"avg_memory"`         // 该语言平均内存使用
}