package model

import (
	"takeaway-go/common/enum"
	"time"

	"gorm.io/gorm"
)

type Category struct {
	ID         uint64    `json:"id" gorm:"primaryKey"`
	Type       int       `json:"type"`
	Name       string    `json:"name"`
	Sort       int       `json:"sort"`
	Status     int       `json:"status" gorm:"default:1"`
	CreateTime time.Time `json:"createTime"`
	UpdateTime time.Time `json:"updateTime"`
	CreateUser uint64    `json:"createUser"`
	UpdateUser uint64    `json:"updateUser"`
}

// TableName 指定表名
func (Category) TableName() string {
	return "category"
}

// BeforeCreate 在创建记录前自动填充字段
func (c *Category) BeforeCreate(tx *gorm.DB) error {
	// 自动填充 创建时间、创建人、更新时间、更新用户
	c.CreateTime = time.Now()
	c.UpdateTime = time.Now()
	// 从上下文获取用户信息
	value := tx.Statement.Context.Value(enum.CurrentId)
	if uid, ok := value.(uint64); ok {
		c.CreateUser = uid
		c.UpdateUser = uid
	}
	return nil
}

// BeforeUpdate 在更新记录前自动填充字段
func (c *Category) BeforeUpdate(tx *gorm.DB) error {
	// 在更新记录时自动填充更新时间
	c.UpdateTime = time.Now()
	// 从上下文获取用户信息
	value := tx.Statement.Context.Value(enum.CurrentId)
	if uid, ok := value.(uint64); ok {
		c.UpdateUser = uid
	}
	return nil
}
