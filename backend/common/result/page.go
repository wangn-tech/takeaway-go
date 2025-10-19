package result

import (
	"gorm.io/gorm"
	"takeaway-go/common/enum"
)

// PageResult 分页结果
type PageResult struct {
	Total   int64       `json:"total"`   //总记录数
	Records interface{} `json:"records"` //当前页数据集合
}

// PageVerify 分页参数校验, 修正非法参数
func PageVerify(page *int, pageSize *int) {
	// 过滤 当前页、单页数量
	if *page < 1 {
		*page = 1
	}
	switch {
	case *pageSize > 100:
		*pageSize = enum.MaxPageSize
	case *pageSize <= 0:
		*pageSize = enum.MinPageSize
	}
}

// Paginate
//
//	// 分页参数校验
//	// 拼接分页 sql
func (p *PageResult) Paginate(page *int, pageSize *int) func(*gorm.DB) *gorm.DB {
	return func(d *gorm.DB) *gorm.DB {
		// 分页校验
		PageVerify(page, pageSize)

		// 拼接分页
		d.Offset((*page - 1) * *pageSize).Limit(*pageSize)
		return d
	}
}
