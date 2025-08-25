package api

import (
	"errors"
	"net/http"
	"time"

	"go_final_project/pkg/db"
)

// POST /api/task/done?id=<number>
func taskDoneHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

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

	// Если правило пустое — просто удаляем
	if stringsTrim(task.Repeat) == "" {
		if err := db.DeleteTask(id); err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, map[string]any{})
		return
	}

	// Иначе вычисляем следующую дату и обновляем её
	now := time.Now()
	next, err := NextDate(now, task.Date, task.Repeat)
	if err != nil {
		writeError(w, err)
		return
	}
	if err := db.UpdateTaskDate(id, next); err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, map[string]any{})
}
