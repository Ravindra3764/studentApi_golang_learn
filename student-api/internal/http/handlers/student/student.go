package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/ravindra3764/student-api/student-api/internal/storage"
	"github.com/ravindra3764/student-api/student-api/internal/types"
	"github.com/ravindra3764/student-api/student-api/internal/utils/response"
)

func New(storage storage.Storage) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Creating a student")

		var student types.Student

		err := json.NewDecoder(r.Body).Decode(&student)

		if errors.Is(err, io.EOF) {
			// response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("Body cant be empty")))
			return
		}

		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		//request validatiion

		if err := validator.New().Struct(student); err != nil {
			validateError := err.(validator.ValidationErrors)
			response.WriteJson(w, http.StatusBadRequest, response.ValidatationError(validateError))
			return
		}

		lastId, err := storage.CreateStudent(

			student.Name,
			student.Email,
			student.Age,
		)

		slog.Info("user created succesfully", slog.String("userId ", fmt.Sprint(lastId)))
		if err != nil {

			response.WriteJson(w, http.StatusInternalServerError, err)
			return

		}

		// w.Write([]byte("Welcome to student apis"))
		response.WriteJson(w, http.StatusCreated, map[string]int64{"id": lastId})

	}

}

func GetById(storage storage.Storage) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		userId := r.PathValue("id")

		slog.Info("Getting a student", slog.String("stuedent id %s", userId))

		intId, err := strconv.ParseInt(userId, 10, 64)

		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		student, err := storage.GetStudentByID(intId)

		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return

		}

		response.WriteJson(w, http.StatusOK, student)

	}

}

func GetAllStudents(storage storage.Storage) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		slog.Info("Getting all student")

		res, err := storage.GetAllStudents()

		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, err)

		}

		response.WriteJson(w, http.StatusOK, res)

	}
}
