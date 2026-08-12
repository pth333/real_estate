//go:build wireinject
// +build wireinject

package wire

import (
	"real_estate_be/internal/controller"
	"real_estate_be/internal/repo"
	"real_estate_be/internal/usecase"

	"github.com/google/wire"
)

func InitializeNotificationHandler() (*controller.NotificationHandler, error) {
	wire.Build(
		providerDB,
		repo.NewNotificationRepository,
		usecase.NewNotificationService,
		controller.NewNotificationHandler,
	)
	return &controller.NotificationHandler{}, nil
}
