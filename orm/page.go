package orm

import (
	"errors"

	"gorm.io/gorm"
)

type PageResponse[T any] struct {
	Total   int64 `json:"total"`   // 总记录数
	Current int64 `json:"current"` // 当前页码，从 1 开始
	Size    int64 `json:"size"`    // 每页条数
	Data    []T   `json:"data"`    // 当前页数据列表
}

func Paginate[T any](db *gorm.DB, page int64, size int64, scopes ...func(db *gorm.DB) *gorm.DB) (PageResponse[T], error) {
	// 分页参数常量配置
	const (
		defaultPage = 1
		defaultSize = 10
		maxSize     = 1000
	)

	// 防御性检查：nil 检查
	if db == nil {
		return PageResponse[T]{}, errors.New("db is nil")
	}

	// 页码容错
	if page <= 0 {
		page = defaultPage
	}
	// 每页条数限制
	switch {
	case size <= 0:
		size = defaultSize
	case size > maxSize:
		size = maxSize
	}
	// ========== 1. 查询总数 ==========
	var total int64
	if err := db.Scopes(scopes...).Count(&total).Error; err != nil {
		return PageResponse[T]{}, err
	}
	// 无数据直接返回空列表，避免无用查询
	if total == 0 {
		return PageResponse[T]{
			Total:   total,
			Current: page,
			Size:    size,
			Data:    make([]T, 0),
		}, nil
	}
	// ========== 2. 分页查询数据 ==========
	offset := (page - 1) * size
	var data []T
	if err := db.Scopes(scopes...).Limit(int(size)).Offset(int(offset)).Find(&data).Error; err != nil {
		return PageResponse[T]{}, err
	}
	return PageResponse[T]{
		Total:   total,
		Current: page,
		Size:    size,
		Data:    data,
	}, nil
}
