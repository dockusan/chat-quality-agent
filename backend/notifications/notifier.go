package notifications

import "context"

// Notifier defines the interface for sending notifications.
type Notifier interface {
	Send(ctx context.Context, subject string, body string) error
	HealthCheck(ctx context.Context) error
}
