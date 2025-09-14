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
	SubmittedAt time.Time          `bson:"submitted_at" json:"submitted_at"`               // 代码提交时间
	JudgedAt    *time.Time         `bson:"judged_at,omitempty" json:"judged_at,omitempty"` // 判题完成时间
}

// CompileInfo 编译信息 - 存储代码编译过程的状态和输出信息
// 编译失败时提供详细错误信息给用户调试
type CompileInfo struct {
	Status  string `bson:"status" json:"status"`   // 编译状态：success(成功)/error(失败)
	Message string `bson:"message" json:"message"` // 编译器输出信息（错误信息或警告）
}

// TestResult 测试结果 - 存储单个测试用例的运行结果
// 为节省存储空间，不直接存储输入输出数据，只存储执行状态和性能指标
type TestResult struct {
	TestCaseID string      `bson:"test_case_id" json:"test_case_id"`   // 测试用例ID，关联TestCase.ID
	Status     string      `bson:"status" json:"status"`               // 运行状态（AC/WA/TLE/MLE/RE等）
	TimeUsed   int         `bson:"time_used" json:"time_used"`         // 该测试用例运行时间（毫秒）
	MemoryUsed int         `bson:"memory_used" json:"memory_used"`     // 该测试用例使用内存（KB）
	Score      int         `bson:"score" json:"score"`                 // 该测试用例得分
	Details    JudgeDetail `bson:"judge_details" json:"judge_details"` // go-judge判题器详细信息
}

// JudgeDetail go-judge判题器详细信息 - 存储底层判题系统的原始数据
// 用于问题排查和性能分析，普通用户不可见
type JudgeDetail struct {
	GoJudgeStatus string `bson:"go_judge_status" json:"go_judge_status"` // go-judge返回的原始状态码
	ExitStatus    int    `bson:"exit_status" json:"exit_status"`         // 程序退出状态码（0表示正常退出）
	RuntimeNS     int64  `bson:"runtime_ns" json:"runtime_ns"`           // 程序运行时间（纳秒级精度）
}