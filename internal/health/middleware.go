package health

import (
	"net/http"
	"strings"
	"time"
	"github.com/DexterDebugs/FloodGate/internal/server"
)

func Middleware(tracker Tracker, routes map[string]string)	func(http.Handler)	http.Handler {
	// Layer 1: factory — runs once at startup
	return func(next http.Handler) http.Handler {	// Layer 2: returns the middleware, actual wrapper function. takes the next handler in chain (the proxy, middleware) and returns a new handler that wraps it.
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {	// Layer 3: per-request
			//Handler func is like an adaptor - it turns plain function into something that satisfies the handler interface.
			var backend string
			for prefix, name := range routes {	//iterate over the map: prefix is the URL path (/users/), name is the value: (users-service)
				if strings.HasPrefix(r.URL.Path, prefix){
					backend = name
					break	//match found, record and break
				}
			}
			start := time.Now()		//capture start time
			rw := &server.ResponseWriter{ResponseWriter: w, Status: 200}	//wrap the writer
			next.ServeHTTP(rw, r)		//hand it to proxy
			if 	backend == "" {		//defensive guard
				return 	//no backend matched, don't record
			}
			isError := rw.Status >= 500
			tracker.Record(backend, time.Since(start), isError)		//record the result
		})
	}
}