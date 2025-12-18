package model

type BlackList struct {
	Token string `json:"token" gorm:"column:token;uniqueIndex;type:varchar(255)"`
}
