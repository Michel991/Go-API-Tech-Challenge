package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/michel991/go-api-tech-challenge/internal/models" // Use your module name
)

func GetAllCourses(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT id, name FROM course") // Change "courses" to "course"
		if err != nil {
			log.Printf("Error retrieving courses: %v", err)
			http.Error(w, "Failed to retrieve courses", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var courses []models.Course
		for rows.Next() {
			var course models.Course
			if err := rows.Scan(&course.ID, &course.Name); err != nil {
				log.Printf("Error scanning course: %v", err)
				http.Error(w, "Failed to scan course", http.StatusInternalServerError)
				return
			}
			courses = append(courses, course)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(courses)
	}
}

func GetCourseByID(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		var course models.Course
		err := db.QueryRow("SELECT id, name FROM course WHERE id = $1", id).Scan(&course.ID, &course.Name) // Change "courses" to "course"
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Course not found", http.StatusNotFound)
			} else {
				http.Error(w, "Failed to retrieve course", http.StatusInternalServerError)
			}
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(course)
	}
}

func UpdateCourse(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		var course models.Course
		if err := json.NewDecoder(r.Body).Decode(&course); err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}

		_, err := db.Exec("UPDATE course SET name = $1 WHERE id = $2", course.Name, id) // Change "courses" to "course"
		if err != nil {
			http.Error(w, "Failed to update course", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func AddCourse(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var course models.Course
		if err := json.NewDecoder(r.Body).Decode(&course); err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}

		err := db.QueryRow("INSERT INTO course (name) VALUES ($1) RETURNING id", course.Name).Scan(&course.ID) // Change "courses" to "course"
		if err != nil {
			http.Error(w, "Failed to add course", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]int{"id": course.ID})
	}
}

func DeleteCourse(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		_, err := db.Exec("DELETE FROM course WHERE id = $1", id) // Change "courses" to "course"
		if err != nil {
			http.Error(w, "Failed to delete course", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
