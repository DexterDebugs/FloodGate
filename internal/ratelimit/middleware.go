package ratelimit

import "net/http"

func Middleware(limiter Limiter) func(next http.Handler) http.Handler{	//takes a handler and returns a handler - Layer 1

	return func(next http.Handler)	http.Handler{		// Layer 2 - middleware (runs once per registsration)
		return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request){	//Layer 3 - always a request handler
			//The Gating logic goes here
			clientID := r.Header.Get("X-API-Key")
			if !limiter.Allow(clientID, r.URL.Path)	{
				http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
				return 
			}
			next.ServeHTTP(w, r)
		})
	}
}