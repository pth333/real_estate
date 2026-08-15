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

func InitializeRealEstateHandler() (*controller.RealEstateHandler, error) {
	wire.Build(
		providerDB,
		repo.NewRealEstateRepository,
		repo.NewCategoryRepository,
		repo.NewImageRepository,
		repo.NewUserRepository,
		repo.NewSearchHistoryRepository,
		kafka.NewProducer,
		usecase.NewRealEstateService,
		controller.NewRealEstateHandler,
	)
	return &controller.RealEstateHandler{}, nil
}
