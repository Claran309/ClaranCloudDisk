package model

// BlackList 黑名单模型
// @Description JWT令牌黑名单，用于存储已注销的令牌
type BlackList struct {
	Token string `json:"token" gorm:"column:token;uniqueIndex;type:varchar(500)"`
}
