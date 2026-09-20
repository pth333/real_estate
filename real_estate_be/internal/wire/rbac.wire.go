//go:build wireinject
// +build wireinject

package wire

import (
	"real_estate_be/internal/controller"
	"real_estate_be/internal/repo"
	"real_estate_be/internal/usecase"

	"github.com/google/wire"
)

// rbacProviderSet — provider cho quản trị role/permission.
var rbacProviderSet = wire.NewSet(
	providerDB,
	repo.NewRbacRepository,
	usecase.NewRbacService,
)

func InitializeAdminRbacHandler() (*controller.AdminRbacHandler, error) {
	wire.Build(
		rbacProviderSet,
		controller.NewAdminRbacHandler,
	)
	return &controller.AdminRbacHandler{}, nil
}
