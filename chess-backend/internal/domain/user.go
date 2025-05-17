package domain

// User 用户领域对象
type User struct {
	UserId   int64
	Username string
	Password string
	// 天梯积分
	Score int
	// 总场数
	TotalCount int
	// 胜场数
	WinCount int
	// 创建时间（毫秒数）
	Ctime int64
	// 更新时间（毫秒数）
	Utime int64
}
