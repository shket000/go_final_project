package db

import (
	"database/sql"
	"os"

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
