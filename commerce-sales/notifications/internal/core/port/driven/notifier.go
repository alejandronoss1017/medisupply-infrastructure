package driven

import "context"

type Notifier interface {
	SendNotification(ctx context.Context, message map[string]interface{}, eventType string) error
}
