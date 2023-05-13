package routes

import (
	"fmt"
	"math"
	"net/http"
	"src/conf"
	"src/models"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq" // PostgreSQL driver
)

type MyHandler struct {
	C conf.Config
}

func (h MyHandler) getLongestDurationMovies(c *gin.Context) error {
	var movies []models.Movie
	rows, err := h.C.Db.Query(`
		SELECT tconst, primaryTitle, runtimeMinutes, genres
		FROM movies 
		ORDER BY runtimeMinutes 
		DESC LIMIT 10
	`)
	if err != nil {
		return h.C.E(err)
	}
	defer rows.Close()

	for rows.Next() {
		var movie models.Movie
		err := rows.Scan(&movie.Tconst, &movie.PrimaryTitle, &movie.RuntimeMinutes, &movie.Genres)
		if err != nil {
			return h.C.E(err)
		}
		movies = append(movies, movie)
	}
	if err = rows.Err(); err != nil {
		return h.C.E(err)
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"results": gin.H{
			"movies": movies,
		},
	})
	return nil
}

func (h MyHandler) createNewMovie(c *gin.Context) error {
	var movie models.InsertMovie
	if err := c.ShouldBindJSON(&movie); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return nil
	}
	tx, err := h.C.Db.Begin()
	if err != nil {
		return h.C.E(err)
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	_, err = tx.Exec(`
		INSERT INTO movies (tconst, titleType, primaryTitle, runtimeMinutes, genres) 
		VALUES ($1, $2, $3, $4, $5)
	`, movie.Tconst, movie.TitleType, movie.PrimaryTitle, movie.RuntimeMinutes, movie.Genres)
	if err != nil {
		return h.C.E(err)
	}

	_, err = tx.Exec(`
		INSERT INTO ratings (tconst, averageRating, numVotes)
		VALUES ($1, $2, $3)
	`, movie.Tconst, movie.AverageRating, movie.NumVotes)
	if err != nil {
		return h.C.E(err)
	}

	err = tx.Commit()
	if err != nil {
		return h.C.E(err)
	}
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"movie":   movie,
	})
	return nil
}

func (h MyHandler) getTopRatedMovies(c *gin.Context) error {
	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("page_size", "10")

	pageNum, err := strconv.Atoi(page)
	if err != nil || pageNum <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page number", "success": false})
		return nil
	}
	size, err := strconv.Atoi(pageSize)
	if err != nil || size <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page size", "success": false})
		return nil
	}

	offset := (pageNum - 1) * size
	limit := size

	rows, err := h.C.Db.Query(`
		SELECT movies.tconst as tconst, primaryTitle, genres, averageRating 
		FROM movies 
		JOIN ratings ON movies.tconst = ratings.tconst 
		WHERE averageRating > 6.0 
		ORDER BY averageRating DESC 
		LIMIT $1 
		OFFSET $2
	`, limit, offset)
	if err != nil {
		h.C.Error.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "success": false})

		return nil
	}
	defer rows.Close()

	var totalCount int
	err = h.C.Db.QueryRow(`
		SELECT COUNT(*) 
		FROM movies 
		JOIN ratings ON movies.tconst = ratings.tconst 
		WHERE averageRating > 6.0
	`).Scan(&totalCount)
	if err != nil {
		h.C.Error.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "success": false})

		return nil
	}

	totalPages := int(math.Ceil(float64(totalCount) / float64(size)))

	nextPage := pageNum + 1
	previousPage := pageNum - 1
	if nextPage > totalPages {
		nextPage = -1
	}
	if previousPage <= 0 {
		previousPage = -1
	}

	var movies []models.InsertMovie
	for rows.Next() {
		var movie models.InsertMovie
		err := rows.Scan(&movie.Tconst, &movie.PrimaryTitle, &movie.Genres, &movie.AverageRating)
		if err != nil {
			h.C.Error.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "success": false})

			return nil
		}
		movies = append(movies, movie)
	}

	if err = rows.Err(); err != nil {
		h.C.Error.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "success": false})
		return nil

	}

	paginatedMovies := models.PaginatedMovies{
		Movies:       movies,
		TotalRecords: totalCount,
		TotalPages:   totalPages,
		CurrentPage:  pageNum,
		NextPage:     nextPage,
		PreviousPage: previousPage,
		Success:      true,
	}
	nextPageURL := ""
	if nextPage != -1 {
		nextPageURL = fmt.Sprintf("%s?page=%d&page_size=%d", c.Request.URL.Path, nextPage, size)
	}
	previousPageURL := ""
	if previousPage != -1 {
		previousPageURL = fmt.Sprintf("%s?page=%d&page_size=%d", c.Request.URL.Path, previousPage, size)
	}
	paginatedMovies.NextPageURL = nextPageURL
	paginatedMovies.PreviousPageURL = previousPageURL

	c.JSON(http.StatusOK, paginatedMovies)
	return nil
}

func (h MyHandler) getGenreMoviesWithSubtotals(c *gin.Context) error {

	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("page_size", "10")

	pageNum, err := strconv.Atoi(page)
	if err != nil || pageNum <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page number", "success": false})
		return nil
	}
	size, err := strconv.Atoi(pageSize)
	if err != nil || size <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page size", "success": false})
		return nil
	}

	offset := (pageNum - 1) * size
	limit := size

	var genres []models.GenreSubtotal
	rows, err := h.C.Db.Query(`
		SELECT m.genres as genres, m.primaryTitle, r.numVotes, sub.total_num_votes as subtotal
		FROM movies AS m
		JOIN ratings AS r ON r.tconst = m.tconst
		JOIN (
			SELECT genres, SUM(numVotes) AS total_num_votes
			FROM movies
			JOIN ratings ON ratings.tconst = movies.tconst
			GROUP BY genres
		) AS sub ON sub.genres = m.genres
		ORDER BY m.genres, m.primaryTitle
		LIMIT $1
		OFFSET $2
	`, limit, offset)
	if err != nil {
		return h.C.E(err)
	}
	defer rows.Close()

	var totalCount int
	err = h.C.Db.QueryRow(`
		SELECT COUNT(*)
		FROM (
			SELECT m.genres, m.primaryTitle, r.numVotes, sub.total_num_votes, COUNT(*) OVER () AS total_records
			FROM movies AS m
			JOIN ratings AS r ON r.tconst = m.tconst
			JOIN (
				SELECT genres, SUM(numVotes) AS total_num_votes
				FROM movies
				JOIN ratings ON ratings.tconst = movies.tconst
				GROUP BY genres
			) AS sub ON sub.genres = m.genres
			ORDER BY m.genres, m.primaryTitle
		) as subquery
	`).Scan(&totalCount)
	if err != nil {
		h.C.Error.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "success": false})
		return nil
	}

	totalPages := int(math.Ceil(float64(totalCount) / float64(size)))

	nextPage := pageNum + 1
	previousPage := pageNum - 1
	if nextPage > totalPages {
		nextPage = -1
	}
	if previousPage <= 0 {
		previousPage = -1
	}

	for rows.Next() {
		var genre models.GenreSubtotal
		err := rows.Scan(&genre.Genre, &genre.PrimaryTitle, &genre.NumVotes, &genre.Subtotal)
		if err != nil {
			h.C.Error.Println(err)
			return err
		}
		genres = append(genres, genre)
	}
	if err = rows.Err(); err != nil {
		h.C.Error.Println(err)
		return err
	}

	paginatedGenres := models.PaginatedGenre{
		Genres:       genres,
		TotalRecords: totalCount,
		TotalPages:   totalPages,
		CurrentPage:  pageNum,
		NextPage:     nextPage,
		PreviousPage: previousPage,
		Success:      true,
	}
	nextPageURL := ""
	if nextPage != -1 {
		nextPageURL = fmt.Sprintf("%s?page=%d&page_size=%d", c.Request.URL.Path, nextPage, size)
	}
	previousPageURL := ""
	if previousPage != -1 {
		previousPageURL = fmt.Sprintf("%s?page=%d&page_size=%d", c.Request.URL.Path, previousPage, size)
	}
	paginatedGenres.NextPageURL = nextPageURL
	paginatedGenres.PreviousPageURL = previousPageURL
	c.JSON(http.StatusOK, paginatedGenres)
	return nil
}

func (h MyHandler) updateRuntimeMinutes(c *gin.Context) error {
	_, err := h.C.Db.Exec(`
	UPDATE movies
	SET runtimeMinutes = CASE
		WHEN genres = 'Documentary' THEN runtimeMinutes + 15
		WHEN genres = 'Animation' THEN runtimeMinutes + 30
		ELSE runtimeMinutes + 45
	END;
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "success": false})
		return nil
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Runtime minutes updated successfully",
		"success": true,
	})
	return nil
}

// func (h MyHandler) updateRuntimeMinutesDec(c *gin.Context) error {
// 	_, err := h.C.Db.Exec(`
// 	UPDATE movies
// 	SET runtimeMinutes = CASE
// 		WHEN genres = 'Documentary' THEN runtimeMinutes - 15
// 		WHEN genres = 'Animation' THEN runtimeMinutes - 30
// 		ELSE runtimeMinutes - 45
// 	END;
// 	`)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return nil
// 	}
// 	c.JSON(http.StatusOK, gin.H{
// 		"message": "Runtime minutes updated successfully",
// 		"status":  "ok",
// 	})
// 	return nil
// }
