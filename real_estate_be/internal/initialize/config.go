package initialize

import (
	"os"
	"path/filepath"
	"real_estate_be/internal/global"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// LoadConfig đọc cấu hình từ:
//  1. File .env (godotenv load vào os.Environ — secret KHÔNG commit lên git)
//  2. File config/config_local.yaml (cấu hình chung, không chứa secret)
//  3. Biến môi trường hệ thống (ưu tiên cao nhất, prefix RE_)
func LoadConfig() {
	// ── Load file .env ─────────────────────────────
	// godotenv load .env vào os.Environ() thật, để viper.AutomaticEnv() đọc được
	loadEnvFiles()

	v := viper.New()

	// ── Load file config_local.yaml ────────────────
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

	// ── Ưu tiên biến môi trường hệ thống ──────────
	// RE_AI_API_KEY -> key "ai.api_key", RE_R2_SECRET_ACCESS_KEY -> "r2.secret_access_key"
	// AutomaticEnv đọc từ os.Getenv — godotenv đã load .env vào đây rồi
	v.SetEnvPrefix("RE")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// BindEnv cho các key KHÔNG có trong file yaml (chỉ có trong .env / env hệ thống),
	// nếu không AutomaticEnv sẽ bỏ qua chúng.
	v.BindEnv("ai.api_key", "RE_AI_API_KEY")
	v.BindEnv("r2.endpoint", "RE_R2_ENDPOINT")
	v.BindEnv("r2.access_key_id", "RE_R2_ACCESS_KEY_ID")
	v.BindEnv("r2.secret_access_key", "RE_R2_SECRET_ACCESS_KEY")
	v.BindEnv("r2.public_url", "RE_R2_PUBLIC_URL")

	if err := v.Unmarshal(&global.Config); err != nil {
		panic(err)
	}
}

// loadEnvFiles tìm và load file .env từ nhiều thư mục làm việc.
// Nếu không tìm thấy thì bỏ qua (secret có thể đến từ env hệ thống).
func loadEnvFiles() {
	// Đường dẫn có khả năng chứa .env nhất, theo thứ tự ưu tiên
	paths := []string{}

	// Thư mục làm việc hiện tại (tăng dần cấp để tìm .env ở root project)
	if cwd, err := os.Getwd(); err == nil {
		for depth := 0; depth <= 3; depth++ {
			dir := cwd
			for i := 0; i < depth; i++ {
				dir = filepath.Dir(dir)
			}
			paths = append(paths, filepath.Join(dir, ".env"))
		}
	}

	// Thư mục của executable
	if exe, err := os.Executable(); err == nil {
		exeDir := filepath.Dir(exe)
		paths = append(paths,
			filepath.Join(exeDir, ".env"),
			filepath.Join(exeDir, "..", ".env"),
			filepath.Join(exeDir, "..", "..", ".env"),
		)
	}

	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			// Load bổ sung, không ghi đè env đã có sẵn
			_ = godotenv.Load(p)
		}
	}
}
