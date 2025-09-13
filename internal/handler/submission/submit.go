package submission

import (
	"net/http"
	"strconv"

	"zhku-oj/internal/middleware"
	"zhku-oj/internal/pkg/utils"
	"zhku-oj/internal/service/interfaces"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Handler 代码提交处理器
type Handler struct {
	service interfaces.SubmissionService
}

// NewSubmissionHandler 创建代码提交处理器
func NewSubmissionHandler(service interfaces.SubmissionService) *Handler {
	return &Handler{
		service: service,
	}
}

// SubmitRequest 代码提交请求
type SubmitRequest struct {
	ProblemID string `json:"problem_id" binding:"required"`
	Code      string `json:"code" binding:"required,max=50000"`
	Language  string `json:"language" binding:"required,oneof=java"`
}

// Submit 提交代码接口
// 用户提交解题代码，系统将进行在线判题
// 请求方法: POST
// 路径: /api/v1/submissions
// 请求体: {"problem_id": "题目ID", "code": "源代码", "language": "编程语言"}
// 响应: {"submission_id": "提交ID", "status": "PENDING"}
func (h *Handler) Submit(c *gin.Context) {
	var req SubmitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	// 获取用户ID
	userIDStr := middleware.GetUserID(c)
	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "用户ID格式错误")
		return
	}

	// 验证题目ID
	problemID, err := primitive.ObjectIDFromHex(req.ProblemID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "题目ID格式错误")
		return
	}

	// 调用服务层处理提交
	submission, err := h.service.Submit(c.Request.Context(), userID, problemID, req.Code, req.Language)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "提交失败: "+err.Error())
		return
	}

	// 返回成功响应
	utils.SuccessResponse(c, gin.H{
		"submission_id": submission.ID.Hex(),
		"status":        submission.Status,
		"submitted_at":  submission.SubmittedAt,
	})
}

// GetSubmission 获取提交详情接口
// 获取指定提交的详细信息和判题结果
// 请求方法: GET
// 路径: /api/v1/submissions/{id}
// 响应: 提交详情包括状态、得分、测试结果等
func (h *Handler) GetSubmission(c *gin.Context) {
	submissionIDStr := c.Param("id")
	submissionID, err := primitive.ObjectIDFromHex(submissionIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "提交ID格式错误")
		return
	}

	// 获取当前用户ID
	userIDStr := middleware.GetUserID(c)
	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "用户ID格式错误")
		return
	}

	// 获取提交详情
	submission, err := h.service.GetSubmission(c.Request.Context(), submissionID, userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "提交不存在或无权限访问")
		return
	}

	utils.SuccessResponse(c, submission)
}

// ListSubmissions 获取提交列表接口
// 获取用户的提交记录列表，支持分页和筛选
// 请求方法: GET
// 路径: /api/v1/submissions?page=1&page_size=20&problem_id=xxx&status=ACCEPTED
// 响应: 分页的提交列表
func (h *Handler) ListSubmissions(c *gin.Context) {
	// 解析查询参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	problemIDStr := c.Query("problem_id")
	status := c.Query("status")
	language := c.Query("language")

	// 限制页面大小
	if pageSize > 100 {
		pageSize = 100
	}
	if pageSize < 1 {
		pageSize = 20
	}
	if page < 1 {
		page = 1
	}

	// 获取当前用户ID
	userIDStr := middleware.GetUserID(c)
	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "用户ID格式错误")
		return
	}

	// 构造查询条件
	filter := map[string]interface{}{
		"user_id": userID,
	}

	if problemIDStr != "" {
		problemID, err := primitive.ObjectIDFromHex(problemIDStr)
		if err != nil {
			utils.ErrorResponse(c, http.StatusBadRequest, "题目ID格式错误")
			return
		}
		filter["problem_id"] = problemID
	}

	if status != "" {
		filter["status"] = status
	}

	if language != "" {
		filter["language"] = language
	}

	// 获取提交列表
	submissions, total, err := h.service.ListSubmissions(c.Request.Context(), filter, page, pageSize)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "获取提交列表失败: "+err.Error())
		return
	}

	// 计算总页数
	totalPages := (int(total) + pageSize - 1) / pageSize

	utils.SuccessResponseWithPagination(c, submissions, utils.Pagination{
		Page:       page,
		PageSize:   pageSize,
		Total:      total,
		TotalPages: totalPages,
	})
}
