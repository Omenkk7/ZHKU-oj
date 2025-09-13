# 校园Java-OJ在线判题系统

## 📋 项目概述

**项目名称**: 校园Java-OJ在线判题系统  
**技术栈**: Go + MongoDB + RabbitMQ + go-judge  
**目标用户**: 100-200名在校学生  
**核心功能**: Java代码在线评测、题库管理、用户管理、成绩统计

## 🏗️ 系统架构

本项目采用六层架构模型：
- **前端层**: 学生Web界面、教师管理端、管理员控制台
- **业务服务层**: Gin Web服务、用户服务、题目服务、提交服务、统计服务
- **消息队列层**: RabbitMQ任务队列、结果通知、统计更新、死信队列
- **判题处理层**: 判题管理器、沙箱负载均衡器、go-judge集群、文件缓存管理器
- **数据存储层**: MongoDB、Redis、MinIO文件存储
- **日志层**: Logrus应用日志、系统监控、链路追踪

## 🚀 快速开始

### 环境要求
- Go 1.18+
- MongoDB 5.0+
- Redis 6.0+
- RabbitMQ 3.9+
- go-judge 沙箱环境

### 安装步骤

1. **克隆项目**
```bash
git clone <repository-url>
cd zhku-oj
```

2. **安装依赖**
```bash
go mod tidy
```

3. **配置文件**
```bash
cp configs/config.yaml configs/config.local.yaml
# 编辑配置文件，修改数据库连接等信息
```

4. **启动基础服务**
```bash
# 使用Docker Compose启动基础服务
docker-compose -f deployments/docker/docker-compose.yml up -d mongodb redis rabbitmq go-judge-1 go-judge-2
```

5. **启动应用服务**
```bash
# 启动Web服务
go run cmd/server/main.go

# 启动判题服务
go run cmd/judger/main.go

# 启动异步任务处理器
go run cmd/worker/main.go
```

## 📁 项目结构

```
zhku-oj/
├── cmd/                     # 应用程序入口
│   ├── server/             # Web服务器
│   ├── judger/             # 判题服务
│   └── worker/             # 异步任务处理器
├── internal/               # 内部包
│   ├── config/             # 配置管理
│   ├── middleware/         # HTTP中间件
│   ├── handler/            # HTTP处理器
│   ├── service/            # 业务逻辑层
│   ├── repository/         # 数据访问层
│   ├── model/              # 数据模型
│   ├── judge/              # 判题相关
│   ├── queue/              # 消息队列
│   └── pkg/                # 工具包
├── configs/                # 配置文件
├── deployments/            # 部署配置
├── docs/                   # 文档
└── tests/                  # 测试
```

## 🔄 核心业务流程

### 判题流程
1. **代码提交阶段** (100-200ms): 参数校验→权限验证→创建记录→任务入队→立即响应
2. **异步判题阶段** (2-30秒): 任务消费→编译Java→运行测试用例→结果计算→资源清理
3. **结果处理阶段** (50-100ms): 结果消费→数据更新→WebSocket通知
4. **异步统计更新**: 用户统计→题目统计→排行榜更新

### go-judge集成
- **编译阶段**: 使用`POST /run`编译Java代码，获取.class文件ID缓存
- **运行阶段**: 使用缓存文件ID逐个执行测试用例
- **资源清理**: 及时调用`DELETE /file/{fileId}`清理缓存文件

## 🛠️ API接口

### 认证相关
- `POST /api/v1/auth/register` - 用户注册
- `POST /api/v1/auth/login` - 用户登录
- `POST /api/v1/auth/logout` - 用户登出

### 题目管理
- `GET /api/v1/problems` - 获取题目列表
- `GET /api/v1/problems/{id}` - 获取题目详情
- `POST /api/v1/problems` - 创建题目(教师权限)

### 代码提交
- `POST /api/v1/submissions` - 提交代码
- `GET /api/v1/submissions/{id}` - 获取提交详情
- `GET /api/v1/submissions` - 获取提交列表

### 用户管理
- `GET /api/v1/users/profile` - 获取用户信息
- `PUT /api/v1/users/profile` - 更新用户信息
- `PUT /api/v1/users/password` - 修改密码

## 🐳 部署配置

### Docker部署
```bash
# 开发环境一键部署
cd deployments/docker
docker-compose up -d

# 生产环境部署
docker-compose -f docker-compose.prod.yml up -d
```

### 健康检查
- Web服务: `GET /health`
- go-judge状态: `GET http://localhost:5050/version`

## 📊 监控和日志

### 日志配置
- 应用日志: 使用Logrus，支持JSON格式
- 日志轮转: 使用lumberjack，自动清理旧日志
- 日志级别: debug, info, warn, error

### 监控指标
- 应用健康状态
- go-judge实例状态
- 判题队列深度
- 数据库连接状态
- 系统资源使用情况

## 🔧 开发指南

### 代码规范
- 遵循Go语言官方代码规范
- 使用统一的错误处理机制
- 所有接口函数必须添加完整的中文注释
- 严格按照项目目录结构组织代码

### 判题模块开发要点
- 沙箱负载均衡: 多个go-judge实例的负载分发
- 文件缓存管理: .class文件的生命周期管理
- 错误重试机制: 沙箱异常时的重试策略
- 资源限制: CPU时间、内存、进程数的严格控制

### 测试
```bash
# 运行单元测试
go test ./tests/unit/...

# 运行集成测试
go test ./tests/integration/...

# 代码覆盖率
go test -cover ./...
```

## 🔒 安全注意事项

- JWT Token安全管理
- 输入验证和XSS防护
- go-judge沙箱安全隔离
- 敏感信息加密存储
- API访问频率限制

## 📈 性能要求

- 支持200-500用户并发
- 判题响应时间 < 30秒
- 接口响应时间 < 200ms
- 支持20-50提交/分钟处理

## 🤝 贡献指南

1. Fork本项目
2. 创建功能分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 创建Pull Request

## 📄 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情

## 📞 联系方式

- 项目维护者: [Your Name]
- 邮箱: [your.email@example.com]
- 问题反馈: [GitHub Issues](https://github.com/your-repo/zhku-oj/issues)

## 🙏 致谢

感谢以下开源项目的支持：
- [go-judge](https://github.com/criyle/go-judge) - 代码执行沙箱
- [Gin](https://github.com/gin-gonic/gin) - Web框架
- [MongoDB](https://www.mongodb.com/) - 数据库
- [RabbitMQ](https://www.rabbitmq.com/) - 消息队列