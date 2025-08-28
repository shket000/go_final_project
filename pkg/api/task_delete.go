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
		writeError(w, errors.New("id is required"), http.StatusBadRequest) // Отсутствует id
		return
	}
	if err := db.DeleteTask(id); err != nil {
		writeError(w, err, http.StatusInternalServerError) // Ошибка при удалении задачи
		return
	}
	writeJSON(w, map[string]any{}, http.StatusOK)
}
