package models

import (
	"PracticeProject/define"
	"PracticeProject/helper"
	"gorm.io/gorm"
)

// 保存数据的结构体名默认与表名对应
type ProblemBasic struct {
	gorm.Model
	Identity          string             `gorm:"column:identity;type:varchar(36);"json:"identity"`                   // 问题表的唯一标识
	ProblemCategories []*ProblemCategory `gorm:"foreignKey:problem_id;references:id" json:"problem_categories"`      // 关联问题分类表
	Title             string             `gorm:"column:title;type:varchar(255);"json:"title"`                        // 文章标题
	Content           string             `gorm:"column:content;type:text;"json:"content"`                            // 文章正文
	MaxMem            int                `gorm:"column:max_mem;type:int(11);"json:"max_mem"`                         // 最大运行内存（kb）
	MaxRuntime        int                `gorm:"column:max_runtime;type:int(11)"json:"max_Runtime"`                  // 最大运行时间（ms）
	PassNum           int64              `gorm:"column:pass_num;type:int(11);" json:"pass_num"`                      // 通过的次数
	SubmitNum         int64              `gorm:"column:submit_num;type:int(11);" json:"submit_num"`                  // 提交次数
	TestCases         []*TestCase        `gorm:"foreignKey:problem_identity;references:identity;" json:"test_cases"` // 关联测试用例表
}

func (table *ProblemBasic) TableName() string {
	return "problem_basic"
}

func GetProblemList(startPoint int, size int, keyword string, categoryIdentity string) ([]*ProblemBasic, int64, error) {
	var count int64
	list := make([]*ProblemBasic, 0) // 保存problemlist
	// 查询表
	tx := DB.Model(new(ProblemBasic)).Preload("ProblemCategories").Preload("ProblemCategories.CategoryBasic"). // 指定要查询的表
															Where("title like ? OR content like ?", "%"+keyword+"%", "%"+keyword+"%") // 查询条件
	var err error = nil
	if categoryIdentity != "" {
		tx.Joins("RIGHT JOIN problem_category pc on pc.problem_id = problem_basic.id"). // 根据外键关联表获得包含指定分类id的问题列（问题表<右关联>问题分类表<右关联>分类表）
												Where("pc.category_id = (SELECT cb.id FROM category_basic cb WHERE cb.identity = ? )", categoryIdentity)
	}
	err = tx.Count(&count).Offset(startPoint).Limit(size).Find(&list).Error
	return list, count, err
	/*data := make([]*Problem, 0)
	DB.Find(&data)
	for _, v := range data {
		fmt.Printf("Problem ==> %v \n", v)
	}*/
}

func GetProblemDetail(identity string) (*gorm.DB, *ProblemBasic) {
	detail := new(ProblemBasic)
	tx := DB.Where("identity = ?", identity).Preload("ProblemCategories").Preload("ProblemCategories.CategoryBasic").First(detail)
	return tx, detail
}

func ProblemCreate(data *ProblemBasic) error {
	err := DB.Create(data).Error
	return err
}

func GetProblemCount(identity string) (int64, error) {
	var cnt int64
	err := DB.Model(new(ProblemBasic)).Where("identity = ?", identity).Count(&cnt).Error
	return cnt, err
}

func ProblemModify(in *define.ProblemBasic) error {
	err := DB.Transaction(func(tx *gorm.DB) error {
		// 问题基础信息保存 problem_basic
		problemBasic := &ProblemBasic{
			Identity:   in.Identity,
			Title:      in.Title,
			Content:    in.Content,
			MaxRuntime: in.MaxRuntime,
			MaxMem:     in.MaxMem,
		}
		err := tx.Where("identity = ?", in.Identity).Updates(problemBasic).Error
		if err != nil {
			return err
		}
		// 查询问题详情
		err = tx.Where("identity = ?", in.Identity).Find(problemBasic).Error
		if err != nil {
			return err
		}

		// 关联问题分类的更新
		// 1、删除已存在的关联关系
		err = tx.Where("problem_id = ?", problemBasic.ID).Delete(new(ProblemCategory)).Error
		if err != nil {
			return err
		}
		// 2、新增新的关联关系
		pcs := make([]*ProblemCategory, 0)
		for _, id := range in.ProblemCategories {
			pcs = append(pcs, &ProblemCategory{
				ProblemId:  problemBasic.ID,
				CategoryId: uint(id),
			})
		}
		err = tx.Create(&pcs).Error
		if err != nil {
			return err
		}
		// 关联测试案例的更新
		// 1、删除已存在的关联关系
		err = tx.Where("problem_identity = ?", in.Identity).Delete(new(TestCase)).Error
		if err != nil {
			return err
		}
		// 2、增加新的关联关系
		tcs := make([]*TestCase, 0)
		for _, v := range in.TestCases {
			// 举个例子 {"input":"1 2\n","output":"3\n"}
			tcs = append(tcs, &TestCase{
				Identity:        helper.GetUUID(),
				ProblemIdentity: in.Identity,
				Input:           v.Input,
				Output:          v.Output,
			})
		}
		err = tx.Create(tcs).Error
		if err != nil {
			return err
		}
		return nil
	})

	return err
}

func GetProblemTestCase(problemIdentity string) (*ProblemBasic, error) {
	problemBasic := new(ProblemBasic)
	return problemBasic, DB.Where("identity = ?", problemIdentity).Preload("TestCases").First(&problemBasic).Error
}
