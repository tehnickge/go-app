package main

import (
	"log"
	"net/http"

	"backend-app/internal/config"
	"backend-app/internal/migrations"
	"backend-app/pkg/db"

	"github.com/go-chi/chi/v5"
)

func main() {
	// Загружаем конфиг
	cfg := config.Load()

	// Подключаемся к базе
	database, err := db.NewPostgres(cfg)
	if err != nil {
		log.Fatal(err)
	}

	// Запускаем миграции
	migrations.RunMigrations(database)

	// Router
	r := chi.NewRouter()
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	log.Println("server running on :8080")
	http.ListenAndServe(":8080", r)
}