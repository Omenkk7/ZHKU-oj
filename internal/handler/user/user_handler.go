package user

import (
	"net/http"
	"strconv"
	"zhku-oj/internal/middleware"
	"zhku-oj/internal/pkg/errors"
	"zhku-oj/internal/pkg/utils"
	"zhku-oj/internal/service/interfaces"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UserHandler 用户控制器 (类似Spring的@RestController)
type UserHandler struct {
	userService interfaces.UserService
}

// NewUserHandler 创建用户控制器实例 (类似Spring的@Autowired构造函数)
func NewUserHandler(userService interfaces.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// CreateUser 创建用户 (类似Spring的@PostMapping)
// POST /api/v1/admin/users
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req interfaces.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendError(c, errors.INVALID_PARAMS)
		return
	}

	user, err := h.userService.CreateUser(c.Request.Context(), &req)
	if err != nil {
		// 如果是业务错误，直接使用统一错误处理
		utils.HandleError(c, err)
		return
	}

	utils.SendSuccess(c, user)
}

// GetUser 获取用户详情 (类似Spring的@GetMapping("/{id}"))
// GET /api/v1/users/{id}
func (h *UserHandler) GetUser(c *gin.Context) {
	idParam := c.Param("id")
	userID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "用户ID格式错误", err.Error())
		return
	}

	user, err := h.userService.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "用户不存在", err.Error())
		return
	}

	utils.SuccessResponse(c, "获取用户成功", user)
}

// UpdateUser 更新用户信息 (类似Spring的@PutMapping("/{id}"))
// PUT /api/v1/users/{id}
func (h *UserHandler) UpdateUser(c *gin.Context) {
	idParam := c.Param("id")
	userID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "用户ID格式错误", err.Error())
		return
	}

	var req interfaces.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "参数错误", err.Error())
		return
	}

	// 权限检查：只能修改自己的信息，或管理员可以修改任何人
	currentUserID := middleware.GetUserID(c)
	currentUserRole := middleware.GetUserRole(c)

	if currentUserID != userID.Hex() && currentUserRole != "admin" {
		utils.ErrorResponse(c, http.StatusForbidden, "无权限修改他人信息", "")
		return
	}

	user, err := h.userService.UpdateUser(c.Request.Context(), userID, &req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "更新用户失败", err.Error())
		return
	}

	utils.SuccessResponse(c, "更新用户成功", user)
}

// DeleteUser 删除用户 (类似Spring的@DeleteMapping("/{id}"))
// DELETE /api/v1/admin/users/{id}
func (h *UserHandler) DeleteUser(c *gin.Context) {
	idParam := c.Param("id")
	userID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "用户ID格式错误", err.Error())
		return
	}

	if err := h.userService.DeleteUser(c.Request.Context(), userID); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "删除用户失败", err.Error())
		return
	}

	utils.SuccessResponse(c, "删除用户成功", nil)
}

// ListUsers 获取用户列表 (类似Spring的@GetMapping with Pageable)
// GET /api/v1/admin/users?page=1&page_size=20&role=student&keyword=张三
func (h *UserHandler) ListUsers(c *gin.Context) {
	var req interfaces.UserListRequest

	// 绑定查询参数 (类似Spring的@RequestParam)
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.SendError(c, errors.INVALID_PARAMS)
		return
	}

	// 设置默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}

	// 处理is_active参数
	if isActiveStr := c.Query("is_active"); isActiveStr != "" {
		if isActive, err := strconv.ParseBool(isActiveStr); err == nil {
			req.IsActive = &isActive
		}
	}

	response, err := h.userService.ListUsers(c.Request.Context(), &req)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	// 使用分页响应
	utils.SendSuccessWithPagination(c, response.Users, response.Page, response.PageSize, response.Total)
}

// GetProfile 获取当前用户信息 (类似Spring Security的@AuthenticationPrincipal)
// GET /api/v1/users/profile
func (h *UserHandler) GetProfile(c *gin.Context) {
	userIDStr := middleware.GetUserID(c)
	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "用户ID格式错误", err.Error())
		return
	}

	user, err := h.userService.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "用户不存在", err.Error())
		return
	}

	utils.SuccessResponse(c, "获取用户信息成功", user)
}

// UpdateProfile 更新当前用户信息
// PUT /api/v1/users/profile
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userIDStr := middleware.GetUserID(c)
	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "用户ID格式错误", err.Error())
		return
	}

	var req interfaces.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "参数错误", err.Error())
		return
	}

	// 普通用户不能修改角色
	req.Role = ""

	user, err := h.userService.UpdateUser(c.Request.Context(), userID, &req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "更新用户信息失败", err.Error())
		return
	}

	utils.SuccessResponse(c, "更新用户信息成功", user)
}

// ChangePassword 修改密码
// PUT /api/v1/users/password
func (h *UserHandler) ChangePassword(c *gin.Context) {
	userIDStr := middleware.GetUserID(c)
	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "用户ID格式错误", err.Error())
		return
	}

	var req interfaces.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "参数错误", err.Error())
		return
	}

	if err := h.userService.ChangePassword(c.Request.Context(), userID, &req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "修改密码失败", err.Error())
		return
	}

	utils.SuccessResponse(c, "修改密码成功", nil)
}

// ActivateUser 激活用户
// PUT /api/v1/admin/users/{id}/activate
func (h *UserHandler) ActivateUser(c *gin.Context) {
	idParam := c.Param("id")
	userID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "用户ID格式错误", err.Error())
		return
	}

	if err := h.userService.ActivateUser(c.Request.Context(), userID); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "激活用户失败", err.Error())
		return
	}

	utils.SuccessResponse(c, "激活用户成功", nil)
}

// DeactivateUser 停用用户
// PUT /api/v1/admin/users/{id}/deactivate
func (h *UserHandler) DeactivateUser(c *gin.Context) {
	idParam := c.Param("id")
	userID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "用户ID格式错误", err.Error())
		return
	}

	if err := h.userService.DeactivateUser(c.Request.Context(), userID); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "停用用户失败", err.Error())
		return
	}

	utils.SuccessResponse(c, "停用用户成功", nil)
}

// GetUserStats 获取用户统计信息
// GET /api/v1/users/{id}/stats
func (h *UserHandler) GetUserStats(c *gin.Context) {
	idParam := c.Param("id")
	userID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "用户ID格式错误", err.Error())
		return
	}

	stats, err := h.userService.GetUserStats(c.Request.Context(), userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "获取用户统计失败", err.Error())
		return
	}

	utils.SuccessResponse(c, "获取用户统计成功", stats)
}
