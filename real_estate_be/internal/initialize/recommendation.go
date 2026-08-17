package initialize

import (
	"log"
	"real_estate_be/internal/global"
	"real_estate_be/pkg/recommendation"
)

// InitRecommendation khởi tạo gRPC client kết nối tới Recommendation Service
func InitRecommendation() {
	addr := global.Config.Recommendation.Addr
	if addr == "" {
		addr = "localhost:50051"
	}

	client, err := recommendation.NewClient(addr)
	if err != nil {
		log.Printf("[gRPC] Khởi tạo kết nối tới Recommendation Service (%s) thất bại: %v", addr, err)
		return
	}

	global.RecommendationClient = client
	log.Printf("[gRPC] Khởi tạo kết nối thành công tới Recommendation Service (%s) (Background)", addr)
}
