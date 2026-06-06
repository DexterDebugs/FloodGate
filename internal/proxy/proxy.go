//its like a waiter. it just does one work: receives a request, forward it to the backend, 
//return the response. that's it. nothing else

package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

func New(targetURL string)	(http.Handler, error){
	target, err := url.Parse(targetURL)	//turns plain string into structured *url.URL 
	if err != nil {
		return nil, err
	}

	return httputil.NewSingleHostReverseProxy(target), nil	//builds an actual reverse proxy
}
/*The proxy itself does this automatically when a request hits it:

Rewrites the request URL to point at the target backend (scheme + host swap)
Forwards the request to that backend
Reads the backend's response
Streams it back to the original client*/