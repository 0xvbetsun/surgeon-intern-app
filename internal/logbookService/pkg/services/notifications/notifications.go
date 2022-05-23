package notifications

import (
	"context"
	"database/sql"
	"time"

	"github.com/friendsofgo/errors"
	commonModel "github.com/vbetsun/surgeon-intern-app/internal/graph/model"
	"github.com/vbetsun/surgeon-intern-app/pkg/models"
	"github.com/vbetsun/surgeon-intern-app/pkg/repo/dbexecutor"
	"github.com/vbetsun/surgeon-intern-app/pkg/repo/notifications"
	"github.com/vbetsun/surgeon-intern-app/pkg/repo/users"
	"github.com/volatiletech/null/v8"
)

type (
	INotificationsService interface {
		GetNotifications(ctx context.Context, notificationType *commonModel.NotificationType, hasSeen bool) ([]*commonModel.Notification, error)
		SetHasSeenNotifications(ctx context.Context, seenNotifications []*commonModel.SeenNotificationInput) ([]*commonModel.Notification, error)
	}

	NotificationsService struct {
		notificationsRepo notifications.INotificationsRepo
		usersRepo         users.IRepo
		dbexecutor.IDBExecutor
	}
)

func NewNotificationsService(notificationsRepo notifications.INotificationsRepo,
	usersRepo users.IRepo, dbExecutor dbexecutor.IDBExecutor) INotificationsService {
	return &NotificationsService{
		notificationsRepo: notificationsRepo,
		usersRepo:         usersRepo,
		IDBExecutor:       dbExecutor,
	}
}

func (s NotificationsService) GetNotifications(ctx context.Context, notificationType *commonModel.NotificationType, hasSeen bool) ([]*commonModel.Notification, error) {
	user, err := s.usersRepo.GetByAuthenticationContext(ctx)
	if err != nil {
		return nil, err
	}

	allNotifications, err := s.notificationsRepo.ListAllNotifications(ctx, user.ID, notificationType, hasSeen)
	if err != nil {
		return nil, err
	}

	mappedNotifications := make([]*commonModel.Notification, 0)

	for _, notification := range allNotifications {
		mappedNotification := mapNotificationToQlModel(notification)
		mappedNotifications = append(mappedNotifications, mappedNotification)
	}
	return mappedNotifications, nil
}

func (s NotificationsService) SetHasSeenNotifications(ctx context.Context, seenNotifications []*commonModel.SeenNotificationInput) ([]*commonModel.Notification, error) {
	user, err := s.usersRepo.GetByAuthenticationContext(ctx)
	if err != nil {
		return nil, err
	}

	notificationIds := make([]string, 0)
	for _, seenNotification := range seenNotifications {
		notificationIds = append(notificationIds, seenNotification.ID)
	}

	notificationsToUpdate, err := s.notificationsRepo.GetNotifications(ctx, notificationIds)
	if err != nil {
		return nil, err
	}
	for _, notificationToUpdate := range notificationsToUpdate {
		if notificationToUpdate.UserID != user.ID {
			return nil, errors.New("Cant set has seen on notification not belonging to user")
		}
		notificationToUpdate.SeenAt = null.TimeFrom(time.Now())
	}
	err = s.RunWithTX(ctx, func(tx *sql.Tx) error {
		err := s.notificationsRepo.UpdateNotifications(ctx, tx, notificationsToUpdate)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	allNotifications, err := s.notificationsRepo.ListAllNotifications(ctx, user.ID, nil, false)
	if err != nil {
		return nil, err
	}

	mappedNotifications := make([]*commonModel.Notification, 0)

	for _, notification := range allNotifications {
		mappedNotification := mapNotificationToQlModel(notification)
		mappedNotifications = append(mappedNotifications, mappedNotification)
	}
	return mappedNotifications, nil
}

func mapNotificationToQlModel(notification *models.Notification) *commonModel.Notification {
	mappedNotification := &commonModel.Notification{
		ID:        notification.ID,
		UserID:    notification.UserID,
		CreatedAt: notification.CreatedAt,
		SeenAt:    notification.SeenAt.Ptr(),
	}
	annotations := &commonModel.NotificationAnnotations{}
	if err := notification.Annotations.Unmarshal(&annotations); err != nil {
		// Todo handle err
	} else {
		mappedNotification.Annotations = annotations
	}
	return mappedNotification
}
