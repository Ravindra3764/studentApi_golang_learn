package sqlite

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"github.com/ravindra3764/student-api/student-api/internal/config"
	"github.com/ravindra3764/student-api/student-api/internal/types"
)

type Sqlite struct {
	Db *sql.DB
}

//intance of the db

func New(cfg *config.Config) (*Sqlite, error) {

	db, err := sql.Open("sqlite3", cfg.StoragePath)

	if err != nil {
		return nil, err
	}

	// Delete table if not required
	// db.Exec(`DROP TABLE IF EXISTS students`)

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS students(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT,
	email TEXT UNIQUE,
	age INTEGER
	)`)

	if err != nil {
		return nil, err
	}

	return &Sqlite{
		Db: db,
	}, nil

}

func (s *Sqlite) CreateStudent(name string, email string, age int) (int64, error) {

	statement, err := s.Db.Prepare("INSERT INTO STUDENTS (name, email, age) VALUES (? ,?, ?)")

	if err != nil {
		return 0, err
	}

	defer statement.Close()

	res, err := statement.Exec(name, email, age)

	if err != nil {
		return 0, err
	}

	lastid, err := res.LastInsertId()

	if err != nil {
		return 0, err
	}

	return lastid, nil

}

func (s *Sqlite) GetStudentByID(id int64) (types.Student, error) {

	statemen, error := s.Db.Prepare("SELECT * FROM STUDENTS WHERE id = ? LIMIT 1")

	if error != nil {
		return types.Student{}, error
	}

	defer statemen.Close()

	var student types.Student

	err := statemen.QueryRow(id).Scan(&student.Id, &student.Name, &student.Email, &student.Age)

	if err != nil {
		if err == sql.ErrNoRows {
			return types.Student{}, fmt.Errorf("No Student found with id : %s", fmt.Sprint(id))

		}
		return types.Student{}, fmt.Errorf("Query Error : %w", err)
	}

	return student, nil

}

func (s *Sqlite) GetAllStudents() ([]types.Student, error) {

	statement, err := s.Db.Prepare("SELECT id, name, email, age FROM STUDENTS")

	if err != nil {

		return nil, err

	}

	defer statement.Close()

	row, err := statement.Query()

	if err != nil {
		return nil, err
	}

	defer row.Close()

	var students []types.Student

	for row.Next() {

		var student types.Student

		err := row.Scan(&student.Id, &student.Name, &student.Email, &student.Age)

		if err != nil {
			return nil, err
		}

		students = append(students, student)
	}

	return students, nil
}

func (s *Sqlite) DeleteStudentByID(id int64) error {

	statement, err := s.Db.Prepare("DELETE FROM STUDENTS WHERE id = ?")

	if err != nil {
		return err
	}

	defer statement.Close()

	result, err := statement.Exec(id)

	if err != nil {
		return err
	}

	rowAffected, err := result.RowsAffected()

	if err != nil {
		return err
	}

	if rowAffected == 0 {
		return fmt.Errorf("no student found with id %d", id)

	}

	return nil

}

func (s *Sqlite) UpdateStudentByID(id int64, name string, email string, age int) (types.Student, error) {

	statement, err := s.Db.Prepare("UPDATE STUDENTS SET name = ?, email = ?, age = ? WHERE id = ?")

	if err != nil {
		return types.Student{}, err
	}

	defer statement.Close()

	result, err := statement.Exec(name, email, age, id)

	if err != nil {
		return types.Student{}, err
	}

	rowAffected, err := result.RowsAffected()

	if err != nil {
		return types.Student{}, err
	}

	if rowAffected == 0 {
		return types.Student{}, fmt.Errorf("no student found with id %d", id)
	}

	return types.Student{
		Id:    id,
		Name:  name,
		Email: email,
		Age:   age,
	}, nil

}
