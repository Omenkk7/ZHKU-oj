/**
 * @Author: omenkk7
 * @Date: 2025/9/14
 * @Desc: 用户相关响应体定义 - API v1规范
 */
package response

import (
	"zhku-oj/pkg/io/constanct"
	"time"
	"zhku-oj/internal/domain"
)

// ============================= 认证相关 - API v1 =============================

// AuthLoginData 登录响应数据结构
type AuthLoginData struct {
	AccessToken  string `json:"access_token"`  // JWT访问令牌
	TokenType    string `json:"token_type"`    // 令牌类型
	ExpiresIn    int64  `json:"expires_in"`    // 过期时间(秒)
	Username     string `json:"username"`      // 用户名
	UserID       string `json:"user_id"`       // 用户ID
	RefreshToken string `json:"refresh_token"` // 刷新令牌
}

/*
API v1规范 - 用户登录
HTTP请求结构：
POST /api/v1/auth/login
Content-Type: application/json
Body: {
  "username": "student001",
  "password": "password123"
}
响应格式：
{
  "code": 0,
  "message": "登录成功",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "token_type": "Bearer",
    "expires_in": 7200,
    "username": "student001",
    "user_id": "uid_12345",
    "refresh_token": "refresh_token_abc"
  }
}
*/
// NewV1LoginResponse API v1 登录成功响应
func NewV1LoginResponse(token, refreshToken, uid, username string, expiresIn int64) Response[AuthLoginData] {
	return NewResponse(0, "登录成功", AuthLoginData{
		AccessToken:  token,
		TokenType:    "Bearer",
		ExpiresIn:    expiresIn,a
		Username:     username,
		UserID:       uid,
		RefreshToken: refreshToken,
	})
}

// AuthRegisterData 注册响应数据结构
type AuthRegisterData struct {
	UserID             string `json:"user_id"`             // 用户ID
	Username           string `json:"username"`            // 用户名
	Email              string `json:"email"`               // 邮箱
	Status             string `json:"status"`              // 账户状态
	CreatedAt          int64  `json:"created_at"`          // 创建时间
	ActivationRequired bool   `json:"activation_required"` // 是否需要激活
}

/*
API v1规范 - 用户注册
HTTP请求结构：
POST /api/v1/auth/register
Content-Type: application/json
Body: {
  "username": "student001",
  "password": "password123",
  "email": "student001@zhku.edu.cn",
  "school": "仲恺农业工程学院",
  "classes": "计科2001",
  "major": "计算机科学与技术"
}
响应格式：
{
  "code": 0,
  "message": "注册成功",
  "data": {
    "user_id": "uid_12345",
    "username": "student001",
    "email": "student001@zhku.edu.cn",
    "status": "active",
    "created_at": 1726286400,
    "activation_required": false
  }
}
*/
// NewV1RegisterResponse API v1 注册成功响应
func NewV1RegisterResponse(uid, username, email string, createdAt int64) Response[AuthRegisterData] {
	return NewResponse(0, "注册成功", AuthRegisterData{
		UserID:             uid,
		Username:           username,
		Email:              email,
		Status:             "active",
		CreatedAt:          createdAt,
		ActivationRequired: false,
	})
}

// ============================= 用户资源管理 - API v1 =============================

// UserProfileData 用户详细信息响应数据结构
type UserProfileData struct {
	UserID           string            `json:"user_id"`           // 用户ID
	Username         string            `json:"username"`          // 用户名
	Email            string            `json:"email"`             // 邮箱
	Profile          UserProfileInfo   `json:"profile"`           // 个人信息
	Statistics       UserStatistics    `json:"statistics"`        // 统计信息
	ExternalAccounts []ExternalAccount `json:"external_accounts"` // 第三方账号
	CreatedAt        int64             `json:"created_at"`        // 注册时间
	UpdatedAt        int64             `json:"updated_at"`        // 更新时间
	LastLoginAt      int64             `json:"last_login_at"`     // 最后登录时间
}

// UserProfileInfo 用户个人信息
type UserProfileInfo struct {
	School      string `json:"school"`       // 学校
	Class       string `json:"class"`        // 班级
	Major       string `json:"major"`        // 专业
	Specialty   string `json:"specialty"`    // 擅长方向
	AvatarURL   string `json:"avatar_url"`   // 头像URL
	DisplayName string `json:"display_name"` // 显示名称
}

// UserStatistics 用户统计信息
type UserStatistics struct {
	Rating           int     `json:"rating"`            // 评分
	TotalSubmissions int64   `json:"total_submissions"` // 总提交数
	SolvedProblems   uint32  `json:"solved_problems"`   // 已解决题目数
	AcceptanceRate   float64 `json:"acceptance_rate"`   // 通过率
	GlobalRanking    int     `json:"global_ranking"`    // 全局排名
}

// ExternalAccount 第三方账号信息
type ExternalAccount struct {
	Platform string `json:"platform"`  // 平台名称
	Username string `json:"username"`  // 用户名
	IsActive bool   `json:"is_active"` // 是否激活
	BindTime int64  `json:"bind_time"` // 绑定时间
}

/*
API v1规范 - 获取用户详细信息
HTTP请求结构：
GET /api/v1/users/{user_id}
Headers: {
  "Authorization": "Bearer <jwt_token>"
}
响应格式：
{
  "code": 0,
  "message": "获取用户信息成功",
  "data": {
    "user_id": "uid_12345",
    "username": "student001",
    "email": "student001@zhku.edu.cn",
    "profile": {...},
    "statistics": {...},
    "external_accounts": [...],
    "created_at": 1726286400,
    "updated_at": 1726286400,
    "last_login_at": 1726286400
  }
}
*/
// NewV1UserProfileResponse API v1 获取用户详细信息响应
func NewV1UserProfileResponse(user *domain.User) Response[UserProfileData] {
	// 计算通过率
	acceptanceRate := 0.0
	if user.Submited > 0 {
		acceptanceRate = float64(user.Solved) / float64(user.Submited) * 100
	}

	// 构建第三方账号列表
	externalAccounts := []ExternalAccount{}
	if user.Vjid != "" {
		externalAccounts = append(externalAccounts, ExternalAccount{
			Platform: "vjudge",
			Username: user.Vjid,
			IsActive: true,
			BindTime: user.RegisterTime,
		})
	}
	if user.CodeForceUser != "" {
		externalAccounts = append(externalAccounts, ExternalAccount{
			Platform: "codeforces",
			Username: user.CodeForceUser,
			IsActive: true,
			BindTime: user.RegisterTime,
		})
	}

	return NewResponse(0, "获取用户信息成功", UserProfileData{
		UserID:   user.UID,
		Username: user.Uname,
		Email:    user.Email,
		Profile: UserProfileInfo{
			School:      user.School,
			Class:       user.Classes,
			Major:       user.Major,
			Specialty:   user.Adept,
			AvatarURL:   user.HeadURL,
			DisplayName: user.Uname,
		},
		Statistics: UserStatistics{
			Rating:           user.Rating,
			TotalSubmissions: user.Submited,
			SolvedProblems:   user.Solved,
			AcceptanceRate:   acceptanceRate,
			GlobalRanking:    0, // 需要从其他服务获取
		},
		ExternalAccounts: externalAccounts,
		CreatedAt:        user.RegisterTime,
		UpdatedAt:        user.RegisterTime,
		LastLoginAt:      user.RegisterTime,
	})
}

// UserUpdateData 更新用户信息响应数据结构
type UserUpdateData struct {
	UserID        string   `json:"user_id"`        // 用户ID
	UpdatedFields []string `json:"updated_fields"` // 更新的字段列表
	UpdatedAt     int64    `json:"updated_at"`     // 更新时间
	Version       int      `json:"version"`        // 版本号
}

/*
API v1规范 - 更新用户信息
HTTP请求结构：
PUT /api/v1/users/{user_id}
Headers: {
  "Authorization": "Bearer <jwt_token>"
}
Content-Type: application/json
Body: {
  "profile": {
    "school": "仲恺农业工程学院",
    "class": "计科2002",
    "major": "软件工程",
    "specialty": "算法竞赛",
    "avatar_url": "https://example.com/avatar.jpg",
    "display_name": "新昵称"
  },
  "email": "newemail@zhku.edu.cn"
}
响应格式：
{
  "code": 0,
  "message": "用户信息更新成功",
  "data": {
    "user_id": "uid_12345",
    "updated_fields": ["profile.school", "profile.class", "email"],
    "updated_at": 1726286400,
    "version": 2
  }
}
*/
// NewV1UserUpdateResponse API v1 更新用户信息成功响应
func NewV1UserUpdateResponse(uid string, updatedFields []string, version int) Response[UserUpdateData] {
	return NewResponse(0, "用户信息更新成功", UserUpdateData{
		UserID:        uid,
		UpdatedFields: updatedFields,
		UpdatedAt:     time.Now().Unix(),
		Version:       version,
	})
}

// PasswordChangeData 修改密码响应数据结构
type PasswordChangeData struct {
	UserID        string `json:"user_id"`        // 用户ID
	ChangedAt     int64  `json:"changed_at"`     // 修改时间
	RequireReauth bool   `json:"require_reauth"` // 是否需要重新登录
	Message       string `json:"message"`        // 提示信息
}

/*
API v1规范 - 修改密码
HTTP请求结构：
PUT /api/v1/users/{user_id}/password
Headers: {
  "Authorization": "Bearer <jwt_token>"
}
Content-Type: application/json
Body: {
  "current_password": "current123",
  "new_password": "newpassword456",
  "confirm_password": "newpassword456"
}
响应格式：
{
  "code": 0,
  "message": "密码修改成功",
  "data": {
    "user_id": "uid_12345",
    "changed_at": 1726286400,
    "require_reauth": true,
    "message": "密码已成功修改，请重新登录"
  }
}
*/
// NewV1PasswordChangeResponse API v1 修改密码成功响应
func NewV1PasswordChangeResponse(uid string) Response[PasswordChangeData] {
	return NewResponse(0, "密码修改成功", PasswordChangeData{
		UserID:        uid,
		ChangedAt:     time.Now().Unix(),
		RequireReauth: true,
		Message:       "密码已成功修改，请重新登录",
	})
}

// ============================= 用户列表管理 - API v1 =============================

// UserListItem 用户列表项数据结构
type UserListItem struct {
	UserID           string  `json:"user_id"`           // 用户ID
	Username         string  `json:"username"`          // 用户名
	DisplayName      string  `json:"display_name"`      // 显示名称
	Email            string  `json:"email"`             // 邮箱
	School           string  `json:"school"`            // 学校
	Class            string  `json:"class"`             // 班级
	Major            string  `json:"major"`             // 专业
	AvatarURL        string  `json:"avatar_url"`        // 头像URL
	Rating           int     `json:"rating"`            // 评分
	TotalSubmissions int64   `json:"total_submissions"` // 总提交数
	SolvedProblems   uint32  `json:"solved_problems"`   // 已解决题目数
	AcceptanceRate   float64 `json:"acceptance_rate"`   // 通过率
	GlobalRanking    int     `json:"global_ranking"`    // 全局排名
	JoinedAt         int64   `json:"joined_at"`         // 加入时间
	IsActive         bool    `json:"is_active"`         // 是否活跃
}

// UserListResponse 用户列表响应数据结构
type UserListResponse struct {
	Users      []UserListItem `json:"users"`      // 用户列表
	Pagination PaginationInfo `json:"pagination"` // 分页信息
	Filters    FilterInfo     `json:"filters"`    // 过滤信息
	Sorting    SortingInfo    `json:"sorting"`    // 排序信息
}

// PaginationInfo 分页信息
type PaginationInfo struct {
	Page       int   `json:"page"`        // 当前页
	PageSize   int   `json:"page_size"`   // 每页大小
	Total      int64 `json:"total"`       // 总数
	TotalPages int   `json:"total_pages"` // 总页数
	HasNext    bool  `json:"has_next"`    // 是否有下一页
	HasPrev    bool  `json:"has_prev"`    // 是否有上一页
}

// FilterInfo 过滤信息
type FilterInfo struct {
	Class    string `json:"class,omitempty"`     // 班级过滤
	School   string `json:"school,omitempty"`    // 学校过滤
	Major    string `json:"major,omitempty"`     // 专业过滤
	Keyword  string `json:"keyword,omitempty"`   // 关键词搜索
	IsActive *bool  `json:"is_active,omitempty"` // 活跃状态过滤
}

// SortingInfo 排序信息
type SortingInfo struct {
	SortBy string `json:"sort_by"` // 排序字段
	Order  string `json:"order"`   // 排序方向
}

/*
API v1规范 - 获取用户列表
HTTP请求结构：
GET /api/v1/users?page=1&page_size=20&class=计科2001&keyword=student&sort_by=rating&order=desc
Headers: {
  "Authorization": "Bearer <jwt_token>"
}
Query参数：
- page: 页码(默认1)
- page_size: 每页大小(默认20, 最大1, 最大100)
- class: 班级过滤(可选)
- school: 学校过滤(可选)
- major: 专业过滤(可选)
- keyword: 关键词搜索(可选)
- is_active: 活跃状态过滤(可选)
- sort_by: 排序字段(rating/solved_problems/total_submissions/joined_at，默认rating)
- order: 排序方向(asc/desc，默认desc)
响应格式：
{
  "code": 0,
  "message": "获取用户列表成功",
  "data": {
    "users": [...],
    "pagination": {...},
    "filters": {...},
    "sorting": {...}
  }
}
*/
// NewV1UserListResponse API v1 获取用户列表响应
func NewV1UserListResponse(users []domain.User, total int64, page, pageSize int, filters FilterInfo, sorting SortingInfo) Response[UserListResponse] {
	userItems := make([]UserListItem, 0, len(users))
	for _, user := range users {
		// 计算通过率
		acceptanceRate := 0.0
		if user.Submited > 0 {
			acceptanceRate = float64(user.Solved) / float64(user.Submited) * 100
		}

		userItems = append(userItems, UserListItem{
			UserID:           user.UID,
			Username:         user.Uname,
			DisplayName:      user.Uname,
			Email:            user.Email,
			School:           user.School,
			Class:            user.Classes,
			Major:            user.Major,
			AvatarURL:        user.HeadURL,
			Rating:           user.Rating,
			TotalSubmissions: user.Submited,
			SolvedProblems:   user.Solved,
			AcceptanceRate:   acceptanceRate,
			GlobalRanking:    0, // 需要从其他服务获取
			JoinedAt:         user.RegisterTime,
			IsActive:         user.Defaulted != "disabled",
		})
	}

	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))

	return NewResponse(0, "获取用户列表成功", UserListResponse{
		Users: userItems,
		Pagination: PaginationInfo{
			Page:       page,
			PageSize:   pageSize,
			Total:      total,
			TotalPages: totalPages,
			HasNext:    page < totalPages,
			HasPrev:    page > 1,
		},
		Filters: filters,
		Sorting: sorting,
	})
}

// ============================= 用户管理操作 - API v1 =============================

// UserDeleteData 删除用户响应数据结构
type UserDeleteData struct {
	UserID       string `json:"user_id"`        // 被删除的用户ID
	DeletedAt    int64  `json:"deleted_at"`     // 删除时间
	DeletedBy    string `json:"deleted_by"`     // 执行删除的管理员ID
	Reason       string `json:"reason"`         // 删除原因
	IsHardDelete bool   `json:"is_hard_delete"` // 是否硬删除
}

/*
API v1规范 - 删除用户
HTTP请求结构：
DELETE /api/v1/users/{user_id}
Headers: {
  "Authorization": "Bearer <admin_jwt_token>"
}
Query参数：
- hard_delete: 是否硬删除(true/false，默认false)
- reason: 删除原因(可选)
注意：只有管理员才能删除用户，软删除只是禁用账户，硬删除会永久删除数据
响应格式：
{
  "code": 0,
  "message": "用户删除成功",
  "data": {
    "user_id": "uid_12345",
    "deleted_at": 1726286400,
    "deleted_by": "admin_uid_001",
    "reason": "违反社区规则",
    "is_hard_delete": false
  }
}
*/
// NewV1UserDeleteResponse API v1 删除用户成功响应
func NewV1UserDeleteResponse(uid, deletedBy, reason string, isHardDelete bool) Response[UserDeleteData] {
	return NewResponse(0, "用户删除成功", UserDeleteData{
		UserID:       uid,
		DeletedAt:    time.Now().Unix(),
		DeletedBy:    deletedBy,
		Reason:       reason,
		IsHardDelete: isHardDelete,
	})
}

// ============================= 用户统计分析 - API v1 =============================

// UserStatsDetailData 用户统计信息响应数据结构
type UserStatsDetailData struct {
	UserID           string           `json:"user_id"`           // 用户ID
	Overview         StatsOverview    `json:"overview"`          // 概览统计
	SubmissionStats  SubmissionStats  `json:"submission_stats"`  // 提交统计
	ProblemStats     ProblemStats     `json:"problem_stats"`     // 题目统计
	RankingInfo      RankingInfo      `json:"ranking_info"`      // 排名信息
	ActivityTimeline []ActivityRecord `json:"activity_timeline"` // 活动时间线
	GeneratedAt      int64            `json:"generated_at"`      // 统计生成时间
}

// StatsOverview 概览统计
type StatsOverview struct {
	Rating           int     `json:"rating"`            // 当前评分
	TotalSubmissions int64   `json:"total_submissions"` // 总提交数
	SolvedProblems   uint32  `json:"solved_problems"`   // 已解决题目数
	AcceptanceRate   float64 `json:"acceptance_rate"`   // 通过率
	GlobalRanking    int     `json:"global_ranking"`    // 全局排名
	DaysActive       int     `json:"days_active"`       // 活跃天数
	ConsecutiveDays  int     `json:"consecutive_days"`  // 连续活跃天数
	LastActivityAt   int64   `json:"last_activity_at"`  // 最后活动时间
}

// SubmissionStats 提交统计
type SubmissionStats struct {
	Accepted            int64   `json:"accepted"`              // AC数量
	WrongAnswer         int64   `json:"wrong_answer"`          // WA数量
	TimeLimitExceeded   int64   `json:"time_limit_exceeded"`   // TLE数量
	MemoryLimitExceeded int64   `json:"memory_limit_exceeded"` // MLE数量
	RuntimeError        int64   `json:"runtime_error"`         // RE数量
	CompileError        int64   `json:"compile_error"`         // CE数量
	Other               int64   `json:"other"`                 // 其他状态
	TotalTime           int64   `json:"total_time"`            // 总用时(毫秒)
	AverageTime         float64 `json:"average_time"`          // 平均用时(毫秒)
	LastSubmitTime      int64   `json:"last_submit_time"`      // 最后提交时间
}

// ProblemStats 题目统计
type ProblemStats struct {
	Easy           int      `json:"easy"`            // 简单题数量
	Medium         int      `json:"medium"`          // 中等题数量
	Hard           int      `json:"hard"`            // 困难题数量
	FavoriteTopics []string `json:"favorite_topics"` // 屡长题型
	WeakTopics     []string `json:"weak_topics"`     // 薄弱题型
}

// RankingInfo 排名信息
type RankingInfo struct {
	GlobalRank int     `json:"global_rank"` // 全局排名
	SchoolRank int     `json:"school_rank"` // 校内排名
	ClassRank  int     `json:"class_rank"`  // 班级排名
	RankChange int     `json:"rank_change"` // 排名变化(正数上升，负数下降)
	Percentile float64 `json:"percentile"`  // 百分位
}

// ActivityRecord 活动记录
type ActivityRecord struct {
	Date        string `json:"date"`        // 日期(YYYY-MM-DD)
	Submissions int    `json:"submissions"` // 当日提交数
	Solved      int    `json:"solved"`      // 当日解决题目数
	TimeSpent   int    `json:"time_spent"`  // 当日用时(分钟)
}

/*
API v1规范 - 获取用户统计信息
HTTP请求结构：
GET /api/v1/users/{user_id}/statistics
Headers: {
  "Authorization": "Bearer <jwt_token>"
}
Query参数：
- period: 统计周期(7d/30d/90d/1y/all，默认all)
- include_timeline: 是否包含活动时间线(true/false，默认true)
响应格式：
{
  "code": 0,
  "message": "获取用户统计成功",
  "data": {
    "user_id": "uid_12345",
    "overview": {...},
    "submission_stats": {...},
    "problem_stats": {...},
    "ranking_info": {...},
    "activity_timeline": [...],
    "generated_at": 1726286400
  }
}
*/
// NewV1UserStatsResponse API v1 获取用户统计信息响应
func NewV1UserStatsResponse(user *domain.User, submissionStats SubmissionStats, ranking RankingInfo, activityTimeline []ActivityRecord) Response[UserStatsDetailData] {
	acceptanceRate := 0.0
	if user.Submited > 0 {
		acceptanceRate = float64(user.Solved) / float64(user.Submited) * 100
	}

	return NewResponse(0, "获取用户统计成功", UserStatsDetailData{
		UserID: user.UID,
		Overview: StatsOverview{
			Rating:           user.Rating,
			TotalSubmissions: user.Submited,
			SolvedProblems:   user.Solved,
			AcceptanceRate:   acceptanceRate,
			GlobalRanking:    ranking.GlobalRank,
			DaysActive:       0, // 需要从活动记录计算
			ConsecutiveDays:  0, // 需要从活动记录计算
			LastActivityAt:   submissionStats.LastSubmitTime,
		},
		SubmissionStats: submissionStats,
		ProblemStats: ProblemStats{
			Easy:           0, // 需要从题目服务获取
			Medium:         0,
			Hard:           0,
			FavoriteTopics: []string{},
			WeakTopics:     []string{},
		},
		RankingInfo:      ranking,
		ActivityTimeline: activityTimeline,
		GeneratedAt:      time.Now().Unix(),
	})
}

// ============================= 第三方账号管理 - API v1 =============================

// ExternalAccountBindData 绑定第三方账号响应数据结构
type ExternalAccountBindData struct {
	UserID           string `json:"user_id"`           // 用户ID
	Platform         string `json:"platform"`          // 平台名称
	ExternalUsername string `json:"external_username"` // 第三方平台用户名
	BindTime         int64  `json:"bind_time"`         // 绑定时间
	Status           string `json:"status"`            // 绑定状态
	IsVerified       bool   `json:"is_verified"`       // 是否已验证
	LastSyncTime     int64  `json:"last_sync_time"`    // 最后同步时间
}

// ExternalAccountUnbindData 解绑第三方账号响应数据结构
type ExternalAccountUnbindData struct {
	UserID       string `json:"user_id"`       // 用户ID
	Platform     string `json:"platform"`      // 平台名称
	UnbindTime   int64  `json:"unbind_time"`   // 解绑时间
	Reason       string `json:"reason"`        // 解绑原因
	DataRetained bool   `json:"data_retained"` // 是否保留数据
}

// ExternalAccountListData 第三方账号列表响应数据结构
type ExternalAccountListData struct {
	UserID    string            `json:"user_id"`    // 用户ID
	Accounts  []ExternalAccount `json:"accounts"`   // 账号列表
	Total     int               `json:"total"`      // 总数
	UpdatedAt int64             `json:"updated_at"` // 更新时间
}

/*
API v1规范 - 绑定第三方账号
HTTP请求结构：
POST /api/v1/users/{user_id}/external-accounts
Headers: {
  "Authorization": "Bearer <jwt_token>"
}
Content-Type: application/json
Body: {
  "platform": "codeforces", // 支持: "codeforces", "vjudge", "atcoder", "leetcode"
  "username": "cf_username",
  "password": "cf_password", // 仅VJ需要，其他平台可选
  "auto_sync": true,  // 是否自动同步
  "verify_method": "api" // 验证方式: "api", "manual"
}
响应格式：
{
  "code": 0,
  "message": "第三方账号绑定成功",
  "data": {
    "user_id": "uid_12345",
    "platform": "codeforces",
    "external_username": "cf_username",
    "bind_time": 1726286400,
    "status": "active",
    "is_verified": true,
    "last_sync_time": 1726286400
  }
}
*/
// NewV1ExternalAccountBindResponse API v1 绑定第三方账号成功响应
func NewV1ExternalAccountBindResponse(uid, platform, externalUsername string, isVerified bool) Response[ExternalAccountBindData] {
	bindTime := time.Now().Unix()
	return NewResponse(0, "第三方账号绑定成功", ExternalAccountBindData{
		UserID:           uid,
		Platform:         platform,
		ExternalUsername: externalUsername,
		BindTime:         bindTime,
		Status:           "active",
		IsVerified:       isVerified,
		LastSyncTime:     bindTime,
	})
}

/*
API v1规范 - 解绑第三方账号
HTTP请求结构：
DELETE /api/v1/users/{user_id}/external-accounts/{platform}
Headers: {
  "Authorization": "Bearer <jwt_token>"
}
Query参数：
- retain_data: 是否保留历史数据(true/false，默认true)
- reason: 解绑原因(可选)
*/
// NewV1ExternalAccountUnbindResponse API v1 解绑第三方账号成功响应
func NewV1ExternalAccountUnbindResponse(uid, platform, reason string, dataRetained bool) Response[ExternalAccountUnbindData] {
	return NewResponse(0, "第三方账号解绑成功", ExternalAccountUnbindData{
		UserID:       uid,
		Platform:     platform,
		UnbindTime:   time.Now().Unix(),
		Reason:       reason,
		DataRetained: dataRetained,
	})
}

/*
API v1规范 - 获取第三方账号列表
HTTP请求结构：
GET /api/v1/users/{user_id}/external-accounts
Headers: {
  "Authorization": "Bearer <jwt_token>"
}
Query参数：
- platform: 过滤平台(可选)
- status: 过滤状态(active/inactive/error，可选)
*/
// NewV1ExternalAccountListResponse API v1 获取第三方账号列表响应
func NewV1ExternalAccountListResponse(uid string, accounts []ExternalAccount) Response[ExternalAccountListData] {
	return NewResponse(0, "获取第三方账号列表成功", ExternalAccountListData{
		UserID:    uid,
		Accounts:  accounts,
		Total:     len(accounts),
		UpdatedAt: time.Now().Unix(),
	})
}

// ============================= 兼容性函数 - 保持与旧版本API兼容 =============================

// 为了保持向后兼容，保留部分旧版本函数别名

// LoginData 旧版本登录数据结构(已弃用，请使用AuthLoginData)
type LoginData = AuthLoginData

// NewLoginResponse 旧版本登录响应(已弃用，请使用NewV1LoginResponse)
func NewLoginResponse(token, uname string) Response[LoginData] {
	return NewV1LoginResponse(token, "", "", uname, 7200)
}
