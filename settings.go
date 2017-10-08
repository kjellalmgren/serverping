package serverping

import (
	"serverping/notify"
)

// Settings
type Settings struct {
	Monitor       *MonitorSettings
	Notifications *NotificationSettings
}

// MonitorSettings
type MonitorSettings struct {
	CheckInterval             int `json:"checkInterval"`
	Timeout                   int `json:"timeout"`
	MaxConnections            int `json:"maxConnections"`
	ExponentialBackoffSeconds int `json:"exponentialBackoffSeconds"`
}

// NotificationSettings
type NotificationSettings struct {
	Email []*notify.EmailSettings `json:"email"`
	Sms   []*notify.SmsSettings   `json:"sms"`
}

// Get Notifiers
func (n *NotificationSettings) GetNotifiers() (notifiers notify.Notifiers) {
	for _, email := range n.Email {
		emailNotifier := &notify.EmailNotifier{Settings: email}
		notifiers = append(notifiers, emailNotifier)
	}
	for _, sms := range n.Sms {
		smsNotifier := &notify.SmsNotifier{Settings: sms}
		notifiers = append(notifiers, smsNotifier)
	}
	return
}
