package router

import (
	"zhku-oj/internal/middleware"

	"github.com/gin-gonic/gin"
)

// setupProblemRoutes 设置题目相关路由
// 题目浏览、搜索、管理等功能
func (rm *RouterManager) setupProblemRoutes(v1 *gin.RouterGroup) {
	problemGroup := v1.Group("/problems")
	problemGroup.Use(middleware.AuthRequired()) // 所有题目接口都需要认证
	{
		// ========== 题目查询接口 ==========

		// 获取题目列表（支持分页、搜索、筛选）
		// GET /api/v1/problems?page=1&page_size=20&keyword=排序&difficulty=easy&tags=算法
		// 响应码: 0-成功, 10002-参数错误
		problemGroup.GET("", rm.problemHandler.ListProblems)

		// 获取题目详情
		// GET /api/v1/problems/{id}
		// 响应码: 0-成功, 10002-参数错误, 30001-题目不存在, 30006-题目访问被拒绝
		problemGroup.GET("/:id", rm.problemHandler.GetProblem)

		// 获取题目统计信息
		// GET /api/v1/problems/{id}/stats
		// 响应码: 0-成功, 10002-参数错误, 30001-题目不存在
		// problemGroup.GET("/:id/stats", rm.problemHandler.GetProblemStats)

		// 获取题目提交记录（分页）
		// GET /api/v1/problems/{id}/submissions?page=1&page_size=20&status=ACCEPTED
		// 响应码: 0-成功, 10002-参数错误, 30001-题目不存在
		// problemGroup.GET("/:id/submissions", rm.submissionHandler.GetProblemSubmissions)

		// ========== 题目管理接口（教师/管理员权限） ==========

		// 创建题目
		// POST /api/v1/problems
		// 权限: teacher, admin
		// 响应码: 0-成功, 10002-参数错误, 10004-权限不足, 30002-题目已存在
		problemGroup.POST("",
			middleware.RoleRequired("teacher", "admin"),
			rm.problemHandler.CreateProblem)

		// 更新题目
		// PUT /api/v1/problems/{id}
		// 权限: teacher, admin
		// 响应码: 0-成功, 10002-参数错误, 10004-权限不足, 30001-题目不存在
		problemGroup.PUT("/:id",
			middleware.RoleRequired("teacher", "admin"),
			rm.problemHandler.UpdateProblem)

		// 删除题目
		// DELETE /api/v1/problems/{id}
		// 权限: admin
		// 响应码: 0-成功, 10002-参数错误, 10004-权限不足, 30001-题目不存在
		problemGroup.DELETE("/:id",
			middleware.RoleRequired("admin"),
			rm.problemHandler.DeleteProblem)

		// 批量导入题目
		// POST /api/v1/problems/import
		// 权限: teacher, admin
		// 响应码: 0-成功, 10002-参数错误, 10004-权限不足
		// problemGroup.POST("/import",
		// 	middleware.RoleRequired("teacher", "admin"),
		// 	rm.problemHandler.ImportProblems)

		// 题目标签管理
		// GET /api/v1/problems/tags
		// 响应码: 0-成功
		// problemGroup.GET("/tags", rm.problemHandler.GetProblemTags)
	}
}
