package seeder

import (
	"encoding/csv"
	"os"
	"src/conf"
	"strconv"

	pg "github.com/lib/pq"
)

type SeedMovie struct {
	Tconst         string
	TitleType      string
	PrimaryTitle   string
	RuntimeMinutes string
	Genres         string
}

type Rating struct {
	Tconst        string
	AverageRating float64
	NumVotes      int
}

var (
	MOVIES_CSV  = "movies.csv"
	RATINGS_CSV = "ratings.csv"
)

func seeder_movies(app *conf.Config) error {
	file, err := os.Open(MOVIES_CSV)
	if err != nil {
		return app.E(err)
	}
	defer file.Close()
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return app.E(err)
	}
	for i, record := range records {
		if i == 0 {
			continue
		}
		movie := SeedMovie{
			Tconst:         record[0],
			TitleType:      record[1],
			PrimaryTitle:   record[2],
			RuntimeMinutes: record[3],
			Genres:         record[4],
		}
		// runTimeMin, err := strconv.ParseInt(movie.RuntimeMinutes, 10, 32)
		// if err != nil {
		// 	log.Println("movie.RuntimeMinutes ", movie.RuntimeMinutes)
		// 	return app.E(err)
		// }
		_, err = app.Db.Exec(`
			INSERT INTO movies (tconst, titleType, primaryTitle, runtimeMinutes, genres)
			VALUES ($1, $2, $3, $4, $5)
		`, movie.Tconst, movie.TitleType, movie.PrimaryTitle, movie.RuntimeMinutes, movie.Genres)
		if err != nil {
			pgErr, ok := err.(*pg.Error)
			if ok && pgErr.Code == "23505" {
				app.Info.Println("Duplicate key violation error:", pgErr)
			} else {
				app.Error.Println(err)
				return app.E(err)
			}
		}

	}
	_, err = app.Db.Exec(`
		INSERT INTO seeder (file_name)
		VALUES($1)
	`, MOVIES_CSV)
	if err != nil {
		return app.E(err)
	}
	// lock_file, err := os.Create(MOVIES_CSV)
	// if err != nil {
	// 	return app.E(err)
	// }
	// defer lock_file.Close()

	app.Info.Println("seed data inserted successfully")
	return nil

}

func seeder_ratings(app *conf.Config) error {
	file, err := os.Open(RATINGS_CSV)
	if err != nil {
		return app.E(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	records, err := reader.ReadAll()
	if err != nil {
		return app.E(err)
	}
	for i, record := range records {
		if i == 0 {
			continue
		}
		averageRating, err := strconv.ParseFloat(record[1], 64)
		if err != nil {
			return app.E(err)
		}

		numVotes, err := strconv.Atoi(record[2])
		if err != nil {
			return app.E(err)
		}

		rating := Rating{
			Tconst:        record[0],
			AverageRating: averageRating,
			NumVotes:      numVotes,
		}

		_, err = app.Db.Exec(`
			INSERT INTO ratings (tconst, averageRating, numVotes)
			VALUES ($1, $2, $3)
		`, rating.Tconst, rating.AverageRating, rating.NumVotes)
		if err != nil {
			pgErr, ok := err.(*pg.Error)
			if ok && pgErr.Code == "23505" {
				app.Info.Println("Duplicate key violation error:", pgErr)
			} else {
				app.Error.Println(err)
				return app.E(err)
			}
		}
	}
	_, err = app.Db.Exec(`
		INSERT INTO seeder (file_name)
		VALUES($1)
	`, RATINGS_CSV)
	if err != nil {
		return app.E(err)
	}

	app.Info.Println("seed data inserted successfully")
	return nil
}

func Seeder_data(app *conf.Config) error {

	var movies_seed_exists bool
	err := app.Db.QueryRow(`
		SELECT EXISTS(
			SELECT 1 
			FROM seeder
			WHERE file_name = $1
		)
	`, MOVIES_CSV).Scan(&movies_seed_exists)
	if err != nil {
		return app.E(err)
	}

	if !movies_seed_exists {
		err := seeder_movies(app)
		if err != nil {
			return app.E(err)
		}
	} else {
		app.Info.Println("movies data exists")
	}

	var ratings_seed_exists bool
	err = app.Db.QueryRow(`
		SELECT EXISTS(
			SELECT 1 
			FROM seeder
			WHERE file_name = $1
		)
	`, RATINGS_CSV).Scan(&ratings_seed_exists)

	if !ratings_seed_exists {
		err := seeder_ratings(app)
		if err != nil {
			return app.E(err)
		}
	} else {
		app.Info.Println("ratings data exists")
	}
	return nil
}
