package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Submission 提交记录模型 - 存储用户代码提交和判题结果的完整信息
// 对应MongoDB集合: submissions
type Submission struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`                        // 提交记录唯一标识ID
	UserID      primitive.ObjectID `bson:"user_id" json:"user_id"`                         // 提交用户ID，关联users集合
	ProblemID   primitive.ObjectID `bson:"problem_id" json:"problem_id"`                   // 题目ID，关联problems集合
	Code        string             `bson:"code" json:"code"`                               // 用户提交的源代码
	Language    string             `bson:"language" json:"language"`                       // 编程语言（java/go/python等）
	Status      string             `bson:"status" json:"status"`                           // 判题状态（PENDING/ACCEPTED/WRONG_ANSWER等）
	Score       int                `bson:"score" json:"score"`                             // 得分（0-100分）
	TimeUsed    int                `bson:"time_used" json:"time_used"`                     // 程序运行时间（毫秒）
	MemoryUsed  int                `bson:"memory_used" json:"memory_used"`                 // 程序使用内存（KB）
	CompileInfo CompileInfo        `bson:"compile_info" json:"compile_info"`               // 编译过程信息
	TestResults []TestResult       `bson:"test_results" json:"test_results"`               // 各测试用例运行结果
	JudgeInfo   JudgeInfo          `bson:"judge_info" json:"judge_info"`                   // 判题过程详细信息
	SubmittedAt time.Time          `bson:"submitted_at" json:"submitted_at"`               // 代码提交时间
	JudgedAt    *time.Time         `bson:"judged_at,omitempty" json:"judged_at,omitempty"` // 判题完成时间
}

// CompileInfo 编译信息 - 存储代码编译过程的状态和输出信息
// 与go-judge沙箱编译阶段的返回结果对应
type CompileInfo struct {
	Status     string `bson:"status" json:"status"`           // 编译状态：Accepted(成功)/Compile Error(失败)
	Message    string `bson:"message" json:"message"`         // 编译器输出信息（错误信息或警告）
	TimeUsed   int64  `bson:"time_used" json:"time_used"`     // 编译耗时（纳秒）
	MemoryUsed int64  `bson:"memory_used" json:"memory_used"` // 编译内存使用（字节）
	ExitStatus int    `bson:"exit_status" json:"exit_status"` // 编译器退出状态码
	FileId     string `bson:"file_id" json:"file_id"`         // go-judge返回的编译文件缓存ID
}

// TestResult 测试结果 - 存储单个测试用例的运行结果
// 与go-judge沙箱运行阶段的返回结果对应
type TestResult struct {
	TestCaseID string      `bson:"test_case_id" json:"test_case_id"`   // 测试用例ID，关联TestCase.ID
	Status     string      `bson:"status" json:"status"`               // 运行状态（Accepted/Wrong Answer/Time Limit Exceeded/Memory Limit Exceeded/Runtime Error等）
	TimeUsed   int         `bson:"time_used" json:"time_used"`         // 该测试用例运行时间（毫秒）
	MemoryUsed int         `bson:"memory_used" json:"memory_used"`     // 该测试用例使用内存（KB）
	Score      int         `bson:"score" json:"score"`                 // 该测试用例得分
	Output     string      `bson:"output" json:"output"`               // 程序实际输出（可选，用于错误分析）
	ErrorMsg   string      `bson:"error_msg" json:"error_msg"`         // 错误信息（如果有）
	Details    JudgeDetail `bson:"judge_details" json:"judge_details"` // go-judge判题器详细信息
}

// JudgeDetail go-judge判题器详细信息 - 存储底层判题系统的原始数据
// 完整保存go-judge返回的原始信息，用于问题排查和性能分析
type JudgeDetail struct {
	GoJudgeStatus string `bson:"go_judge_status" json:"go_judge_status"` // go-judge返回的原始状态码
	ExitStatus    int    `bson:"exit_status" json:"exit_status"`         // 程序退出状态码（0表示正常退出）
	RuntimeNS     int64  `bson:"runtime_ns" json:"runtime_ns"`           // 程序运行时间（纳秒级精度）
	MemoryBytes   int64  `bson:"memory_bytes" json:"memory_bytes"`       // 内存使用量（字节级精度）
	ProcPeak      int    `bson:"proc_peak" json:"proc_peak"`             // 峰值进程数
	FileError     string `bson:"file_error" json:"file_error"`           // 文件操作错误信息（如果有）
}

// JudgeInfo 判题信息 - 存储整个判题过程的元信息
// 用于跟踪判题流程和WebSocket实时通知
type JudgeInfo struct {
	JudgeVersion   string    `bson:"judge_version" json:"judge_version"`       // 判题器版本
	JudgeStartTime time.Time `bson:"judge_start_time" json:"judge_start_time"` // 判题开始时间
	JudgeEndTime   time.Time `bson:"judge_end_time" json:"judge_end_time"`     // 判题结束时间
	TotalTime      int64     `bson:"total_time" json:"total_time"`             // 总判题耗时（毫秒）
	CompileStage   string    `bson:"compile_stage" json:"compile_stage"`       // 编译阶段状态
	TestStage      string    `bson:"test_stage" json:"test_stage"`             // 测试阶段状态
	PassedTests    int       `bson:"passed_tests" json:"passed_tests"`         // 通过的测试用例数
	TotalTests     int       `bson:"total_tests" json:"total_tests"`           // 总测试用例数
	JudgeMode      string    `bson:"judge_mode" json:"judge_mode"`             // 判题模式（interface/acm）
	SystemError    string    `bson:"system_error" json:"system_error"`         // 系统错误信息（如果有）
}
