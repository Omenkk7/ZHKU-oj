package sandbox

import (
	"context"
	"fmt"
	"strings"
	"time"
)

// ExampleUsage 沙箱客户端使用示例
// 演示如何使用SandboxClient进行Java代码的编译和运行
func ExampleUsage() {
	// 1. 创建沙箱配置
	config := &SandboxConfig{
		URL:                 "http://localhost:5050",
		Weight:              1,
		MaxConcurrent:       10,
		Timeout:             30 * time.Second,
		HealthCheckInterval: 10 * time.Second,
		Enabled:             true,
		RetryTimes:          3,
		RetryInterval:       1 * time.Second,
	}

	// 2. 创建沙箱客户端
	client, err := NewSandboxClient(config)
	if err != nil {
		fmt.Printf("创建沙箱客户端失败: %v\n", err)
		return
	}
	defer client.Close()

	// 3. 检查沙箱健康状态
	if !client.IsHealthy() {
		fmt.Println("沙箱服务不健康")
		return
	}

	// 4. 准备Java代码
	javaCode := `
import java.util.*;

public class Main {
    public static void main(String[] args) {
        Scanner scanner = new Scanner(System.in);
        int a = scanner.nextInt();
        int b = scanner.nextInt();
        System.out.println(a + b);
        scanner.close();
    }
}
`

	// 5. 构建编译请求
	compileReq := &CompileRequest{
		SourceCode:  javaCode,
		SourceFile:  "Main.java",
		CompileCmd:  []string{"/usr/bin/javac", "-encoding", "UTF-8", "-cp", "/w", "Main.java"},
		CompileEnv:  []string{"PATH=/usr/bin:/bin", "JAVA_HOME=/usr/lib/jvm/java-17-openjdk-amd64", "CLASSPATH=.", "LANG=C.UTF-8"},
		CPULimit:    10000000000, // 10秒
		MemoryLimit: 268435456,   // 256MB
		StackLimit:  134217728,   // 128MB
		ProcLimit:   50,
		OutputLimit: 10240, // 10KB
	}

	// 6. 执行编译
	ctx := context.Background()
	compileResult, err := client.CompileJava(ctx, compileReq)
	if err != nil {
		fmt.Printf("编译失败: %v\n", err)
		return
	}

	fmt.Printf("编译状态: %s\n", compileResult.Status)
	fmt.Printf("编译时间: %d ms\n", ConvertTimeToMS(compileResult.Time))
	fmt.Printf("编译内存: %d KB\n", ConvertMemoryToKB(compileResult.Memory))

	// 7. 检查编译是否成功
	if compileResult.Status != "Accepted" {
		fmt.Printf("编译失败: %s\n", compileResult.Files["stderr"])
		return
	}

	// 8. 获取编译后的class文件ID
	classFileID, exists := compileResult.FileIDs["Main.class"]
	if !exists {
		fmt.Println("未找到编译后的class文件")
		return
	}

	fmt.Printf("Class文件ID: %s\n", classFileID)

	// 9. 准备测试用例
	testCases := []struct {
		Input    string
		Expected string
	}{
		{"1 2\n", "3"},
		{"10 20\n", "30"},
		{"100 200\n", "300"},
	}

	// 10. 执行测试用例
	for i, testCase := range testCases {
		fmt.Printf("\n执行测试用例 %d:\n", i+1)

		runReq := &RunTestRequest{
			ClassFileID: classFileID,
			Input:       testCase.Input,
			RunCmd:      []string{"/usr/bin/java", "-cp", "/w", "Main"},
			RunEnv:      []string{"PATH=/usr/bin:/bin", "JAVA_HOME=/usr/lib/jvm/java-17-openjdk-amd64", "CLASSPATH=.", "LANG=C.UTF-8"},
			CPULimit:    5000000000, // 5秒
			MemoryLimit: 134217728,  // 128MB
			StackLimit:  67108864,   // 64MB
			ProcLimit:   1,
			OutputLimit: 10240, // 10KB
		}

		runResult, err := client.RunJava(ctx, runReq)
		if err != nil {
			fmt.Printf("运行失败: %v\n", err)
			continue
		}

		fmt.Printf("输入: %s", testCase.Input)
		fmt.Printf("期望输出: %s\n", testCase.Expected)
		fmt.Printf("实际输出: %s\n", runResult.Files["stdout"])
		fmt.Printf("运行状态: %s\n", runResult.Status)
		fmt.Printf("运行时间: %d ms\n", ConvertTimeToMS(runResult.Time))
		fmt.Printf("运行内存: %d KB\n", ConvertMemoryToKB(runResult.Memory))

		// 判断结果
		actualOutput := strings.TrimSpace(runResult.Files["stdout"])
		if actualOutput == testCase.Expected {
			fmt.Println("✅ 测试通过")
		} else {
			fmt.Println("❌ 测试失败")
		}
	}

	// 11. 清理缓存文件
	err = client.DeleteFile(ctx, classFileID)
	if err != nil {
		fmt.Printf("删除缓存文件失败: %v\n", err)
	}

	fmt.Println("\n判题完成！")
}

// ExampleBatchJudge 批量判题示例
// 演示如何批量处理多个代码提交
func ExampleBatchJudge() {
	config := &SandboxConfig{
		URL:                 "http://localhost:5050",
		Weight:              1,
		MaxConcurrent:       10,
		Timeout:             30 * time.Second,
		HealthCheckInterval: 10 * time.Second,
		Enabled:             true,
		RetryTimes:          3,
		RetryInterval:       1 * time.Second,
	}

	client, err := NewSandboxClient(config)
	if err != nil {
		fmt.Printf("创建沙箱客户端失败: %v\n", err)
		return
	}
	defer client.Close()

	// 模拟多个代码提交
	submissions := []struct {
		ID   string
		Code string
	}{
		{
			ID: "sub001",
			Code: `
import java.util.*;
public class Main {
    public static void main(String[] args) {
        Scanner sc = new Scanner(System.in);
        int a = sc.nextInt();
        int b = sc.nextInt();
        System.out.println(a + b);
    }
}`,
		},
		{
			ID: "sub002",
			Code: `
import java.util.*;
public class Main {
    public static void main(String[] args) {
        Scanner sc = new Scanner(System.in);
        int a = sc.nextInt();
        int b = sc.nextInt();
        System.out.println(a * b);  // 错误的逻辑，应该是相加
    }
}`,
		},
	}

	testInput := "10 20\n"
	expectedOutput := "30"

	for _, submission := range submissions {
		fmt.Printf("\n处理提交: %s\n", submission.ID)

		// 编译配置
		compileConfig := JavaCompileConfig{
			Command:     []string{"/usr/bin/javac", "-encoding", "UTF-8", "-cp", "/w", "Main.java"},
			Env:         []string{"PATH=/usr/bin:/bin", "JAVA_HOME=/usr/lib/jvm/java-17-openjdk-amd64", "CLASSPATH=.", "LANG=C.UTF-8"},
			CPULimit:    10000000000,
			MemoryLimit: 268435456,
			StackLimit:  134217728,
			ProcLimit:   50,
			OutputLimit: 10,
		}

		// 运行配置
		runtimeConfig := JavaRuntimeConfig{
			Command:     []string{"/usr/bin/java", "-cp", "/w", "Main"},
			Env:         []string{"PATH=/usr/bin:/bin", "JAVA_HOME=/usr/lib/jvm/java-17-openjdk-amd64", "CLASSPATH=.", "LANG=C.UTF-8"},
			CPULimit:    5000000000,
			MemoryLimit: 134217728,
			StackLimit:  67108864,
			ProcLimit:   1,
			OutputLimit: 10,
		}

		// 构建请求
		compileReq := BuildJavaCompileRequest(submission.Code, compileConfig)

		// 执行编译
		ctx := context.Background()
		compileResult, err := client.CompileJava(ctx, compileReq)
		if err != nil {
			fmt.Printf("编译失败: %v\n", err)
			continue
		}

		// 检查编译结果
		if compileResult.Status != "Accepted" {
			fmt.Printf("编译错误: %s\n", compileResult.Files["stderr"])
			continue
		}

		classFileID := compileResult.FileIDs["Main.class"]

		// 构建运行请求
		runReq := BuildJavaRunRequest(classFileID, testInput, runtimeConfig)

		// 执行运行
		runResult, err := client.RunJava(ctx, runReq)
		if err != nil {
			fmt.Printf("运行失败: %v\n", err)
			client.DeleteFile(ctx, classFileID)
			continue
		}

		// 比较结果
		actualOutput := strings.TrimSpace(runResult.Files["stdout"])
		if actualOutput == expectedOutput {
			fmt.Printf("✅ 提交 %s: ACCEPTED\n", submission.ID)
		} else {
			fmt.Printf("❌ 提交 %s: WRONG_ANSWER (期望: %s, 实际: %s)\n", submission.ID, expectedOutput, actualOutput)
		}

		// 清理文件
		client.DeleteFile(ctx, classFileID)
	}
}

// ExampleErrorHandling 错误处理示例
// 演示如何处理各种错误情况
func ExampleErrorHandling() {
	config := &SandboxConfig{
		URL:                 "http://localhost:5050",
		Timeout:             30 * time.Second,
		MaxConcurrent:       10,
		HealthCheckInterval: 10 * time.Second,
		Enabled:             true,
	}

	client, err := NewSandboxClient(config)
	if err != nil {
		fmt.Printf("创建客户端失败: %v\n", err)
		return
	}
	defer client.Close()

	// 1. 编译错误示例
	fmt.Println("1. 编译错误示例:")
	badCode := `
public class Main {
    public static void main(String[] args) {
        System.out.println("Hello World"  // 缺少分号和右括号
    }
`

	compileReq := &CompileRequest{
		SourceCode:  badCode,
		SourceFile:  "Main.java",
		CompileCmd:  []string{"/usr/bin/javac", "-encoding", "UTF-8", "-cp", "/w", "Main.java"},
		CompileEnv:  []string{"PATH=/usr/bin:/bin", "JAVA_HOME=/usr/lib/jvm/java-17-openjdk-amd64"},
		CPULimit:    10000000000,
		MemoryLimit: 268435456,
		ProcLimit:   50,
		OutputLimit: 10240,
	}

	ctx := context.Background()
	compileResult, err := client.CompileJava(ctx, compileReq)
	if err != nil {
		fmt.Printf("编译请求失败: %v\n", err)
	} else {
		fmt.Printf("编译状态: %s\n", compileResult.Status)
		if compileResult.Status != "Accepted" {
			fmt.Printf("编译错误信息: %s\n", compileResult.Files["stderr"])
			// 使用状态映射
			ojStatus := MapStatus(compileResult.Status, compileResult.ExitStatus, true)
			fmt.Printf("OJ状态: %s\n", ojStatus)
			// 格式化错误信息
			userMsg := FormatErrorMessage(compileResult.Status, compileResult.Files["stderr"], compileResult.ExitStatus)
			fmt.Printf("用户提示: %s\n", userMsg)
		}
	}

	// 2. 运行时错误示例
	fmt.Println("\n2. 运行时错误示例:")
	runtimeErrorCode := `
public class Main {
    public static void main(String[] args) {
        int[] arr = new int[5];
        System.out.println(arr[10]); // 数组越界
    }
}
`

	compileReq.SourceCode = runtimeErrorCode
	compileResult, err = client.CompileJava(ctx, compileReq)
	if err == nil && compileResult.Status == "Accepted" {
		classFileID := compileResult.FileIDs["Main.class"]

		runReq := &RunTestRequest{
			ClassFileID: classFileID,
			Input:       "",
			RunCmd:      []string{"/usr/bin/java", "-cp", "/w", "Main"},
			RunEnv:      []string{"PATH=/usr/bin:/bin", "JAVA_HOME=/usr/lib/jvm/java-17-openjdk-amd64"},
			CPULimit:    5000000000,
			MemoryLimit: 134217728,
			ProcLimit:   1,
			OutputLimit: 10240,
		}

		runResult, err := client.RunJava(ctx, runReq)
		if err != nil {
			fmt.Printf("运行请求失败: %v\n", err)
		} else {
			fmt.Printf("运行状态: %s\n", runResult.Status)
			fmt.Printf("退出码: %d\n", runResult.ExitStatus)
			if runResult.Files["stderr"] != "" {
				fmt.Printf("错误输出: %s\n", runResult.Files["stderr"])
			}
			// 使用状态映射
			ojStatus := MapStatus(runResult.Status, runResult.ExitStatus, false)
			fmt.Printf("OJ状态: %s\n", ojStatus)
			// 格式化错误信息
			userMsg := FormatErrorMessage(runResult.Status, runResult.Files["stderr"], runResult.ExitStatus)
			fmt.Printf("用户提示: %s\n", userMsg)
		}

		client.DeleteFile(ctx, classFileID)
	}

	// 3. 获取版本信息
	fmt.Println("\n3. 获取沙箱版本信息:")
	version, err := client.GetVersion(ctx)
	if err != nil {
		fmt.Printf("获取版本失败: %v\n", err)
	} else {
		fmt.Printf("版本信息: %v\n", version)
	}
}

// ExampleCompileAndRunJava 简化Java代码执行示例
// 演示如何使用CompileAndRunJava方法一次性完成编译和运行
func ExampleCompileAndRunJava() {
	// 创建沙箱配置
	config := &SandboxConfig{
		URL:                 "http://localhost:5050",
		Timeout:             30 * time.Second,
		MaxConcurrent:       10,
		HealthCheckInterval: 30 * time.Second,
		Enabled:             true,
	}

	// 创建客户端
	client, err := NewSandboxClient(config)
	if err != nil {
		fmt.Printf("创建沙箱客户端失败: %v\n", err)
		return
	}
	defer client.Close()

	// Java源代码 - 两数之和
	javaCode := `
import java.util.*;

public class Main {
    public static void main(String[] args) {
        Scanner scanner = new Scanner(System.in);
        
        // 读取数组长度
        int n = scanner.nextInt();
        int[] nums = new int[n];
        
        // 读取数组元素
        for (int i = 0; i < n; i++) {
            nums[i] = scanner.nextInt();
        }
        
        // 读取目标值
        int target = scanner.nextInt();
        
        // 调用解题方法
        int[] result = twoSum(nums, target);
        
        // 输出结果
        System.out.println(result[0] + " " + result[1]);
        
        scanner.close();
    }
    
    public static int[] twoSum(int[] nums, int target) {
        Map<Integer, Integer> map = new HashMap<>();
        for (int i = 0; i < nums.length; i++) {
            int complement = target - nums[i];
            if (map.containsKey(complement)) {
                return new int[]{map.get(complement), i};
            }
            map.put(nums[i], i);
        }
        throw new IllegalArgumentException("No two sum solution");
    }
}
`

	// 测试用例输入: 4个数字[2, 7, 11, 15]，目标值9
	inputData := "4\n2 7 11 15\n9\n"

	ctx := context.Background()

	// 使用CompileAndRunJava方法一次性完成编译和运行
	result, err := client.CompileAndRunJava(ctx, javaCode, inputData, nil, nil)
	if err != nil {
		fmt.Printf("执行失败: %v\n", err)
		return
	}

	// 输出结果
	fmt.Printf("执行状态: %s\n", result.Status)
	fmt.Printf("执行成功: %t\n", result.Success)
	fmt.Printf("程序输出: %s\n", result.Output)
	fmt.Printf("总耗时: %d ms\n", result.TimeUsed)
	fmt.Printf("内存使用: %d KB\n", result.MemoryUsed)
	fmt.Printf("编译时间: %d ms\n", result.CompileTime)
	fmt.Printf("运行时间: %d ms\n", result.RunTime)

	if !result.Success {
		fmt.Printf("错误信息: %s\n", result.CompileError)
		fmt.Printf("错误输出: %s\n", result.ErrorOutput)
	}
}

// ExampleMultiLanguage 多语言代码执行示例
// 演示如何使用RunCode方法执行不同编程语言的代码
func ExampleMultiLanguage() {
	// 创建客户端
	config := &SandboxConfig{
		URL:     "http://localhost:5050",
		Timeout: 30 * time.Second,
		Enabled: true,
	}

	client, err := NewSandboxClient(config)
	if err != nil {
		fmt.Printf("创建沙箱客户端失败: %v\n", err)
		return
	}
	defer client.Close()

	ctx := context.Background()

	// 1. Java代码执行
	fmt.Println("=== Java代码执行 ===")
	javaReq := &CodeExecutionRequest{
		Language:           "java",
		SourceCode:         `public class Main { public static void main(String[] args) { System.out.println("Hello from Java!"); } }`,
		Input:              "",
		TimeLimit:          1000,
		MemoryLimit:        128,
		CompileTimeLimit:   10000,
		CompileMemoryLimit: 256,
		OutputLimit:        10240,
	}

	result, err := client.RunCode(ctx, javaReq)
	if err != nil {
		fmt.Printf("Java执行失败: %v\n", err)
	} else {
		fmt.Printf("Java输出: %s\n", result.Output)
	}

	// 2. Python代码执行
	fmt.Println("\n=== Python代码执行 ===")
	pythonReq := &CodeExecutionRequest{
		Language: "python",
		SourceCode: `print("Hello from Python!")
print(1 + 2)`,
		Input:       "",
		TimeLimit:   1000,
		MemoryLimit: 128,
		OutputLimit: 10240,
	}

	result, err = client.RunCode(ctx, pythonReq)
	if err != nil {
		fmt.Printf("Python执行失败: %v\n", err)
	} else {
		fmt.Printf("Python输出: %s\n", result.Output)
	}

	// 3. C++代码执行
	fmt.Println("\n=== C++代码执行 ===")
	cppReq := &CodeExecutionRequest{
		Language: "cpp",
		SourceCode: `#include<iostream>
using namespace std;
int main() { cout << "Hello from C++!" << endl; return 0; }`,
		Input:              "",
		TimeLimit:          1000,
		MemoryLimit:        128,
		CompileTimeLimit:   10000,
		CompileMemoryLimit: 256,
		OutputLimit:        10240,
	}

	result, err = client.RunCode(ctx, cppReq)
	if err != nil {
		fmt.Printf("C++执行失败: %v\n", err)
	} else {
		fmt.Printf("C++输出: %s\n", result.Output)
	}
}
