package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Problem 题目模型 - 存储编程题目的完整信息
// 对应MongoDB集合: problems
type Problem struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`            // 题目唯一标识ID
	Title        string             `bson:"title" json:"title"`                 // 题目标题
	Description  string             `bson:"description" json:"description"`     // 题目描述（支持Markdown格式）
	InputFormat  string             `bson:"input_format" json:"input_format"`   // 输入格式说明
	OutputFormat string             `bson:"output_format" json:"output_format"` // 输出格式说明
	SampleInput  string             `bson:"sample_input" json:"sample_input"`   // 示例输入
	SampleOutput string             `bson:"sample_output" json:"sample_output"` // 示例输出
	TimeLimit    int                `bson:"time_limit" json:"time_limit"`       // 时间限制（毫秒）
	MemoryLimit  int                `bson:"memory_limit" json:"memory_limit"`   // 内存限制（MB）
	Difficulty   string             `bson:"difficulty" json:"difficulty"`       // 难度级别: easy(简单), medium(中等), hard(困难)
	Tags         []string           `bson:"tags" json:"tags"`                   // 题目标签（如：算法、数据结构、动态规划等）
	TestCases    []TestCase         `bson:"test_cases" json:"test_cases"`       // 测试用例集合
	Stats        ProblemStats       `bson:"stats" json:"stats"`                 // 题目统计信息
	IsPublic     bool               `bson:"is_public" json:"is_public"`         // 是否公开可见
	CreatedBy    primitive.ObjectID `bson:"created_by" json:"created_by"`       // 创建者用户ID
	CreatedAt    time.Time          `bson:"created_at" json:"created_at"`       // 创建时间
	UpdatedAt    time.Time          `bson:"updated_at" json:"updated_at"`       // 最后更新时间
}

// TestCase 测试用例模型 - 存储测试用例的元数据信息
// 注意：不直接存储大体量输入输出数据，实际文件存放在磁盘或对象存储中
type TestCase struct {
	ID       string `bson:"id" json:"id"`               // 测试用例唯一标识
	FilePath string `bson:"file_path" json:"file_path"` // 测试数据文件存放路径（磁盘或OSS）
	Score    int    `bson:"score" json:"score"`         // 该测试用例的分值
	IsPublic bool   `bson:"is_public" json:"is_public"` // 是否公开给学生查看
}

// ProblemStats 题目统计信息 - 记录题目的提交和通过数据
// 由后台异步任务定时更新，用于题目难度分析和推荐
type ProblemStats struct {
	TotalSubmissions int     `bson:"total_submissions" json:"total_submissions"` // 总提交次数
	AcceptedCount    int     `bson:"accepted_count" json:"accepted_count"`       // 通过次数
	AcceptanceRate   float64 `bson:"acceptance_rate" json:"acceptance_rate"`     // 通过率（百分比）
}
