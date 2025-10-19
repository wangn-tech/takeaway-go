package repository

import (
	"gorm.io/gorm"
	"takeaway-go/common/result"
	"takeaway-go/internal/api/request"
	"takeaway-go/internal/model"
)

type EmployeeDao struct {
	db *gorm.DB
}

func NewEmployeeDao(db *gorm.DB) *EmployeeDao {
	return &EmployeeDao{
		db: db,
	}
}

// GetByUsername 根据用户名查询员工信息
func (e *EmployeeDao) GetByUsername(username string) (*model.Employee, error) {
	var employee model.Employee
	// 使用 First 来获取第一条匹配的记录，如果没有找到会返回 gorm.ErrRecordNotFound 错误
	err := e.db.Where("username = ?", username).First(&employee).Error
	if err != nil {
		return nil, err
	}
	return &employee, nil
}

// Save 新增员工信息
func (e *EmployeeDao) Save(employee *model.Employee) error {
	// Create 方法用于插入一条新记录
	return e.db.Create(employee).Error
}

// PageQuery 分页查询员工信息
//
//	// select count(*) from employee where name like %name% limit x,y
//	// select count(*) from employee where name = ? limit x,y
func (e *EmployeeDao) PageQuery(dto request.EmployeePageQueryDTO) (*result.PageResult, error) {
	var pageResult result.PageResult
	var employeeList []model.Employee
	var err error

	// 构建查询
	query := e.db.Model(&model.Employee{})

	// 如果提供了 name，则添加模糊查询条件
	if dto.Name != "" {
		query = query.Where("name LIKE ?", "%"+dto.Name+"%")
	}

	// 计算总数 total
	if err = query.Count(&pageResult.Total).Error; err != nil {
		return nil, err
	}

	// 执行分页查询
	err = query.Scopes(pageResult.Paginate(&dto.Page, &dto.PageSize)).Find(&employeeList).Error
	if err != nil {
		return nil, err
	}
	pageResult.Records = employeeList

	return &pageResult, nil
}
