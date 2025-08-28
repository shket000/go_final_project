package api

import "net/http"

// Init регистрирует все API-хендлеры.
func Init() {
	http.HandleFunc("/api/nextdate", nextDateHandler)

	http.HandleFunc("/api/task", taskHandler) // mux по методам
	http.HandleFunc("/api/tasks", tasksHandler)

	http.HandleFunc("/api/task/done", taskDoneHandler)
}
