package model

import (
	"takeaway-go/common/enum"
	"time"

	"gorm.io/gorm"
)

type Employee struct {
	ID         uint64    `json:"id" gorm:"primaryKey"`
	Name       string    `json:"name"`
	Username   string    `json:"username" gorm:"unique"`
	Password   string    `json:"password"`
	Phone      string    `json:"phone"`
	Sex        string    `json:"sex"`
	IdNumber   string    `json:"idNumber"`
	Status     int       `json:"status" gorm:"default:1"`
	CreateTime time.Time `json:"createTime"`
	UpdateTime time.Time `json:"updateTime"`
	CreateUser uint64    `json:"createUser"`
	UpdateUser uint64    `json:"updateUser"`
}

// TableName 指定表名
func (Employee) TableName() string {
	return "employee"
}

func (e *Employee) BeforeCreate(tx *gorm.DB) error {
	// 自动填充 创建时间、创建人、更新时间、更新用户
	e.CreateTime = time.Now()
	e.UpdateTime = time.Now()
	// 从上下文获取用户信息
	value := tx.Statement.Context.Value(enum.CurrentId)
	if uid, ok := value.(uint64); ok {
		e.CreateUser = uid
		e.UpdateUser = uid
	}
	return nil
}
