package db

import (
	"database/sql"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func AddTask(task *Task) (int64, error) {
	const q = `INSERT INTO scheduler(date, title, comment, repeat) VALUES (?, ?, ?, ?)`
	res, err := DB.Exec(q, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// Tasks выбирает ближайшие задачи, возможно с фильтром поиска (LIKE) и лимитом.
func Tasks(limit int, like string) ([]*Task, error) {
	var rows *sql.Rows
	var err error

	if like != "" {
		const q = `
			SELECT id, date, title, comment, repeat
			FROM scheduler
			WHERE title LIKE ? OR comment LIKE ?
			ORDER BY date
			LIMIT ?`
		rows, err = DB.Query(q, like, like, limit)
	} else {
		const q = `
			SELECT id, date, title, comment, repeat
			FROM scheduler
			ORDER BY date
			LIMIT ?`
		rows, err = DB.Query(q, limit)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*Task
	for rows.Next() {
		t := &Task{}
		if err := rows.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	if tasks == nil {
		tasks = []*Task{}
	}
	return tasks, rows.Err()
}

// TasksByDate выбирает задачи ровно на указанную дату (YYYYMMDD), с лимитом.
func TasksByDate(limit int, date string) ([]*Task, error) {
	const q = `
		SELECT id, date, title, comment, repeat
		FROM scheduler
		WHERE date = ?
		ORDER BY date
		LIMIT ?`
	rows, err := DB.Query(q, date, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*Task
	for rows.Next() {
		t := &Task{}
		if err := rows.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	if tasks == nil {
		tasks = []*Task{}
	}
	return tasks, rows.Err()
}

func GetTask(id string) (*Task, error) {
	t := &Task{}
	const q = `
		SELECT id, date, title, comment, repeat
		FROM scheduler WHERE id = ?`
	if err := DB.QueryRow(q, id).Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat); err != nil {
		return nil, err
	}
	return t, nil
}

func UpdateTask(task *Task) error {
	const q = `
		UPDATE scheduler
		SET date = ?, title = ?, comment = ?, repeat = ?
		WHERE id = ?`
	res, err := DB.Exec(q, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func UpdateTaskDate(id string, newDate string) error {
	const q = `UPDATE scheduler SET date = ? WHERE id = ?`
	res, err := DB.Exec(q, newDate, id)
	if err != nil {
		return err
	}
	_, err = res.RowsAffected()
	return err
}

func DeleteTask(id string) error {
	const q = `DELETE FROM scheduler WHERE id = ?`
	_, err := DB.Exec(q, id)
	return err
}
