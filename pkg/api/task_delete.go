package api

import (
	"errors"
	"net/http"

	"go_final_project/pkg/db"
)

// DELETE /api/task?id=<number>
func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := stringsTrim(r.URL.Query().Get("id"))
	if id == "" {
		writeError(w, errors.New("id is required"))
		return
	}
	if err := db.DeleteTask(id); err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, map[string]any{})
}
