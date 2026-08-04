package initialize

import (
	"fmt"
	"real_estate_be/internal/global"
	model "real_estate_be/internal/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitMysql() {
	m := global.Config.Mysql
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		m.Username,
		m.Password,
		m.Host,
		m.Port,
		m.DBName,
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	// db.AutoMigrate(&model.User{})
	// db.AutoMigrate(&model.Category{})
	// db.AutoMigrate(&model.Image{})
	// db.AutoMigrate(&model.RealEstate{})
	db.AutoMigrate(&model.SearchHistory{})

	if err != nil {
		panic(err)
	}

	global.DB = db
}
