package main

import (
	"fmt"
	"net/http"
	"time"
)

func securityHeaderMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//Preprocessing
		w.Header().Set("X-XSS-Protection", "1: mode=block")
		w.Header().Set("X_Frame-Options", "deny")
		//continue the chain
		next.ServeHTTP(w, r)
	})
}

func (app *application) logRequestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//Preprocessing
		start := time.Now()
		//continue the chain
		next.ServeHTTP(w, r)
		//preprosessing
		app.infoLog.Printf("%s %s %s %s %s",
			r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI(), time.Since(start))
	})
}

func (app *application) recoverPanicMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "Close")
				app.serverError(w, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}
