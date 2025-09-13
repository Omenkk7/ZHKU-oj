package main

import (
	"fmt"
	"net/http"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
)

// 模拟Repository层
type UserRepository struct {
	name string
}

func NewUserRepository() *UserRepository {
	repo := &UserRepository{name: "UserRepository实例"}
	fmt.Printf("🗃️  创建Repository: %s (地址: %p)\n", repo.name, repo)
	return repo
}

func (r *UserRepository) FindByID(id string) string {
	return fmt.Sprintf("用户%s的数据", id)
}

// 模拟Service层
type UserService struct {
	name string
	repo *UserRepository
}

func NewUserService(repo *UserRepository) *UserService {
	service := &UserService{
		name: "UserService实例",
		repo: repo,
	}
	fmt.Printf("⚙️  创建Service: %s (地址: %p) -> 依赖Repository: %p\n",
		service.name, service, service.repo)
	return service
}

func (s *UserService) GetUser(id string) string {
	return s.repo.FindByID(id)
}

// 模拟Handler层
type UserHandler struct {
	name    string
	service *UserService
}

func NewUserHandler(service *UserService) *UserHandler {
	handler := &UserHandler{
		name:    "UserHandler实例",
		service: service,
	}
	fmt.Printf("🎮 创建Handler: %s (地址: %p) -> 依赖Service: %p\n",
		handler.name, handler, handler.service)
	return handler
}

func (h *UserHandler) GetUser(c *gin.Context) {
	id := c.Param("id")
	user := h.service.GetUser(id)

	fmt.Printf("📝 处理请求 /users/%s - Handler地址: %p\n", id, h)

	c.JSON(http.StatusOK, gin.H{
		"user":            user,
		"handler_addr":    fmt.Sprintf("%p", h),
		"service_addr":    fmt.Sprintf("%p", h.service),
		"repository_addr": fmt.Sprintf("%p", h.service.repo),
	})
}

// 垃圾回收测试函数
func printGCStats() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("🗑️  GC运行次数: %d, 堆内存: %d KB\n", m.NumGC, m.HeapAlloc/1024)
}

func main() {
	fmt.Println("🚀 开始创建应用对象...")

	// ========== 依赖注入链 ==========
	// 1. 创建Repository层
	userRepo := NewUserRepository()

	// 2. 创建Service层（注入Repository）
	userService := NewUserService(userRepo)

	// 3. 创建Handler层（注入Service）
	userHandler := NewUserHandler(userService)

	fmt.Println("\n📋 对象引用关系:")
	fmt.Printf("Handler(%p) -> Service(%p) -> Repository(%p)\n",
		userHandler, userService, userRepo)

	// ========== 路由注册 ==========
	router := gin.New()
	router.GET("/users/:id", userHandler.GetUser) // 🔗 关键：路由持有Handler方法引用

	// 添加GC测试接口
	router.GET("/gc", func(c *gin.Context) {
		runtime.GC() // 手动触发垃圾回收
		printGCStats()
		c.JSON(200, gin.H{"message": "GC executed"})
	})

	// 添加对象地址查看接口
	router.GET("/objects", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"handler_addr":    fmt.Sprintf("%p", userHandler),
			"service_addr":    fmt.Sprintf("%p", userService),
			"repository_addr": fmt.Sprintf("%p", userRepo),
			"note":            "这些地址在整个应用运行期间保持不变",
		})
	})

	fmt.Println("\n🌐 路由注册完成，对象被路由表引用")
	fmt.Println("📊 初始GC状态:")
	printGCStats()

	// ========== 启动服务器 ==========
	server := &http.Server{
		Addr:    ":8080",
		Handler: router, // 🔗 服务器持有路由引用
	}

	fmt.Println("\n🎯 服务器启动在 :8080")
	fmt.Println("🧪 测试接口:")
	fmt.Println("  GET /users/123     - 获取用户信息")
	fmt.Println("  GET /objects       - 查看对象地址")
	fmt.Println("  GET /gc           - 手动触发GC")

	// 定时打印GC状态
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				fmt.Println("\n⏰ 定时GC检查:")
				runtime.GC()
				printGCStats()
				fmt.Printf("   Handler仍在: %p, Service仍在: %p, Repository仍在: %p\n",
					userHandler, userService, userRepo)
			}
		}
	}()

	// 启动HTTP服务器（阻塞运行）
	// 🔗 关键：ListenAndServe()保持服务器运行，维持整个引用链
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		fmt.Printf("❌ 服务器启动失败: %v\n", err)
	}
}

/*
运行示例输出:

🚀 开始创建应用对象...
🗃️  创建Repository: UserRepository实例 (地址: 0xc000010240)
⚙️  创建Service: UserService实例 (地址: 0xc000010260) -> 依赖Repository: 0xc000010240
🎮 创建Handler: UserHandler实例 (地址: 0xc000010280) -> 依赖Service: 0xc000010260

📋 对象引用关系:
Handler(0xc000010280) -> Service(0xc000010260) -> Repository(0xc000010240)

🌐 路由注册完成，对象被路由表引用
📊 初始GC状态:
🗑️  GC运行次数: 1, 堆内存: 1024 KB

🎯 服务器启动在 :8080
🧪 测试接口:
  GET /users/123     - 获取用户信息
  GET /objects       - 查看对象地址
  GET /gc           - 手动触发GC

[GIN] 2024/01/15 - 10:30:00 | 200 |     125.5µs |       127.0.0.1 | GET      "/users/123"
📝 处理请求 /users/123 - Handler地址: 0xc000010280

⏰ 定时GC检查:
🗑️  GC运行次数: 5, 堆内存: 1156 KB
   Handler仍在: 0xc000010280, Service仍在: 0xc000010260, Repository仍在: 0xc000010240

说明：即使经过多次GC，对象地址始终不变，证明对象没有被回收！
*/
