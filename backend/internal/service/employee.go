package service

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"takeaway-go/common/enum"
	"takeaway-go/common/result"
	"takeaway-go/internal/api/request"
	"takeaway-go/internal/api/response"
	"takeaway-go/internal/app/config"
	"takeaway-go/internal/app/initializer/database"
	"takeaway-go/internal/model"
	"takeaway-go/internal/repository"
	utils2 "takeaway-go/internal/utils"
)

type EmployeeService struct {
	repo *repository.EmployeeDao
}

func NewEmployeeService() *EmployeeService {
	return &EmployeeService{
		repo: repository.NewEmployeeDao(database.DB),
	}
}

const DefaultPassword = "123456" // 定义默认密码

// Login 处理员工登录逻辑
func (s *EmployeeService) Login(loginDTO request.EmployeeLoginDTO) (*response.EmployeeLoginVO, error) {
	// 查询员工是否存在
	employee, err := s.repo.GetByUsername(loginDTO.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("账号不存在")
		}
		return nil, err // 其他数据库错误
	}

	// 检查账号状态
	if employee.Status == 0 {
		return nil, errors.New("账号已被锁定")
	}

	// 校验密码 (使用 bcrypt)
	if !utils2.CheckPasswordHash(loginDTO.Password, employee.Password) {
		return nil, errors.New("密码错误")
	}

	// 生成 JWT token
	jwtConfig := config.AppConf.JWT.Admin
	tokenResponse, err := utils2.GenerateTokenPair(employee.ID, jwtConfig.Name, jwtConfig.Secret)
	if err != nil {
		return nil, errors.New("生成令牌失败")
	}

	// 封装返回结果
	loginVO := &response.EmployeeLoginVO{
		ID:       employee.ID,
		UserName: employee.Username,
		Name:     employee.Name,
		Token:    tokenResponse.AccessToken,
	}

	return loginVO, nil
}

// AddEmployee 新增员工
func (s *EmployeeService) AddEmployee(dto request.EmployeeAddDTO) error {
	// 检查用户名是否已存在
	if _, err := s.repo.GetByUsername(dto.Username); err == nil {
		return errors.New("用户名已存在")
	}

	// 对默认密码进行加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(DefaultPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("密码加密失败")
	}

	// 创建 model.Employee
	employee := model.Employee{
		Username: dto.Username,
		Name:     dto.Name,
		Password: string(hashedPassword),
		Phone:    dto.Phone,
		Sex:      dto.Sex,
		IdNumber: dto.IdNumber,
		Status:   enum.ENABLE, // 默认为启用状态
	}

	// 保存到数据库
	return s.repo.Save(&employee)
}

func (s *EmployeeService) PageQuery(dto request.EmployeePageQueryDTO) (*result.PageResult, error) {
	// 调用 repository 进行分页查询
	pageResult, err := s.repo.PageQuery(dto)
	if err != nil {
		return nil, err
	}

	// 数据转换：将 []model.Employee 转换为 []response.EmployeePageVO
	// 获取查询到的员工列表
	employees, ok := pageResult.Records.([]model.Employee)
	if !ok {
		return nil, errors.New("数据类型转换失败")
	}

	// 创建一个新的 VO 列表
	var employeeVOs []response.EmployeePageVO
	for _, emp := range employees {
		employeeVOs = append(employeeVOs, response.EmployeePageVO{
			ID:         emp.ID,
			Name:       emp.Name,
			Username:   emp.Username,
			Phone:      emp.Phone,
			Sex:        emp.Sex,
			Status:     emp.Status,
			UpdateTime: emp.UpdateTime,
		})
	}

	// 将转换后的数据放回 PageResult
	pageResult.Records = employeeVOs

	return pageResult, nil
}
