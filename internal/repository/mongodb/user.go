package mongodb

import (
	"context"
	"fmt"
	"time"
	"zhku-oj/internal/model"
	"zhku-oj/internal/repository/interfaces"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// 用户仓储层
type userRepository struct {
	collection *mongo.Collection
}

// NewUserRepository
func NewUserRepository(client *mongo.Client, database string) interfaces.UserRepository {
	return &userRepository{
		collection: client.Database(database).Collection("users"),
	}
}

// Create 创建用户 (类似JpaRepository的save方法)
func (r *userRepository) Create(ctx context.Context, user *model.User) error {
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	result, err := r.collection.InsertOne(ctx, user)
	if err != nil {
		return fmt.Errorf("创建用户失败: %w", err)
	}

	user.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// GetByID 根据ID获取用户 (类似JpaRepository的findById)
func (r *userRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*model.User, error) {
	var user model.User
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("用户不存在")
		}
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}
	return &user, nil
}

// GetByUsername 根据用户名获取用户 (类似Spring的findByUsername)
func (r *userRepository) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	err := r.collection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("用户不存在")
		}
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}
	return &user, nil
}

// GetByStudentID 根据学号获取用户
func (r *userRepository) GetByStudentID(ctx context.Context, studentID string) (*model.User, error) {
	var user model.User
	err := r.collection.FindOne(ctx, bson.M{"student_id": studentID}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("用户不存在")
		}
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}
	return &user, nil
}

// GetByEmail 根据邮箱获取用户 (类似Spring的findByEmail)
func (r *userRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("用户不存在")
		}
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}
	return &user, nil
}

// Update 更新用户 (类似JpaRepository的save方法)
func (r *userRepository) Update(ctx context.Context, user *model.User) error {
	user.UpdatedAt = time.Now()

	update := bson.M{
		"$set": bson.M{
			"username":   user.Username,
			"email":      user.Email,
			"real_name":  user.RealName,
			"role":       user.Role,
			"class":      user.Class,
			"grade":      user.Grade,
			"avatar":     user.Avatar,
			"is_active":  user.IsActive,
			"updated_at": user.UpdatedAt,
		},
	}

	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": user.ID}, update)
	if err != nil {
		return fmt.Errorf("更新用户失败: %w", err)
	}
	return nil
}

// UpdatePassword 更新密码
func (r *userRepository) UpdatePassword(ctx context.Context, userID primitive.ObjectID, hashedPassword string) error {
	update := bson.M{
		"$set": bson.M{
			"password":   hashedPassword,
			"updated_at": time.Now(),
		},
	}

	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": userID}, update)
	if err != nil {
		return fmt.Errorf("更新密码失败: %w", err)
	}
	return nil
}

// Delete 删除用户 (类似JpaRepository的deleteById)
func (r *userRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return fmt.Errorf("删除用户失败: %w", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("用户不存在")
	}
	return nil
}

// List 分页查询用户列表 (类似Spring的Page<User> findAll(Pageable))
func (r *userRepository) List(ctx context.Context, page, pageSize int, filters map[string]interface{}) ([]*model.User, int64, error) {
	// 构建查询条件 (类似Spring Data JPA的Specification)
	filter := bson.M{}
	for key, value := range filters {
		switch key {
		case "role":
			filter["role"] = value
		case "class":
			filter["class"] = value
		case "grade":
			filter["grade"] = value
		case "is_active":
			filter["is_active"] = value
		case "keyword": // 关键词搜索
			if keyword, ok := value.(string); ok && keyword != "" {
				filter["$or"] = []bson.M{
					{"username": bson.M{"$regex": keyword, "$options": "i"}},
					{"real_name": bson.M{"$regex": keyword, "$options": "i"}},
					{"student_id": bson.M{"$regex": keyword, "$options": "i"}},
				}
			}
		}
	}

	// 计算跳过的文档数 (类似Spring的PageRequest.of(page, size))
	skip := (page - 1) * pageSize

	// 查询选项
	opts := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(pageSize)).
		SetSort(bson.D{{"created_at", -1}}) // 按创建时间倒序

	// 执行查询
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, fmt.Errorf("查询用户列表失败: %w", err)
	}
	defer cursor.Close(ctx)

	// 解析结果
	var users []*model.User
	if err = cursor.All(ctx, &users); err != nil {
		return nil, 0, fmt.Errorf("解析用户数据失败: %w", err)
	}

	// 计算总数 (类似Spring的Page.getTotalElements())
	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("统计用户总数失败: %w", err)
	}

	return users, total, nil
}

// UpdateStats 更新用户统计信息
func (r *userRepository) UpdateStats(ctx context.Context, userID primitive.ObjectID, stats model.UserStats) error {
	update := bson.M{
		"$set": bson.M{
			"stats":      stats,
			"updated_at": time.Now(),
		},
	}

	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": userID}, update)
	if err != nil {
		return fmt.Errorf("更新用户统计失败: %w", err)
	}
	return nil
}

// UpdateLastLogin 更新最后登录时间
func (r *userRepository) UpdateLastLogin(ctx context.Context, userID primitive.ObjectID) error {
	now := time.Now()
	update := bson.M{
		"$set": bson.M{
			"last_login": &now,
			"updated_at": now,
		},
	}

	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": userID}, update)
	if err != nil {
		return fmt.Errorf("更新最后登录时间失败: %w", err)
	}
	return nil
}

// ExistsByUsername 检查用户名是否存在 (类似Spring的existsByUsername)
func (r *userRepository) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	count, err := r.collection.CountDocuments(ctx, bson.M{"username": username})
	if err != nil {
		return false, fmt.Errorf("检查用户名失败: %w", err)
	}
	return count > 0, nil
}

// ExistsByEmail 检查邮箱是否存在
func (r *userRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	count, err := r.collection.CountDocuments(ctx, bson.M{"email": email})
	if err != nil {
		return false, fmt.Errorf("检查邮箱失败: %w", err)
	}
	return count > 0, nil
}

// ExistsByStudentID 检查学号是否存在
func (r *userRepository) ExistsByStudentID(ctx context.Context, studentID string) (bool, error) {
	count, err := r.collection.CountDocuments(ctx, bson.M{"student_id": studentID})
	if err != nil {
		return false, fmt.Errorf("检查学号失败: %w", err)
	}
	return count > 0, nil
}
