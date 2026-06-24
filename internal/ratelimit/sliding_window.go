package ratelimit

import (
	"context"
	"fmt"
	"time"
	"github.com/redis/go-redis/v9"
)

type SlidingWindow struct {
	client *redis.Client
	limit int
	window time.Duration
}

func NewSlidingWindow (client *redis.Client, limit int, window time.Duration)	*SlidingWindow {
	return &SlidingWindow{
		client: client,
		limit: limit,
		window: window,
	}
}

func (sw *SlidingWindow) Allow(clientID, route string)	bool {
	ctx := context.Background()

	currentBucket := time.Now().Unix()/int64(sw.window.Seconds())
	currKey := fmt.Sprintf("ratelimit:%s:%s:%d", clientID, route, currentBucket)

	prevBucket := currentBucket - 1
	prevKey := fmt.Sprintf("ratelimit:%s:%s:%d", clientID, route, prevBucket)

	count, err := sw.client.Incr(ctx, currKey).Result()
	if err != nil {
		return true		//nothing its intentional
	}
	if count == 1 {
		sw.client.Expire(ctx, currKey, sw.window)
	}

	prevCount, err := sw.client.Get(ctx, prevKey).Int64()
	if err == redis.Nil {
		prevCount = 0
	} else if err != nil {
		return true   // fail open
	}
	windowSeconds := int64(sw.window.Seconds())
	elapsedSec := time.Now().Unix() - (currentBucket * windowSeconds)
	elapsed := float64(elapsedSec)	/	float64(windowSeconds)
	estimate := float64(prevCount)*(1-elapsed) + float64(count)

	return estimate <= float64(sw.limit)
}

func (sw *SlidingWindow) SetLimit(limit int){
	sw.limit = limit
}