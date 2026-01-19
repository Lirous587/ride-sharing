package main

import "net/http"

func enableCORS(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Request-Method", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Request-Headers", "Content-Type, Authorization, Fuck")

		// allow preflight requests from the browser API
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		handler(w, r)
	}
}
