package wire

import (
	"real_estate_be/internal/global"

	"gorm.io/gorm"
)

func providerDB() *gorm.DB {
	return global.DB
}
