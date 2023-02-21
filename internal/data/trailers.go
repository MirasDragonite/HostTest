package data

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/shynggys9219/greenlight/internal/validator"
)

type Trailer struct {
	ID int64 `json:"id"` // Unique integer ID for the movie

	Name     string `json:"trailersname"` // Movie title
	Duration int64  `json:"duration"`     // Movie release year, "omitempty" - hide from response if empty
	Date     string `json:"premierdate"`  // Movie title
	// time the movie information is updated
}

type TrailerModel struct {
	DB *sql.DB
}

func (m TrailerModel) Insert(trailer *Trailer) error {
	query := `
		INSERT INTO trailers(trailersname, duration, premierdate)
		VALUES ($1, $2, $3)
		RETURNING id`

	return m.DB.QueryRow(query, &trailer.Name, &trailer.Duration, &trailer.Date).Scan(&trailer.ID)
}

func ValidateTrailer(v *validator.Validator, trailer *Trailer) {
	v.Check(trailer.Name != "", "trailersname", "must be provided")
	v.Check(len(trailer.Name) <= 500, "trailersname", "must not be more than 500 bytes long")
	v.Check(trailer.Date != "", "premierdate", "must be provided")
	v.Check(trailer.Duration != 0, "duration", "must be provided")
	v.Check(trailer.Duration > 0, "duration", "must be a positive integer")
}

// arguments.
func (t TrailerModel) SearchByName(title string, filters Filters) ([]*Trailer, error) {

	query := fmt.Sprintf(`
SELECT id, trailersname, duration, premierdate
FROM trailers
WHERE (to_tsvector('simple', trailersname) @@ plainto_tsquery('simple', $1) OR $1 = '')
ORDER BY %s %s, id ASC`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := t.DB.QueryContext(ctx, query, title)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	trailers := []*Trailer{}
	for rows.Next() {
		var trailer Trailer
		err := rows.Scan(
			&trailer.ID,
			&trailer.Name,
			&trailer.Duration,
			&trailer.Date,
		)
		if err != nil {
			return nil, err
		}
		trailers = append(trailers, &trailer)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return trailers, nil
}

// func (m TrailerModel) Searchname(title string, filters Filters) ([]*Trailer, error) {
// 	// Update the SQL query to include the filter conditions.
// 	query := `
// 	SELECT id,trailersname, duration, premierdate
// 	FROM trailers
// 	WHERE (LOWER(trailersname) = LOWER($1) OR $1 = '')
// 	ORDER BY id`
// 	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
// 	defer cancel()
// 	// Pass the title and genres as the placeholder parameter values.
// 	rows, err := m.DB.QueryContext(ctx, query, title)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()
// 	trailers := []*Trailer{}
// 	for rows.Next() {
// 		var trailer Trailer
// 		err := rows.Scan(
// 			&trailer.ID,
// 			&trailer.Name,
// 			&trailer.Duration,
// 			&trailer.Date,
// 		)
// 		if err != nil {
// 			return nil, err
// 		}
// 		trailers = append(trailers, &trailer)
// 	}
// 	if err = rows.Err(); err != nil {
// 		return nil, err
// 	}
// 	return trailers, nil
// }

//////////
