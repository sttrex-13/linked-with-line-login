package main

import (
	"api/pkg/line"
	"api/pkg/util"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), (1 * time.Minute))
	defer cancel()

	// handler
	lineHandler := line.NewLINEHandler()

	r := mux.NewRouter()
	r.Use(CORS)

	r.Methods("OPTIONS").HandlerFunc(HandlerPreflight)
	r.HandleFunc("/health-check", func(w http.ResponseWriter, r *http.Request) {
		util.WriteJSONResponse(w, http.StatusOK, nil)
	}).Methods(http.MethodGet)

	linerouter := r.PathPrefix("/line").Subrouter()
	linerouter.HandleFunc("/webhook", lineHandler.WebHook).Methods(http.MethodPost)
	linerouter.HandleFunc("/request-login", lineHandler.RequestLogin).Methods(http.MethodGet)
	linerouter.HandleFunc("/login-callback", lineHandler.LoginCallback).Methods(http.MethodGet)

	srv := &http.Server{
		Handler:      r,
		Addr:         "0.0.0.0:8000",
		WriteTimeout: 1 * time.Minute,
		ReadTimeout:  1 * time.Minute,
	}

	go srv.ListenAndServe()
	log.Println("service is ready...")

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	<-c
	srv.Shutdown(ctx)
	log.Println("shutting down...")
	os.Exit(0)
}

func SetHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

func HandlerPreflight(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w)
	w.WriteHeader(http.StatusNoContent)
}

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		SetHeaders(w)
		next.ServeHTTP(w, r)
	})
}
