package web

// UserVO 用户视图对象
type UserVO struct {
	UserId   int64  `json:"userId"`
	Username string `json:"username"`
	// 天梯积分
	Score int `json:"score"`
	// 总场数
	TotalCount int `json:"totalCount"`
	// 胜场数
	WinCount int `json:"winCount"`
	// 创建时间（毫秒数）
	Ctime int64 `json:"ctime"`
	// 更新时间（毫秒数）
	Utime int64 `json:"utime"`
}
