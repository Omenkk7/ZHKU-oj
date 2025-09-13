# API 响应体使用示例

## 📋 概述

本文档展示了校园Java-OJ系统中统一响应体的使用方法和示例。

## 🔧 统一响应格式

### 基础响应结构
```json
{
    "code": 0,           // 错误码，0表示成功，非0表示各种错误
    "message": "成功",    // 响应消息
    "data": {},          // 响应数据，可选
    "trace_id": "..."    // 链路追踪ID，可选
}
```

### 分页响应结构
```json
{
    "code": 0,
    "message": "成功",
    "data": [],
    "pagination": {
        "page": 1,
        "page_size": 20,
        "total": 100,
        "total_pages": 5
    }
}
```

## 📋 错误码规范

### 错误码分类
- **通用错误码**: 10000-10999
- **用户模块**: 20000-20999  
- **题目模块**: 30000-30999
- **提交模块**: 40000-40999
- **判题模块**: 50000-50999
- **管理模块**: 60000-60999
- **竞赛模块**: 70000-70999

### 主要错误码
```javascript
const ERROR_CODES = {
    // 通用错误
    SUCCESS: 0,
    SYSTEM_ERROR: 10001,
    INVALID_PARAMS: 10002,
    UNAUTHORIZED: 10003,
    FORBIDDEN: 10004,
    NOT_FOUND: 10005,
    
    // 用户模块
    USER_NOT_FOUND: 20001,
    USERNAME_ALREADY_EXISTS: 20003,
    EMAIL_ALREADY_EXISTS: 20004,
    INVALID_PASSWORD: 20006,
    USER_DISABLED: 20008,
    LOGIN_FAILED: 20010,
    
    // 提交模块
    SUBMISSION_NOT_FOUND: 40001,
    CODE_TOO_LONG: 40004,
    CODE_EMPTY: 40005,
    DUPLICATE_SUBMISSION: 40007,
    SUBMISSION_TOO_FREQUENT: 40010,
    
    // 判题模块
    JUDGE_SYSTEM_ERROR: 50001,
    COMPILE_ERROR: 50004,
    RUNTIME_ERROR: 50005,
    TIME_LIMIT_EXCEEDED: 50006,
    MEMORY_LIMIT_EXCEEDED: 50007,
    WRONG_ANSWER: 50009
};
```

## 🎯 API 示例

### 1. 用户登录
**请求**:
```bash
POST /api/v1/auth/login
Content-Type: application/json

{
    "username": "zhangsan",
    "password": "123456"
}
```

**成功响应**:
```json
{
    "code": 0,
    "message": "成功",
    "data": {
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
        "user": {
            "id": "507f1f77bcf86cd799439011",
            "username": "zhangsan",
            "email": "zhangsan@example.com",
            "role": "student",
            "real_name": "张三"
        }
    }
}
```

**失败响应**:
```json
{
    "code": 20010,
    "message": "登录失败",
    "data": {
        "detail": "用户名或密码错误"
    }
}
```

### 2. 获取用户列表
**请求**:
```bash
GET /api/v1/admin/users?page=1&page_size=20&role=student&keyword=张
Authorization: Bearer {token}
```

**成功响应**:
```json
{
    "code": 0,
    "message": "成功",
    "data": [
        {
            "id": "507f1f77bcf86cd799439011",
            "username": "zhangsan",
            "email": "zhangsan@example.com",
            "real_name": "张三",
            "role": "student",
            "class": "软件工程1班",
            "is_active": true,
            "created_at": "2024-01-15T10:30:00Z"
        }
    ],
    "pagination": {
        "page": 1,
        "page_size": 20,
        "total": 1,
        "total_pages": 1
    }
}
```

### 3. 代码提交
**请求**:
```bash
POST /api/v1/submissions
Authorization: Bearer {token}
Content-Type: application/json

{
    "problem_id": "507f1f77bcf86cd799439012",
    "language": "java",
    "code": "public class Main {\n    public static void main(String[] args) {\n        System.out.println(\"Hello World\");\n    }\n}"
}
```

**成功响应**:
```json
{
    "code": 0,
    "message": "成功",
    "data": {
        "submission_id": "507f1f77bcf86cd799439013",
        "status": "PENDING"
    }
}
```

**重复提交错误**:
```json
{
    "code": 40007,
    "message": "请勿重复提交相同代码",
    "data": {
        "detail": "相同代码在10分钟内已提交"
    }
}
```

### 4. 获取判题结果
**请求**:
```bash
GET /api/v1/submissions/507f1f77bcf86cd799439013
Authorization: Bearer {token}
```

**判题中响应**:
```json
{
    "code": 0,
    "message": "成功",
    "data": {
        "id": "507f1f77bcf86cd799439013",
        "problem_id": "507f1f77bcf86cd799439012",
        "status": "JUDGING",
        "score": 0,
        "time_used": 0,
        "memory_used": 0,
        "submitted_at": "2024-01-15T10:35:00Z"
    }
}
```

**判题完成响应**:
```json
{
    "code": 0,
    "message": "成功",
    "data": {
        "id": "507f1f77bcf86cd799439013",
        "problem_id": "507f1f77bcf86cd799439012",
        "status": "ACCEPTED",
        "score": 100,
        "time_used": 125,
        "memory_used": 2048,
        "test_results": [
            {
                "test_case_id": "1",
                "status": "ACCEPTED",
                "time_used": 125,
                "memory_used": 2048,
                "score": 100
            }
        ],
        "submitted_at": "2024-01-15T10:35:00Z",
        "judged_at": "2024-01-15T10:35:03Z"
    }
}
```

### 5. 参数验证错误
**请求**:
```bash
POST /api/v1/admin/users
Authorization: Bearer {token}
Content-Type: application/json

{
    "username": "",
    "email": "invalid-email",
    "password": "123"
}
```

**响应**:
```json
{
    "code": 10002,
    "message": "参数错误",
    "data": {
        "errors": {
            "username": "用户名不能为空",
            "email": "邮箱格式不正确", 
            "password": "密码长度至少6位"
        }
    }
}
```

## 🚀 前端集成示例

### JavaScript/TypeScript
```typescript
// 定义响应接口
interface ApiResponse<T = any> {
    code: number;
    message: string;
    data?: T;
    trace_id?: string;
}

interface PaginationResponse<T> extends ApiResponse<T[]> {
    pagination: {
        page: number;
        page_size: number;
        total: number;
        total_pages: number;
    };
}

// 统一响应处理函数
function handleApiResponse<T>(response: ApiResponse<T>): T {
    if (response.code === 0) {
        return response.data!;
    }
    
    // 错误处理
    const errorMessage = response.message;
    const errorDetail = response.data?.detail;
    
    // 根据错误码进行不同处理
    switch (response.code) {
        case 10003: // UNAUTHORIZED
            // 清除本地token，跳转登录页
            localStorage.removeItem('token');
            window.location.href = '/login';
            break;
            
        case 20001: // USER_NOT_FOUND
            showError('用户不存在');
            break;
            
        case 40007: // DUPLICATE_SUBMISSION
            showWarning('请勿重复提交相同代码');
            break;
            
        case 50006: // TIME_LIMIT_EXCEEDED
            showInfo('代码执行超时');
            break;
            
        default:
            showError(errorMessage + (errorDetail ? ': ' + errorDetail : ''));
    }
    
    throw new Error(errorMessage);
}

// API 调用示例
class ApiClient {
    private baseUrl = '/api/v1';
    
    async get<T>(url: string): Promise<T> {
        const response = await fetch(`${this.baseUrl}${url}`, {
            headers: {
                'Authorization': `Bearer ${localStorage.getItem('token')}`
            }
        });
        const data = await response.json();
        return handleApiResponse<T>(data);
    }
    
    async post<T>(url: string, body: any): Promise<T> {
        const response = await fetch(`${this.baseUrl}${url}`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${localStorage.getItem('token')}`
            },
            body: JSON.stringify(body)
        });
        const data = await response.json();
        return handleApiResponse<T>(data);
    }
}

// 使用示例
const api = new ApiClient();

// 获取用户信息
try {
    const user = await api.get<User>('/users/profile');
    console.log('用户信息:', user);
} catch (error) {
    console.error('获取用户信息失败:', error.message);
}

// 提交代码
try {
    const result = await api.post<SubmissionResult>('/submissions', {
        problem_id: problemId,
        language: 'java',
        code: codeContent
    });
    console.log('提交成功:', result);
} catch (error) {
    console.error('提交失败:', error.message);
}
```

### Vue.js 示例
```vue
<template>
    <div>
        <el-table :data="users" :loading="loading">
            <el-table-column prop="username" label="用户名" />
            <el-table-column prop="real_name" label="真实姓名" />
            <el-table-column prop="email" label="邮箱" />
        </el-table>
        
        <el-pagination
            :current-page="currentPage"
            :page-size="pageSize"
            :total="total"
            @current-change="handlePageChange"
        />
    </div>
</template>

<script>
export default {
    data() {
        return {
            users: [],
            loading: false,
            currentPage: 1,
            pageSize: 20,
            total: 0
        };
    },
    
    async mounted() {
        await this.loadUsers();
    },
    
    methods: {
        async loadUsers() {
            this.loading = true;
            try {
                const url = `/admin/users?page=${this.currentPage}&page_size=${this.pageSize}`;
                const data = await this.$api.get(url);
                
                this.users = data.data;
                this.total = data.pagination.total;
                
                this.$message.success('获取用户列表成功');
            } catch (error) {
                // 错误已在 handleApiResponse 中处理
                console.error('加载用户列表失败:', error.message);
            } finally {
                this.loading = false;
            }
        },
        
        async handlePageChange(page) {
            this.currentPage = page;
            await this.loadUsers();
        }
    }
};
</script>
```

## 📝 后端开发者指南

### Go Handler 示例
```go
package user

import (
    "zhku-oj/internal/pkg/errors"
    "zhku-oj/internal/pkg/utils"
    "github.com/gin-gonic/gin"
)

func (h *UserHandler) CreateUser(c *gin.Context) {
    var req CreateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        // 参数验证错误
        utils.SendError(c, errors.INVALID_PARAMS)
        return
    }
    
    user, err := h.userService.CreateUser(c.Request.Context(), &req)
    if err != nil {
        // 业务错误，统一处理
        utils.HandleError(c, err)
        return
    }
    
    // 成功响应
    utils.SendSuccess(c, user)
}
```

### Service 层错误处理
```go
package impl

import "zhku-oj/internal/pkg/errors"

func (s *userService) CreateUser(ctx context.Context, req *CreateUserRequest) (*User, error) {
    // 检查用户名是否存在
    exists, err := s.userRepo.ExistsByUsername(ctx, req.Username)
    if err != nil {
        return nil, errors.Wrap(errors.DATABASE_ERROR, err)
    }
    if exists {
        return nil, errors.NewUsernameAlreadyExists()
    }
    
    // 其他业务逻辑...
    
    return user, nil
}
```

这个统一的响应体系统为前后端提供了清晰、一致的API交互规范，大大提升了开发效率和用户体验。