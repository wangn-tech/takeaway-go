package repository

import (
	"context"
	"takeaway-go/common/result"
	"takeaway-go/internal/api/request"
	"takeaway-go/internal/model"
	"time"

	"gorm.io/gorm"
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
func (e *EmployeeDao) GetByUsername(ctx context.Context, username string) (*model.Employee, error) {
	var employee model.Employee
	// 使用 First 来获取第一条匹配的记录，如果没有找到会返回 gorm.ErrRecordNotFound 错误
	err := e.db.WithContext(ctx).Where("username = ?", username).First(&employee).Error
	if err != nil {
		return nil, err
	}
	return &employee, nil
}

// Save 新增员工信息
func (e *EmployeeDao) Save(ctx context.Context, employee *model.Employee) error {
	// Create 方法用于插入一条新记录
	return e.db.WithContext(ctx).Create(employee).Error
}

// PageQuery 分页查询员工信息
//
//	// select count(*) from employee where name like %name% limit x,y
//	// select count(*) from employee where name = ? limit x,y
func (e *EmployeeDao) PageQuery(ctx context.Context, dto request.EmployeePageQueryDTO) (*result.PageResult, error) {
	var pageResult result.PageResult
	var employeeList []model.Employee
	var err error

	// 构建查询
	query := e.db.WithContext(ctx).Model(&model.Employee{})

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

// Update 更新员工信息
func (e *EmployeeDao) Update(ctx context.Context, employee *model.Employee) error {
	err := e.db.WithContext(ctx).Model(&model.Employee{}).Where("id = ?", employee.ID).Updates(employee).Error
	if err != nil {
		return err
	}
	return nil
}

// UpdateStatus 更新员工状态
func (e *EmployeeDao) UpdateStatus(ctx context.Context, employee *model.Employee) error {
	err := e.db.WithContext(ctx).Model(&model.Employee{}).Where("id = ?", employee.ID).Update("status", employee.Status).Error
	if err != nil {
		return err
	}
	return nil
}

// GetByID 根据ID获取员工信息
func (e *EmployeeDao) GetByID(ctx context.Context, id uint64) (*model.Employee, error) {
	var employee model.Employee
	err := e.db.WithContext(ctx).Where("id = ?", id).First(&employee).Error
	if err != nil {
		return nil, err
	}
	return &employee, nil
}

// UpdatePassword 更新员工密码
func (r *EmployeeDao) UpdatePassword(ctx context.Context, empId uint64, hashedPassword string) error {
	return r.db.WithContext(ctx).
		Model(&model.Employee{}).
		Where("id = ?", empId).
		Updates(map[string]interface{}{
			"password":    hashedPassword,
			"update_time": time.Now(),
		}).Error
}
