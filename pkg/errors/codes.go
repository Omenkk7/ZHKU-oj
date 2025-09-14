package errors

// 错误码定义 - 校园Java-OJ系统统一错误码
// 采用5位数字编码：AABBB
// AA: 模块代码 (10:通用, 20:用户, 30:题目, 40:提交, 50:判题, 60:管理)
// BBB: 具体错误代码 (001-999)

const (
	// ========== 通用错误码 (10000-10999) ==========
	SUCCESS             = 0     // 成功
	SYSTEM_ERROR        = 10001 // 系统内部错误
	INVALID_PARAMS      = 10002 // 参数错误
	UNAUTHORIZED        = 10003 // 未授权
	FORBIDDEN           = 10004 // 无权限
	NOT_FOUND           = 10005 // 资源不存在
	METHOD_NOT_ALLOWED  = 10006 // 请求方法不允许
	TOO_MANY_REQUESTS   = 10007 // 请求过于频繁
	SERVICE_UNAVAILABLE = 10008 // 服务不可用
	INVALID_TOKEN       = 10009 // Token无效或过期
	REQUEST_TIMEOUT     = 10010 // 请求超时

	// ========== 用户模块错误码 (20000-20999) ==========
	USER_NOT_FOUND            = 20001 // 用户不存在
	USER_ALREADY_EXISTS       = 20002 // 用户已存在
	USERNAME_ALREADY_EXISTS   = 20003 // 用户名已存在
	EMAIL_ALREADY_EXISTS      = 20004 // 邮箱已存在
	STUDENT_ID_ALREADY_EXISTS = 20005 // 学号已存在
	INVALID_PASSWORD          = 20006 // 密码错误
	PASSWORD_TOO_WEAK         = 20007 // 密码强度不足
	USER_DISABLED             = 20008 // 用户已被禁用
	USER_NOT_VERIFIED         = 20009 // 用户未验证
	LOGIN_FAILED              = 20010 // 登录失败
	LOGOUT_FAILED             = 20011 // 登出失败
	REGISTER_FAILED           = 20012 // 注册失败
	OLD_PASSWORD_INCORRECT    = 20013 // 旧密码不正确
	PASSWORD_CHANGE_FAILED    = 20014 // 密码修改失败
	PROFILE_UPDATE_FAILED     = 20015 // 用户信息更新失败
	INSUFFICIENT_PERMISSION   = 20016 // 权限不足
	USER_STATS_ERROR          = 20017 // 用户统计信息错误

	// ========== 题目模块错误码 (30000-30999) ==========
	PROBLEM_NOT_FOUND      = 30001 // 题目不存在
	PROBLEM_ALREADY_EXISTS = 30002 // 题目已存在
	PROBLEM_CREATE_FAILED  = 30003 // 题目创建失败
	PROBLEM_UPDATE_FAILED  = 30004 // 题目更新失败
	PROBLEM_DELETE_FAILED  = 30005 // 题目删除失败
	PROBLEM_ACCESS_DENIED  = 30006 // 题目访问被拒绝
	PROBLEM_NOT_PUBLIC     = 30007 // 题目未公开
	TESTCASE_NOT_FOUND     = 30008 // 测试用例不存在
	TESTCASE_INVALID       = 30009 // 测试用例无效
	PROBLEM_STATS_ERROR    = 30010 // 题目统计错误

	// ========== 提交模块错误码 (40000-40999) ==========
	SUBMISSION_NOT_FOUND      = 40001 // 提交记录不存在
	SUBMISSION_CREATE_FAILED  = 40002 // 提交创建失败
	SUBMISSION_UPDATE_FAILED  = 40003 // 提交更新失败
	CODE_TOO_LONG             = 40004 // 代码过长
	CODE_EMPTY                = 40005 // 代码为空
	LANGUAGE_NOT_SUPPORTED    = 40006 // 编程语言不支持
	DUPLICATE_SUBMISSION      = 40007 // 重复提交
	SUBMISSION_ACCESS_DENIED  = 40008 // 提交访问被拒绝
	SUBMISSION_LIMIT_EXCEEDED = 40009 // 提交次数超限
	SUBMISSION_TOO_FREQUENT   = 40010 // 提交过于频繁

	// ========== 判题模块错误码 (50000-50999) ==========
	JUDGE_SYSTEM_ERROR        = 50001 // 判题系统错误
	JUDGE_TIMEOUT             = 50002 // 判题超时
	JUDGE_QUEUE_FULL          = 50003 // 判题队列已满
	COMPILE_ERROR             = 50004 // 编译错误
	RUNTIME_ERROR             = 50005 // 运行时错误
	TIME_LIMIT_EXCEEDED       = 50006 // 时间超限
	MEMORY_LIMIT_EXCEEDED     = 50007 // 内存超限
	OUTPUT_LIMIT_EXCEEDED     = 50008 // 输出超限
	WRONG_ANSWER              = 50009 // 答案错误
	PRESENTATION_ERROR        = 50010 // 格式错误
	SANDBOX_ERROR             = 50011 // 沙箱错误
	JUDGE_SERVICE_UNAVAILABLE = 50012 // 判题服务不可用
	FILE_CACHE_ERROR          = 50013 // 文件缓存错误
	TESTCASE_EXECUTION_ERROR  = 50014 // 测试用例执行错误

	// ========== 管理模块错误码 (60000-60999) ==========
	ADMIN_PERMISSION_DENIED = 60001 // 管理员权限不足
	SYSTEM_MAINTENANCE      = 60002 // 系统维护中
	CONFIG_ERROR            = 60003 // 配置错误
	DATABASE_ERROR          = 60004 // 数据库错误
	CACHE_ERROR             = 60005 // 缓存错误
	MESSAGE_QUEUE_ERROR     = 60006 // 消息队列错误
	BACKUP_FAILED           = 60007 // 备份失败
	RESTORE_FAILED          = 60008 // 恢复失败
	SYSTEM_STATS_ERROR      = 60009 // 系统统计错误

	// ========== 竞赛模块错误码 (70000-70999) ==========
	CONTEST_NOT_FOUND           = 70001 // 竞赛不存在
	CONTEST_NOT_STARTED         = 70002 // 竞赛未开始
	CONTEST_ENDED               = 70003 // 竞赛已结束
	CONTEST_ACCESS_DENIED       = 70004 // 竞赛访问被拒绝
	CONTEST_REGISTRATION_FAILED = 70005 // 竞赛报名失败
)

// 错误码到消息的映射
var errorMessages = map[int]string{
	// 通用错误码
	SUCCESS:             "成功",
	SYSTEM_ERROR:        "系统内部错误",
	INVALID_PARAMS:      "参数错误",
	UNAUTHORIZED:        "未授权访问",
	FORBIDDEN:           "访问被禁止",
	NOT_FOUND:           "资源不存在",
	METHOD_NOT_ALLOWED:  "请求方法不允许",
	TOO_MANY_REQUESTS:   "请求过于频繁，请稍后重试",
	SERVICE_UNAVAILABLE: "服务暂时不可用",
	INVALID_TOKEN:       "Token无效或已过期",
	REQUEST_TIMEOUT:     "请求超时",

	// 用户模块
	USER_NOT_FOUND:            "用户不存在",
	USER_ALREADY_EXISTS:       "用户已存在",
	USERNAME_ALREADY_EXISTS:   "用户名已存在",
	EMAIL_ALREADY_EXISTS:      "邮箱已被注册",
	STUDENT_ID_ALREADY_EXISTS: "学号已被注册",
	INVALID_PASSWORD:          "密码错误",
	PASSWORD_TOO_WEAK:         "密码强度不足",
	USER_DISABLED:             "用户已被禁用",
	USER_NOT_VERIFIED:         "用户未验证",
	LOGIN_FAILED:              "登录失败",
	LOGOUT_FAILED:             "登出失败",
	REGISTER_FAILED:           "注册失败",
	OLD_PASSWORD_INCORRECT:    "旧密码不正确",
	PASSWORD_CHANGE_FAILED:    "密码修改失败",
	PROFILE_UPDATE_FAILED:     "用户信息更新失败",
	INSUFFICIENT_PERMISSION:   "权限不足",
	USER_STATS_ERROR:          "用户统计信息获取失败",

	// 题目模块
	PROBLEM_NOT_FOUND:      "题目不存在",
	PROBLEM_ALREADY_EXISTS: "题目已存在",
	PROBLEM_CREATE_FAILED:  "题目创建失败",
	PROBLEM_UPDATE_FAILED:  "题目更新失败",
	PROBLEM_DELETE_FAILED:  "题目删除失败",
	PROBLEM_ACCESS_DENIED:  "题目访问被拒绝",
	PROBLEM_NOT_PUBLIC:     "题目未公开",
	TESTCASE_NOT_FOUND:     "测试用例不存在",
	TESTCASE_INVALID:       "测试用例无效",
	PROBLEM_STATS_ERROR:    "题目统计信息获取失败",

	// 提交模块
	SUBMISSION_NOT_FOUND:      "提交记录不存在",
	SUBMISSION_CREATE_FAILED:  "提交创建失败",
	SUBMISSION_UPDATE_FAILED:  "提交更新失败",
	CODE_TOO_LONG:             "代码长度超过限制",
	CODE_EMPTY:                "代码不能为空",
	LANGUAGE_NOT_SUPPORTED:    "不支持的编程语言",
	DUPLICATE_SUBMISSION:      "请勿重复提交相同代码",
	SUBMISSION_ACCESS_DENIED:  "提交记录访问被拒绝",
	SUBMISSION_LIMIT_EXCEEDED: "提交次数超过限制",
	SUBMISSION_TOO_FREQUENT:   "提交过于频繁，请稍后重试",

	// 判题模块
	JUDGE_SYSTEM_ERROR:        "判题系统错误",
	JUDGE_TIMEOUT:             "判题超时",
	JUDGE_QUEUE_FULL:          "判题队列已满，请稍后重试",
	COMPILE_ERROR:             "编译错误",
	RUNTIME_ERROR:             "运行时错误",
	TIME_LIMIT_EXCEEDED:       "时间超限",
	MEMORY_LIMIT_EXCEEDED:     "内存超限",
	OUTPUT_LIMIT_EXCEEDED:     "输出超限",
	WRONG_ANSWER:              "答案错误",
	PRESENTATION_ERROR:        "格式错误",
	SANDBOX_ERROR:             "沙箱执行错误",
	JUDGE_SERVICE_UNAVAILABLE: "判题服务不可用",
	FILE_CACHE_ERROR:          "文件缓存错误",
	TESTCASE_EXECUTION_ERROR:  "测试用例执行错误",

	// 管理模块
	ADMIN_PERMISSION_DENIED: "管理员权限不足",
	SYSTEM_MAINTENANCE:      "系统正在维护中",
	CONFIG_ERROR:            "系统配置错误",
	DATABASE_ERROR:          "数据库错误",
	CACHE_ERROR:             "缓存错误",
	MESSAGE_QUEUE_ERROR:     "消息队列错误",
	BACKUP_FAILED:           "备份失败",
	RESTORE_FAILED:          "恢复失败",
	SYSTEM_STATS_ERROR:      "系统统计信息获取失败",

	// 竞赛模块
	CONTEST_NOT_FOUND:           "竞赛不存在",
	CONTEST_NOT_STARTED:         "竞赛尚未开始",
	CONTEST_ENDED:               "竞赛已结束",
	CONTEST_ACCESS_DENIED:       "竞赛访问被拒绝",
	CONTEST_REGISTRATION_FAILED: "竞赛报名失败",
}

// GetErrorMessage 根据错误码获取错误消息
func GetErrorMessage(code int) string {
	if message, exists := errorMessages[code]; exists {
		return message
	}
	return "未知错误"
}

// IsUserError 判断是否为用户模块错误
func IsUserError(code int) bool {
	return code >= 20000 && code < 21000
}

// IsProblemError 判断是否为题目模块错误
func IsProblemError(code int) bool {
	return code >= 30000 && code < 31000
}

// IsSubmissionError 判断是否为提交模块错误
func IsSubmissionError(code int) bool {
	return code >= 40000 && code < 41000
}

// IsJudgeError 判断是否为判题模块错误
func IsJudgeError(code int) bool {
	return code >= 50000 && code < 51000
}

// IsAdminError 判断是否为管理模块错误
func IsAdminError(code int) bool {
	return code >= 60000 && code < 61000
}

// IsContestError 判断是否为竞赛模块错误
func IsContestError(code int) bool {
	return code >= 70000 && code < 71000
}
