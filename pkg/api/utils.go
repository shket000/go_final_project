package api

import (
	"encoding/json"
	"net/http"
)

const DateFmt = "20060102"

func writeJSON(w http.ResponseWriter, v any, statusCode int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(statusCode) // Устанавливаем нужный код статуса
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, err error, statusCode int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(statusCode) // Устанавливаем код HTTP статуса
	_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
}

// afterNow возвращает true, если a > b (по дате без времени).
func afterNow(a, b any) bool {
	switch x := a.(type) {
	case string:
		switch y := b.(type) {
		case string:
			return x > y
		}
	}
	return false
}
