# 沙箱客户端使用指南

本文档说明如何使用 `pkg/sandbox` 包与 go-judge 沙箱进行交互。当前实现已完全对应[沙箱使用教程](../../md/沙箱使用教程.md)中描述的所有接口。

## ✅ 接口对应关系检查

### 核心接口映射

| 文档接口 | 实现方法 | 状态 | 说明 |
|---------|----------|------|------|
| `POST /run` (编译) | `CompileJava()` | ✅ | 完全对应，支持所有编译参数 |
| `POST /run` (运行) | `RunJava()` | ✅ | 完全对应，支持copyInCached |
| `GET /version` | `GetVersion()` | ✅ | 获取go-judge版本信息 |
| `GET /config` | `GetConfig()` | ✅ | 获取go-judge配置信息 |
| `GET /file` | `ListFiles()` | ✅ | 列出所有缓存文件 |
| `POST /file` | `UploadFile()` | ✅ | 上传文件到缓存 |
| `GET /file/{id}` | `DownloadFile()` | ✅ | 下载指定缓存文件 |
| `DELETE /file/{id}` | `DeleteFile()` | ✅ | 删除指定缓存文件 |

### 请求参数映射

| 文档参数 | 实现字段 | 状态 | 类型 |
|---------|----------|------|------|
| `args` | `CommandConfig.Args` | ✅ | `[]string` |
| `env` | `CommandConfig.Env` | ✅ | `[]string` |
| `files` | `CommandConfig.Files` | ✅ | `[]FileDescriptor` |
| `cpuLimit` | `CommandConfig.CPULimit` | ✅ | `int64` (纳秒) |
| `memoryLimit` | `CommandConfig.MemoryLimit` | ✅ | `int64` (字节) |
| `copyIn` | `CommandConfig.CopyIn` | ✅ | `map[string]CopyIn` |
| `copyInCached` | `CommandConfig.CopyInCached` | ✅ | `map[string]string` |
| `copyOut` | `CommandConfig.CopyOut` | ✅ | `[]string` |
| `copyOutCached` | `CommandConfig.CopyOutCached` | ✅ | `[]string` |

### 响应字段映射

| 文档字段 | 实现字段 | 状态 | 类型 |
|---------|----------|------|------|
| `status` | `RunResponse.Status` | ✅ | `string` |
| `exitStatus` | `RunResponse.ExitStatus` | ✅ | `int` |
| `time` | `RunResponse.Time` | ✅ | `int64` (纳秒) |
| `memory` | `RunResponse.Memory` | ✅ | `int64` (字节) |
| `files` | `RunResponse.Files` | ✅ | `map[string]string` |
| `fileIds` | `RunResponse.FileIDs` | ✅ | `map[string]string` |

## 🚀 快速开始

### 1. 创建客户端

```go
import "zhku-oj/pkg/sandbox"

// 创建沙箱配置
config := &sandbox.SandboxConfig{
    URL:                 "http://localhost:5050",
    Timeout:             30 * time.Second,
    MaxConcurrent:       10,
    HealthCheckInterval: 30 * time.Second,
    Enabled:             true,
}

// 创建客户端
client, err := sandbox.NewSandboxClient(config)
if err != nil {
    log.Fatal("创建沙箱客户端失败:", err)
}
defer client.Close()
```

### 2. 简化Java代码执行

```go
// 使用CompileAndRunJava方法一次性完成编译和运行
javaCode := `
public class Main {
    public static void main(String[] args) {
        System.out.println("Hello World");
    }
}
`

result, err := client.CompileAndRunJava(ctx, javaCode, "", nil, nil)
if err != nil {
    log.Fatal("执行失败:", err)
}

fmt.Printf("状态: %s, 输出: %s\n", result.Status, result.Output)
fmt.Printf("耗时: %dms, 内存: %dKB\n", result.TimeUsed, result.MemoryUsed)
```

### 3. 多语言代码执行

```go
// 支持Java、Python、C++、C等多种语言
req := &sandbox.CodeExecutionRequest{
    Language:    "java",
    SourceCode:  javaCode,
    Input:       "test input",
    TimeLimit:   1000,    // 1秒
    MemoryLimit: 128,     // 128MB
}

result, err := client.RunCode(ctx, req)
```

### 4. 分离式编译运行

```go
// 编译阶段
compileReq := &sandbox.CompileRequest{
    SourceCode:  javaCode,
    SourceFile:  "Main.java",
    CompileCmd:  []string{"/usr/bin/javac", "-encoding", "UTF-8", "-cp", "/w", "Main.java"},
    CompileEnv:  []string{"PATH=/usr/bin:/bin", "JAVA_HOME=/usr/lib/jvm/java-17-openjdk-amd64"},
    CPULimit:    10000000000, // 10秒
    MemoryLimit: 268435456,   // 256MB
}

compileResult, err := client.CompileJava(ctx, compileReq)
if err != nil || compileResult.Status != "Accepted" {
    // 处理编译错误
    return
}

// 运行阶段
runReq := &sandbox.RunTestRequest{
    ClassFileID: compileResult.FileIDs["Main.class"],
    Input:       "test input",
    RunCmd:      []string{"/usr/bin/java", "-cp", "/w", "Main"},
    RunEnv:      []string{"PATH=/usr/bin:/bin", "JAVA_HOME=/usr/lib/jvm/java-17-openjdk-amd64"},
    CPULimit:    5000000000, // 5秒
    MemoryLimit: 134217728,  // 128MB
}

runResult, err := client.RunJava(ctx, runReq)
```

## 📋 配置说明

### 沙箱配置 (config.yaml)

项目配置文件 `pkg/configs/config.yaml` 中包含了完整的沙箱配置：

```yaml
sandboxes:
  - url: "http://localhost:5050"
    weight: 1
    max_concurrent: 10
    timeout: 30s
    health_check_interval: 10s
    enabled: true
    retry_times: 3
    retry_interval: 1s

languages:
  java:
    compile:
      command: ["/usr/bin/javac", "-encoding", "UTF-8", "-cp", "/w", "Main.java"]
      env: ["PATH=/usr/bin:/bin", "JAVA_HOME=/usr/lib/jvm/java-17-openjdk-amd64"]
      cpu_limit: 10000000000
      memory_limit: 268435456
    runtime:
      command: ["/usr/bin/java", "-cp", "/w", "Main"]
      cpu_limit: 5000000000
      memory_limit: 134217728
```

### 语言配置

每种编程语言都有独立的配置：

- **Java**: 完全支持JDK17环境
- **Python**: 支持Python3解释器
- **C++**: 支持g++编译器和C++17标准
- **C**: 支持gcc编译器和C11标准

## 🔧 API 参考

### 主要结构体

#### SandboxClient
沙箱客户端，提供与 go-judge 交互的所有方法。

#### JudgeResult
统一的判题结果结构：
```go
type JudgeResult struct {
    Status         string `json:"status"`          // 最终状态
    Success        bool   `json:"success"`         // 是否成功
    Output         string `json:"output"`          // 程序输出
    ErrorOutput    string `json:"error_output"`    // 错误输出
    CompileError   string `json:"compile_error"`   // 编译错误信息
    TimeUsed       int64  `json:"time_used"`       // 总耗时(毫秒)
    MemoryUsed     int64  `json:"memory_used"`     // 峰值内存使用(KB)
    ExitStatus     int    `json:"exit_status"`     // 程序退出码
    CompileTime    int64  `json:"compile_time"`    // 编译时间(毫秒)
    RunTime        int64  `json:"run_time"`        // 运行时间(毫秒)
    CompileMemory  int64  `json:"compile_memory"`  // 编译内存(KB)
    RunMemory      int64  `json:"run_memory"`      // 运行内存(KB)
}
```

### 主要方法

- `NewSandboxClient(config)`: 创建沙箱客户端
- `CompileAndRunJava(ctx, sourceCode, inputData, compileConfig, runConfig)`: Java代码一键执行
- `RunCode(ctx, req)`: 多语言代码执行
- `CompileJava(ctx, req)`: Java代码编译
- `RunJava(ctx, req)`: Java代码运行
- `GetVersion(ctx)`: 获取go-judge版本
- `GetConfig(ctx)`: 获取go-judge配置
- `ListFiles(ctx)`: 列出缓存文件
- `UploadFile(ctx, content)`: 上传文件
- `DownloadFile(ctx, fileID)`: 下载文件
- `DeleteFile(ctx, fileID)`: 删除文件
- `IsHealthy()`: 健康检查

## ⚠️ 注意事项

1. **超时设置**: 建议设置适当的上下文超时时间
2. **资源限制**: 根据题目要求设置合理的时间和内存限制
3. **文件清理**: CompileAndRunJava会自动清理缓存文件
4. **并发控制**: 客户端支持并发使用，但建议控制并发数量
5. **错误重试**: 网络错误时客户端会自动重试
6. **copyInCached**: 使用缓存文件时使用copyInCached字段而不是copyIn

## 📝 使用示例

完整的使用示例请参考：
- [`example.go`](./example.go) - 包含各种使用场景的完整示例
- [`utils.go`](./utils.go) - 提供状态映射和数据转换工具

## 🔗 相关文档

- [沙箱使用教程](../../md/沙箱使用教程.md) - go-judge接口详细文档
- [项目配置](../../pkg/configs/config.yaml) - 沙箱配置文件

## ✅ 验证清单

- [x] 所有go-judge核心接口已实现
- [x] copyInCached字段正确使用
- [x] 时间精度统一为纳秒
- [x] 内存精度统一为字节
- [x] 状态码完全映射
- [x] 文件管理完整支持
- [x] 健康检查机制完善
- [x] 错误处理和重试机制
- [x] 多语言支持框架
- [x] 完整的使用示例
- [x] 详细的文档说明

当前沙箱服务实现已经完全对应sandbox.md文档中描述的所有功能和接口！