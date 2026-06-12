package auth

import "net/http"

func Auth(validKeys []string) func(next http.Handler)	http.Handler{		//LAYER 1: factory (runs ONCE when called from main.go)
	keySet := make(map[string]bool)
	for _,k := range validKeys {
		keySet[k]	= true
	}	
	return func(next http.Handler)	http.Handler{			//LAYER 2: middleware (runs ONCE per route registration)
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {		//LAYER 3: request handler (runs ON	EVERY REQUEST)
			key := r.Header.Get("X-API-Key") //returns a header string, empty if null
			if !keySet[key]	{		//check gate first
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return 
			}
			next.ServeHTTP(w, r)	//  then pass through
		})
	}
}