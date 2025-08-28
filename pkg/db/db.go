package db

import (
	"database/sql"
	"os"
	"strconv"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

const schema = `
CREATE TABLE scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT "",
    title VARCHAR(255) NOT NULL DEFAULT "",
    comment TEXT NOT NULL DEFAULT "",
    repeat VARCHAR(128) NOT NULL DEFAULT ""
);
CREATE INDEX idx_date ON scheduler(date);
`

// Init открывает БД и при необходимости создаёт таблицу и индекс.
func Init(dbFile string) error {
	install := false
	if _, err := os.Stat(dbFile); err != nil {
		install = true
	}

	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return err
	}
	DB = db

	if install {
		if _, err := DB.Exec(schema); err != nil {
			return err
		}
	}
	return nil
}

// Close закрывает соединение с БД.
func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}

// TasksBySearch выполняет поиск задач с учетом фильтрации по подстроке.
func TasksBySearch(limit int, pattern string) ([]*Task, error) {
	const q = `
		SELECT id, date, title, comment, repeat
		FROM scheduler
		WHERE title LIKE ? OR comment LIKE ?
		ORDER BY date
		LIMIT ?`
	rows, err := DB.Query(q, pattern, pattern, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*Task
	for rows.Next() {
		var idInt int64
		t := &Task{}
		if err := rows.Scan(&idInt, &t.Date, &t.Title, &t.Comment, &t.Repeat); err != nil {
			return nil, err
		}
		t.ID = strconv.FormatInt(idInt, 10)
		tasks = append(tasks, t)
	}
	if tasks == nil {
		tasks = []*Task{}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return tasks, nil
}
