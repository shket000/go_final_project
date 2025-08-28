package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"go_final_project/pkg/db"
)

// taskHandler — общий роутер для /api/task по методам.
func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		addTaskHandler(w, r)
	case http.MethodGet:
		getTaskHandler(w, r)
	case http.MethodPut:
		updateTaskHandler(w, r)
	case http.MethodDelete:
		deleteTaskHandler(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

type addTaskReq struct {
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var req addTaskReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, err, http.StatusBadRequest) // Неверный формат тела запроса
		return
	}

	if stringsTrim(req.Title) == "" {
		writeError(w, errors.New("title is required"), http.StatusBadRequest) // Отсутствует заголовок
		return
	}

	now := time.Now()
	// Если дата пустая — ставим сегодняшнюю
	if stringsTrim(req.Date) == "" {
		req.Date = now.Format(DateFmt)
	} else {
		if _, err := time.Parse(DateFmt, req.Date); err != nil {
			writeError(w, errors.New("invalid date format"), http.StatusBadRequest) // Неверный формат даты
			return
		}
	}

	// Проверка повторений
	if stringsTrim(req.Repeat) != "" {
		next, err := NextDate(now, req.Date, req.Repeat)
		if err != nil {
			writeError(w, err, http.StatusBadRequest) // Ошибка в вычислении повторений
			return
		}
		t, _ := time.Parse(DateFmt, req.Date)
		if !t.After(stripTime(now)) {
			req.Date = next
		}
	} else {
		t, _ := time.Parse(DateFmt, req.Date)
		if t.Before(stripTime(now)) {
			req.Date = now.Format(DateFmt)
		}
	}

	id, err := db.AddTask(&db.Task{
		Date:    req.Date,
		Title:   req.Title,
		Comment: req.Comment,
		Repeat:  req.Repeat,
	})
	if err != nil {
		writeError(w, err, http.StatusInternalServerError) // Ошибка при добавлении задачи
		return
	}

	writeJSON(w, map[string]string{"id": itoa(id)}, http.StatusCreated)
}

func stringsTrim(s string) string { return strings.TrimSpace(s) }

func itoa(id int64) string {
	return strconv.FormatInt(id, 10)
}
