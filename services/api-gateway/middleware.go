package main

import "net/http"

func enableCORS(next http.Handler) http.Handler {
	// return func(w http.ResponseWriter, r *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Request-Method", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Request-Headers", "Content-Type, Authorization, Cookie")

		// allow preflight requests from the browser API
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
