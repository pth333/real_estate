package initialize

import (
	"os"
	"path/filepath"
	"real_estate_be/internal/global"

	"github.com/spf13/viper"
)

func LoadConfig() {
	v := viper.New()

	// Allow loading config from different working directories.
	v.AddConfigPath("./config")
	v.AddConfigPath("../config")
	v.AddConfigPath("../../config")

	if exe, err := os.Executable(); err == nil {
		exeDir := filepath.Dir(exe)
		v.AddConfigPath(filepath.Join(exeDir, "config"))
		v.AddConfigPath(filepath.Join(exeDir, "..", "config"))
		v.AddConfigPath(filepath.Join(exeDir, "..", "..", "config"))
	}

	v.SetConfigName("config_local")
	v.SetConfigType("yaml")

	err := v.ReadInConfig()
	if err != nil {
		panic(err)
	}

	if err := v.Unmarshal(&global.Config); err != nil {
		panic(err)
	}
}
