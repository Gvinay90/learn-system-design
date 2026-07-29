// Package factory implements the Factory Method pattern: a
// NotificationFactory that produces Email/SMS/Push notifications behind
// a common Notification interface, so callers never construct concrete
// types directly.
package factory

import (
	"errors"
	"fmt"
)

type NotificationType int

const (
	Email NotificationType = iota
	SMS
	Push
)

type Notification interface {
	Send(recipient, message string) string
}

type EmailNotification struct{}

func (EmailNotification) Send(recipient, message string) string {
	return fmt.Sprintf("Email to %s: %s", recipient, message)
}

type SMSNotification struct{}

func (SMSNotification) Send(recipient, message string) string {
	return fmt.Sprintf("SMS to %s: %s", recipient, message)
}

type PushNotification struct{}

func (PushNotification) Send(recipient, message string) string {
	return fmt.Sprintf("Push to %s: %s", recipient, message)
}

var ErrUnknownNotificationType = errors.New("unknown notification type")

// CreateNotification is the factory method: it hides construction of the
// concrete Notification type behind a single switch, so adding a new
// channel means adding one case here rather than touching every call site.
func CreateNotification(t NotificationType) (Notification, error) {
	switch t {
	case Email:
		return EmailNotification{}, nil
	case SMS:
		return SMSNotification{}, nil
	case Push:
		return PushNotification{}, nil
	default:
		return nil, ErrUnknownNotificationType
	}
}
