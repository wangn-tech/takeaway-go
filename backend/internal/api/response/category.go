package response

import "time"

// CategoryVO 分类信息返回
type CategoryVO struct {
	ID         uint64    `json:"id"`
	Type       int       `json:"type"`   // 1:菜品分类 2:套餐分类
	Name       string    `json:"name"`   // 分类名称
	Sort       int       `json:"sort"`   // 排序
	Status     int       `json:"status"` // 0:禁用 1:启用
	CreateTime time.Time `json:"createTime"`
	UpdateTime time.Time `json:"updateTime"`
	CreateUser uint64    `json:"createUser"`
	UpdateUser uint64    `json:"updateUser"`
}
