# 校园Java-OJ系统数据库设计

## 📊 MongoDB 集合设计

### 1. users 集合 - 用户信息
```json
{
  "_id": ObjectId("64f8a123b45c6789d0123456"),
  "student_id": "2021001001",
  "username": "zhang_san",
  "password": "$2a$10$hashed_password",
  "email": "zhangsan@school.edu.cn",
  "real_name": "张三",
  "role": "student", // student, teacher, admin
  "class": "计算机科学与技术2021-1班",
  "grade": "2021",
  "avatar": "http://example.com/avatar.jpg",
  "is_active": true,
  "stats": {
    "total_submissions": 45,
    "accepted_count": 23,
    "problems_solved": 18,
    "ranking": 12,
    "total_score": 1250,
    "max_streak": 7,
    "current_streak": 3
  },
  "preferences": {
    "language": "java",
    "theme": "light",
    "notifications": true
  },
  "created_at": ISODate("2024-01-15T10:30:00Z"),
  "updated_at": ISODate("2024-01-20T14:20:00Z"),
  "last_login": ISODate("2024-01-20T09:15:00Z"),
  "login_count": 156
}
```

### 2. problems 集合 - 题目信息
```json
{
  "_id": ObjectId("64f8a123b45c6789d0123457"),
  "title": "两数之和",
  "description": "给定一个整数数组 nums 和一个整数目标值 target...",
  "input_format": "第一行包含数组长度n和目标值target...",
  "output_format": "输出两个整数的下标...",
  "sample_input": "4 9\n2 7 11 15",
  "sample_output": "0 1",
  "time_limit": 1000, // 毫秒
  "memory_limit": 128, // MB
  "difficulty": "easy", // easy, medium, hard
  "tags": ["array", "hash-table", "two-pointers"],
  "category": "algorithm", // algorithm, data-structure, math
  "source": "LeetCode", // 题目来源
  "test_cases": [
    {
      "id": "case1",
      "input": "4 9\n2 7 11 15",
      "output": "0 1",
      "score": 20,
      "is_public": true,
      "description": "基础测试用例"
    },
    {
      "id": "case2", 
      "input": "3 6\n3 2 4",
      "output": "1 2",
      "score": 30,
      "is_public": false,
      "description": "边界测试用例"
    }
  ],
  "stats": {
    "total_submissions": 1250,
    "accepted_count": 812,
    "acceptance_rate": 0.65,
    "average_time": 245,
    "average_memory": 8192,
    "difficulty_rating": 4.2
  },
  "constraints": {
    "max_code_length": 50000,
    "allowed_languages": ["java"],
    "special_judge": false
  },
  "is_public": true,
  "is_active": true,
  "created_by": ObjectId("64f8a123b45c6789d0123460"),
  "created_at": ISODate("2024-01-10T08:00:00Z"),
  "updated_at": ISODate("2024-01-15T10:30:00Z"),
  "publish_time": ISODate("2024-01-12T00:00:00Z")
}
```

### 3. submissions 集合 - 提交记录
```json
{
  "_id": ObjectId("64f8a123b45c6789d0123458"),
  "user_id": ObjectId("64f8a123b45c6789d0123456"),
  "problem_id": ObjectId("64f8a123b45c6789d0123457"),
  "code": "public class Main {\n    public static void main(String[] args) {\n        // Java解题代码\n    }\n}",
  "language": "java",
  "status": "ACCEPTED", // PENDING, JUDGING, ACCEPTED, WRONG_ANSWER, etc.
  "score": 100,
  "time_used": 245, // 毫秒
  "memory_used": 8192, // KB
  "code_length": 450,
  "compile_info": {
    "status": "SUCCESS", // SUCCESS, FAILED
    "time_used": 2456, // 纳秒
    "memory_used": 59801600, // 字节
    "message": "", // 编译错误信息
    "compiler_output": {
      "stdout": "",
      "stderr": ""
    }
  },
  "test_results": [
    {
      "test_case_id": "case1",
      "status": "ACCEPTED",
      "time_used": 45, // 毫秒
      "memory_used": 2048, // KB
      "score": 20,
      "input": "4 9\n2 7 11 15",
      "expected_output": "0 1",
      "actual_output": "0 1",
      "judge_details": {
        "go_judge_status": "Accepted",
        "exit_status": 0,
        "runtime_ns": 45000000,
        "sandbox_instance": "go-judge-1",
        "file_cache_used": true,
        "execution_id": "exec_abc123"
      }
    }
  ],
  "judge_info": {
    "judge_server": "go-judge-1",
    "total_judge_time": "12.5s",
    "compile_time_ns": 870867000,
    "avg_run_time_ns": 123456789,
    "file_cache_cleared": true,
    "retry_count": 0,
    "queue_time": 1500 // 排队等待时间(毫秒)
  },
  "ip_address": "192.168.1.100",
  "user_agent": "Mozilla/5.0...",
  "submitted_at": ISODate("2024-01-15T14:30:00Z"),
  "judged_at": ISODate("2024-01-15T14:30:15Z"),
  "contest_id": null // 如果是竞赛提交
}
```

### 4. contests 集合 - 竞赛信息 (预留)
```json
{
  "_id": ObjectId("64f8a123b45c6789d0123459"),
  "title": "2024春季编程竞赛",
  "description": "面向大一大二学生的编程竞赛",
  "type": "public", // public, private, class
  "problem_ids": [
    ObjectId("64f8a123b45c6789d0123457"),
    ObjectId("64f8a123b45c6789d0123461")
  ],
  "participant_ids": [
    ObjectId("64f8a123b45c6789d0123456")
  ],
  "settings": {
    "duration": 7200, // 秒
    "penalty_time": 1200, // 错误提交惩罚时间
    "freeze_time": 3600, // 封榜时间
    "max_submissions": 50
  },
  "start_time": ISODate("2024-03-15T09:00:00Z"),
  "end_time": ISODate("2024-03-15T11:00:00Z"),
  "freeze_time": ISODate("2024-03-15T10:00:00Z"),
  "status": "upcoming", // upcoming, running, ended
  "created_by": ObjectId("64f8a123b45c6789d0123460"),
  "created_at": ISODate("2024-03-01T10:00:00Z"),
  "updated_at": ISODate("2024-03-10T15:30:00Z")
}
```

### 5. user_stats 集合 - 用户统计详情
```json
{
  "_id": ObjectId("64f8a123b45c6789d0123462"),
  "user_id": ObjectId("64f8a123b45c6789d0123456"),
  "period": "2024-01", // 统计周期
  "stats": {
    "submissions_by_status": {
      "ACCEPTED": 23,
      "WRONG_ANSWER": 15,
      "TIME_LIMIT_EXCEEDED": 5,
      "COMPILE_ERROR": 2
    },
    "problems_by_difficulty": {
      "easy": {"solved": 8, "attempted": 12},
      "medium": {"solved": 7, "attempted": 18},
      "hard": {"solved": 3, "attempted": 8}
    },
    "daily_submissions": {
      "2024-01-15": 5,
      "2024-01-16": 3,
      "2024-01-17": 7
    },
    "average_time": 280,
    "average_memory": 6800,
    "best_streak": 7,
    "total_score": 1250
  },
  "created_at": ISODate("2024-01-31T23:59:59Z"),
  "updated_at": ISODate("2024-01-31T23:59:59Z")
}
```

### 6. problem_stats 集合 - 题目统计详情
```json
{
  "_id": ObjectId("64f8a123b45c6789d0123463"),
  "problem_id": ObjectId("64f8a123b45c6789d0123457"),
  "period": "2024-01",
  "stats": {
    "submissions_by_status": {
      "ACCEPTED": 812,
      "WRONG_ANSWER": 320,
      "TIME_LIMIT_EXCEEDED": 85,
      "COMPILE_ERROR": 33
    },
    "submissions_by_language": {
      "java": {"total": 1250, "accepted": 812}
    },
    "time_distribution": {
      "0-100ms": 245,
      "100-500ms": 420,
      "500-1000ms": 147
    },
    "memory_distribution": {
      "0-50MB": 680,
      "50-100MB": 132
    },
    "daily_submissions": {
      "2024-01-15": 45,
      "2024-01-16": 32
    }
  },
  "created_at": ISODate("2024-01-31T23:59:59Z"),
  "updated_at": ISODate("2024-01-31T23:59:59Z")
}
```

### 7. judge_queue 集合 - 判题队列状态
```json
{
  "_id": ObjectId("64f8a123b45c6789d0123464"),
  "submission_id": ObjectId("64f8a123b45c6789d0123458"),
  "status": "PROCESSING", // PENDING, PROCESSING, COMPLETED, FAILED
  "priority": 5, // 1-10, 10最高
  "retry_count": 0,
  "assigned_judge": "go-judge-1",
  "progress": {
    "stage": "RUNNING", // COMPILING, RUNNING, COMPLETED
    "current_test_case": 3,
    "total_test_cases": 5,
    "start_time": ISODate("2024-01-15T14:30:05Z"),
    "estimated_finish": ISODate("2024-01-15T14:30:20Z")
  },
  "error_info": {
    "error_type": "SANDBOX_ERROR",
    "error_message": "Connection timeout to go-judge",
    "error_count": 1
  },
  "created_at": ISODate("2024-01-15T14:30:00Z"),
  "updated_at": ISODate("2024-01-15T14:30:05Z")
}
```

### 8. system_logs 集合 - 系统日志
```json
{
  "_id": ObjectId("64f8a123b45c6789d0123465"),
  "level": "INFO", // DEBUG, INFO, WARN, ERROR, FATAL
  "service": "judge-service",
  "action": "submit_code",
  "user_id": ObjectId("64f8a123b45c6789d0123456"),
  "submission_id": ObjectId("64f8a123b45c6789d0123458"),
  "message": "用户提交代码成功",
  "details": {
    "ip_address": "192.168.1.100",
    "user_agent": "Mozilla/5.0...",
    "execution_time": 150,
    "memory_usage": 45678912
  },
  "trace_id": "abc123def456",
  "timestamp": ISODate("2024-01-15T14:30:00Z")
}
```

## 🔍 索引设计

### 用户集合索引
```javascript
db.users.createIndex({ "username": 1 }, { unique: true })
db.users.createIndex({ "email": 1 }, { unique: true })
db.users.createIndex({ "student_id": 1 }, { unique: true })
db.users.createIndex({ "class": 1, "stats.ranking": 1 })
db.users.createIndex({ "role": 1, "is_active": 1 })
db.users.createIndex({ "created_at": 1 })
```

### 题目集合索引
```javascript
db.problems.createIndex({ "difficulty": 1, "is_public": 1 })
db.problems.createIndex({ "tags": 1, "is_active": 1 })
db.problems.createIndex({ "created_by": 1, "created_at": -1 })
db.problems.createIndex({ "stats.acceptance_rate": -1 })
db.problems.createIndex({ "title": "text", "description": "text" })
```

### 提交记录索引
```javascript
db.submissions.createIndex({ "user_id": 1, "submitted_at": -1 })
db.submissions.createIndex({ "problem_id": 1, "status": 1 })
db.submissions.createIndex({ "status": 1, "submitted_at": -1 })
db.submissions.createIndex({ "user_id": 1, "problem_id": 1, "status": 1 })
db.submissions.createIndex({ "contest_id": 1, "submitted_at": 1 })
db.submissions.createIndex({ "judged_at": -1 })
```

### 判题队列索引
```javascript
db.judge_queue.createIndex({ "status": 1, "priority": -1, "created_at": 1 })
db.judge_queue.createIndex({ "submission_id": 1 }, { unique: true })
db.judge_queue.createIndex({ "assigned_judge": 1, "status": 1 })
```

### 日志集合索引
```javascript
db.system_logs.createIndex({ "timestamp": -1 })
db.system_logs.createIndex({ "level": 1, "timestamp": -1 })
db.system_logs.createIndex({ "service": 1, "action": 1, "timestamp": -1 })
db.system_logs.createIndex({ "user_id": 1, "timestamp": -1 })
db.system_logs.createIndex({ "trace_id": 1 })
```

## 📈 Redis 缓存设计

### 1. 用户会话缓存
```
Key: session:{user_id}
Value: {
  "token": "jwt_token_string",
  "expires_at": "2024-01-16T10:30:00Z",
  "login_ip": "192.168.1.100",
  "last_activity": "2024-01-15T14:30:00Z"
}
TTL: 24小时
```

### 2. 排行榜缓存
```
Key: ranking:class:{class_name}
Type: Sorted Set
Score: problems_solved
Member: user_id
TTL: 1小时
```

### 3. 题目列表缓存
```
Key: problems:list:{difficulty}:{page}
Value: [problem_list_json]
TTL: 30分钟
```

### 4. 提交状态缓存
```
Key: submission:{submission_id}:status
Value: {
  "status": "JUDGING",
  "progress": 60,
  "current_test": 3,
  "total_tests": 5
}
TTL: 30分钟
```

### 5. 题目统计缓存
```
Key: problem:{problem_id}:stats
Value: {
  "total_submissions": 1250,
  "accepted_count": 812,
  "acceptance_rate": 0.65
}
TTL: 1小时
```