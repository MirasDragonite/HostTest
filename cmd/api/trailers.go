package main

import (
	"fmt"
	"net/http"

	"github.com/shynggys9219/greenlight/internal/data"
	"github.com/shynggys9219/greenlight/internal/validator"
)

func (app *application) createTrailerHandler(w http.ResponseWriter, r *http.Request) {

	var input struct {
		Name     string `json:"trailersname"`
		Duration int64  `json:"duration"`
		Date     string `json:"premierdate"`
	}

	// if there is error with decoding, we are sending corresponding message
	err := app.readJSON(w, r, &input) //non-nil pointer as the target decode destination
	if err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, err.Error())
	}

	trailer := &data.Trailer{
		Name:     input.Name,
		Duration: input.Duration,
		Date:     input.Date,
	}

	err = app.models.Trailer.Insert(trailer)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/trailers/%d", trailer.ID))
	err = app.writeJSON(w, http.StatusCreated, envelope{"trailers": trailer}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	// // Dump the contents of the input struct in a HTTP response.
	// fmt.Fprintf(w, "%+v\n", input) //+v here is adding the field name of a value // https://pkg.go.dev/fmt
}

func (app *application) listTrailerHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title string
		data.Filters
	}
	v := validator.New()
	qs := r.URL.Query()
	input.Title = app.readString(qs, "trailersname", "")

	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)
	input.Filters.Sort = app.readString(qs, "sort", "id")

	input.Filters.SortSafelist = []string{"id", "trailersname", "premierdate", "duration", "-premierdate", "-duration", "-id"}
	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	fmt.Fprintf(w, "%+v\n", input)

	trailers, err := app.models.Trailer.SearchByName(input.Title, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"trailers": trailers}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// func (app *application) listTrailerHandler(w http.ResponseWriter, r *http.Request) {
// 	var input struct {
// 		Title string

// 		data.Filters
// 	}
// 	v := validator.New()
// 	qs := r.URL.Query()
// 	input.Title = app.readString(qs, "title", "")

// 	input.Filters.Page = app.readInt(qs, "page", 1, v)
// 	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)
// 	input.Filters.Sort = app.readString(qs, "sort", "id")
// 	input.Filters.SortSafelist = []string{"id", "trailersname", "duration", "premierdate", "-id", "-trailersname", "-duration", "-premierdate"}

// 	if data.ValidateFilters(v, input.Filters); !v.Valid() {
// 		app.failedValidationResponse(w, r, v.Errors)
// 		return
// 	}
// 	// Call the GetAll() method to retrieve the movies, passing in the various filter
// 	// parameters.
// 	trailers, err := app.models.Trailer.GetAll(input.Title, input.Filters)
// 	if err != nil {
// 		app.serverErrorResponse(w, r, err)
// 		return
// 	}
// 	// Send a JSON response containing the movie data.
// 	err = app.writeJSON(w, http.StatusOK, envelope{"trailers": trailers}, nil)
// 	if err != nil {
// 		app.serverErrorResponse(w, r, err)
// 	}
// }

// func (app *application) listTrailerHandler(w http.ResponseWriter, r *http.Request) {
// 	var input struct {
// 		Name string
// 		data.Filters
// 	}
// 	v := validator.New()
// 	qs := r.URL.Query()
// 	input.Name = app.readString(qs, "title", "")

// 	input.Filters.Page = app.readInt(qs, "page", 1, v)
// 	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)
// 	input.Filters.Sort = app.readString(qs, "sort", "id")
// 	input.Filters.SortSafelist = []string{"id", "trailersname", "duration", "premierdate", "-trailersname", "-duration", "-premierdate"}
// 	if data.ValidateFilters(v, input.Filters); !v.Valid() {
// 		app.failedValidationResponse(w, r, v.Errors)
// 		return
// 	}
// 	// Call the GetAll() method to retrieve the movies, passing in the various filter
// 	// parameters.
// 	trailers, err := app.models.Trailer.Searchname(input.Name, input.Filters)
// 	if err != nil {
// 		app.serverErrorResponse(w, r, err)
// 		return
// 	}
// 	// Send a JSON response containing the movie data.
// 	err = app.writeJSON(w, http.StatusOK, envelope{"trailers": trailers}, nil)
// 	if err != nil {
// 		app.serverErrorResponse(w, r, err)
// 	}
// }
