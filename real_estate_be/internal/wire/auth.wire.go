//go:build wireinject
// +build wireinject

package wire

import (
	"real_estate_be/internal/controller"
	"real_estate_be/internal/repo"
	"real_estate_be/internal/usecase"

	"github.com/google/wire"
)

func InitializeAuthHandler() (*controller.UserHandler, error) {
	// Auth
	wire.Build(
		providerDB,
		repo.NewUserRepository,
		usecase.NewAuthService,
		controller.NewUserHandler,
	)
	return &controller.UserHandler{}, nil
}
