package storage

import (
	"github.com/ravindra3764/student-api/student-api/internal/types"
)

type Storage interface {
	CreateStudent(name string, email string, age int) (int64, error)

	GetStudentByID(id int64) (types.Student, error)

	GetAllStudents() ([]types.Student, error)
	DeleteStudentByID(id int64) error

	UpdateStudentByID(id int64, name string, email string, age int) (types.Student, error)
}
