package api

import (
	"errors"
	"net/http"

	"go_final_project/pkg/db"
)

// GET /api/task?id=<number>
func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := stringsTrim(r.URL.Query().Get("id"))
	if id == "" {
		writeError(w, errors.New("id is required"))
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, task)
}
