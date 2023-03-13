package models

import (
	"PracticeProject/helper"
	"gorm.io/gorm"
)

type CategoryBasic struct {
	gorm.Model
	Identity string `gorm:"column:identity;type:varchar(36);"json:"identity"` // 分类的唯一标识
	Name     string `gorm:"column:name;type:varchar(100);"json:"name"`        // 分类名称
	ParentId int    `gorm:"column:parent_id;type:int(11);"json:"parent_id"`   // 父级id
}

func (table *CategoryBasic) TableName() string {
	return "category_basic"
}

func GetCategoryList(startPoint int, size int, keyword string) ([]*CategoryBasic, int64, error) {
	var count int64
	categoryList := make([]*CategoryBasic, 0)
	err := DB.Model(new(CategoryBasic)).Where("name like ?", "%"+keyword+"%").
		Count(&count).Limit(size).Offset(startPoint).Find(&categoryList).Error
	return categoryList, count, err
}

func CategoryCreate(name string, parentId int) error {
	category := &CategoryBasic{
		Identity: helper.GetUUID(),
		Name:     name,
		ParentId: parentId,
	}
	err := DB.Create(category).Error
	return err
}

func CategoryModify(identity string, name string, parentId int) error {
	category := &CategoryBasic{
		Identity: identity,
		Name:     name,
		ParentId: parentId,
	}
	err := DB.Model(new(CategoryBasic)).Where("identity = ?", identity).Updates(category).Error
	return err
}

func GetCategoryCount(identity string) (int64, error) {
	var cnt int64
	err := DB.Model(new(CategoryBasic)).Where("identity = ?", identity).Count(&cnt).Error
	return cnt, err
}
func CountProblemByCategory(identity string) (int64, error) {
	var cnt int64
	err := DB.Model(new(ProblemCategory)).Where("category_id = (SELECT id FROM category_basic WHERE identity = ? LIMIT 1)", identity).Count(&cnt).Error
	return cnt, err
}

func CategoryDelete(identity string) error {
	err := DB.Where("identity = ?", identity).Delete(new(CategoryBasic)).Error
	return err
}
