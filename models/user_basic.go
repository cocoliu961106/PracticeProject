package models

import (
	"errors"
	"gorm.io/gorm"
)

type UserBasic struct {
	gorm.Model
	Identity  string `gorm:"column:identity;type:varchar(36);"json:"identity"`  // 用户的唯一标识
	Name      string `gorm:"column:name;type:varchar(100);"json:"name"`         // 用户名
	Password  string `gorm:"column:password;type:varchar(32);"json:"password"`  // 密码
	Phone     string `gorm:"column:phone;type:varchar(20);"json:"phone"`        // 电话
	Mail      string `gorm:"column:mail;type:varchar(20);"json:"phone"`         // 邮箱
	PassNum   int64  `gorm:"column:pass_num;type:int(11);" json:"pass_num"`     // 通过的次数
	SubmitNum int64  `gorm:"column:submit_num;type:int(11);" json:"submit_num"` // 提交次数
	IsAdmin   int    `gorm:"column:is_admin;type:tinyint(1);"json:"is_admin"`   // 是否为管理员[0-否，1-是]
}

func (table *UserBasic) TableName() string {
	return "user_basic"
}

func GetUserDetail(identity string) (*UserBasic, error) {
	data := new(UserBasic)
	err := DB.Omit("password").Where("identity = ?", identity).Find(&data).Error
	if data.Identity == "" {
		err = errors.New("不存在该用户")
	}
	return data, err
}

func GetLoginResult(username string, password string) (*UserBasic, error) {
	data := new(UserBasic)
	err := DB.Where("name = ? AND password = ?", username, password).First(&data).Error
	return data, err
}

func FindRepeatEmail(mail string) (int64, error) {
	var cnt int64
	err := DB.Where("mail = ?", mail).Model(new(UserBasic)).Count(&cnt).Error
	return cnt, err
}

func GetRankList(startPoint int, size int) ([]*UserBasic, int64, error) {
	var count int64
	list := make([]*UserBasic, 0)
	err := DB.Model(new(UserBasic)).Count(&count).Order("pass_num DESC, submit_num ASC").
		Offset(startPoint).Limit(size).Find(&list).Error
	return list, count, err
}
