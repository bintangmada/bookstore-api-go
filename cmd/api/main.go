package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/bintangmada/bookstore-api-go/internal/book"
	"github.com/bintangmada/bookstore-api-go/internal/config"
	"github.com/bintangmada/bookstore-api-go/internal/database"
)

func main() {
	cfg := config.Load()

	db, err := database.NewMySQL(cfg)
	if err != nil {
		log.Fatalf("gagal terhubung ke MySQL: %v", err)
	}
	defer db.Close()

	log.Printf("berhasil terhubung ke database MySQL %q (%s:%s)", cfg.DBName, cfg.DBHost, cfg.DBPort)

	bookRepo := book.NewRepository(db)
	bookHandler := book.NewHandler(bookRepo)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		pingCtx, pingCancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer pingCancel()

		w.Header().Set("Content-Type", "application/json")
		if err := db.PingContext(pingCtx); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			_ = json.NewEncoder(w).Encode(map[string]string{
				"status":   "error",
				"database": "disconnected",
			})
			return
		}

		_ = json.NewEncoder(w).Encode(map[string]string{
			"status":   "ok",
			"database": "connected",
		})
	})

	bookHandler.RegisterRoutes(mux)

	addr := ":" + cfg.AppPort
	server := &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("server berjalan di http://localhost%s", addr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("server berhenti: %v", err)
	}
}
