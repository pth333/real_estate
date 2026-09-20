package global

import (
	"real_estate_be/internal/sse"
	"real_estate_be/pkg/recommendation"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/redis/go-redis/v9"

	"gorm.io/gorm"
)

var (
	Config               ConfigSettings
	DB                   *gorm.DB
	SSEHub               *sse.Hub
	S3Client             *s3.Client
	RedisClient          *redis.Client
	RecommendationClient *recommendation.Client
)

type ConfigSettings struct {
	Server         ServerConfig         `mapstructure:"server"`
	Mysql          MysqlConfig          `mapstructure:"mysql"`
	Kafka          KafkaConfig          `mapstructure:"kafka"`
	Redis          RedisConfig          `mapstructure:"redis"`
	R2             R2Config             `mapstructure:"r2"`
	AI             AIConfig             `mapstructure:"ai"`
	Recommendation RecommendationConfig `mapstructure:"recommendation"`
	Infobip        InfobipConfig        `mapstructure:"infobip"`
	Payment        PaymentConfig        `mapstructure:"payment"`
	Deposit        DepositConfig        `mapstructure:"deposit"`
	Admin          AdminConfig          `mapstructure:"admin"`
}

// PaymentConfig — cấu hình cổng thanh toán cho luồng đặt cọc escrow
type PaymentConfig struct {
	// URL FE nhận kết quả thanh toán (return URL của cổng)
	ReturnURL string `mapstructure:"return_url"`
	// URL backend nhận IPN (cổng gọi server-to-server)
	IPNURL string       `mapstructure:"ipn_url"`
	VNPay  VNPaySetting `mapstructure:"vnpay"`
}

type VNPaySetting struct {
	TmnCode       string `mapstructure:"tmn_code"`
	HashSecret    string `mapstructure:"hash_secret"`
	PaymentURL    string `mapstructure:"payment_url"`
	Locale        string `mapstructure:"locale"`
	ExpireMinutes int    `mapstructure:"expire_minutes"`
}

// DepositConfig — các mốc thời gian nghiệp vụ đặt cọc (đọc từ plan mục 2 & 6)
type DepositConfig struct {
	// Số tiền cọc mặc định (VNĐ) khi FE không gửi lên
	DefaultAmount float64 `mapstructure:"default_amount"`
	// Phí môi giới mặc định (VNĐ)
	DefaultBrokerFee float64 `mapstructure:"default_broker_fee"`
	// Thời hạn môi giới phải xác nhận lịch (giờ)
	BrokerConfirmHours int `mapstructure:"broker_confirm_hours"`
	// Hiệu lực OTP check-in (phút)
	OTPValidMinutes int `mapstructure:"otp_valid_minutes"`
	// Thời gian ân hạn sau viewing_start trước khi mở cửa sổ báo cáo (giờ)
	CheckinGraceHours int `mapstructure:"checkin_grace_hours"`
	// Cửa sổ để 2 bên tự báo cáo kết quả (giờ)
	ReportWindowHours int `mapstructure:"report_window_hours"`
	// Hạn upload bằng chứng khi có dispute (giờ)
	DisputeEvidenceHours int `mapstructure:"dispute_evidence_hours"`
	// Thời hạn thanh toán của 1 deposit trước khi tự huỷ (phút)
	PaymentTimeoutMinutes int `mapstructure:"payment_timeout_minutes"`
	// Chu kỳ chạy cron (giây)
	CronIntervalSeconds int `mapstructure:"cron_interval_seconds"`
}

// AdminConfig — email tài khoản được nâng quyền ADMIN khi khởi động
type AdminConfig struct {
	Email string `mapstructure:"email"`
}

type InfobipConfig struct {
	ApiKey string `mapstructure:"api_key"`
}

type RecommendationConfig struct {
	Addr string `mapstructure:"addr"`
}

type AIConfig struct {
	APIKey string `mapstructure:"api_key"`
	Model  string `mapstructure:"model"`
}

type R2Config struct {
	Endpoint        string `mapstructure:"endpoint"`
	Region          string `mapstructure:"region"`
	AccessKeyID     string `mapstructure:"access_key_id"`
	SecretAccessKey string `mapstructure:"secret_access_key"`
	Bucket          string `mapstructure:"bucket"`
	PublicURL       string `mapstructure:"public_url"`
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
	Brokers     []string    `mapstructure:"brokers"`
	ClientID    string      `mapstructure:"client_id"`
	GroupPrefix string      `mapstructure:"group_prefix"`
	Topics      KafkaTopics `mapstructure:"topics"`
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
