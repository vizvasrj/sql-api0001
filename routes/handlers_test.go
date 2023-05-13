package routes

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"src/conf"
	"src/database"
	"src/models"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	_ "github.com/lib/pq" // PostgreSQL driver
)

type getLongestDurationMoviesResponse struct {
	Results Results `json:"results"`
	Success bool    `json:"success"`
}

type Results struct {
	Movies []models.Movie `json:"movies"`
}

func getApp() conf.Config {
	db, err := database.Get_database()
	if err != nil {
		log.Fatalln("error in opening postgres connection", err)
	}
	r := gin.Default()

	var infoLog *log.Logger = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	var errorLog *log.Logger = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Llongfile)
	quit := make(chan bool)

	app := conf.Config{
		R:     r,
		Db:    db,
		Quit:  quit,
		Error: errorLog,
		Info:  infoLog,
		E:     conf.AddCallerInfo,
	}
	return app
}

func TestGetLongestDurationMovies(t *testing.T) {
	app := getApp()
	router := Setup_router(&app)

	recorder := httptest.NewRecorder()

	req, err := http.NewRequest(
		http.MethodGet, "/api/v1/longest-duration-movies", nil,
	)
	if err != nil {
		log.Fatal("error in http.NewRequest", err)
	}
	router.ServeHTTP(recorder, req)
	// log.Printf("%#v", req)
	assert.Equal(t, 200, recorder.Code)
	responseData, _ := ioutil.ReadAll(recorder.Body)
	var result_data getLongestDurationMoviesResponse
	err = json.Unmarshal(responseData, &result_data)
	if err != nil {
		log.Fatal("Error while unmarshing", err)
	}
	assert.Equal(t, result_data.Success, true)

	rows, err := app.Db.Query(`
		SELECT tconst
		FROM movies
		ORDER BY runtimeMinutes
		DESC LIMIT 10
	`)
	if err != nil {
		log.Fatal("error while getting data from db", err)
	}
	defer rows.Close()

	var movies []models.Movie
	for rows.Next() {
		var movie models.Movie
		err := rows.Scan(&movie.Tconst)
		if err != nil {
			log.Fatal("error while scaning rows", err)
		}
		movies = append(movies, movie)
	}

	if err = rows.Err(); err != nil {
		log.Fatal("error", err)
	}

	assert.Equal(t, result_data.Results.Movies[0].Tconst, movies[0].Tconst)
	assert.Equal(t, result_data.Results.Movies[1].Tconst, movies[1].Tconst)
	assert.Equal(t, result_data.Results.Movies[2].Tconst, movies[2].Tconst)
	assert.Equal(t, result_data.Results.Movies[3].Tconst, movies[3].Tconst)
	assert.Equal(t, result_data.Results.Movies[4].Tconst, movies[4].Tconst)
	assert.Equal(t, result_data.Results.Movies[5].Tconst, movies[5].Tconst)
	assert.Equal(t, result_data.Results.Movies[6].Tconst, movies[6].Tconst)
	assert.Equal(t, result_data.Results.Movies[7].Tconst, movies[7].Tconst)
	assert.Equal(t, result_data.Results.Movies[8].Tconst, movies[8].Tconst)
	assert.Equal(t, result_data.Results.Movies[9].Tconst, movies[9].Tconst)
}

func TestCreateNewMovie(t *testing.T) {
	app := getApp()
	router := Setup_router(&app)
	recorder := httptest.NewRecorder()
	payload := models.InsertMovie{
		Tconst:         "tt999999",
		TitleType:      "Short",
		PrimaryTitle:   "Test Movies",
		RuntimeMinutes: 55,
		Genres:         "Comedy",
		AverageRating:  5.9,
		NumVotes:       110,
	}
	payloadBytes, _ := json.Marshal(payload)

	req, err := http.NewRequest(
		http.MethodPost, "/api/v1/new-movie", bytes.NewBuffer(payloadBytes),
	)
	if err != nil {
		log.Fatal(err)
	}
	router.ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusCreated, recorder.Code)

	var get_movie models.InsertMovie
	err = app.Db.QueryRow(`
	SELECT 
		ratings.tconst, titleType, primaryTitle, 
		runtimeMinutes, genres, averageRating, numVotes
	FROM movies
	JOIN ratings
	ON ratings.tconst = movies.tconst
	WHERE ratings.tconst = $1
	`, payload.Tconst).Scan(
		&get_movie.Tconst, &get_movie.TitleType,
		&get_movie.PrimaryTitle, &get_movie.RuntimeMinutes,
		&get_movie.Genres, &get_movie.AverageRating, &get_movie.NumVotes,
	)

	if err != nil {
		log.Fatal("error while geting from manual db", err)
	}

	assert.Equal(t, payload.Tconst, get_movie.Tconst)
	assert.Equal(t, payload.TitleType, get_movie.TitleType)
	assert.Equal(t, payload.PrimaryTitle, get_movie.PrimaryTitle)
	assert.Equal(t, payload.RuntimeMinutes, get_movie.RuntimeMinutes)
	assert.Equal(t, payload.Genres, get_movie.Genres)
	assert.Equal(t, payload.AverageRating, get_movie.AverageRating)
	assert.Equal(t, payload.NumVotes, get_movie.NumVotes)

	_, err = app.Db.Exec(`
		DELETE FROM movies
		WHERE tconst = $1
	`, payload.Tconst)
	if err != nil {
		log.Fatal(err)
	}
}

func TestGetTopRatedMovies(t *testing.T) {
	app := getApp()
	router := Setup_router(&app)
	recorder := httptest.NewRecorder()

	req, err := http.NewRequest(
		http.MethodGet, "/api/v1/top-rated-movies", nil,
	)
	if err != nil {
		log.Fatal("error in http.NewRequest", err)
	}
	router.ServeHTTP(recorder, req)
	assert.Equal(t, 200, recorder.Code)

}

func TestGetGenreMoviesWithSubtotals(t *testing.T) {
	app := getApp()
	router := Setup_router(&app)
	recorder := httptest.NewRecorder()
	req, err := http.NewRequest(
		http.MethodGet, "/api/v1/genre-movies-with-subtotals", nil,
	)
	if err != nil {
		log.Fatal("error in http.NewRequest", err)
	}
	router.ServeHTTP(recorder, req)
	assert.Equal(t, 200, recorder.Code)

}

func TestUpdateRuntimeMinutes(t *testing.T) {
	app := getApp()
	router := Setup_router(&app)
	recorder := httptest.NewRecorder()
	req, err := http.NewRequest(
		http.MethodPost, "/api/v1/update-runtime-minutes", nil,
	)
	if err != nil {
		log.Fatal("error in http.NewRequest", err)
	}
	router.ServeHTTP(recorder, req)
	assert.Equal(t, 200, recorder.Code)

	_, err = app.Db.Exec(`
		UPDATE movies
		SET runtimeMinutes = CASE
			WHEN genres = 'Documentary' THEN runtimeMinutes - 15
			WHEN genres = 'Animation' THEN runtimeMinutes - 30
			ELSE runtimeMinutes - 45
		END
		`)
	if err != nil {
		log.Fatal("Database error", err)
	}

}
