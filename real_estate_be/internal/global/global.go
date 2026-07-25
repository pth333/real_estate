package global

import (
	"real_estate_be/internal/sse"

	"gorm.io/gorm"
)

var (
	Config ConfigSettings
	DB     *gorm.DB
	SSEHub *sse.Hub
)

type ConfigSettings struct {
	Server ServerConfig `mapstructure:"server"`
	Mysql  MysqlConfig  `mapstructure:"mysql"`
	Kafka  KafkaConfig  `mapstructure:"kafka"`
	Redis  RedisConfig  `mapstructure:"redis"`
}

type ServerConfig struct {
	Port int `mapstructure:"port"`
}

type MysqlConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"database"`
}

type KafkaConfig struct {
	Brokers     []string      `mapstructure:"brokers"`
	ClientID    string        `mapstructure:"client_id"`
	GroupPrefix string        `mapstructure:"group_prefix"`
	Topics      KafkaTopics   `mapstructure:"topics"`
}

type KafkaTopics struct {
	RealEstateCrawled  string `mapstructure:"real_estate_crawled"`
	RealEstateEnriched string `mapstructure:"real_estate_enriched"`
	RealEstateNotified string `mapstructure:"real_estate_notified"`
}

type RedisConfig struct {
	Addr string `mapstructure:"addr"`
	DB   int    `mapstructure:"db"`
}
