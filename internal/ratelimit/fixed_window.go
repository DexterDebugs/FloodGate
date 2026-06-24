package ratelimit

import (
	"context"		//for redis cancellation /timeouts
	"fmt"			//for spintf to build the key
	"time"
	"github.com/redis/go-redis/v9"		//redis client library
)
//consumers
type FixedWindow struct {	
	client *redis.Client	//to talk to redis - the actual connection
	limit int			//how many requests per how long
	window time.Duration
}

//constructor
//This pattern is called dependency injection — the struct receives its dependencies instead of creating them.
func NewFixedWindow (client *redis.Client, limit int, window time.Duration)	*FixedWindow{
	return &FixedWindow{
		client: client,
		limit: limit,
		window: window,
	}
}

func (fw *FixedWindow)	Allow(clientID, route string)	bool {
	ctx := context.Background()		//through context from redis is how go propogates timeouts and cancellations

	bucket := time.Now().Unix()/int64(fw.window.Seconds())		//converts float to int64 for integer division, division is to floor the result
	key := fmt.Sprintf("ratelimit:%s:%s:%d", clientID, route, bucket)		// s,s for clientID and route, d for the bucket
	//Resulting key: ratelimit:dev-key-shanksss-001:/users/1:29166666
	count, err := fw.client.Incr(ctx, key).Result()		//atomic incrementation through redis
	/*Two layers happening here:

	fw.client.Incr(ctx, key) — schedules the Redis INCR command, 
		returns a *redis.IntCmd (a result handle, not the actual result yet)
	.Result() — actually executes and returns (int64, error) — the new count, plus any error
	*/

	if err != nil {
		return false	// fail open: don't take the gateway down if Redis flickers
	}
	if count == 1 {
		fw.client.Expire(ctx, key, fw.window)
	}
	return count <= int64(fw.limit)
}

func (fw *FixedWindow) SetLimit(limit int){
	fw.limit = limit
}