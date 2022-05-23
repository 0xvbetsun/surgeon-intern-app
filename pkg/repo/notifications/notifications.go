package notifications

import (
	"context"
	"database/sql"

	commonModel "github.com/vbetsun/surgeon-intern-app/internal/graph/model"
	"github.com/vbetsun/surgeon-intern-app/pkg/models"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type (
	INotificationsRepo interface {
		ListAllNotifications(ctx context.Context, userId string, notificationType *commonModel.NotificationType, hasSeen bool) ([]*models.Notification, error)
		GetNotifications(ctx context.Context, notificationIds []string) ([]*models.Notification, error)
		UpdateNotifications(ctx context.Context, tx *sql.Tx, notificationsToUpdate []*models.Notification) error
	}
	NotificationsRepo struct {
		db *sql.DB
	}
)

func NewNotificationsRepo(db *sql.DB) INotificationsRepo {
	return &NotificationsRepo{db: db}
}

func (r NotificationsRepo) ListAllNotifications(ctx context.Context, userId string, notificationType *commonModel.NotificationType, hasSeen bool) ([]*models.Notification, error) {
	filters := []qm.QueryMod{
		models.NotificationWhere.UserID.EQ(userId),
	}
	if hasSeen {
		filters = append(filters, models.NotificationWhere.SeenAt.IsNotNull())
	} else {
		filters = append(filters, models.NotificationWhere.SeenAt.IsNull())
	}
	notifications, err := models.Notifications(filters...).All(ctx, r.db)

	if err != nil {
		return nil, err
	}
	return notifications, nil
}

func (r NotificationsRepo) GetNotifications(ctx context.Context, notificationIds []string) ([]*models.Notification, error) {
	notifications, err := models.Notifications(models.NotificationWhere.ID.IN(notificationIds)).All(ctx, r.db)

	if err != nil {
		return nil, err
	}
	return notifications, nil
}

func (r NotificationsRepo) UpdateNotifications(ctx context.Context, tx *sql.Tx, notificationsToUpdate []*models.Notification) error {
	for _, notificationToUpdate := range notificationsToUpdate {
		_, err := notificationToUpdate.Update(ctx, tx, boil.Infer())
		if err != nil {
			return err
		}
	}
	return nil
}

func AddNotification(ctx context.Context, tx *sql.Tx, userId string, annotations commonModel.NotificationAnnotations) error {
	notification := models.Notification{
		UserID: userId,
	}

	notification.Annotations.Marshal(annotations)
	err := notification.Insert(ctx, tx, boil.Infer())
	if err != nil {
		return err
	}
	return nil
}
