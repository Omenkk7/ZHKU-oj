package interfaces

import (
	"context"
	"zhku-oj/internal/model"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CreateUserRequest 创建用户请求 (类似Spring的DTO)
type CreateUserRequest struct {
	StudentID string `json:"student_id" binding:"required"`
	Username  string `json:"username" binding:"required,min=3,max=20"`
	Password  string `json:"password" binding:"required,min=6"`
	Email     string `json:"email" binding:"required,email"`
	RealName  string `json:"real_name" binding:"required"`
	Role      string `json:"role" binding:"required,oneof=student teacher admin"`
	Class     string `json:"class"`
	Grade     string `json:"grade"`
}

// UpdateUserRequest 更新用户请求 (类似Spring的DTO)
type UpdateUserRequest struct {
	Username string `json:"username" binding:"omitempty,min=3,max=20"`
	Email    string `json:"email" binding:"omitempty,email"`
	RealName string `json:"real_name" binding:"omitempty"`
	Role     string `json:"role" binding:"omitempty,oneof=student teacher admin"`
	Class    string `json:"class"`
	Grade    string `json:"grade"`
	Avatar   string `json:"avatar"`
	IsActive *bool  `json:"is_active"`
}

// ChangePasswordRequest 修改密码请求
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

// UserListRequest 用户列表查询请求 (类似Spring的Specification)
type UserListRequest struct {
	Page     int    `form:"page,default=1" binding:"min=1"`
	PageSize int    `form:"page_size,default=20" binding:"min=1,max=100"`
	Role     string `form:"role"`
	Class    string `form:"class"`
	Grade    string `form:"grade"`
	IsActive *bool  `form:"is_active"`
	Keyword  string `form:"keyword"`
}

// UserListResponse 用户列表响应 (类似Spring的Page<UserDTO>)
type UserListResponse struct {
	Users      []*model.User `json:"users"`
	Total      int64         `json:"total"`
	Page       int           `json:"page"`
	PageSize   int           `json:"page_size"`
	TotalPages int           `json:"total_pages"`
}

// UserService 用户业务服务接口 (类似Spring的@Service接口)
type UserService interface {
	// CreateUser 创建用户 (类似Spring的@Transactional方法)
	CreateUser(ctx context.Context, req *CreateUserRequest) (*model.User, error)

	// GetUserByID 根据ID获取用户
	GetUserByID(ctx context.Context, id primitive.ObjectID) (*model.User, error)

	// GetUserByUsername 根据用户名获取用户
	GetUserByUsername(ctx context.Context, username string) (*model.User, error)

	// UpdateUser 更新用户信息
	UpdateUser(ctx context.Context, id primitive.ObjectID, req *UpdateUserRequest) (*model.User, error)

	// ChangePassword 修改密码
	ChangePassword(ctx context.Context, userID primitive.ObjectID, req *ChangePasswordRequest) error

	// DeleteUser 删除用户
	DeleteUser(ctx context.Context, id primitive.ObjectID) error

	// ListUsers 分页查询用户列表 (类似Spring的Page<User> findAll())
	ListUsers(ctx context.Context, req *UserListRequest) (*UserListResponse, error)

	// ActivateUser 激活用户
	ActivateUser(ctx context.Context, id primitive.ObjectID) error

	// DeactivateUser 停用用户
	DeactivateUser(ctx context.Context, id primitive.ObjectID) error

	// GetUserStats 获取用户统计信息
	GetUserStats(ctx context.Context, userID primitive.ObjectID) (*model.UserStats, error)

	// ValidateUser 验证用户存在性和权限
	ValidateUser(ctx context.Context, userID primitive.ObjectID, requiredRole string) (*model.User, error)
}
