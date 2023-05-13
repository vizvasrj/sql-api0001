package seeder

import (
	"log"
	"os"
	"src/conf"
)

func seed_movies(app *conf.Config) error {
	filename := ".movies_seed.lock"
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		log.Printf("File '%s' does not exist\n", filename)
		_, err := app.Db.Exec(`COPY movies(tconst, titleType, primaryTitle, runtimeMinutes, genres) FROM 'movies.csv' DELIMITER ',' CSV HEADER`)
		if err != nil {
			return err
		}
		file, err := os.Create(".movies_seed.lock")
		if err != nil {
			return err
		}
		defer file.Close()
	} else {
		log.Println("movies data exists.")
	}
	return nil
}

func seed_ratings(app *conf.Config) error {
	filename := ".rating_seed.lock"
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		_, err := app.Db.Exec(`COPY movies(tconst, titleType, primaryTitle, runtimeMinutes, genres) FROM 'movies.csv' DELIMITER ',' CSV HEADER`)
		if err != nil {
			return err
		}
		file, err := os.Create(".rating_seed.lock")
		if err != nil {
			return err
		}
		defer file.Close()
	} else {
		log.Println("rating data exists.")
	}
	return nil
}

func seed_data(app *conf.Config) error {
	err := seed_movies(app)
	if err != nil {
		return err
	}

	err = seed_ratings(app)
	if err != nil {
		return err
	}
	return nil
}
