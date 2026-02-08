package models

import "time"

type entityType string

type RequestContext struct {
	UserID    string
	ApiKey    string
	IpAddress string
}

func (r *RequestContext) SetUserID(id string) {
	r.UserID = id
}

func (r *RequestContext) SetAPIKey(key string) {
	r.ApiKey = key
}

func (r *RequestContext) SetIPAddress(ip string) {
	r.IpAddress = ip
}

const (
	User entityType = "User"
	IP   entityType = "IP"
)

type LimitPolicy struct {
	Requests           int
	Timeframe          time.Duration
	MaxBurst           int
	ConcurrentRequests int
	Entity             entityType
}

func (l *LimitPolicy) SetRequests(r int) {
	l.Requests = r
}

func (l *LimitPolicy) SetTimeframe(t time.Duration) {
	l.Timeframe = t
}

func (l *LimitPolicy) SetMaxBurst(b int) {
	l.MaxBurst = b
}

func (l *LimitPolicy) SetConcurrentRequests(c int) {
	l.ConcurrentRequests = c
}

func (l *LimitPolicy) SetEntity(e entityType) {
	l.Entity = e
}
