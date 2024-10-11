package api

import (
	"database/sql"

	"github.com/go-chi/chi/v5"
	"github.com/michel991/go-api-tech-challenge/internal/api/handlers"
	"github.com/michel991/go-api-tech-challenge/internal/api/middleware"
)

func SetupRoutes(r *chi.Mux, db *sql.DB) {
	r.Use(middleware.Logging)

	r.Route("/api", func(r chi.Router) {
		r.Route("/course", func(r chi.Router) {
			r.Get("/", handlers.GetAllCourses(db))
			r.Get("/{id}", handlers.GetCourseByID(db))
			r.Put("/{id}", handlers.UpdateCourse(db))
			r.Post("/", handlers.AddCourse(db))
			r.Delete("/{id}", handlers.DeleteCourse(db))
		})

		r.Route("/person", func(r chi.Router) {
			r.Get("/", handlers.GetAllPersons(db))
			r.Get("/{name}", handlers.GetPersonByName(db))
			r.Put("/{id}", handlers.UpdatePerson(db))
			r.Post("/", handlers.AddPerson(db))
			r.Delete("/{id}", handlers.DeletePerson(db))
		})
	})
}
