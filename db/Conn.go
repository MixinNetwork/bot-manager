package db

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var Conn *gorm.DB

func Connect() (err error) {
	host := beego.AppConfig.String("dbHost")
	user := beego.AppConfig.String("dbUser")
	database := beego.AppConfig.String("dbName")
	pass := beego.AppConfig.String("dbPass")
	dbConnString := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", host, user, database, pass)
	db, err := gorm.Open("postgres", dbConnString)
	Conn = db
	if err != nil {
		return err
	}
	return nil
}
