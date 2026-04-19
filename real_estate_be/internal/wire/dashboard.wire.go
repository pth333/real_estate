//go:build wireinject
// +build wireinject

package wire

import (
	"real_estate_be/internal/controller"
	"real_estate_be/internal/repo"
	"real_estate_be/internal/usecase"

	"github.com/google/wire"
)

func InitializeDashboardHandler() (*controller.DashboardHandler, error) {
	// Dashboard
	wire.Build(
		providerDB,
		repo.NewDashboardRepository,
		usecase.NewDashboardService,
		controller.NewDashboardHandler,
	)
	return &controller.DashboardHandler{}, nil
}
