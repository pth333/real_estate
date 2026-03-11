package global

import "gorm.io/gorm"

var (
	Config ConfigSettings
	DB     *gorm.DB
)

type ConfigSettings struct {
	Server ServerConfig `mapstructure:"server"`
	Mysql  MysqlConfig  `mapstructure:"mysql"`
}

type ServerConfig struct {
	Port int `mapstructure:"port"`
}

type MysqlConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	UserName string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
}
