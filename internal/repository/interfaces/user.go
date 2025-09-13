package interfaces

import (
	"context"
	"zhku-oj/internal/model"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UserRepository 用户数据访问接口 (类似Spring的@Repository)
type UserRepository interface {
	// Create 创建用户 (类似Spring的save方法)
	Create(ctx context.Context, user *model.User) error

	// GetByID 根据ID获取用户 (类似Spring的findById)
	GetByID(ctx context.Context, id primitive.ObjectID) (*model.User, error)

	// GetByUsername 根据用户名获取用户 (类似Spring的findByUsername)
	GetByUsername(ctx context.Context, username string) (*model.User, error)

	// GetByStudentID 根据学号获取用户
	GetByStudentID(ctx context.Context, studentID string) (*model.User, error)

	// GetByEmail 根据邮箱获取用户 (类似Spring的findByEmail)
	GetByEmail(ctx context.Context, email string) (*model.User, error)

	// Update 更新用户 (类似Spring的save方法)
	Update(ctx context.Context, user *model.User) error

	// UpdatePassword 更新密码
	UpdatePassword(ctx context.Context, userID primitive.ObjectID, hashedPassword string) error

	// Delete 删除用户 (类似Spring的deleteById)
	Delete(ctx context.Context, id primitive.ObjectID) error

	// List 分页查询用户列表 (类似Spring的findAll with Pageable)
	List(ctx context.Context, page, pageSize int, filters map[string]interface{}) ([]*model.User, int64, error)

	// UpdateStats 更新用户统计信息
	UpdateStats(ctx context.Context, userID primitive.ObjectID, stats model.UserStats) error

	// UpdateLastLogin 更新最后登录时间
	UpdateLastLogin(ctx context.Context, userID primitive.ObjectID) error

	// ExistsByUsername 检查用户名是否存在 (类似Spring的existsByUsername)
	ExistsByUsername(ctx context.Context, username string) (bool, error)

	// ExistsByEmail 检查邮箱是否存在
	ExistsByEmail(ctx context.Context, email string) (bool, error)

	// ExistsByStudentID 检查学号是否存在
	ExistsByStudentID(ctx context.Context, studentID string) (bool, error)
}
