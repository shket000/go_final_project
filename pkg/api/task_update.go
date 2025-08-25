package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"go_final_project/pkg/db"
)

type updateTaskReq struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

// PUT /api/task
func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var req updateTaskReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, err)
		return
	}
	if stringsTrim(req.ID) == "" {
		writeError(w, errors.New("id is required"))
		return
	}
	if stringsTrim(req.Title) == "" {
		writeError(w, errors.New("title is required"))
		return
	}
	if stringsTrim(req.Date) == "" {
		writeError(w, errors.New("date is required"))
		return
	}
	if _, err := time.Parse(DateFmt, req.Date); err != nil {
		writeError(w, errors.New("invalid date format"))
		return
	}

	// Если repeat указан — проверяем правило (и получим next, но менять дату по автологике при update не обязаны).
	if stringsTrim(req.Repeat) != "" {
		if _, err := NextDate(time.Now(), req.Date, req.Repeat); err != nil {
			writeError(w, err)
			return
		}
	}

	err := db.UpdateTask(&db.Task{
		ID:      req.ID,
		Date:    req.Date,
		Title:   req.Title,
		Comment: req.Comment,
		Repeat:  req.Repeat,
	})
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, map[string]any{})
}
