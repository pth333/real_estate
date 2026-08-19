//go:build wireinject
// +build wireinject

package wire

import (
	"real_estate_be/internal/controller"
	"real_estate_be/internal/repo"
	"real_estate_be/internal/usecase"
	"real_estate_be/pkg/kafka"

	"github.com/google/wire"
)

func InitializeManagerPostHandler() (*controller.ManagerPostHandler, error) {
	wire.Build(
		providerDB,
		repo.NewManagerPostRepository,
		repo.NewRealEstateRepository, // Thêm dòng này
		repo.NewImageRepository,
		repo.NewUserRepository,
		kafka.NewProducer,
		usecase.NewManagerPostUseCase,
		controller.NewManagerPostHandler,
	)
	return &controller.ManagerPostHandler{}, nil
}
