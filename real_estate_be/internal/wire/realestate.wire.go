//go:build wireinject
// +build wireinject

package wire

import (
	"real_estate_be/internal/controller"
	"real_estate_be/internal/repo"
	"real_estate_be/internal/usecase"

	"github.com/google/wire"
)

func InitializeRealEstateHandler() (*controller.RealEstateHandler, error) {
	wire.Build(
		providerDB,
		repo.NewRealEstateRepository,
		repo.NewCategoryRepository,
		usecase.NewRealEstateService,
		controller.NewRealEstateHandler,
	)
	return &controller.RealEstateHandler{}, nil
}
