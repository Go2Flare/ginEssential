package common

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"oceanlearn.learn/ginessential/model"
)

//声明了要记得初始化
var DB *gorm.DB

func InitDB() *gorm.DB {
	host := "localhost"
	port := "3306"
	database := "ginessential"
	username := "root"
	password := "4.234.23123"
	charset := "utf8"
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true",
		username,
		password,
		host,
		port,
		database,
		charset)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database, err:" + err.Error())
	}

	//按结构体User自动合并到表User
	db.AutoMigrate(&model.User{})
	//赋值
	DB = db
	return db
}

func GetDB() *gorm.DB{
	return DB
}
