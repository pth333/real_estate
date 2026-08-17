package recommendation

import (
	"context"
	"fmt"
	"strconv"
	"time"

	pb "real_estate_be/proto/recommendation"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Client đại diện cho gRPC client kết nối tới Recommendation Service
type Client struct {
	conn   *grpc.ClientConn
	client pb.RecommendationServiceClient
}

// NewClient khởi tạo kết nối gRPC và trả về đối tượng Client
func NewClient(addr string) (*Client, error) {
	// Thiết lập các option cho kết nối gRPC (Insecure kết nối nội bộ)
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	// Sử dụng gRPC NewClient hiện đại (không chặn, tự động kết nối lại)
	conn, err := grpc.NewClient(addr, opts...)
	if err != nil {
		return nil, fmt.Errorf("không thể khởi tạo kết nối gRPC tới %s: %v", addr, err)
	}

	client := pb.NewRecommendationServiceClient(conn)

	return &Client{
		conn:   conn,
		client: client,
	}, nil
}

// Close đóng kết nối gRPC an toàn
func (c *Client) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// GetRecommendations gọi API gRPC để lấy danh sách gợi ý bất động sản
func (c *Client) GetRecommendations(
	ctx context.Context,
	userID uint64,
	sessionID string,
	realEstateID uint64,
	lat float64,
	lon float64,
	limit int32,
) ([]uint64, string, error) {
	// Chuẩn bị các tham số cho request gRPC
	req := &pb.RecommendRequest{
		UserId:       "",
		SessionId:    sessionID,
		RealEstateId: "",
		Latitude:     lat,
		Longitude:    lon,
		Limit:        limit,
	}

	if userID > 0 {
		req.UserId = strconv.FormatUint(userID, 10)
	}

	if realEstateID > 0 {
		req.RealEstateId = strconv.FormatUint(realEstateID, 10)
	}

	// Gọi gRPC Server với timeout là 500ms
	grpcCtx, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
	defer cancel()

	resp, err := c.client.GetRecommendations(grpcCtx, req)
	if err != nil {
		return nil, "", fmt.Errorf("lỗi gọi gRPC GetRecommendations: %v", err)
	}

	// Parse danh sách ID trả về từ string sang uint64
	propertyIDs := make([]uint64, 0, len(resp.GetPropertyIds()))
	for _, idStr := range resp.GetPropertyIds() {
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err == nil {
			propertyIDs = append(propertyIDs, id)
		}
	}

	return propertyIDs, resp.GetStrategy(), nil
}
