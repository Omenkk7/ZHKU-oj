/**
 * @Author: omenkk7
 * @Date: 2025/9/14
 * @Desc: 用户相关响应体定义 - API v1规范
 */
package response

import (
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
		ExpiresIn:    expiresIn,
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
	UserID       string `json:"user_id"`        // 用户ID
	ChangedAt    int64  `json:"changed_at"`     // 修改时间
	RequireReauth bool   `json:"require_reauth"` // 是否需要重新登录
	Message      string `json:"message"`        // 提示信息
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

// ============================= 用户列表查询 =============================

// UserListItem 用户列表项数据结构
type UserListItem struct {
	UID      string `json:"uid"`      // 用户ID
	Uname    string `json:"uname"`    // 用户名
	School   string `json:"school"`   // 学校
	Classes  string `json:"classes"`  // 班级
	Major    string `json:"major"`    // 专业
	Rating   int    `json:"rating"`   // 评分
	Submited int64  `json:"submited"` // 提交数
	Solved   uint32 `json:"solved"`   // 已解决题目数
}

// UserListData 用户列表响应数据结构
type UserListData struct {
	Users      []UserListItem `json:"users"`       // 用户列表
	Total      int64          `json:"total"`       // 总数
	Page       int            `json:"page"`        // 当前页
	PageSize   int            `json:"page_size"`   // 每页大小
	TotalPages int            `json:"total_pages"` // 总页数
}

/*
HTTP请求结构：
GET /api/users?page=1&page_size=20&class=计科2001&keyword=student
Headers: {
  "Authorization": "Bearer <jwt_token>"
}
Query参数：
- page: 页码(默认1)
- page_size: 每页大小(默认20)
- class: 班级过滤(可选)
- keyword: 关键词搜索(可选)
- sort_by: 排序字段(rating/solved/submited，默认rating)
- order: 排序方向(asc/desc，默认desc)
*/
// NewUserListResponse 获取用户列表响应
func NewUserListResponse(users []domain.User, total int64, page, pageSize int) Response[UserListData] {
	userItems := make([]UserListItem, 0, len(users))
	for _, user := range users {
		userItems = append(userItems, UserListItem{
			UID:      user.UID,
			Uname:    user.Uname,
			School:   user.School,
			Classes:  user.Classes,
			Major:    user.Major,
			Rating:   user.Rating,
			Submited: user.Submited,
			Solved:   user.Solved,
		})
	}

	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))

	return NewResponse(0, "获取用户列表成功", UserListData{
		Users:      userItems,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	})
}

// ============================= 删除用户 =============================

// DeleteUserData 删除用户响应数据结构
type DeleteUserData struct {
	UID     string `json:"uid"`     // 被删除的用户ID
	Message string `json:"message"` // 删除消息
}

/*
HTTP请求结构：
DELETE /api/users/{uid}
Headers: {
  "Authorization": "Bearer <jwt_token>"
}
注意：只有管理员才能删除用户
*/
// NewDeleteUserResponse 删除用户成功响应
func NewDeleteUserResponse(uid string) Response[DeleteUserData] {
	return NewResponse(0, "用户删除成功", DeleteUserData{
		UID:     uid,
		Message: "用户账户已被成功删除",
	})
}

// ============================= 用户统计信息 =============================

// UserStatsData 用户统计信息响应数据结构
type UserStatsData struct {
	UID              string  `json:"uid"`               // 用户ID
	TotalSubmissions int64   `json:"total_submissions"` // 总提交数
	SolvedProblems   uint32  `json:"solved_problems"`   // 已解决题目数
	AcceptanceRate   float64 `json:"acceptance_rate"`   // 通过率
	Rating           int     `json:"rating"`            // 当前评分
	Ranking          int     `json:"ranking"`           // 排名
	LastSubmitTime   int64   `json:"last_submit_time"`  // 最后提交时间
}

/*
HTTP请求结构：
GET /api/users/{uid}/stats
Headers: {
  "Authorization": "Bearer <jwt_token>"
}
*/
// NewUserStatsResponse 获取用户统计信息响应
func NewUserStatsResponse(user *domain.User, ranking int, lastSubmitTime int64) Response[UserStatsData] {
	acceptanceRate := 0.0
	if user.Submited > 0 {
		acceptanceRate = float64(user.Solved) / float64(user.Submited) * 100
	}

	return NewResponse(0, "获取用户统计成功", UserStatsData{
		UID:              user.UID,
		TotalSubmissions: user.Submited,
		SolvedProblems:   user.Solved,
		AcceptanceRate:   acceptanceRate,
		Rating:           user.Rating,
		Ranking:          ranking,
		LastSubmitTime:   lastSubmitTime,
	})
}

// ============================= 绑定第三方账号 =============================

// BindAccountData 绑定第三方账号响应数据结构
type BindAccountData struct {
	UID         string `json:"uid"`          // 用户ID
	AccountType string `json:"account_type"` // 账号类型
	Message     string `json:"message"`      // 绑定消息
}

/*
HTTP请求结构：
POST /api/users/{uid}/bind
Headers: {
  "Authorization": "Bearer <jwt_token>"
}
Content-Type: application/json
Body: {
  "type": "codeforces", // 或 "vjudge"
  "username": "cf_username",
  "password": "cf_password" // 仅VJ需要
}
*/
// NewBindAccountResponse 绑定第三方账号成功响应
func NewBindAccountResponse(uid, accountType string) Response[BindAccountData] {
	return NewResponse(0, "账号绑定成功", BindAccountData{
		UID:         uid,
		AccountType: accountType,
		Message:     "第三方账号绑定成功",
	})
}
