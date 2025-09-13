# 路由映射集中管理设计文档

## 🎯 设计目标

将原本散落在 `main.go` 中的所有路由配置集中管理，提高代码的可维护性和可扩展性。

## 📁 文件结构

```
internal/router/
├── router.go       # 路由管理器主文件
├── auth.go         # 认证相关路由
├── user.go         # 用户相关路由  
├── problem.go      # 题目相关路由
├── submission.go   # 提交相关路由
├── admin.go        # 管理员相关路由
└── health.go       # 健康检查路由
```

## 🔧 核心设计

### 1. 路由管理器模式

```go
type RouterManager struct {
    // Handler依赖注入
    authHandler       *auth.AuthHandler
    userHandler       *user.UserHandler
    problemHandler    *problem.ProblemHandler
    submissionHandler *submission.SubmissionHandler
    adminHandler      *admin.AdminHandler
}
```

### 2. 模块化路由配置

每个业务模块的路由独立管理：

- **认证模块** ([auth.go](file://g:\code-project\zhku-oj\internal\router\auth.go)): 注册、登录、登出、Token验证
- **用户模块** ([user.go](file://g:\code-project\zhku-oj\internal\router\user.go)): 用户信息、个人资料、统计数据
- **题目模块** ([problem.go](file://g:\code-project\zhku-oj\internal\router\problem.go)): 题目CRUD、搜索、标签管理
- **提交模块** ([submission.go](file://g:\code-project\zhku-oj\internal\router\submission.go)): 代码提交、判题结果查询
- **管理模块** ([admin.go](file://g:\code-project\zhku-oj\internal\router\admin.go)): 系统管理、用户管理、数据统计

### 3. 统一的路由注册

```go
func (rm *RouterManager) SetupRoutes(router *gin.Engine) {
    // 全局中间件
    router.Use(middleware.Logger())
    router.Use(middleware.Recovery()) 
    router.Use(middleware.CORS())

    // 健康检查
    rm.setupHealthRoutes(router)

    // API路由组
    v1 := router.Group("/api/v1")
    {
        rm.setupAuthRoutes(v1)      // 认证路由
        rm.setupUserRoutes(v1)      // 用户路由
        rm.setupProblemRoutes(v1)   // 题目路由
        rm.setupSubmissionRoutes(v1) // 提交路由
        rm.setupAdminRoutes(v1)     // 管理路由
    }
}
```

## 📋 完整的API路由清单

### 🔐 认证相关 (/api/v1/auth)

| 方法 | 路径 | 功能 | 权限 | 响应码 |
|------|------|------|------|--------|
| POST | `/register` | 用户注册 | 无 | 0,10002,20002 |
| POST | `/login` | 用户登录 | 无 | 0,10002,20010 |
| POST | `/logout` | 用户登出 | 认证 | 0,10003 |
| PUT | `/password` | 修改密码 | 认证 | 0,10002,20013 |
| GET | `/verify` | 验证Token | 认证 | 0,10003,10009 |

### 👤 用户相关 (/api/v1/users)

| 方法 | 路径 | 功能 | 权限 | 响应码 |
|------|------|------|------|--------|
| GET | `/profile` | 获取个人信息 | 认证 | 0,10002,20001 |
| PUT | `/profile` | 更新个人信息 | 认证 | 0,10002,20001 |
| GET | `/:id` | 获取用户信息 | 认证 | 0,10002,20001 |
| GET | `/:id/stats` | 获取用户统计 | 认证 | 0,10002,20001 |

### 📚 题目相关 (/api/v1/problems)

| 方法 | 路径 | 功能 | 权限 | 响应码 |
|------|------|------|------|--------|
| GET | `` | 题目列表 | 认证 | 0,10002 |
| GET | `/:id` | 题目详情 | 认证 | 0,10002,30001,30006 |
| POST | `` | 创建题目 | 教师/管理员 | 0,10002,10004,30002 |
| PUT | `/:id` | 更新题目 | 教师/管理员 | 0,10002,10004,30001 |
| DELETE | `/:id` | 删除题目 | 管理员 | 0,10002,10004,30001 |

### 📝 提交相关 (/api/v1/submissions)

| 方法 | 路径 | 功能 | 权限 | 响应码 |
|------|------|------|------|--------|
| POST | `` | 提交代码 | 认证 | 0,10002,30001,40004,40007 |
| GET | `/:id` | 提交详情 | 认证 | 0,10002,40001,40008 |
| GET | `` | 提交列表 | 认证 | 0,10002 |

### 🛡️ 管理相关 (/api/v1/admin)

| 方法 | 路径 | 功能 | 权限 | 响应码 |
|------|------|------|------|--------|
| GET | `/dashboard` | 管理仪表板 | 管理员 | 0,10004 |
| GET | `/system/status` | 系统状态 | 管理员 | 0,10004 |
| POST | `/users` | 创建用户 | 管理员 | 0,10002,20002 |
| GET | `/users` | 用户列表 | 管理员 | 0,10002 |
| PUT | `/users/:id` | 更新用户 | 管理员 | 0,10002,20001 |
| DELETE | `/users/:id` | 删除用户 | 管理员 | 0,10002,20001 |
| PUT | `/users/:id/activate` | 激活用户 | 管理员 | 0,10002,20001 |
| PUT | `/users/:id/deactivate` | 停用用户 | 管理员 | 0,10002,20001 |

### 🏥 健康检查

| 方法 | 路径 | 功能 | 权限 | 响应码 |
|------|------|------|------|--------|
| GET | `/health` | 基础健康检查 | 无 | 200 |
| GET | `/health/detailed` | 详细健康检查 | 无 | 200 |
| GET | `/info` | 服务信息 | 无 | 200 |

## 🚀 使用方式

### 在main.go中使用

```go
func main() {
    // ... 初始化Handler等

    // 创建路由器
    router := gin.New()

    // 创建路由管理器并设置所有路由
    routerManager := router.NewRouterManager(
        authHandler,
        userHandler, 
        problemHandler,
        submissionHandler,
        adminHandler,
    )
    routerManager.SetupRoutes(router)

    // 启动服务器
    server := &http.Server{
        Addr:    cfg.Server.Port,
        Handler: router,
    }
    server.ListenAndServe()
}
```

## ✅ 优势对比

### 改进前 (main.go中管理)

❌ **问题**:
- main函数过于庞大 (150+ 行路由配置)
- 路由散落各处，难以维护
- 添加新路由需要修改main.go
- 路由逻辑与启动逻辑混合

### 改进后 (集中式管理)

✅ **优势**:
- **模块化**: 每个业务模块路由独立管理
- **可维护**: 路由变更只需修改对应模块文件
- **可扩展**: 新增模块只需添加对应路由文件
- **清晰性**: main.go专注于应用启动逻辑
- **文档化**: 每个路由都有详细的注释说明

## 🔮 扩展建议

### 1. 路由中间件扩展
```go
// 为特定路由组添加专用中间件
problemGroup.Use(middleware.ProblemAccessControl())
submissionGroup.Use(middleware.SubmissionRateLimit())
```

### 2. 版本管理
```go
// 支持API版本控制
v1 := router.Group("/api/v1")
v2 := router.Group("/api/v2") 
```

### 3. 路由缓存
```go
// 为查询接口添加缓存中间件
userGroup.GET("/:id", middleware.Cache(5*time.Minute), userHandler.GetUser)
```

### 4. 自动化文档生成
```go
// 集成Swagger自动生成API文档
router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
```

这种集中式路由管理方案大大提升了项目的可维护性和可扩展性，符合Go项目的最佳实践。