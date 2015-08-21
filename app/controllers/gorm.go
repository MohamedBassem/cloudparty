package controllers

import (
	"fmt"

	"github.com/MohamedBassem/cloudparty/app/models"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/revel/revel"
)

var DB *gorm.DB
var Test string

func init() {
	revel.OnAppStart(InitDB)
}

func InitDB() {
	databaseHost := fmt.Sprintf("%v:%v@tcp(%v)/%v?charset=utf8&parseTime=True&loc=Local", revel.Config.StringDefault("db.user", ""), revel.Config.StringDefault("db.password", ""), revel.Config.StringDefault("db.host", ""), revel.Config.StringDefault("db.name", ""))
	_DB, err := gorm.Open("mysql", databaseHost)
	DB = &_DB
	if err != nil {
		revel.INFO.Println("DB Error", err)
	}
	revel.INFO.Println("DB Connected", err)
	DB.LogMode(true)

	DB.AutoMigrate(&models.Playlist{})
	DB.AutoMigrate(&models.Song{})
}
