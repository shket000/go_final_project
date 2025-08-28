package api

import (
	"net/http"
	"strings"
	"time"

	"go_final_project/pkg/db"
)

type tasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

// GET /api/tasks[?search=...]
// Задание со звёздочкой: поиск по подстроке title/comment (LIKE) и по дате "02.01.2006".
func tasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	const defaultLimit = 50

	search := strings.TrimSpace(r.URL.Query().Get("search"))
	if search == "" {
		tasks, err := db.Tasks(defaultLimit, "")
		if err != nil {
			writeError(w, err, http.StatusBadRequest)
			return
		}
		writeJSON(w, tasksResp{Tasks: tasks}, http.StatusOK)
		return
	}

	// Попробуем распарсить дату вида 02.01.2006
	if _, err := time.Parse("02.01.2006", search); err == nil {
		// Если дата корректная, ищем задачи по дате
		d, _ := time.Parse("02.01.2006", search) // игнорируем ошибку, она уже была проверена
		date := d.Format(DateFmt)
		tasks, err := db.TasksByDate(defaultLimit, date)
		if err != nil {
			writeError(w, err, http.StatusBadRequest)
			return
		}
		writeJSON(w, tasksResp{Tasks: tasks}, http.StatusOK)
		return
	}

	// Используем новый метод для LIKE поиска
	pattern := "%" + search + "%"
	tasks, err := db.TasksBySearch(defaultLimit, pattern)
	if err != nil {
		writeError(w, err, http.StatusBadRequest)
		return
	}
	writeJSON(w, tasksResp{Tasks: tasks}, http.StatusOK)
}
