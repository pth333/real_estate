package usecase

import (
	model "real_estate_be/internal/models"
	"real_estate_be/internal/repo"
)

type INotificationService interface {
	GetNotifications() ([]model.Notification, error)
}

type NotificationService struct {
	repo repo.INotificationRepository
}

func NewNotificationService(repo repo.INotificationRepository) INotificationService {
	return &NotificationService{repo: repo}
}

func (s *NotificationService) GetNotifications() ([]model.Notification, error) {
	return s.repo.GetLatest(50)
}
