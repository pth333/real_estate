package wire

import (
	"real_estate_be/internal/global"
	"real_estate_be/pkg/sms"

	"gorm.io/gorm"
)

func providerDB() *gorm.DB {
	return global.DB
}

func providerSMS() sms.Provider {
	return sms.NewConsoleProvider()
}
