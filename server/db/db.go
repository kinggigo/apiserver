package db

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

const (
	DB_HOST = ""
	DB_USER = ""
	DB_PW   = "1q2w3e4r5t"
	DB_NAME = ""
	DB_PORT = ""
)

var GormDB *gorm.DB

var (
	dbConnUrl = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True", DB_USER, DB_PW, DB_HOST, DB_PORT, DB_NAME)
)

func GetDBConnURL() string {
	return dbConnUrl
}

//func NewDB() ( error) {
//	db, err := gorm.Open("mysql", dbConnUrl)
//	if err != nil {
//		return nil, err
//	}
//	//defer db.Close()
//	db.LogMode(true)
//	GormDB := &WooriDB{db}
//	return GormDB, nil
//}
