package mysql

import (
	"database/sql"

	"github.com/shankar7042/students-api/internal/config"

	_ "github.com/go-sql-driver/mysql"
)

type MySql struct {
	Db *sql.DB
}

func New(cfg *config.Config) (*MySql, error) {
	db, err := sql.Open("mysql", cfg.StoragePath)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS students(
		id INTEGER AUTO_INCREMENT PRIMARY KEY,
		name TEXT,
		email TEXT,
		age INTEGER
	);`)

	if err != nil {
		return nil, err
	}

	return &MySql{
		Db: db,
	}, nil

}

func (m *MySql) CreateStudent(name string, email string, age int) (int64, error) {
	stmt, err := m.Db.Prepare("INSERT INTO students (name, email, age) VALUES (?, ?, ?)")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(name, email, age)
	if err != nil {
		return 0, err
	}

	lastId, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return lastId, nil
}
