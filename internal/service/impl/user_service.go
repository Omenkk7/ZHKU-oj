package impl

import (
	"context"
	"fmt"
	"math"
	"zhku-oj/internal/model"
	repoInterface "zhku-oj/internal/repository/interfaces"
	serviceInterface "zhku-oj/internal/service/interfaces"
	errors2 "zhku-oj/pkg/errors"

	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// userService 用户服务实现 (类似Spring的@Service实现类)
type userService struct {
	userRepo    repoInterface.UserRepository
	redisClient *redis.Client
}

// NewUserService 创建用户服务实例 (类似Spring的@Autowired构造函数)
func NewUserService(userRepo repoInterface.UserRepository, redisClient *redis.Client) serviceInterface.UserService {
	return &userService{
		userRepo:    userRepo,
		redisClient: redisClient,
	}
}

// CreateUser 创建用户 (类似Spring的@Transactional方法)
func (s *userService) CreateUser(ctx context.Context, req *serviceInterface.CreateUserRequest) (*model.User, error) {
	// 1. 验证唯一性约束 (类似Spring的@Valid + 自定义验证)
	exists, err := s.userRepo.ExistsByUsername(ctx, req.Username)
	if err != nil {
		return nil, fmt.Errorf("检查用户名失败: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("用户名已存在")
	}

	exists, err = s.userRepo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("检查邮箱失败: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("邮箱已存在")
	}

	exists, err = s.userRepo.ExistsByStudentID(ctx, req.StudentID)
	if err != nil {
		return nil, fmt.Errorf("检查学号失败: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("学号已存在")
	}

	// 2. 密码加密 (类似Spring Security的PasswordEncoder)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors2.Wrap(errors2.SYSTEM_ERROR, err)
	}

	// 3. 创建用户模型
	user := &model.User{
		StudentID: req.StudentID,
		Username:  req.Username,
		Password:  string(hashedPassword),
		Email:     req.Email,
		RealName:  req.RealName,
		Role:      req.Role,
		Class:     req.Class,
		Grade:     req.Grade,
		IsActive:  true,
		Stats: model.UserStats{
			TotalSubmissions: 0,
			AcceptedCount:    0,
			ProblemsSolved:   0,
			Ranking:          0,
		},
	}

	// 4. 保存到数据库 (类似Spring的@Transactional)
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, errors2.Wrap(errors2.DATABASE_ERROR, err)
	}

	// 5. 清除密码字段 (安全考虑)
	user.Password = ""

	return user, nil
}

// GetUserByID 根据ID获取用户
func (s *userService) GetUserByID(ctx context.Context, id primitive.ObjectID) (*model.User, error) {
	// 先尝试从Redis缓存获取 (类似Spring Cache)
	cacheKey := fmt.Sprintf("user:%s", id.Hex())

	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// 清除密码字段
	user.Password = ""

	// 缓存到Redis (TTL: 1小时)
	// 这里省略了Redis序列化代码，实际项目中需要实现

	return user, nil
}

// GetUserByUsername 根据用户名获取用户
func (s *userService) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	user, err := s.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	// 清除密码字段
	user.Password = ""
	return user, nil
}

// UpdateUser 更新用户信息
func (s *userService) UpdateUser(ctx context.Context, id primitive.ObjectID, req *serviceInterface.UpdateUserRequest) (*model.User, error) {
	// 1. 获取现有用户
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// 2. 检查唯一性约束
	if req.Username != "" && req.Username != user.Username {
		exists, err := s.userRepo.ExistsByUsername(ctx, req.Username)
		if err != nil {
			return nil, fmt.Errorf("检查用户名失败: %w", err)
		}
		if exists {
			return nil, fmt.Errorf("用户名已存在")
		}
		user.Username = req.Username
	}

	if req.Email != "" && req.Email != user.Email {
		exists, err := s.userRepo.ExistsByEmail(ctx, req.Email)
		if err != nil {
			return nil, fmt.Errorf("检查邮箱失败: %w", err)
		}
		if exists {
			return nil, fmt.Errorf("邮箱已存在")
		}
		user.Email = req.Email
	}

	// 3. 更新其他字段
	if req.RealName != "" {
		user.RealName = req.RealName
	}
	if req.Role != "" {
		user.Role = req.Role
	}
	if req.Class != "" {
		user.Class = req.Class
	}
	if req.Grade != "" {
		user.Grade = req.Grade
	}
	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}
	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}

	// 4. 保存更新
	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("更新用户失败: %w", err)
	}

	// 5. 清除缓存
	cacheKey := fmt.Sprintf("user:%s", id.Hex())
	s.redisClient.Del(ctx, cacheKey)

	// 清除密码字段
	user.Password = ""
	return user, nil
}

// ChangePassword 修改密码
func (s *userService) ChangePassword(ctx context.Context, userID primitive.ObjectID, req *serviceInterface.ChangePasswordRequest) error {
	// 1. 获取用户（包含密码）
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	// 2. 验证旧密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)); err != nil {
		return fmt.Errorf("旧密码不正确")
	}

	// 3. 加密新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("密码加密失败: %w", err)
	}

	// 4. 更新密码
	if err := s.userRepo.UpdatePassword(ctx, userID, string(hashedPassword)); err != nil {
		return fmt.Errorf("更新密码失败: %w", err)
	}

	return nil
}

// DeleteUser 删除用户 (类似Spring的软删除或硬删除)
func (s *userService) DeleteUser(ctx context.Context, id primitive.ObjectID) error {
	// 检查用户是否存在
	_, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// 执行删除 (这里是硬删除，实际项目中可能需要软删除)
	if err := s.userRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("删除用户失败: %w", err)
	}

	// 清除缓存
	cacheKey := fmt.Sprintf("user:%s", id.Hex())
	s.redisClient.Del(ctx, cacheKey)

	return nil
}

// ListUsers 分页查询用户列表 (类似Spring的Page<User> findAll())
func (s *userService) ListUsers(ctx context.Context, req *serviceInterface.UserListRequest) (*serviceInterface.UserListResponse, error) {
	// 构建过滤条件
	filters := make(map[string]interface{})
	if req.Role != "" {
		filters["role"] = req.Role
	}
	if req.Class != "" {
		filters["class"] = req.Class
	}
	if req.Grade != "" {
		filters["grade"] = req.Grade
	}
	if req.IsActive != nil {
		filters["is_active"] = *req.IsActive
	}
	if req.Keyword != "" {
		filters["keyword"] = req.Keyword
	}

	// 查询数据
	users, total, err := s.userRepo.List(ctx, req.Page, req.PageSize, filters)
	if err != nil {
		return nil, fmt.Errorf("查询用户列表失败: %w", err)
	}

	// 清除密码字段
	for _, user := range users {
		user.Password = ""
	}

	// 计算总页数
	totalPages := int(math.Ceil(float64(total) / float64(req.PageSize)))

	return &serviceInterface.UserListResponse{
		Users:      users,
		Total:      total,
		Page:       req.Page,
		PageSize:   req.PageSize,
		TotalPages: totalPages,
	}, nil
}

// ActivateUser 激活用户
func (s *userService) ActivateUser(ctx context.Context, id primitive.ObjectID) error {
	req := &serviceInterface.UpdateUserRequest{
		IsActive: &[]bool{true}[0],
	}
	_, err := s.UpdateUser(ctx, id, req)
	return err
}

// DeactivateUser 停用用户
func (s *userService) DeactivateUser(ctx context.Context, id primitive.ObjectID) error {
	req := &serviceInterface.UpdateUserRequest{
		IsActive: &[]bool{false}[0],
	}
	_, err := s.UpdateUser(ctx, id, req)
	return err
}

// GetUserStats 获取用户统计信息
func (s *userService) GetUserStats(ctx context.Context, userID primitive.ObjectID) (*model.UserStats, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return &user.Stats, nil
}

// ValidateUser 验证用户存在性和权限
func (s *userService) ValidateUser(ctx context.Context, userID primitive.ObjectID, requiredRole string) (*model.User, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if !user.IsActive {
		return nil, fmt.Errorf("用户已被停用")
	}

	if requiredRole != "" {
		switch requiredRole {
		case model.RoleAdmin:
			if user.Role != model.RoleAdmin {
				return nil, fmt.Errorf("需要管理员权限")
			}
		case model.RoleTeacher:
			if user.Role != model.RoleTeacher && user.Role != model.RoleAdmin {
				return nil, fmt.Errorf("需要教师或管理员权限")
			}
		}
	}

	// 清除密码字段
	user.Password = ""
	return user, nil
}
