//go:build wireinject
// +build wireinject

package wire

import (
	"real_estate_be/internal/controller"
	"real_estate_be/internal/repo"
	"real_estate_be/internal/usecase"

	"github.com/google/wire"
)

func InitializeTrackingHandler() (*controller.TrackingHandler, error) {
	wire.Build(
		providerDB,
		repo.NewSearchHistoryRepository,
		repo.NewViewHistoryRepository,
		usecase.NewTrackingService,
		controller.NewTrackingHandler,
	)
	return &controller.TrackingHandler{}, nil
}
