package store

import (
	"database/sql"
	"errors"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"

	"final/task"
)

const LimitOfTasks = 50

type DataBase struct {
	database *sql.DB
}

func InitializeDataBase() (DataBase, error) {
	appPath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
		return DataBase{database: nil}, err
	}

	dbFile := filepath.Join(filepath.Dir(appPath), "scheduler.db")
	_, err = os.Stat(dbFile)

	var install bool
	if err != nil {
		install = true
	}

	db, err := sql.Open("sqlite3", "scheduler.db")

	if err != nil {
		log.Fatal(err)
		return DataBase{database: nil}, err
	}
	if install {
		createTableSql := `CREATE TABLE IF NOT EXISTS scheduler (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			date CHAR(8) NOT NULL DEFAULT "",
			title VARCHAR(128) NOT NULL DEFAULT "",
			comment TEXT NOT NULL DEFAULT "",
			repeat VARCHAR(128) NOT NULL DEFAULT ""
			);
			CREATE INDEX IF NOT EXISTS scheduler_date ON scheduler (date);`
		_, err = db.Exec(createTableSql)
		if err != nil {
			log.Fatal(err)
			return DataBase{database: nil}, err
		}
		return DataBase{database: db}, nil
	}
	return DataBase{database: db}, nil
}

func (db *DataBase) AddTaskToDataBase(task task.Task) (int64, error) {
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
	res, err := db.database.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, errors.New("Ошибка добавления задачи")
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, errors.New("Ошибка добавления задачи")
	}
	return id, nil
}

func (db *DataBase) GetTasks() ([]task.Task, error) {
	rows, err := db.database.Query(`SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT :limit`, sql.Named("limit", LimitOfTasks))
	if err != nil {
		return nil, errors.New("Ошибка выполнения запроса: ")
	}
	defer rows.Close()

	tasks := make([]task.Task, 0, 0)

	for rows.Next() {
		var task task.Task
		if err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
			return nil, errors.New("Ошибка чтения строки: ")
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.New("Ошибка обработки результата: ")
	}
	return tasks, nil
}

func (db *DataBase) FindTask(id string) (task.Task, string) {
	var task task.Task
	query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`
	err := db.database.QueryRow(query, id).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return task, "задача не найдена"
		}
		return task, "ошибка выполнения запроса"
	}

	return task, ""
}

func (db *DataBase) UpdateTask(task task.Task) string {
	query := `UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`
	res, err := db.database.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return "ошибка выполнения запроса"
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil || rowsAffected == 0 {
		return "задача не найдена"
	}

	return ""
}

func (db *DataBase) DeleteTask(id string) string {
	query := `DELETE FROM scheduler WHERE id = ?`
	res, err := db.database.Exec(query, id)
	if err != nil {
		return "ошибка выполнения запроса"
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil || rowsAffected == 0 {
		return "запись не найдена"
	}

	return ""
}

func (db *DataBase) UpdateTaskDate(date, id string) string {
	query := `UPDATE scheduler SET date = ? WHERE id = ?`
	res, err := db.database.Exec(query, date, id)
	if err != nil {
		return "ошибка обновления задачи"
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil || rowsAffected == 0 {
		return "задача не найдена"
	}

	return ""
}
