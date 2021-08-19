package model

import "gorm.io/gorm"

//用gorm连接数据库
type User struct {
	gorm.Model
	Name      string `gorm:"type:varchar(20);notnull"`
	Telephone string `gorm:"type:varchar(110);notnull;unique"`
	Password  string `gorm:"size:255;not null"`
}
