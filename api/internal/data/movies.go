package data

import (
	"time"

	"github.com/ryanbyrne30/htmx/api/internal/validator"
)

type Movie struct {
	ID        int64				`json:"id"`
	CreatedAt time.Time 	`json:"-"`
	Title 		string			`json:"title"`
	Year 			int32 			`json:"year,omitempty"`
	Runtime 	Runtime				`json:"runtime,omitempty"`
	Genres 		[]string 		`json:"genres,omitempty"`
	Version 	int32				`json:"version"`
}

func ValidateMovie(v *validator.Validator, movie *Movie) {
	v.Check(movie.Title != "", "title", "must be provided")
	v.Check(len(movie.Title) <= 500, "title", "must not exceed 500 characters")

	v.Check(movie.Year != 0, "year", "must be provided")
	v.Check(movie.Year >= 1800, "year", "must be after 1800")
	v.Check(movie.Year <= int32(time.Now().Year()), "year", "year must not be in the future")

	v.Check(movie.Runtime != 0, "runtime", "must be provided")
	v.Check(movie.Runtime > 0, "runtime", "must be a positive integer")

	v.Check(movie.Genres != nil, "genres", "must be provided")
	v.Check(len(movie.Genres) >= 1, "genres", "must contain at least 1 genre")
	v.Check(len(movie.Genres) <= 5, "genres", "cannot contain more than 5 genres")
	v.Check(validator.Unique(movie.Genres), "genres", "must not contain duplicate values")

} 