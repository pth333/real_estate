//go:build wireinject
// +build wireinject

package wire

import (
	"real_estate_be/internal/controller"
	"real_estate_be/internal/repo"
	"real_estate_be/internal/usecase"

	"github.com/google/wire"
)

func InitializeManagerPostHandler() (*controller.ManagerPostHandler, error) {
	wire.Build(
		providerDB,
		repo.NewManagerPostRepository,
		repo.NewUserRepository,
		usecase.NewManagerPostUseCase,
		controller.NewManagerPostHandler,
	)
	return &controller.ManagerPostHandler{}, nil
}
