//go:build wireinject
// +build wireinject

package wire

import (
	"real_estate_be/internal/controller"
	"real_estate_be/internal/repo"
	"real_estate_be/internal/usecase"

	"github.com/google/wire"
)

func InitializeAIHandler() (*controller.AIHandler, error) {
	wire.Build(
		repo.NewAIRepository,
		usecase.NewAIService,
		controller.NewAIHandler,
	)
	return &controller.AIHandler{}, nil
}
