package ratelimit

type Limiter interface{
	Allow(clientID string, route string) bool
	SetLimit(limit int)
}