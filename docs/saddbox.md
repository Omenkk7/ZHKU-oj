# go-judge 沙箱 Go 项目集成指南

## 📋 概述

本指南将详细说明如何在 Go 项目中集成和使用 go-judge 沙箱服务。根据项目架构信息，go-judge 提供了多种集成方式，包括直接导入、HTTP 客户端调用和 gRPC 调用。

## 🔧 集成方式选择

### 1. HTTP REST API 调用（推荐）
- **适用场景**：独立部署的沙箱服务
- **优势**：服务隔离、易于扩展、语言无关
- **部署方式**：Docker 容器或二进制文件

### 2. 直接导入模块
- **适用场景**：单体应用、性能要求极高
- **优势**：无网络开销、直接函数调用
- **限制**：仅支持 Linux 环境

### 3. gRPC 调用
- **适用场景**：高性能、大批量处理
- **优势**：二进制协议、高性能、类型安全

## 🚀 方式一：HTTP REST API 集成（推荐）

### 1.1 启动沙箱服务

#### Docker 启动（推荐）
```bash
# 基础启动
docker run -it --rm --privileged --shm-size=256m \
  -p 5050:5050 --name=go-judge \
  criyle/go-judge

# 带 JDK 17 的版本
docker run -it --rm --privileged --shm-size=256m \
  -p 5050:5050 --name=go-judge-java \
  go-judge-java-wsl

# 生产环境启动（带配置）
docker run -it --rm --privileged --shm-size=256m \
  -p 5050:5050 \
  -v $(pwd)/mount.yaml:/mount.yaml \
  -v $(pwd)/seccomp.yaml:/seccomp.yaml \
  --name=go-judge \
  criyle/go-judge \
  -mount-conf=/mount.yaml \
  -seccomp-conf=/seccomp.yaml \
  -parallelism=4
```

#### 二进制文件启动
```bash
# 下载二进制文件
wget https://github.com/criyle/go-judge/releases/latest/download/go-judge-linux-amd64

# 设置权限
chmod +x go-judge-linux-amd64

# 启动服务
./go-judge-linux-amd64 -http-addr=:5050 -parallelism=4
```

### 1.2 Go 客户端实现

#### 基础客户端结构
```go
package judge

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "time"
)

// JudgeClient Go-Judge 沙箱客户端
type JudgeClient struct {
    BaseURL    string
    HTTPClient *http.Client
}

// NewJudgeClient 创建新的判题客户端
func NewJudgeClient(baseURL string) *JudgeClient {
    return &JudgeClient{
        BaseURL: baseURL,
        HTTPClient: &http.Client{
            Timeout: 30 * time.Second,
            Transport: &http.Transport{
                MaxIdleConns:        100,
                MaxIdleConnsPerHost: 10,
                IdleConnTimeout:     90 * time.Second,
            },
        },
    }
}

// RunRequest 执行请求结构
type RunRequest struct {
    Cmd []Command `json:"cmd"`
}

// Command 单个命令配置
type Command struct {
    Args         []string          `json:"args"`
    Env          []string          `json:"env,omitempty"`
    Files        []File            `json:"files"`
    CPULimit     int64             `json:"cpuLimit"`
    MemoryLimit  int64             `json:"memoryLimit"`
    ProcLimit    int               `json:"procLimit"`
    CopyIn       map[string]Input  `json:"copyIn,omitempty"`
    CopyOut      []string          `json:"copyOut,omitempty"`
    CopyOutCached []string         `json:"copyOutCached,omitempty"`
    CopyInCached map[string]string `json:"copyInCached,omitempty"`
    CopyOutMax   int64             `json:"copyOutMax,omitempty"`
}

// File 文件描述符配置
type File struct {
    Content string `json:"content,omitempty"`
    Name    string `json:"name,omitempty"`
    Max     int64  `json:"max,omitempty"`
}

// Input 输入文件配置
type Input struct {
    Content string `json:"content"`
}

// RunResponse 执行响应结构
type RunResponse struct {
    Status     string            `json:"status"`
    ExitStatus int               `json:"exitStatus"`
    Time       int64             `json:"time"`
    Memory     int64             `json:"memory"`
    RunTime    int64             `json:"runTime"`
    Files      map[string]string `json:"files"`
    FileIds    map[string]string `json:"fileIds,omitempty"`
}
```

#### 核心执行方法
```go
// RunCode 执行代码
func (c *JudgeClient) RunCode(ctx context.Context, req *RunRequest) ([]RunResponse, error) {
    reqBody, err := json.Marshal(req)
    if err != nil {
        return nil, fmt.Errorf("marshal request: %w", err)
    }

    httpReq, err := http.NewRequestWithContext(ctx, "POST", c.BaseURL+"/run", bytes.NewReader(reqBody))
    if err != nil {
        return nil, fmt.Errorf("create request: %w", err)
    }
    httpReq.Header.Set("Content-Type", "application/json")

    resp, err := c.HTTPClient.Do(httpReq)
    if err != nil {
        return nil, fmt.Errorf("execute request: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        body, _ := io.ReadAll(resp.Body)
        return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
    }

    var results []RunResponse
    if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
        return nil, fmt.Errorf("decode response: %w", err)
    }

    return results, nil
}

// CompileAndRun 编译并运行 Java 代码
func (c *JudgeClient) CompileAndRun(ctx context.Context, sourceCode, inputData string, timeLimit, memoryLimit int64) (*JudgeResult, error) {
    // 步骤1: 编译
    compileReq := &RunRequest{
        Cmd: []Command{{
            Args: []string{"/usr/bin/javac", "-encoding", "UTF-8", "-cp", "/w", "Main.java"},
            Env:  []string{"PATH=/usr/bin:/bin"},
            Files: []File{
                {Content: ""},
                {Name: "stdout", Max: 10240},
                {Name: "stderr", Max: 10240},
            },
            CPULimit:    10000000000, // 10秒编译时间
            MemoryLimit: 134217728,   // 128MB
            ProcLimit:   50,
            CopyIn: map[string]Input{
                "Main.java": {Content: sourceCode},
            },
            CopyOut:       []string{"stdout", "stderr"},
            CopyOutCached: []string{"Main.class"},
            CopyOutMax:    10240,
        }},
    }

    compileResults, err := c.RunCode(ctx, compileReq)
    if err != nil {
        return nil, fmt.Errorf("compile request failed: %w", err)
    }

    compileResult := compileResults[0]
    if compileResult.Status != "Accepted" {
        return &JudgeResult{
            Status:       "Compile Error",
            CompileError: compileResult.Files["stderr"],
            Success:      false,
        }, nil
    }

    // 获取编译后的 class 文件 ID
    classFileId := compileResult.FileIds["Main.class"]

    // 步骤2: 运行
    runReq := &RunRequest{
        Cmd: []Command{{
            Args: []string{"/usr/bin/java", "-cp", "/w", "Main"},
            Env:  []string{"PATH=/usr/bin:/bin"},
            Files: []File{
                {Content: inputData},
                {Name: "stdout", Max: 10240},
                {Name: "stderr", Max: 10240},
            },
            CPULimit:    timeLimit * 1000000,   // 转换为纳秒
            MemoryLimit: memoryLimit * 1024 * 1024, // 转换为字节
            ProcLimit:   50,
            CopyInCached: map[string]string{
                "Main.class": classFileId,
            },
            CopyOut:    []string{"stdout", "stderr"},
            CopyOutMax: 10240,
        }},
    }

    runResults, err := c.RunCode(ctx, runReq)
    if err != nil {
        return nil, fmt.Errorf("run request failed: %w", err)
    }

    runResult := runResults[0]
    return &JudgeResult{
        Status:      runResult.Status,
        Output:      runResult.Files["stdout"],
        Error:       runResult.Files["stderr"],
        TimeUsed:    runResult.Time / 1000000,    // 转换为毫秒
        MemoryUsed:  runResult.Memory / 1024,     // 转换为KB
        ExitCode:    runResult.ExitStatus,
        Success:     runResult.Status == "Accepted",
    }, nil
}

// JudgeResult 判题结果
type JudgeResult struct {
    Status       string `json:"status"`
    Output       string `json:"output"`
    Error        string `json:"error"`
    CompileError string `json:"compile_error,omitempty"`
    TimeUsed     int64  `json:"time_used"`     // 毫秒
    MemoryUsed   int64  `json:"memory_used"`   // KB
    ExitCode     int    `json:"exit_code"`
    Success      bool   `json:"success"`
}
```

#### 使用示例
```go
package main

import (
    "context"
    "fmt"
    "log"
    "strings"
)

func main() {
    // 创建客户端
    client := NewJudgeClient("http://localhost:5050")

    // Java 代码示例
    javaCode := `
import java.util.Scanner;

public class Main {
    public static void main(String[] args) {
        Scanner sc = new Scanner(System.in);
        int a = sc.nextInt();
        int b = sc.nextInt();
        System.out.println(a + b);
        sc.close();
    }
}
`

    // 测试用例
    inputData := "1 2\n"
    expectedOutput := "3"

    // 执行判题
    ctx := context.Background()
    result, err := client.CompileAndRun(ctx, javaCode, inputData, 1000, 128)
    if err != nil {
        log.Fatalf("Judge failed: %v", err)
    }

    // 处理结果
    if result.Success {
        actualOutput := strings.TrimSpace(result.Output)
        if actualOutput == expectedOutput {
            fmt.Println("✅ 测试通过")
            fmt.Printf("输出: %s\n", actualOutput)
            fmt.Printf("时间: %dms\n", result.TimeUsed)
            fmt.Printf("内存: %dKB\n", result.MemoryUsed)
        } else {
            fmt.Println("❌ 输出不匹配")
            fmt.Printf("期望: %s\n", expectedOutput)
            fmt.Printf("实际: %s\n", actualOutput)
        }
    } else {
        fmt.Printf("❌ 执行失败: %s\n", result.Status)
        if result.CompileError != "" {
            fmt.Printf("编译错误: %s\n", result.CompileError)
        }
        if result.Error != "" {
            fmt.Printf("运行错误: %s\n", result.Error)
        }
    }
}
```

## 🔧 方式二：直接模块导入集成

根据项目源码分析，go-judge 支持直接导入使用，但需要正确配置环境和依赖。

### 2.1 项目依赖配置

#### go.mod 配置
```go
module yourproject

go 1.21

require (
    github.com/criyle/go-judge v1.8.8
    github.com/criyle/go-sandbox v0.10.3
    go.uber.org/zap v1.26.0
)
```

### 2.2 环境初始化

根据源码 `cmd/go-judge/main.go` 的初始化流程：

```go
package main

import (
    "context"
    "log"
    "os"
    "runtime"
    "time"

    "github.com/criyle/go-judge/env"
    "github.com/criyle/go-judge/env/pool"
    "github.com/criyle/go-judge/filestore"
    "github.com/criyle/go-judge/worker"
    "github.com/criyle/go-sandbox/container"
    "go.uber.org/zap"
)

// SandboxService 沙箱服务
type SandboxService struct {
    envPool   worker.EnvironmentPool
    fileStore filestore.FileStore
    worker    *worker.Worker
    logger    *zap.Logger
}

// NewSandboxService 创建沙箱服务
func NewSandboxService() (*SandboxService, error) {
    // 初始化日志
    logger, err := zap.NewProduction()
    if err != nil {
        return nil, err
    }

    // Linux 环境初始化
    if runtime.GOOS == "linux" {
        container.Init()
    }

    // 环境配置
    envConfig := env.Config{
        ContainerInitPath:  "/usr/local/bin/go-judge-init",
        MountConf:          "mount.yaml",
        TmpFsParam:         "size=128m,nr_inodes=4k",
        NetShare:           false,
        CgroupPrefix:       "gojudge",
        ContainerCredStart: 10000,
        EnableCPURate:      false,
        CPUCfsPeriod:       100 * time.Millisecond,
    }

    // 创建环境构建器
    builder, _, err := env.NewBuilder(envConfig, logger)
    if err != nil {
        return nil, err
    }

    // 创建环境池
    envPool := pool.NewPool(builder)

    // 创建文件存储
    fileStoreDir := "/dev/shm/go-judge"
    if runtime.GOOS != "linux" {
        fileStoreDir = os.TempDir() + "/go-judge"
    }
    os.MkdirAll(fileStoreDir, 0755)
    fileStore := filestore.NewFileLocalStore(fileStoreDir)

    // 创建工作器
    workerInstance := worker.New(worker.Config{
        FileStore:              fileStore,
        EnvironmentPool:        envPool,
        Parallelism:            runtime.NumCPU(),
        WorkDir:                "/w",
        TimeLimitCheckInterval: time.Millisecond,
        OutputLimit:            64 << 20, // 64MB
        CopyOutLimit:           64 << 20, // 64MB
        OpenFileLimit:          256,
        ExtraMemoryLimit:       16 << 20, // 16MB
    })

    return &SandboxService{
        envPool:   envPool,
        fileStore: fileStore,
        worker:    workerInstance,
        logger:    logger,
    }, nil
}

// ExecuteCode 执行代码
func (s *SandboxService) ExecuteCode(ctx context.Context, req *worker.Request) (*worker.Response, error) {
    responseChan := s.worker.Submit(ctx, req)
    response := <-responseChan
    return response, nil
}

// Close 关闭服务
func (s *SandboxService) Close() {
    s.envPool.Destroy()
}
```

### 2.3 使用示例

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/criyle/go-judge/worker"
)

func main() {
    // 创建沙箱服务
    service, err := NewSandboxService()
    if err != nil {
        log.Fatalf("Failed to create sandbox service: %v", err)
    }
    defer service.Close()

    // 构建执行请求
    javaCode := `
public class Main {
    public static void main(String[] args) {
        System.out.println("Hello from go-judge!");
    }
}
`

    // 编译请求
    compileRequest := &worker.Request{
        RequestID: "compile-001",
        Cmd: []worker.Cmd{
            {
                Args: []string{"/usr/bin/javac", "-cp", "/w", "Main.java"},
                Env:  []string{"PATH=/usr/bin:/bin"},
                Files: []worker.CmdFile{
                    {Content: ""},
                    {Name: "stdout", Max: 10240},
                    {Name: "stderr", Max: 10240},
                },
                CPULimit:    10 * time.Second,
                MemoryLimit: 128 << 20, // 128MB
                ProcLimit:   50,
                CopyIn: map[string]worker.CmdCopyInFile{
                    "Main.java": {Content: javaCode},
                },
                CopyOut:       []string{"stdout", "stderr"},
                CopyOutCached: []string{"Main.class"},
            },
        },
    }

    // 执行编译
    ctx := context.Background()
    compileResponse, err := service.ExecuteCode(ctx, compileRequest)
    if err != nil {
        log.Fatalf("Compile failed: %v", err)
    }

    if compileResponse.Results[0].Status != worker.StatusAccepted {
        log.Fatalf("Compile error: %s", compileResponse.Results[0].Files["stderr"])
    }

    fmt.Println("编译成功！")
    fmt.Printf("编译结果: %+v\n", compileResponse.Results[0])
}
```

## 🔗 方式三：gRPC 集成

### 3.1 gRPC 客户端

根据项目中的 gRPC 支持：

```go
package judge

import (
    "context"
    "time"

    "github.com/criyle/go-judge/pb"
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
)

// GRPCJudgeClient gRPC 判题客户端
type GRPCJudgeClient struct {
    conn   *grpc.ClientConn
    client pb.ExecutorClient
}

// NewGRPCJudgeClient 创建 gRPC 客户端
func NewGRPCJudgeClient(addr string) (*GRPCJudgeClient, error) {
    conn, err := grpc.Dial(addr, 
        grpc.WithTransportCredentials(insecure.NewCredentials()),
        grpc.WithTimeout(30*time.Second),
    )
    if err != nil {
        return nil, err
    }

    return &GRPCJudgeClient{
        conn:   conn,
        client: pb.NewExecutorClient(conn),
    }, nil
}

// Close 关闭连接
func (c *GRPCJudgeClient) Close() error {
    return c.conn.Close()
}

// Execute 执行代码
func (c *GRPCJudgeClient) Execute(ctx context.Context, req *pb.Request) (*pb.Response, error) {
    return c.client.Exec(ctx, req)
}
```

### 3.2 启动带 gRPC 的沙箱服务

```bash
# 启动支持 gRPC 的服务
docker run -it --rm --privileged --shm-size=256m \
  -p 5050:5050 -p 5051:5051 \
  --name=go-judge \
  criyle/go-judge \
  -grpc-addr=:5051
```

## 📊 配置文件示例

### Mount 配置 (mount.yaml)

```yaml
version: 1
mounts:
  - source: tmpfs
    target: /tmp
    type: tmpfs
    data: size=128m,nr_inodes=4k
  - source: /usr/lib/jvm/java-17-openjdk-amd64
    target: /usr/lib/jvm/java-17-openjdk-amd64
    type: bind
    options: [ro, rbind]
  - source: /usr/bin/java
    target: /usr/bin/java
    type: bind
    options: [ro]
  - source: /usr/bin/javac
    target: /usr/bin/javac
    type: bind
    options: [ro]
workDir: /w
hostName: go-judge
domainName: go-judge
uid: 1000
gid: 1000
```

### Seccomp 配置 (seccomp.yaml)

```yaml
version: 1
defaultAction: allow
syscalls:
  - names: [ptrace, mount, umount, reboot, setsid]
    action: errno
    errno: 1
  - names: [socket, connect, bind, listen, accept]
    action: errno
    errno: 1
  - names: [execve]
    action: allow
```

## 🛡️ 安全机制集成

根据项目记忆中的安全机制，集成时需要注意：

### 1. Namespace 隔离
```go
// 在环境配置中启用隔离
envConfig := env.Config{
    NetShare: false, // 网络隔离
    // 其他配置...
}
```

### 2. Cgroup 资源控制
```go
// 资源限制配置
cmd := worker.Cmd{
    CPULimit:    time.Second,     // CPU 时间限制
    MemoryLimit: 128 << 20,       // 内存限制 128MB
    ProcLimit:   50,              // 进程数限制
}
```

### 3. Seccomp 系统调用过滤
```go
// 启用 Seccomp 过滤器
envConfig := env.Config{
    SeccompConf: "seccomp.yaml",
}
```

## 📋 最佳实践

### 1. 资源配置
```go
// Java 执行的推荐配置
type JavaJudgeConfig struct {
    CompileTimeLimit time.Duration // 编译时间限制：10秒
    CompileMemoryLimit int64       // 编译内存限制：256MB
    RunTimeLimit time.Duration     // 运行时间限制：1-5秒
    RunMemoryLimit int64           // 运行内存限制：128-512MB
    ProcLimit int                  // 进程限制：50
    OutputLimit int64              // 输出限制：10MB
}

func DefaultJavaConfig() *JavaJudgeConfig {
    return &JavaJudgeConfig{
        CompileTimeLimit:   10 * time.Second,
        CompileMemoryLimit: 256 << 20, // 256MB
        RunTimeLimit:       2 * time.Second,
        RunMemoryLimit:     128 << 20, // 128MB
        ProcLimit:          50,
        OutputLimit:        10 << 20, // 10MB
    }
}
```

### 2. 错误处理
```go
// 完整的错误处理
func handleJudgeResult(result *worker.Result) *JudgeStatus {
    switch result.Status {
    case worker.StatusAccepted:
        return &JudgeStatus{Verdict: "Accepted"}
    case worker.StatusMemoryLimitExceeded:
        return &JudgeStatus{Verdict: "Memory Limit Exceeded"}
    case worker.StatusTimeLimitExceeded:
        return &JudgeStatus{Verdict: "Time Limit Exceeded"}
    case worker.StatusRuntimeError:
        return &JudgeStatus{
            Verdict: "Runtime Error",
            Error:   result.Files["stderr"],
        }
    default:
        return &JudgeStatus{
            Verdict: "System Error",
            Error:   fmt.Sprintf("Unknown status: %s", result.Status),
        }
    }
}
```

## 📋 总结

根据 go-judge 项目架构，推荐的集成方式：

1. **HTTP REST API（推荐）**：
    - 服务独立部署，易于维护
    - 支持多语言客户端
    - 适合大多数应用场景

2. **直接导入模块**：
    - 性能最优，无网络开销
    - 仅适用于 Linux 环境
    - 需要正确配置 cgroup 权限

3. **gRPC 调用**：
    - 高性能二进制协议
    - 适合高并发场景
    - 类型安全的接口定义

每种方式都完整支持 go-judge 的安全机制，包括 Namespace 隔离、Cgroup 资源控制和 Seccomp 系统调用过滤，可以根据具体需求选择合适的集成方案。