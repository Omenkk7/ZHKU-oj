# Go响应封装系统指导文档 - 面向Java开发者

## 📋 概述

本文档面向有Java/Spring Boot开发经验的开发者，帮助快速理解和使用Go项目中的统一响应封装系统。

## 🔄 Java vs Go 对比

### Java Spring Boot 风格
```java
// Spring Boot 典型响应
@RestController
public class UserController {
    @GetMapping("/users/{id}")
    public ResponseEntity<ApiResponse<User>> getUser(@PathVariable String id) {
        try {
            User user = userService.findById(id);
            return ResponseEntity.ok(ApiResponse.success(user));
        } catch (UserNotFoundException e) {
            return ResponseEntity.ok(ApiResponse.error(20001, "用户不存在"));
        }
    }
}

// 响应数据结构
public class ApiResponse<T> {
    private int code;
    private String message;
    private T data;
    // getters/setters...
}
```

### Go Gin 风格（本项目）
```go
// Go Gin 等价实现
func (h *UserHandler) GetUser(c *gin.Context) {
    userID := c.Param("id")
    
    user, err := h.userService.GetByID(userID)
    if err != nil {
        if errors.IsBusinessError(err) {
            utils.SendBusinessError(c, err.(*errors.BusinessError))
            return
        }
        utils.SendError(c, errors.SYSTEM_ERROR)
        return
    }
    
    utils.SendSuccess(c, user)
}
```

## 🏗️ 系统架构

### 文件结构对应关系

| Go文件 | Java等价概念 | 作用 |
|--------|-------------|------|
| `codes.go` | `ErrorCodeEnum.java` | 错误码常量定义 |
| `errors.go` | `BusinessException.java` | 业务异常类 |
| `response.go` | `ResponseUtil.java` | 响应工具类 |

### 核心组件

```
响应封装系统
├── 错误码管理 (codes.go)
│   ├── 分层错误码定义
│   ├── 错误消息映射
│   └── 错误分类判断
├── 业务异常处理 (errors.go)
│   ├── BusinessError结构体
│   ├── 异常构造函数
│   └── 异常包装方法
└── 响应工具 (response.go)
    ├── 统一响应结构
    ├── 成功响应方法
    ├── 错误响应方法
    └── 分页响应支持
```

## 📊 错误码体系 (codes.go)

### 分层设计

```go
// 5位数字编码：AABBB
// AA: 模块代码 (10:通用, 20:用户, 30:题目, 40:提交, 50:判题, 60:管理)
// BBB: 具体错误代码 (001-999)

const (
    // 通用错误码 (10000-10999)
    SUCCESS      = 0     // 成功
    SYSTEM_ERROR = 10001 // 系统内部错误
    
    // 用户模块错误码 (20000-20999)
    USER_NOT_FOUND = 20001 // 用户不存在
    
    // 题目模块错误码 (30000-30999)
    PROBLEM_NOT_FOUND = 30001 // 题目不存在
)
```

### Java对比
```java
// Java枚举方式
public enum ErrorCode {
    SUCCESS(0, "成功"),
    SYSTEM_ERROR(10001, "系统内部错误"),
    USER_NOT_FOUND(20001, "用户不存在");
    
    private final int code;
    private final String message;
    
    // 构造函数和getter方法...
}
```

### 使用方法
```go
// 获取错误消息
message := errors.GetErrorMessage(errors.USER_NOT_FOUND)
// 输出: "用户不存在"

// 判断错误类型
if errors.IsUserError(20001) {
    // 处理用户模块错误
}
```

## 🚨 业务异常处理 (errors.go)

### BusinessError结构体

```go
// Go版本
type BusinessError struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
    Detail  string `json:"detail,omitempty"`
}

// 实现error接口
func (e *BusinessError) Error() string {
    return fmt.Sprintf("Code: %d, Message: %s", e.Code, e.Message)
}
```

### Java对比
```java
// Java版本
public class BusinessException extends RuntimeException {
    private int code;
    private String message;
    private String detail;
    
    public BusinessException(int code, String message) {
        super(message);
        this.code = code;
        this.message = message;
    }
    
    // getters/setters...
}
```

### 异常创建方式

#### 1. 通用创建方法
```go
// Go
err := errors.New(errors.USER_NOT_FOUND, "用户ID无效")
err := errors.Newf(errors.USER_NOT_FOUND, "用户ID %s 无效", userID)
err := errors.Wrap(errors.SYSTEM_ERROR, originalErr)
```

```java
// Java
throw new BusinessException(ErrorCode.USER_NOT_FOUND, "用户ID无效");
throw new BusinessException(ErrorCode.USER_NOT_FOUND, 
    String.format("用户ID %s 无效", userId));
```

#### 2. 预定义便捷方法
```go
// Go - 预定义方法
err := errors.NewUserNotFound("用户不存在")
err := errors.NewInvalidPassword("密码错误")
err := errors.NewSystemError("数据库连接失败")
```

```java
// Java - 静态方法
public class BusinessExceptions {
    public static BusinessException userNotFound(String detail) {
        return new BusinessException(ErrorCode.USER_NOT_FOUND, detail);
    }
    
    public static BusinessException invalidPassword(String detail) {
        return new BusinessException(ErrorCode.INVALID_PASSWORD, detail);
    }
}
```

### 异常检测和转换
```go
// 检测是否为业务异常
if errors.IsBusinessError(err) {
    bizErr, _ := errors.GetBusinessError(err)
    fmt.Printf("业务错误码: %d", bizErr.GetCode())
}
```

## 📤 响应工具 (response.go)

### 响应结构体

```go
// Go版本
type Response struct {
    Code    int         `json:"code"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
    TraceID string      `json:"trace_id,omitempty"`
}

// 分页响应
type PageResponse struct {
    Response
    Pagination Pagination `json:"pagination"`
}
```

### Java对比
```java
// Java版本
public class ApiResponse<T> {
    private int code;
    private String message;
    private T data;
    private String traceId;
    
    // 静态工厂方法
    public static <T> ApiResponse<T> success(T data) {
        return new ApiResponse<>(0, "成功", data);
    }
    
    public static ApiResponse<Void> error(int code, String message) {
        return new ApiResponse<>(code, message, null);
    }
}
```

### 核心响应方法

#### 1. 成功响应
```go
// Go - 简化版本（用户需求）
utils.SendSuccess(c, userData)
// 输出: {"code":0,"message":"成功","data":{...}}

// Go - 完整版本
utils.SuccessResponse(c, "获取用户信息成功", userData)
```

```java
// Java
return ResponseEntity.ok(ApiResponse.success(userData));
```

#### 2. 错误响应
```go
// Go - 通过错误码
utils.SendError(c, errors.USER_NOT_FOUND)
// 输出: {"code":20001,"message":"用户不存在"}

// Go - 业务异常
bizErr := errors.NewUserNotFound("用户ID不存在")
utils.SendBusinessError(c, bizErr)

// Go - 带详情
utils.SendErrorWithDetail(c, errors.USER_NOT_FOUND, "用户ID格式错误")
```

```java
// Java
return ResponseEntity.ok(ApiResponse.error(20001, "用户不存在"));
```

#### 3. 分页响应
```go
// Go
utils.SendSuccessWithPagination(c, userList, page, pageSize, total)
```

```java
// Java
PageInfo<User> pageInfo = new PageInfo<>(userList);
return ResponseEntity.ok(ApiResponse.success(pageInfo));
```

## 🎯 实际使用示例

### 完整的Controller示例

#### Go版本
```go
// UserHandler
func (h *UserHandler) GetUserList(c *gin.Context) {
    // 1. 参数验证
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
    
    if page < 1 || pageSize < 1 || pageSize > 100 {
        utils.SendError(c, errors.INVALID_PARAMS)
        return
    }
    
    // 2. 业务逻辑调用
    users, total, err := h.userService.GetUserList(page, pageSize)
    if err != nil {
        // 统一错误处理
        utils.HandleError(c, err)
        return
    }
    
    // 3. 成功响应
    utils.SendSuccessWithPagination(c, users, page, pageSize, total)
}

func (h *UserHandler) CreateUser(c *gin.Context) {
    var req model.CreateUserRequest
    
    // 1. 参数绑定和验证
    if err := c.ShouldBindJSON(&req); err != nil {
        utils.ValidationErrorResponse(c, map[string]string{
            "request": "请求参数格式错误",
        })
        return
    }
    
    // 2. 业务逻辑
    user, err := h.userService.CreateUser(&req)
    if err != nil {
        // 自动处理业务异常和系统异常
        utils.HandleError(c, err)
        return
    }
    
    // 3. 成功响应
    utils.SendSuccess(c, user)
}
```

#### Java对比版本
```java
@RestController
@RequestMapping("/api/users")
public class UserController {
    
    @GetMapping
    public ResponseEntity<ApiResponse<PageInfo<User>>> getUserList(
            @RequestParam(defaultValue = "1") int page,
            @RequestParam(defaultValue = "10") int pageSize) {
        
        if (page < 1 || pageSize < 1 || pageSize > 100) {
            return ResponseEntity.ok(ApiResponse.error(10002, "参数错误"));
        }
        
        try {
            PageInfo<User> pageInfo = userService.getUserList(page, pageSize);
            return ResponseEntity.ok(ApiResponse.success(pageInfo));
        } catch (BusinessException e) {
            return ResponseEntity.ok(ApiResponse.error(e.getCode(), e.getMessage()));
        }
    }
    
    @PostMapping
    public ResponseEntity<ApiResponse<User>> createUser(@RequestBody CreateUserRequest request) {
        try {
            User user = userService.createUser(request);
            return ResponseEntity.ok(ApiResponse.success(user));
        } catch (BusinessException e) {
            return ResponseEntity.ok(ApiResponse.error(e.getCode(), e.getMessage()));
        }
    }
}
```

### Service层错误处理

#### Go版本
```go
func (s *UserService) GetByID(id string) (*model.User, error) {
    // 参数验证
    if id == "" {
        return nil, errors.NewInvalidParams("用户ID不能为空")
    }
    
    // 数据库查询
    user, err := s.userRepo.FindByID(id)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            return nil, errors.NewUserNotFound("用户不存在")
        }
        return nil, errors.NewDatabaseError("查询用户失败")
    }
    
    return user, nil
}
```

#### Java对比
```java
@Service
public class UserService {
    
    public User getById(String id) {
        if (StringUtils.isEmpty(id)) {
            throw BusinessExceptions.invalidParams("用户ID不能为空");
        }
        
        User user = userRepository.findById(id);
        if (user == null) {
            throw BusinessExceptions.userNotFound("用户不存在");
        }
        
        return user;
    }
}
```

## 🔧 最佳实践

### 1. 错误处理流程
```go
// 推荐的错误处理模式
func (h *Handler) SomeMethod(c *gin.Context) {
    // 1. 参数验证 - 使用INVALID_PARAMS
    if invalidParam {
        utils.SendError(c, errors.INVALID_PARAMS)
        return
    }
    
    // 2. 权限检查 - 使用FORBIDDEN
    if !hasPermission {
        utils.SendError(c, errors.FORBIDDEN)
        return
    }
    
    // 3. 业务逻辑调用
    result, err := h.service.DoSomething()
    if err != nil {
        // 4. 统一错误处理
        utils.HandleError(c, err)
        return
    }
    
    // 5. 成功响应
    utils.SendSuccess(c, result)
}
```

### 2. 业务异常抛出
```go
// Service层推荐做法
func (s *Service) BusinessMethod() error {
    // 使用预定义的异常构造函数
    if userNotExists {
        return errors.NewUserNotFound("用户不存在")
    }
    
    // 包装底层错误
    if dbErr != nil {
        return errors.Wrap(errors.DATABASE_ERROR, dbErr)
    }
    
    return nil
}
```

### 3. 响应格式一致性
```go
// ✅ 推荐：使用统一的响应方法
utils.SendSuccess(c, data)
utils.SendError(c, errorCode)
utils.SendBusinessError(c, businessError)

// ❌ 不推荐：直接使用gin的JSON方法
c.JSON(200, gin.H{"code": 0, "data": data})
```

## 📋 常用错误码参考

| 错误码 | 常量名 | 含义 | Java等价 |
|--------|--------|------|----------|
| 0 | SUCCESS | 成功 | 200 OK |
| 10001 | SYSTEM_ERROR | 系统错误 | 500 Internal Server Error |
| 10002 | INVALID_PARAMS | 参数错误 | 400 Bad Request |
| 10003 | UNAUTHORIZED | 未授权 | 401 Unauthorized |
| 10004 | FORBIDDEN | 无权限 | 403 Forbidden |
| 20001 | USER_NOT_FOUND | 用户不存在 | 自定义业务错误 |
| 20006 | INVALID_PASSWORD | 密码错误 | 自定义业务错误 |

## 🚀 快速上手清单

1. **引入响应工具**
   ```go
   import "zhku-oj/internal/pkg/utils"
   import "zhku-oj/internal/pkg/errors"
   ```

2. **Handler层统一格式**
   ```go
   // 成功响应
   utils.SendSuccess(c, data)
   
   // 错误响应
   utils.SendError(c, errors.ERROR_CODE)
   
   // 统一错误处理
   utils.HandleError(c, err)
   ```

3. **Service层异常抛出**
   ```go
   // 使用预定义异常
   return errors.NewUserNotFound("详细信息")
   
   // 通用异常构造
   return errors.New(errors.USER_NOT_FOUND, "详细信息")
   ```

4. **前端响应处理**
   ```javascript
   // 前端JavaScript处理
   axios.post('/api/users', data)
     .then(response => {
       if (response.data.code === 0) {
         // 成功处理
         console.log(response.data.data);
       } else {
         // 错误处理
         alert(response.data.message);
       }
     });
   ```

通过以上指导，Java开发者可以快速理解和使用Go项目中的响应封装系统，保持与Java Spring Boot开发习惯的一致性。