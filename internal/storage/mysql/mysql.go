package mysql

import (
	"database/sql"
	"fmt"

	"github.com/shankar7042/students-api/internal/config"
	"github.com/shankar7042/students-api/internal/types"

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

func (m *MySql) GetStudentById(id int64) (types.Student, error) {
	stmt, err := m.Db.Prepare("SELECT id, name, email, age FROM students WHERE id = ? LIMIT 1")
	if err != nil {
		return types.Student{}, err
	}

	defer stmt.Close()

	var student types.Student

	err = stmt.QueryRow(id).Scan(&student.Id, &student.Name, &student.Email, &student.Age)
	if err != nil {
		if err == sql.ErrNoRows {
			return types.Student{}, fmt.Errorf("no student found with id=%s", fmt.Sprint(id))
		}
		return types.Student{}, fmt.Errorf("query error: %w", err)
	}

	return student, nil
}

func (m *MySql) GetStudents() ([]types.Student, error) {
	stmt, err := m.Db.Prepare("SELECT id, name, email, age FROM students")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var students []types.Student

	for rows.Next() {
		var student types.Student

		err := rows.Scan(&student.Id, &student.Name, &student.Email, &student.Age)
		if err != nil {
			return nil, err
		}

		students = append(students, student)
	}

	return students, nil
}
