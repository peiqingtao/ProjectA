package controller

import (
	"config"
	"fmt"
	"github.com/jinzhu/gorm"
	"log"
)

//基础控制器代码

var orm *gorm.DB

func InitDB() (*gorm.DB,error){
	//初始化 gorm
	gorm.DefaultTableNameHandler = func (db *gorm.DB, defaultTableName string) string  {
		return config.App["DB_TABLE_PREFIX"] + defaultTableName
	}
	// 连接数据库
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&loc=%s&parseTime=%s",
		config.App["MYSQL_USER"],
		config.App["MYSQL_PASSWORD"],
		config.App["MYSQL_HOST"],
		config.App["MYSQL_PORT"],
		config.App["MYSQL_DBNAME"],
		config.App["MYSQL_CHARSET"],
		config.App["MYSQL_LOC"],
		config.App["MYSQL_PARSETIME"],
	)
	db, err := gorm.Open(config.App["DB_DRIVER"], dsn)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	orm = db
	return db, nil
}