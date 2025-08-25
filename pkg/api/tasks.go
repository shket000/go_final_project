package api

import (
	"net/http"
	"regexp"
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
			writeError(w, err)
			return
		}
		writeJSON(w, tasksResp{Tasks: tasks})
		return
	}

	// Матчим дату вида 02.01.2006
	dateRe := regexp.MustCompile(`^\d{2}\.\d{2}\.\d{4}$`)
	if dateRe.MatchString(search) {
		d, err := time.Parse("02.01.2006", search)
		if err != nil {
			writeError(w, err)
			return
		}
		date := d.Format(DateFmt)
		tasks, err := db.TasksByDate(defaultLimit, date)
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, tasksResp{Tasks: tasks})
		return
	}

	// LIKE поиск
	pattern := "%" + search + "%"
	tasks, err := db.Tasks(defaultLimit, pattern)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, tasksResp{Tasks: tasks})
}
