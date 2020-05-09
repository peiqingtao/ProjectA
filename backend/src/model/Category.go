package model

import (
	"github.com/jinzhu/gorm"
	"time"
)


// 模型结构体(类)定义
type Category struct {
	// 嵌套结构体
	gorm.Model
	ParentId uint
	Name string
	Logo string
	Description string
	SortOrder int
	MetaTitle string
	MetaKeywords string
	MetaDescription string
	CreatedAt *time.Time
	DeletedAt *time.Time
	UpdatedAt *time.Time
	// has many
	//Products []Product // Product.CategoryID

	// has one

	// belongs to

	// many to many
}