package server

import (
	"log"
	"net/http"
	"time"
)

type ResponseWriter struct{
	http.ResponseWriter			//Go generates a hidden field accessible as rw.ResponseWriter
	Status int
}


//Overriding method
func (rw *ResponseWriter) WriteHeader(code int){		//wraps http.responsewriter and remembers the status code when its written
	rw.Status = code		//captures the status code into your field
	rw.ResponseWriter.WriteHeader(code)		//Without this line,
	//  the status code would never actually get written to the HTTP response — your wrapper would intercept it and silently drop it
}


//middleware function
func Logging(next http.Handler)	http.Handler {	//takes a handler, returns a handler	
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {		//it is a type that converts such a function into something that satisfies http.Handler. It's an adapter.
		start := time.Now()
		rw := &ResponseWriter{ResponseWriter: w, Status: 200}	//creates a wrapper
		next.ServeHTTP(rw, r)	//pass wrapper through the next chain
		log.Printf("%s %s %d %v", r.Method, r.URL.Path, rw.Status, time.Since(start))
	}) 
}