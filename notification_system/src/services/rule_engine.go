package services

import "notification_system/src/models"

type RuleEngine struct{}

func NewRuleEngine() *RuleEngine { return &RuleEngine{} }

var defaultChannels = map[models.NotificationCategory][]models.ChannelType{
	models.CategoryTransaction: {models.ChannelEmail, models.ChannelPush, models.ChannelInApp},
	models.CategoryMarketing:   {models.ChannelEmail, models.ChannelInApp},
	models.CategorySecurity:    {models.ChannelEmail, models.ChannelSMS, models.ChannelPush},
	models.CategorySystem:      {models.ChannelInApp},
}

func (r *RuleEngine) ResolveChannels(n models.Notification, prefs *models.UserPreference) []models.ChannelType {
	channels := defaultChannels[n.Category]
	if prefs == nil || len(prefs.EnabledChannels) == 0 {
		return channels
	}
	var allowed []models.ChannelType
	for _, ch := range channels {
		if prefs.IsChannelAllowed(ch) {
			allowed = append(allowed, ch)
		}
	}
	return allowed
}
