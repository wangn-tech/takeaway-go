package response

import "time"

// EmployeeLoginVO 员工登录成功后返回的视图对象
type EmployeeLoginVO struct {
	ID       uint64 `json:"id"`
	Name     string `json:"name"`
	UserName string `json:"userName"`
	Token    string `json:"token"`
}

// EmployeePageVO 员工分页查询返回的视图对象
type EmployeePageVO struct {
	ID         uint64    `json:"id"`
	Name       string    `json:"name"`
	Username   string    `json:"username"`
	Phone      string    `json:"phone"`
	Sex        string    `json:"sex"`
	Status     int       `json:"status"`
	UpdateTime time.Time `json:"updateTime"`
}
