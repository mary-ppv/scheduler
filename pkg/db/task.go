package db

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	_ "modernc.org/sqlite"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func AddTask(task *Task) (int64, error) {
	res, err := DB.Exec("INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)",
		task.Date,
		task.Title,
		task.Comment,
		task.Repeat)
	if err != nil {
		return 0, fmt.Errorf("failed to insert task: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last inserted ID: %w", err)
	}

	return id, nil
}

func Tasks(limit int, searchFilter string) ([]*Task, error) {
	if DB == nil {
		return nil, fmt.Errorf("database not initialized")
	}

	var rows *sql.Rows
	var query string
	var args []interface{}

	if timeStr, err := time.Parse("02.01.2006", searchFilter); err == nil {
		formattedDate := timeStr.Format("20060102")
		query = `SELECT id, date, title, comment, repeat FROM scheduler WHERE date = ? ORDER BY date ASC LIMIT ?`
		args = []interface{}{formattedDate, limit}
	} else if searchFilter != "" {
		searchPattern := "%" + searchFilter + "%"
		query = `SELECT id, date, title, comment, repeat FROM scheduler WHERE title LIKE ? OR comment LIKE ? ORDER BY date ASC LIMIT ?`
		args = []interface{}{searchPattern, searchPattern, limit}
	} else {
		query = `SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date ASC LIMIT ?`
		args = []interface{}{limit}
	}

	rows, err := DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*Task
	for rows.Next() {
		var t Task
		var id int64
		err := rows.Scan(&id, &t.Date, &t.Title, &t.Comment, &t.Repeat)
		if err != nil {
			return nil, err
		}
		t.ID = strconv.FormatInt(id, 10)
		tasks = append(tasks, &t)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	if tasks == nil {
		tasks = make([]*Task, 0)
	}

	return tasks, nil
}

func GetTask(id string) (*Task, error) {
	var task Task

	err := DB.QueryRow(`SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`, id).Scan(
		&task.ID,
		&task.Date,
		&task.Title,
		&task.Comment,
		&task.Repeat,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("task not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get task: %w", err)
	}

	return &task, nil
}

func UpdateTask(task *Task) error {
	err := CheckDate(task, time.Now())
	if err != nil {
		return fmt.Errorf("invalid task date: %w", err)
	}

	var updates []string
	var args []interface{}

	if task.Date != "" {
		updates = append(updates, "date = ?")
		args = append(args, task.Date)
	}
	if task.Title != "" {
		updates = append(updates, "title = ?")
		args = append(args, task.Title)
	}
	if task.Comment != "" {
		updates = append(updates, "comment = ?")
		args = append(args, task.Comment)
	}
	if task.Repeat != "" {
		updates = append(updates, "repeat = ?")
		args = append(args, task.Repeat)
	}

	if len(updates) == 0 {
		return fmt.Errorf("no fields to update")
	}

	id, err := strconv.ParseInt(task.ID, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid task ID: %w", err)
	}
	args = append(args, id)

	query := fmt.Sprintf(
		"UPDATE scheduler SET %s WHERE id = ?",
		strings.Join(updates, ", "),
	)

	res, err := DB.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to execute update query: %w", err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if count == 0 {
		return fmt.Errorf("no rows updated — invalid ID or no changes made")
	}

	return nil
}

func DeleteTask(id string) error {
	if id == "" || id == "nil" {
		return fmt.Errorf("empty id")
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		return fmt.Errorf("invalid task ID: %w", err)
	}

	res, err := DB.Exec("DELETE FROM scheduler WHERE id = ?", idInt)
	if err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if count == 0 {
		return fmt.Errorf("no rows deleted — invalid ID or task already deleted")
	}

	return nil
}

func UpdateDate(next string, id string) error {
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid task ID: %w", err)
	}

	date, err := time.Parse("20060102", next)
	if err != nil {
		return fmt.Errorf("invalid date format: %w", err)
	}

	_, err = DB.Exec("UPDATE scheduler SET date = ? WHERE id = ?", date.Format("20060102"), idInt)
	if err != nil {
		return fmt.Errorf("failed to update task date: %w", err)
	}

	return nil
}
