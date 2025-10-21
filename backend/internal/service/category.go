package service

import (
	"context"
	"errors"
	"takeaway-go/common/enum"
	"takeaway-go/common/result"
	"takeaway-go/internal/api/request"
	"takeaway-go/internal/api/response"
	"takeaway-go/internal/model"
	"takeaway-go/internal/repository"
	"takeaway-go/pkg/logger"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ICategoryService interface {
	// Update 修改分类
	Update(ctx context.Context, dto request.CategoryDTO) error
	// PageQuery 分类分页查询
	PageQuery(ctx context.Context, dto request.CategoryPageQueryDTO) (*result.PageResult, error)
	// List 根据类型查询分类
	List(ctx context.Context, cate int) ([]*response.CategoryVO, error)
	// DeleteById 根据 id 删除分类
	DeleteById(ctx context.Context, id uint64) error
	// AddCategory 新增分类
	AddCategory(ctx context.Context, dto request.CategoryDTO) error
	// SetStatus 启用、禁用分类
	SetStatus(ctx context.Context, id uint64, status int) error
}

type CategoryService struct {
	repo *repository.CategoryDao
}

func NewCategoryService(repo *repository.CategoryDao) ICategoryService {
	return &CategoryService{
		repo: repo,
	}
}

// Update 修改分类
func (s *CategoryService) Update(ctx context.Context, dto request.CategoryDTO) error {

	// 更新分类
	category := &model.Category{
		ID:   dto.Id,
		Type: dto.Type,
		Name: dto.Name,
		Sort: dto.Sort,
	}

	if err := s.repo.Update(ctx, category); err != nil {
		logger.Log.Error("Update: 更新分类失败", zap.Uint64("id", dto.Id), zap.Error(err))
		return errors.New("更新分类失败")
	}

	logger.Log.Info("Update: 更新分类成功", zap.Uint64("id", dto.Id))
	return nil
}

// SetStatus 启用、禁用分类
func (s *CategoryService) SetStatus(ctx context.Context, id uint64, status int) error {
	err := s.repo.SetStatus(ctx, &model.Category{
		ID:     id,
		Status: status,
	})
	if err != nil {
		logger.Log.Error("SetStatus: 设置分类状态失败", zap.Uint64("id", id), zap.Int("status", status), zap.Error(err))
		return errors.New("设置分类状态失败")
	}

	logger.Log.Info("SetStatus: 设置分类状态成功", zap.Uint64("id", id), zap.Int("status", status))
	return nil
}

// Delete 删除分类
func (s *CategoryService) DeleteById(ctx context.Context, id uint64) error {

	// TODO: 检查是否有关联的菜品或套餐
	// 这里暂时直接删除，后续添加菜品和套餐模块后再完善
	if err := s.repo.DeleteById(ctx, id); err != nil {
		logger.Log.Error("Delete: 删除分类失败", zap.Uint64("id", id), zap.Error(err))
		return errors.New("删除分类失败")
	}

	logger.Log.Info("Delete: 删除分类成功", zap.Uint64("id", id))
	return nil
}

// GetByID 根据ID查询分类
func (s *CategoryService) GetByID(ctx context.Context, id uint64) (*response.CategoryVO, error) {
	category, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("分类不存在")
		}
		return nil, errors.New("查询分类失败")
	}

	return &response.CategoryVO{
		ID:         category.ID,
		Type:       category.Type,
		Name:       category.Name,
		Sort:       category.Sort,
		Status:     category.Status,
		CreateTime: category.CreateTime,
		UpdateTime: category.UpdateTime,
	}, nil
}

// List 查询分类列表
func (s *CategoryService) List(ctx context.Context, cate int) ([]*response.CategoryVO, error) {
	categories, err := s.repo.List(ctx, &cate)
	if err != nil {
		logger.Log.Error("List: 查询分类列表失败", zap.Error(err))
		return nil, errors.New("查询分类列表失败")
	}

	result := make([]*response.CategoryVO, 0, len(categories))
	for _, cat := range categories {
		result = append(result, &response.CategoryVO{
			ID:         cat.ID,
			Type:       cat.Type,
			Name:       cat.Name,
			Sort:       cat.Sort,
			Status:     cat.Status,
			CreateTime: cat.CreateTime,
			UpdateTime: cat.UpdateTime,
		})
	}

	return result, nil
}

// PageQuery 分页查询分类
func (s *CategoryService) PageQuery(ctx context.Context, dto request.CategoryPageQueryDTO) (*result.PageResult, error) {
	pageResult, err := s.repo.Page(ctx, dto)
	if err != nil {
		logger.Log.Error("PageQuery: 查询分类失败", zap.Error(err))
		return nil, errors.New("查询分类失败")
	}

	categories, ok := pageResult.Records.([]model.Category)
	if !ok {
		return nil, errors.New("数据类型转换失败")
	}
	// 将结果转换为 VO
	var categoryVOs []response.CategoryVO
	for _, cat := range categories {
		categoryVOs = append(categoryVOs, response.CategoryVO{
			ID:         cat.ID,
			Type:       cat.Type,
			Name:       cat.Name,
			Sort:       cat.Sort,
			Status:     cat.Status,
			CreateTime: cat.CreateTime,
			UpdateTime: cat.UpdateTime,
			CreateUser: cat.CreateUser,
			UpdateUser: ctx.Value("userID").(uint64),
		})
	}

	pageResult.Records = categoryVOs

	return pageResult, nil
}

func (s *CategoryService) AddCategory(ctx context.Context, dto request.CategoryDTO) error {
	// 检查分类名称是否已存在
	exists, err := s.repo.CheckNameExists(ctx, dto.Name, dto.Type)
	if err != nil {
		logger.Log.Error("Add: 检查分类名称失败", zap.Error(err))
		return errors.New("检查分类名称失败")
	}
	if exists {
		return errors.New("分类名称已存在")
	}

	// 新增分类
	err = s.repo.Insert(ctx, &model.Category{
		Name:   dto.Name,
		Sort:   dto.Sort,
		Type:   dto.Type,
		Status: enum.ENABLE,
	})
	if err != nil {
		logger.Log.Error("Add: 新增分类失败", zap.Error(err))
		return errors.New("新增分类失败")
	}

	return nil
}
