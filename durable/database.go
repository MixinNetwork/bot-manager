package durable

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var Conn *gorm.DB

func Connect(host string, database string, user string, pass string) (db *gorm.DB, err error) {
	dbConnString := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", host, user, database, pass)
	db, err = gorm.Open("postgres", dbConnString)
	Conn = db
	return
}

var models []interface{}
var migrations []string

func RegisterModel(model interface{}) {
	models = append(models, model)
}

func RegisterMigration(migration string) {
	migrations = append(migrations, migration)
}

func AutoMigrate() {
	for _, migration := range migrations {
		Conn.Exec(migration)
	}
	Conn.AutoMigrate(models...)
}

func init() {
	conn := getDBConnection()
	defer conn.Close()
}

func getDBConnection() *gorm.DB {
	dbHost := beego.AppConfig.String("dbHost")
	dbUser := beego.AppConfig.String("dbUser")
	dbName := beego.AppConfig.String("dbName")
	dbPass := beego.AppConfig.String("dbPass")
	conn, err := Connect(dbHost, dbName, dbUser, dbPass)
	if err != nil {
		panic(err.Error())
	}
	return conn
}
