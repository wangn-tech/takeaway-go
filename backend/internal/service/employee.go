package service

import (
	"context"
	"errors"
	"takeaway-go/common/enum"
	"takeaway-go/common/result"
	"takeaway-go/internal/api/request"
	"takeaway-go/internal/api/response"
	"takeaway-go/internal/app/config"
	"takeaway-go/internal/model"
	"takeaway-go/internal/repository"
	"takeaway-go/internal/utils"
	"takeaway-go/pkg/logger"
	"time"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type IEmployeeService interface {
	Login(context.Context, request.EmployeeLoginDTO) (*response.EmployeeLoginVO, error)
	Logout(ctx context.Context, userID uint64) error
	AddEmployee(ctx context.Context, dto request.EmployeeAddDTO) error
	PageQuery(ctx context.Context, dto request.EmployeePageQueryDTO) (*result.PageResult, error)
	UpdateStatus(ctx context.Context, id uint64, status int) error
	UpdateEmployee(ctx context.Context, dto *request.EmployeeEditDTO) error
	GetByID(ctx context.Context, id uint64) (*model.Employee, error)
	EditPassword(context.Context, request.EmployeeEditPasswordDTO) error
}

type EmployeeService struct {
	repo *repository.EmployeeDao
}

func NewEmployeeService(repo *repository.EmployeeDao) IEmployeeService {
	// return &EmployeeService{
	// 	repo: repository.NewEmployeeDao(database.DB),
	// }
	return &EmployeeService{repo: repo}
}

const DefaultPassword = "123456" // 定义默认密码

// GetByID 根据ID获取员工信息
func (s *EmployeeService) GetByID(ctx context.Context, id uint64) (*model.Employee, error) {
	employee, err := s.repo.GetByID(ctx, id)
	if err != nil {
		logger.Log.Error("GetByID: 获取员工信息失败", zap.Uint64("id", id), zap.Error(err))
		return nil, err
	}
	employee.Password = "***"
	return employee, nil
}

func (s *EmployeeService) EditPassword(ctx context.Context, dto request.EmployeeEditPasswordDTO) error {
	employee, err := s.repo.GetByID(ctx, dto.EmpId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("员工不存在")
		}
		return err
	}

	// 校验用户密码
	if !utils.CheckPasswordHash(dto.OldPassword, employee.Password) {
		return errors.New("旧密码错误")
	}

	// 对新密码进行加密
	hashedPassword, err := utils.HashPassword(dto.NewPassword)
	if err != nil {
		logger.Log.Error("EditPassword: 密码加密失败", zap.Error(err))
		return errors.New("密码加密失败")
	}

	// 更新员工密码
	err = s.repo.UpdatePassword(ctx, dto.EmpId, hashedPassword)
	if err != nil {
		logger.Log.Error("EditPassword: 更新密码失败", zap.Uint64("empId", dto.EmpId), zap.Error(err))
		return errors.New("更新密码失败")
	}
	// 密码修改成功后，撤销该用户的所有 token
	if err := utils.RevokeToken(dto.EmpId); err != nil {
		logger.Log.Error("EditPassword: 撤销 token 失败", zap.Uint64("userId", dto.EmpId), zap.Error(err))
		return errors.New("修改密码成功，但登出失败，请重新登录")
	}
	logger.Log.Info("EditPassword: 修改密码成功", zap.Uint64("empId", dto.EmpId))
	return nil
}

// Login 处理员工登录逻辑
func (s *EmployeeService) Login(ctx context.Context, loginDTO request.EmployeeLoginDTO) (*response.EmployeeLoginVO, error) {
	// 查询员工是否存在
	employee, err := s.repo.GetByUsername(ctx, loginDTO.Username)
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
	if !utils.CheckPasswordHash(loginDTO.Password, employee.Password) {
		return nil, errors.New("密码错误")
	}

	// 生成 JWT token
	jwtConfig := config.AppConf.JWT.Admin
	tokenResponse, err := utils.GenerateTokenPair(employee.ID, jwtConfig.Name, jwtConfig.Secret)
	if err != nil {
		return nil, errors.New("生成令牌失败")
	}

	// ✅ 将 token 存储到 Redis（7天过期）
	if err := utils.StoreTokenInRedis(employee.ID, tokenResponse.AccessToken, ""); err != nil {
		logger.Log.Error("Login: 存储 token 失败", zap.Uint64("userID", employee.ID), zap.Error(err))
		// 可以选择不返回错误，因为 token 已经生成
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

// Logout 处理员工登出逻辑
func (s *EmployeeService) Logout(ctx context.Context, userID uint64) error {
	// 调用 RevokeToken 函数，从 Redis 中删除用户的 token
	if err := utils.RevokeToken(userID); err != nil {
		logger.Log.Error("Logout: 撤销 token 失败", zap.Uint64("userID", userID), zap.Error(err))
		return errors.New("登出失败")
	}

	logger.Log.Info("Logout: 登出成功", zap.Uint64("userID", userID))
	return nil
}

// AddEmployee 新增员工
func (s *EmployeeService) AddEmployee(ctx context.Context, dto request.EmployeeAddDTO) error {
	// 检查用户名是否已存在
	if _, err := s.repo.GetByUsername(ctx, dto.Username); err == nil {
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
	return s.repo.Save(ctx, &employee)
}

// PageQuery 分页查询员工信息
func (s *EmployeeService) PageQuery(ctx context.Context, dto request.EmployeePageQueryDTO) (*result.PageResult, error) {
	// 调用 repository 进行分页查询
	pageResult, err := s.repo.PageQuery(ctx, dto)
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

// Update 更新员工信息
func (s *EmployeeService) UpdateEmployee(ctx context.Context, dto *request.EmployeeEditDTO) error {
	// 检查员工是否存在
	existEmployee, err := s.repo.GetByID(ctx, dto.ID)
	if err != nil {
		return errors.New("员工不存在")
	}

	// 如果修改了用户名,检查用户名是否已被使用
	if dto.Username != existEmployee.Username {
		_, err := s.repo.GetByUsername(ctx, dto.Username)
		if err == nil {
			return errors.New("用户名已存在")
		}
	}

	// 更新员工信息
	employee := &model.Employee{
		ID:         dto.ID,
		Username:   dto.Username,
		Name:       dto.Name,
		Phone:      dto.Phone,
		Sex:        dto.Sex,
		IdNumber:   dto.IdNumber,
		UpdateTime: time.Now(),
	}

	return s.repo.Update(ctx, employee)
}

// UpdateStatus 更新员工状态
func (s *EmployeeService) UpdateStatus(ctx context.Context, id uint64, status int) error {
	// 检查员工是否存在
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return errors.New("员工不存在")
	}

	// 更新状态
	err = s.repo.UpdateStatus(ctx, &model.Employee{
		ID:     id,
		Status: status,
	})
	if err != nil {
		return err
	}

	return nil
}
