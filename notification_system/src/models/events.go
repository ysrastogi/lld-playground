package models

type NotificationEvent struct {
	Notification Notification
	Channels     []ChannelType
}
