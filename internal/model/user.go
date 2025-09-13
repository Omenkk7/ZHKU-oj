package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User 用户模型
type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	StudentID string             `bson:"student_id" json:"student_id"`
	Username  string             `bson:"username" json:"username"`
	Password  string             `bson:"password" json:"-"`
	Email     string             `bson:"email" json:"email"`
	RealName  string             `bson:"real_name" json:"real_name"`
	Role      string             `bson:"role" json:"role"` // student, teacher, admin
	Class     string             `bson:"class" json:"class"`
	Grade     string             `bson:"grade" json:"grade"`
	Avatar    string             `bson:"avatar" json:"avatar"`
	IsActive  bool               `bson:"is_active" json:"is_active"`
	Stats     UserStats          `bson:"stats" json:"stats"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
	LastLogin *time.Time         `bson:"last_login,omitempty" json:"last_login,omitempty"`
}

// UserStats 用户统计信息
type UserStats struct {
	TotalSubmissions int `bson:"total_submissions" json:"total_submissions"`
	AcceptedCount    int `bson:"accepted_count" json:"accepted_count"`
	ProblemsSolved   int `bson:"problems_solved" json:"problems_solved"`
	Ranking          int `bson:"ranking" json:"ranking"`
}

// Problem 题目模型
type Problem struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title        string             `bson:"title" json:"title"`
	Description  string             `bson:"description" json:"description"`
	InputFormat  string             `bson:"input_format" json:"input_format"`
	OutputFormat string             `bson:"output_format" json:"output_format"`
	SampleInput  string             `bson:"sample_input" json:"sample_input"`
	SampleOutput string             `bson:"sample_output" json:"sample_output"`
	TimeLimit    int                `bson:"time_limit" json:"time_limit"`     // 毫秒
	MemoryLimit  int                `bson:"memory_limit" json:"memory_limit"` // MB
	Difficulty   string             `bson:"difficulty" json:"difficulty"`     // easy, medium, hard
	Tags         []string           `bson:"tags" json:"tags"`
	TestCases    []TestCase         `bson:"test_cases" json:"test_cases"`
	Stats        ProblemStats       `bson:"stats" json:"stats"`
	IsPublic     bool               `bson:"is_public" json:"is_public"`
	CreatedBy    primitive.ObjectID `bson:"created_by" json:"created_by"`
	CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt    time.Time          `bson:"updated_at" json:"updated_at"`
}

// TestCase 测试用例
type TestCase struct {
	ID       string `bson:"id" json:"id"`
	Input    string `bson:"input" json:"input"`
	Output   string `bson:"output" json:"output"`
	Score    int    `bson:"score" json:"score"`
	IsPublic bool   `bson:"is_public" json:"is_public"`
}

// ProblemStats 题目统计信息
type ProblemStats struct {
	TotalSubmissions int     `bson:"total_submissions" json:"total_submissions"`
	AcceptedCount    int     `bson:"accepted_count" json:"accepted_count"`
	AcceptanceRate   float64 `bson:"acceptance_rate" json:"acceptance_rate"`
	AverageTime      int     `bson:"average_time" json:"average_time"`
	AverageMemory    int     `bson:"average_memory" json:"average_memory"`
}

// Submission 提交记录模型
type Submission struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID      primitive.ObjectID `bson:"user_id" json:"user_id"`
	ProblemID   primitive.ObjectID `bson:"problem_id" json:"problem_id"`
	Code        string             `bson:"code" json:"code"`
	Language    string             `bson:"language" json:"language"`
	Status      string             `bson:"status" json:"status"`
	Score       int                `bson:"score" json:"score"`
	TimeUsed    int                `bson:"time_used" json:"time_used"`     // 毫秒
	MemoryUsed  int                `bson:"memory_used" json:"memory_used"` // KB
	CompileInfo CompileInfo        `bson:"compile_info" json:"compile_info"`
	TestResults []TestResult       `bson:"test_results" json:"test_results"`
	SubmittedAt time.Time          `bson:"submitted_at" json:"submitted_at"`
	JudgedAt    *time.Time         `bson:"judged_at,omitempty" json:"judged_at,omitempty"`
}

// CompileInfo 编译信息
type CompileInfo struct {
	Status     string `bson:"status" json:"status"`
	TimeUsed   int64  `bson:"time_used" json:"time_used"`
	MemoryUsed int64  `bson:"memory_used" json:"memory_used"`
	Message    string `bson:"message" json:"message"`
}

// TestResult 测试结果
type TestResult struct {
	TestCaseID     string      `bson:"test_case_id" json:"test_case_id"`
	Status         string      `bson:"status" json:"status"`
	TimeUsed       int         `bson:"time_used" json:"time_used"`
	MemoryUsed     int         `bson:"memory_used" json:"memory_used"`
	Score          int         `bson:"score" json:"score"`
	Input          string      `bson:"input,omitempty" json:"input,omitempty"`
	ExpectedOutput string      `bson:"expected_output,omitempty" json:"expected_output,omitempty"`
	ActualOutput   string      `bson:"actual_output,omitempty" json:"actual_output,omitempty"`
	JudgeDetails   JudgeDetail `bson:"judge_details" json:"judge_details"`
}

// JudgeDetail go-judge详细信息
type JudgeDetail struct {
	GoJudgeStatus string `bson:"go_judge_status" json:"go_judge_status"`
	ExitStatus    int    `bson:"exit_status" json:"exit_status"`
	RuntimeNS     int64  `bson:"runtime_ns" json:"runtime_ns"`
}

// Contest 竞赛模型 (预留)
type Contest struct {
	ID          primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	Title       string               `bson:"title" json:"title"`
	Description string               `bson:"description" json:"description"`
	ProblemIDs  []primitive.ObjectID `bson:"problem_ids" json:"problem_ids"`
	StartTime   time.Time            `bson:"start_time" json:"start_time"`
	EndTime     time.Time            `bson:"end_time" json:"end_time"`
	CreatedBy   primitive.ObjectID   `bson:"created_by" json:"created_by"`
	CreatedAt   time.Time            `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time            `bson:"updated_at" json:"updated_at"`
}

// 提交状态常量
const (
	StatusPending             = "PENDING"
	StatusJudging             = "JUDGING"
	StatusAccepted            = "ACCEPTED"
	StatusWrongAnswer         = "WRONG_ANSWER"
	StatusTimeLimitExceeded   = "TIME_LIMIT_EXCEEDED"
	StatusMemoryLimitExceeded = "MEMORY_LIMIT_EXCEEDED"
	StatusRuntimeError        = "RUNTIME_ERROR"
	StatusCompileError        = "COMPILE_ERROR"
	StatusSystemError         = "SYSTEM_ERROR"
	StatusDangerousSyscall    = "DANGEROUS_SYSCALL"
	StatusOutputLimitExceeded = "OUTPUT_LIMIT_EXCEEDED"
)

// 用户角色常量
const (
	RoleStudent = "student"
	RoleTeacher = "teacher"
	RoleAdmin   = "admin"
)

// 题目难度常量
const (
	DifficultyEasy   = "easy"
	DifficultyMedium = "medium"
	DifficultyHard   = "hard"
)

// 编程语言常量
const (
	LanguageJava = "java"
)
