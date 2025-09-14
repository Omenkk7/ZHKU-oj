package domain

/*
@Author: omenkk7
@Date: 2025/9/14 10:44
@Description: 用户持久层数据
*/

type User struct {
	UID           string `gorm:"column:UID"`           // 用户唯一 ID（主键，可以是学号/系统生成的 UUID）
	Uname         string `gorm:"column:UserName"`      // 用户名（登录名或展示用昵称）
	PassWord      string `gorm:"column:Pass"`          // 登录密码（加密存储，不能明文）
	School        string `gorm:"column:School"`        // 学校名称（所属学校）
	Classes       string `gorm:"column:Classes"`       // 班级（所属班级，比如“计科2001”）
	Major         string `gorm:"column:Major"`         // 专业（所属专业，比如“计算机科学与技术”）
	Adept         string `gorm:"column:Adept"`         // 擅长方向（可选字段，比如“算法”、“数据库”）
	Vjid          string `gorm:"column:Vjid"`          // Virtual Judge 用户名（用于第三方 OJ 账号绑定）
	Vjpwd         string `gorm:"column:Vjpwd"`         // Virtual Judge 密码（第三方 OJ 账号密码，需加密存储）
	Email         string `gorm:"column:Email"`         // 邮箱（用户注册/找回密码/通知使用）
	CodeForceUser string `gorm:"column:CodeForceUser"` // Codeforces 用户名（用于绑定第三方平台）
	HeadURL       string `gorm:"column:HeadUrl"`       // 头像 URL（用户头像地址）
	Rating        int    `gorm:"column:Rating"`        // 评分（用户评级，通常用于 OJ 积分/比赛 Elo）
	LoginIP       string `gorm:"column:LoginIP"`       // 最近登录 IP（记录用户最后一次登录的 IP 地址）
	RegisterTime  int64  `gorm:"column:RegisterTime"`  // 注册时间（Unix 时间戳，记录用户创建时间）
	Submited      int64  `gorm:"column:Submited"`      // 提交题目数（总共提交过多少次）
	Solved        uint32 `gorm:"column:Solved"`        // 已解决题目数（AC 的题目数）
	Defaulted     string `gorm:"column:Defaulted"`     // 默认状态（可能表示是否启用/封禁用户）
}
