package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/michel991/go-api-tech-challenge/internal/models" //Use your module name
)

func GetAllPersons(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT id, first_name, last_name, type, age FROM person")
		if err != nil {
			log.Printf("Error retrieving persons: %v", err)
			http.Error(w, "Failed to retrieve persons", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var persons []models.Person
		for rows.Next() {
			var person models.Person
			if err := rows.Scan(&person.ID, &person.FirstName, &person.LastName, &person.Type, &person.Age); err != nil {
				log.Printf("Error scanning person: %v", err)
				http.Error(w, "Failed to scan person", http.StatusInternalServerError)
				return
			}
			person.Courses, err = getPersonCourses(db, person.ID)
			if err != nil {
				log.Printf("Error retrieving courses for person %d: %v", person.ID, err)
				http.Error(w, "Failed to retrieve courses for person", http.StatusInternalServerError)
				return
			}
			persons = append(persons, person)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(persons)
	}
}

func GetPersonByID(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		var person models.Person
		err := db.QueryRow("SELECT id, first_name, last_name, type, age FROM person WHERE id = $1", id).
			Scan(&person.ID, &person.FirstName, &person.LastName, &person.Type, &person.Age)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Person not found", http.StatusNotFound)
			} else {
				http.Error(w, "Failed to retrieve person", http.StatusInternalServerError)
			}
			return
		}

		person.Courses, err = getPersonCourses(db, person.ID)
		if err != nil {
			log.Printf("Error retrieving courses for person %d: %v", person.ID, err)
			http.Error(w, "Failed to retrieve courses for person", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(person)
	}
}

func UpdatePerson(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		if id == "" {
			log.Printf("No ID provided for update")
			http.Error(w, "ID is required", http.StatusBadRequest)
			return
		}

		var person models.Person
		if err := json.NewDecoder(r.Body).Decode(&person); err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}

		// Update the person using the ID
		_, err := db.Exec("UPDATE person SET first_name = $1, last_name = $2, type = $3, age = $4 WHERE id = $5",
			person.FirstName, person.LastName, person.Type, person.Age, id)
		if err != nil {
			log.Printf("Error updating person with ID %s: %v", id, err)
			http.Error(w, "Failed to update person", http.StatusInternalServerError)
			return
		}

		// Update the person's courses
		if err := updatePersonCourses(db, person.ID, person.Courses); err != nil {
			log.Printf("Error updating courses for person %d: %v", person.ID, err)
			http.Error(w, "Failed to update courses for person", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func AddPerson(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var person models.Person
		if err := json.NewDecoder(r.Body).Decode(&person); err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}

		err := db.QueryRow("INSERT INTO person (first_name, last_name, type, age) VALUES ($1, $2, $3, $4) RETURNING id",
			person.FirstName, person.LastName, person.Type, person.Age).Scan(&person.ID)
		if err != nil {
			http.Error(w, "Failed to add person", http.StatusInternalServerError)
			return
		}

		if err := updatePersonCourses(db, person.ID, person.Courses); err != nil {
			log.Printf("Error adding courses for person %d: %v", person.ID, err)
			http.Error(w, "Failed to add courses for person", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]int{"id": person.ID})
	}
}

func DeletePerson(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		if id == "" {
			log.Printf("No ID provided for deletion")
			http.Error(w, "ID is required", http.StatusBadRequest)
			return
		}

		// Delete related entries in person_course
		_, err := db.Exec("DELETE FROM person_course WHERE person_id = $1", id)
		if err != nil {
			log.Printf("Error deleting courses for person with ID %s: %v", id, err)
			http.Error(w, "Failed to delete courses for person", http.StatusInternalServerError)
			return
		}

		// Delete the person
		_, err = db.Exec("DELETE FROM person WHERE id = $1", id)
		if err != nil {
			log.Printf("Error deleting person with ID %s: %v", id, err)
			http.Error(w, "Failed to delete person", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

// Helper function to get courses for a person
func getPersonCourses(db *sql.DB, personID int) ([]int, error) {
	rows, err := db.Query("SELECT course_id FROM person_course WHERE person_id = $1", personID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var courses []int
	for rows.Next() {
		var courseID int
		if err := rows.Scan(&courseID); err != nil {
			return nil, err
		}
		courses = append(courses, courseID)
	}
	return courses, nil
}

// Helper function to update courses for a person
func updatePersonCourses(db *sql.DB, personID int, courses []int) error {
	// Delete existing courses
	_, err := db.Exec("DELETE FROM person_course WHERE person_id = $1", personID)
	if err != nil {
		return err
	}

	// Insert new courses
	for _, courseID := range courses {
		_, err := db.Exec("INSERT INTO person_course (person_id, course_id) VALUES ($1, $2)", personID, courseID)
		if err != nil {
			return err
		}
	}
	return nil
}

func GetPersonByName(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := chi.URLParam(r, "name")
		var person models.Person
		err := db.QueryRow("SELECT id, first_name, last_name, type, age FROM person WHERE first_name = $1 OR last_name = $1", name).
			Scan(&person.ID, &person.FirstName, &person.LastName, &person.Type, &person.Age)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Person not found", http.StatusNotFound)
			} else {
				http.Error(w, "Failed to retrieve person", http.StatusInternalServerError)
			}
			return
		}

		person.Courses, err = getPersonCourses(db, person.ID)
		if err != nil {
			log.Printf("Error retrieving courses for person %d: %v", person.ID, err)
			http.Error(w, "Failed to retrieve courses for person", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(person)
	}
}
