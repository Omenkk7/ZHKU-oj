# Auth工具类使用教程

## 📖 架构概述

本auth包采用了**策略模式 + 工厂模式 + 引擎模式**的设计，提供了灵活、可扩展的Token认证解决方案。

### 核心组件

```
auth/
├── token.go          # Token抽象接口
├── jwt.go            # JWT Token实现
├── tokenFactory.go   # Token工厂
├── tokenEngine.go    # Token引擎调度器
└── defaultToken.go   # 默认Token单例(如果存在)
```

### 设计模式分析

1. **策略模式**: `Token`接口定义统一规范，不同实现(`JWTToken`, `RedisToken`等)可以互换
2. **工厂模式**: `TokenFactory`根据类型创建对应的Token实现
3. **引擎模式**: `TokenEngine`作为调度器，统一管理Token实例

## 🏗️ 业务层使用教程

### 1. 基础使用方式

#### 1.1 直接使用具体实现

```go
package service

import (
    "time"
    "zhku-oj/internal/pkg/utils/auth"
)

// 在Service层使用JWT Token
type AuthService struct {
    tokenEngine *auth.TokenEngine
}

func NewAuthService() *AuthService {
    // 创建JWT Token引擎
    engine := auth.NewTokenEngine(
        auth.TokenJWT,
        "your-secret-key",      // 密钥
        time.Hour * 24 * 7,     // 7天过期
    )
    
    return &AuthService{
        tokenEngine: engine,
    }
}

// 用户登录，生成Token
func (s *AuthService) Login(userID string, additionalClaims map[string]interface{}) (string, error) {
    // 添加额外的声明
    claims := map[string]interface{}{
        "role":     "user",
        "permissions": []string{"read", "write"},
    }
    
    // 合并传入的额外声明
    for k, v := range additionalClaims {
        claims[k] = v
    }
    
    return s.tokenEngine.Generate(userID, claims)
}

// 验证Token并获取用户信息
func (s *AuthService) ValidateToken(tokenString string) (map[string]interface{}, error) {
    return s.tokenEngine.Parse(tokenString)
}

// 刷新Token
func (s *AuthService) RefreshToken(tokenString string) (string, error) {
    return s.tokenEngine.Refresh(tokenString)
}
```

#### 1.2 使用工厂模式

```go
package service

import (
    "time"
    "zhku-oj/internal/pkg/utils/auth"
)

// 根据配置动态选择Token类型
type TokenConfig struct {
    Type       auth.TokenType `yaml:"type"`       // "jwt" 或 "redis"
    Secret     string         `yaml:"secret"`
    Expiration time.Duration  `yaml:"expiration"`
}

func CreateTokenService(config TokenConfig) *AuthService {
    // 使用工厂创建Token实例
    token := auth.TokenFactory(config.Type, config.Secret, config.Expiration)
    
    return &AuthService{
        token: token,
    }
}

type AuthService struct {
    token auth.Token
}

func (s *AuthService) GenerateUserToken(userID string) (string, error) {
    claims := map[string]interface{}{
        "iat": time.Now().Unix(),
        "type": "access_token",
    }
    return s.token.Generate(userID, claims)
}
```

### 2. 在Handler层的集成使用

#### 2.1 认证中间件

```go
package middleware

import (
    "net/http"
    "strings"
    "zhku-oj/internal/pkg/utils/auth"
    "github.com/gin-gonic/gin"
)

// JWT认证中间件
func JWTAuthMiddleware(tokenEngine *auth.TokenEngine) gin.HandlerFunc {
    return func(c *gin.Context) {
        // 从Header中获取Token
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing authorization header"})
            c.Abort()
            return
        }
        
        // 解析Bearer Token
        tokenString := strings.TrimPrefix(authHeader, "Bearer ")
        if tokenString == authHeader {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format"})
            c.Abort()
            return
        }
        
        // 验证Token
        claims, err := tokenEngine.Parse(tokenString)
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
            c.Abort()
            return
        }
        
        // 将用户信息存储到Context
        c.Set("user_id", claims["sub"])
        c.Set("user_claims", claims)
        
        c.Next()
    }
}
```

#### 2.2 Handler中的使用

```go
package handler

import (
    "net/http"
    "zhku-oj/internal/pkg/utils/auth"
    "zhku-oj/internal/service"
    "github.com/gin-gonic/gin"
)

type AuthHandler struct {
    authService *service.AuthService
    tokenEngine *auth.TokenEngine
}

func NewAuthHandler(authService *service.AuthService, tokenEngine *auth.TokenEngine) *AuthHandler {
    return &AuthHandler{
        authService: authService,
        tokenEngine: tokenEngine,
    }
}

// 登录接口
func (h *AuthHandler) Login(c *gin.Context) {
    var req struct {
        Username string `json:"username" binding:"required"`
        Password string `json:"password" binding:"required"`
    }
    
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    // 验证用户凭据
    user, err := h.authService.ValidateCredentials(req.Username, req.Password)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
        return
    }
    
    // 生成Token
    claims := map[string]interface{}{
        "username": user.Username,
        "role":     user.Role,
        "school":   user.School,
    }
    
    token, err := h.tokenEngine.Generate(user.ID, claims)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{
        "access_token": token,
        "token_type":   "Bearer",
        "expires_in":   7 * 24 * 3600, // 7天
        "user_info": gin.H{
            "id":       user.ID,
            "username": user.Username,
            "role":     user.Role,
        },
    })
}
```

## 🚀 Client端注册教程

### 1. 依赖注入配置

#### 1.1 在main.go中初始化

```go
package main

import (
    "log"
    "time"
    "zhku-oj/internal/pkg/utils/auth"
    "zhku-oj/internal/service"
    "zhku-oj/internal/handler"
    "github.com/gin-gonic/gin"
)

func main() {
    // 1. 创建Token引擎
    tokenEngine := auth.NewTokenEngine(
        auth.TokenJWT,
        "your-super-secret-key-2024", // 生产环境从配置文件读取
        time.Hour * 24 * 7,           // 7天过期
    )
    
    // 2. 创建Service层
    authService := service.NewAuthService(/* 依赖注入 */)
    userService := service.NewUserService(/* 依赖注入 */)
    
    // 3. 创建Handler层
    authHandler := handler.NewAuthHandler(authService, tokenEngine)
    userHandler := handler.NewUserHandler(userService)
    
    // 4. 创建中间件
    jwtMiddleware := middleware.JWTAuthMiddleware(tokenEngine)
    
    // 5. 设置路由
    router := gin.Default()
    
    // 公开路由
    public := router.Group("/api/v1")
    {
        public.POST("/auth/login", authHandler.Login)
        public.POST("/auth/register", authHandler.Register)
        public.POST("/auth/refresh", authHandler.RefreshToken)
    }
    
    // 需要认证的路由
    protected := router.Group("/api/v1")
    protected.Use(jwtMiddleware) // 应用JWT中间件
    {
        protected.GET("/auth/me", authHandler.GetCurrentUser)
        protected.GET("/users/profile", userHandler.GetProfile)
        protected.PUT("/users/profile", userHandler.UpdateProfile)
        protected.POST("/problems", problemHandler.CreateProblem)
    }
    
    log.Fatal(router.Run(":8080"))
}
```

## 📝 总结

该auth工具类提供了：

1. **灵活的架构**: 支持多种Token实现(JWT、Redis等)
2. **简单的API**: 统一的Generate/Parse/Refresh接口
3. **易于集成**: 可直接用于Service、Handler、Middleware
4. **可扩展性**: 通过工厂模式轻松添加新的Token类型
5. **生产就绪**: 包含安全最佳实践和错误处理

使用时建议按照**配置 → 初始化 → 注入 → 使用**的流程进行集成。