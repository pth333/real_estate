//go:build wireinject
// +build wireinject

package wire

import (
	"real_estate_be/internal/controller"
	"real_estate_be/internal/repo"
	"real_estate_be/internal/usecase"

	"github.com/google/wire"
)

// depositProviderSet — provider dùng chung cho mọi handler của luồng đặt cọc.
var depositProviderSet = wire.NewSet(
	providerDB,
	providerPaymentGateway,
	providerMailer,
	repo.NewDepositRepository,
	repo.NewDisputeRepository,
	repo.NewTransactionRepository,
	repo.NewBrokerRatingRepository,
	repo.NewNotificationLogRepository,
	repo.NewDepositPolicyRepository,
	repo.NewRealEstateRepository,
	repo.NewUserRepository,
	usecase.NewDepositService,
)

func InitializeDepositHandler() (*controller.DepositHandler, error) {
	wire.Build(
		depositProviderSet,
		controller.NewDepositHandler,
	)
	return &controller.DepositHandler{}, nil
}

func InitializeAdminDepositHandler() (*controller.AdminDepositHandler, error) {
	wire.Build(
		depositProviderSet,
		controller.NewAdminDepositHandler,
	)
	return &controller.AdminDepositHandler{}, nil
}

// InitializeDepositService — dùng cho scheduler chạy nền ở initialize.Run.
func InitializeDepositService() (usecase.IDepositService, error) {
	wire.Build(depositProviderSet)
	return nil, nil
}
