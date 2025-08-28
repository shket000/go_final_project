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
		writeError(w, errors.New("method not allowed"), http.StatusMethodNotAllowed) // Неверный метод
		return
	}

	id := stringsTrim(r.URL.Query().Get("id"))
	if id == "" {
		writeError(w, errors.New("id is required"), http.StatusBadRequest) // Отсутствует id
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeError(w, err, http.StatusNotFound) // Задача не найдена
		return
	}

	if stringsTrim(task.Repeat) == "" {
		if err := db.DeleteTask(id); err != nil {
			writeError(w, err, http.StatusInternalServerError) // Ошибка при удалении задачи
			return
		}
		writeJSON(w, map[string]any{}, http.StatusOK)
		return
	}

	now := time.Now()
	next, err := NextDate(now, task.Date, task.Repeat)
	if err != nil {
		writeError(w, err, http.StatusBadRequest) // Ошибка в вычислении следующей даты
		return
	}
	if err := db.UpdateTaskDate(id, next); err != nil {
		writeError(w, err, http.StatusInternalServerError) // Ошибка при обновлении даты задачи
		return
	}
	writeJSON(w, map[string]any{}, http.StatusOK)
}
