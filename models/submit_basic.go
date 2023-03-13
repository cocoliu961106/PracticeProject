package models

import (
	"PracticeProject/helper"
	"errors"
	"gorm.io/gorm"
)

type SubmitBasic struct {
	gorm.Model
	Identity        string        `gorm:"column:identity;type:varchar(36);"json:"identity"`                 // 提交唯一标识
	ProblemIdentity string        `gorm:"column:problem_identity;type:varchar(36);"json:"problem_identity"` // 问题唯一标识
	ProblemBasic    *ProblemBasic `gorm:"foreignKey:identity;references:problem_identity"`                  // 关联问题表
	UserIdentity    string        `gorm:"column:user_identity;type:varchar(36);"json:"user_identity"`       // 用户唯一标识
	UserBasic       *UserBasic    `gorm:"foreignKey:identity;references:user_identity"`                     // 关联用户表
	Path            string        `gorm:"column:path;type:varchar(255);"json:"path"`                        // 代码存放路径
	Status          int           `gorm:"column:status;type:tinyint(1);"json:"status"`                      // [-1-待判断，1-答案正确，2-答案错误， 3-运行超时， 4-运行超内存，5-编译错误]
}

func (table *SubmitBasic) TableName() string {
	return "submit_basic"
}

func GetSubmitList(problemIdentity string, userIdentity string, status int) *gorm.DB {
	tx := DB.Model(new(SubmitBasic)).Preload("ProblemBasic", func(db *gorm.DB) *gorm.DB {
		return db.Omit("content")
	}).Preload("UserBasic")
	if problemIdentity != "" {
		tx.Where("problem_identity = ?", problemIdentity)
	}
	if userIdentity != "" {
		tx.Where("user_identity = ?", userIdentity)
	}
	if status != 0 {
		tx.Where("status = ?", status)
	}
	return tx
}

func Submit(sb *SubmitBasic, userClaim *helper.UserClaims, problemIdentity string) error {
	err := DB.Transaction(func(tx *gorm.DB) error {
		err := tx.Create(sb).Error
		if err != nil {
			return errors.New("SubmitBasic Save Error:" + err.Error())
		}

		// 用户表与问题表的提交次数+1，若提交结果正确，通过次数+1
		m := make(map[string]interface{})
		m["submit_num"] = gorm.Expr("submit_num + ?", 1)
		if sb.Status == 1 {
			m["pass_num"] = gorm.Expr("pass_num + ?", 1)
		}
		// 更新 user_basic
		err = tx.Model(new(UserBasic)).Where("identity = ?", userClaim.Identity).Updates(m).Error
		if err != nil {
			return errors.New("UserBasic Modify Error:" + err.Error())
		}
		// 更新 problem_basic
		err = tx.Model(new(ProblemBasic)).Where("identity = ?", problemIdentity).Updates(m).Error
		if err != nil {
			return errors.New("ProblemBasic Modify Error:" + err.Error())
		}
		return nil
	})
	return err
}
