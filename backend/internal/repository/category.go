package repository

import (
	"context"
	"takeaway-go/common/result"
	"takeaway-go/internal/api/request"
	"takeaway-go/internal/model"

	"gorm.io/gorm"
)

type CategoryDao struct {
	db *gorm.DB
}

func NewCategoryDao(db *gorm.DB) *CategoryDao {
	return &CategoryDao{
		db: db,
	}
}

// UpdateStatus 更新分类状态
func (r *CategoryDao) SetStatus(ctx context.Context, category *model.Category) error {
	return r.db.WithContext(ctx).
		Model(&model.Category{}).
		Where("id = ?", category.ID).
		Update("status", category.Status).Error
}

// Create 创建分类
func (r *CategoryDao) Create(ctx context.Context, category *model.Category) error {
	return r.db.WithContext(ctx).Create(category).Error
}

// Update 更新分类
func (r *CategoryDao) Update(ctx context.Context, category *model.Category) error {
	return r.db.WithContext(ctx).
		Model(&model.Category{}).
		Where("id = ?", category.ID).
		Updates(category).Error
}

// GetByID 根据ID查询分类
func (r *CategoryDao) GetByID(ctx context.Context, id uint64) (*model.Category, error) {
	var category model.Category
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&category).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

// DeleteById 根据ID删除分类
func (r *CategoryDao) DeleteById(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.Category{}).Error
}

// Insert 新增分类
func (r *CategoryDao) Insert(ctx context.Context, category *model.Category) error {
	return r.db.WithContext(ctx).Create(category).Error
}

// List 查询分类列表（按类型）
func (r *CategoryDao) List(ctx context.Context, categoryType *int) ([]*model.Category, error) {
	var categories []*model.Category

	query := r.db.WithContext(ctx).Model(&model.Category{}).
		Where("status = ?", 1) // 只查询启用的分类

	if categoryType != nil {
		query = query.Where("type = ?", *categoryType)
	}

	if err := query.Order("sort ASC, create_time DESC").
		Find(&categories).Error; err != nil {
		return nil, err
	}

	return categories, nil
}

// Page 分类分页查询
func (r *CategoryDao) Page(ctx context.Context, dto request.CategoryPageQueryDTO) (*result.PageResult, error) {
	var pageResult result.PageResult
	var categoryList []model.Category

	query := r.db.WithContext(ctx).Model(&model.Category{})

	if dto.Name != "" {
		query = query.Where("name LIKE ?", "%"+dto.Name+"%")
	}
	if dto.Type != 0 {
		query = query.Where("type = ?", dto.Type)
	}
	// 计算总数
	err := query.Count(&pageResult.Total).Error
	if err != nil {
		return nil, err
	}
	// 分页查询
	err = query.Scopes(pageResult.Paginate(&dto.Page, &dto.PageSize)).
		Order("create_time desc").
		Find(&categoryList).
		Error
	if err != nil {
		return nil, err
	}

	pageResult.Records = categoryList
	return &pageResult, nil
}

// CheckNameExists 检查分类名称是否存在
func (r *CategoryDao) CheckNameExists(ctx context.Context, name string, categoryType int) (bool, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&model.Category{}).
		Where("name = ? AND type = ?", name, categoryType)

	err := query.Count(&count).Error
	return count > 0, err
}
