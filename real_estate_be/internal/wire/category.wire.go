//go:build wireinject
// +build wireinject

package wire

import (
	"real_estate_be/internal/controller"
	"real_estate_be/internal/repo"
	"real_estate_be/internal/usecase"

	"github.com/google/wire"
)

func InitializeCategoryHandler() (*controller.CategoryHandler, error) {
	wire.Build(
		providerDB,
		repo.NewCategoryRepository,
		usecase.NewCategoryService,
		controller.NewCategoryHandler,
	)
	return &controller.CategoryHandler{}, nil
}
