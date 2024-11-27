package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi/v5"

	"final/handlers"
	"final/store"
	"final/tests"
)

func main() {
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = strconv.Itoa(tests.Port)
	}

	db, err := store.InitializeDataBase()
	if err != nil {
		log.Fatalf("Ошибка инициализации базы данных: %v", err)
	}

	handlers := handlers.Handlers{db}

	r := chi.NewRouter()
	r.Handle("/*", handlers.GetWebDir())
	r.Get("/api/nextdate", handlers.GetnextDateHandler)
	r.Post("/api/task", handlers.AddTask())
	r.Get("/api/tasks", handlers.GetTasks())
	r.Get("/api/task", handlers.GetTask())
	r.Put("/api/task", handlers.EditTask())
	r.Post("/api/task/done", handlers.MarkTask())
	r.Delete("/api/task", handlers.DeleteTask())

	log.Printf("Сервер слушает порт %s", port)

	if err := http.ListenAndServe(":"+port, r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
