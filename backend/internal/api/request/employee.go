package request

// EmployeeLoginDTO 员工登录请求的数据传输对象
type EmployeeLoginDTO struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// EmployeeAddDTO 定义了新增员工时需要传入的参数
type EmployeeAddDTO struct {
	Username string `json:"username" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Phone    string `json:"phone" binding:"required"`
	Sex      string `json:"sex" binding:"required"`
	IdNumber string `json:"idNumber" binding:"required"`
}

// EmployeePageQueryDTO 定义了分页查询员工时可以传入的参数
type EmployeePageQueryDTO struct {
	Name     string `form:"name"`     // 员工姓名 (模糊插叙, 可选)
	Page     int    `form:"page"`     // 页数
	PageSize int    `form:"pageSize"` // 页容量
}

// EmployeeEditDTO 编辑员工信息的数据传输对象
type EmployeeEditDTO struct {
	ID       uint64 `json:"id" binding:"required"`
	Username string `json:"username" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Phone    string `json:"phone" binding:"required,len=11"`
	Sex      string `json:"sex" binding:"required"`
	IdNumber string `json:"idNumber" binding:"required,len=18"`
}

// EmployeeUpdateStatusDTO 更新员工状态的数据传输对象
type EmployeeUpdateStatusDTO struct {
	ID     uint64 `json:"id" binding:"required"`
	Status int    `json:"status" binding:"required,oneof=0 1"`
}

// EmployeeEditPassword 编辑员工密码的数据传输对象
type EmployeeEditPasswordDTO struct {
	EmpId       uint64 `json:"empId"`
	NewPassword string `json:"newPassword" binding:"required"`
	OldPassword string `json:"oldPassword" binding:"required"`
}
