package server

import (
	"fmt"
	"net/http"
	"os"

	"go_final_project/pkg/api"
)

const defaultPort = "7540"

func Run() error {
	// Порт из окружения (задача со звёздочкой)
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = defaultPort
	}

	// Раздача фронтенда
	webDir := "web"
	http.Handle("/", http.FileServer(http.Dir(webDir)))

	// Регистрация API
	api.Init()

	addr := fmt.Sprintf(":%s", port)
	fmt.Printf("Server listening on http://localhost%s\n", addr)
	return http.ListenAndServe(addr, nil)
}
