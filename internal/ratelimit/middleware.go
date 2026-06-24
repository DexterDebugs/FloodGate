package ratelimit

import (
	"net/http"
	"strings"
)

func Middleware(limiters map[string]Limiter, routes map[string]string) func(next http.Handler) http.Handler{	//takes a handler and returns a handler - Layer 1
/*	Which route does this request path match? — answered by walking the routes map (path → name)
	Which limiter is responsible for that route? — answered by looking up the name in limiters*/


	return func(next http.Handler)	http.Handler{		// Layer 2 - middleware (runs once per registsration)
		return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request){	//Layer 3 - always a request handler
			//The Gating logic goes here
			var routeName string		//walk through the routes to find which route this request path matches
			for path, name := range routes {
				if strings.HasPrefix(r.URL.Path, path) {
					routeName = name
					break
				}
			}
			// if no route matched, pass through (no rate limiting)
			if routeName == "" {
				next.ServeHTTP(w, r)
				return 
			}

			limiter, ok := limiters[routeName]		//look up the limiter for this route
			if !ok {
				next.ServeHTTP(w, r)
				return 
			}
			//gating logic
			clientID := r.Header.Get("X-API-Key")
			if !limiter.Allow(clientID, r.URL.Path)	{
				http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
				return 
			}
			next.ServeHTTP(w, r)
		})
	}
}