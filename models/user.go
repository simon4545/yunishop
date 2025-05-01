package models

import "time"

// 用户基础结构
type BaseUser struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	UserName   string    `gorm:"size:50"`
	Password   string    `gorm:"size:50"`
	Name       string    `gorm:"size:100;not null" json:"name"`
	BirthDate  time.Time `json:"birth_date"`
	Province   string    `gorm:"size:50" json:"province"`
	City       string    `gorm:"size:50" json:"city"`
	Address    string    `gorm:"size:255" json:"address"`
	Gender     string    `gorm:"size:10" json:"gender"` // male, female, other
	ReferrerID uint      `json:"referrer_id"`           // 介绍人ID
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// 客户
type Customer struct {
	BaseUser
	Level      int
	SessionKey string `gorm:"size:255"`
	OpenId     string `gorm:"size:255" json:"openid"`
	AgentID    uint   `json:"agent_id"` // 归属代理商编号
}

// 代理商
type Agent struct {
	BaseUser
}

// 管理员
type Admin struct {
	BaseUser
	Role string `gorm:"size:50" json:"role"` // 管理员角色
}
