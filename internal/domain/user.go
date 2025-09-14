package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

/*
@Author: omenkk7
@Date: 2025/9/14 10:44
@Description: 用户持久层数据 x
*/

type User struct {
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"id"`                          // 用户唯一标识ID
	StudentID         string             `bson:"student_id" json:"student_id"`                     // 学号，用于学校系统集成
	Username          string             `bson:"username" json:"username"`                         // 用户名，用于登录，全局唯一
	Password          string             `bson:"password" json:"-"`                                // 加密后，不返回给前端
	Email             string             `bson:"email" json:"email"`                               // 邮箱地址，用于找回密码等功能
	RealName          string             `bson:"real_name" json:"real_name"`                       // 真实姓名
	Role              string             `bson:"role" json:"role"`                                 // 用户角色: student(学生), teacher(教师), admin(管理员)
	Class             string             `bson:"class" json:"class"`                               // 所属班级
	Grade             string             `bson:"grade" json:"grade"`                               // 所属年级
	Avatar            string             `bson:"avatar" json:"avatar"`                             // 头像URL
	IsActive          bool               `bson:"is_active" json:"is_active"`                       // 账户是否激活
	Stats             UserStats          `bson:"stats" json:"stats"`                               // 用户统计信息
	Preferences       UserPreferences    `bson:"preferences" json:"preferences"`                   // 用户偏好设置
	CreatedAt         time.Time          `bson:"created_at" json:"created_at"`                     // 账户创建时间
	UpdatedAt         time.Time          `bson:"updated_at" json:"updated_at"`                     // 账户最后更新时间
	LastLogin         *time.Time         `bson:"last_login,omitempty" json:"last_login,omitempty"` // 最后登录时间
	PasswordUpdatedAt *time.Time         `bson:"password_updated_at,omitempty" json:"-"`           // 密码最后修改时间
	LoginCount        int                `bson:"login_count" json:"login_count"`                   // 登录次数统计
}

// UserStats 用户统计信息 - 记录用户在系统中的各项表现数据
// 由后台定时任务或异步队列更新，提升查询性能
type UserStats struct {
	TotalSubmissions int `bson:"total_submissions" json:"total_submissions"` // 总提交次数
	AcceptedCount    int `bson:"accepted_count" json:"accepted_count"`       // AC（通过）次数
	ProblemsSolved   int `bson:"problems_solved" json:"problems_solved"`     // 解决的题目数量（去重）
	TotalScore       int `bson:"total_score" json:"total_score"`             // 总得分
	MaxStreak        int `bson:"max_streak" json:"max_streak"`               // 最大连续AC天数
	CurrentStreak    int `bson:"current_streak" json:"current_streak"`       // 当前连续AC天数
}

// UserPreferences 用户偏好设置 - 存储用户的个性化配置
type UserPreferences struct {
	Language      string `bson:"language" json:"language"`           // 默认编程语言偏好
	Theme         string `bson:"theme" json:"theme"`                 // 界面主题: light(浅色)/dark(深色)
	Notifications bool   `bson:"notifications" json:"notifications"` // 是否开启系统通知
}
