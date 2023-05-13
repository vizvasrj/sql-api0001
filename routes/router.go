package routes

import (
	"net/http"
	"src/conf"

	"github.com/gin-gonic/gin"
)

// errors.Wrap(conf.ErrNotFound, "index is not found here")

func Setup_router(app *conf.Config) *gin.Engine {

	h := MyHandler{
		C: *app,
	}

	app.R.GET("/", Handle(h.index))
	v1 := app.R.Group("/api/v1")
	{
		v1.GET("/longest-duration-movies", Handle(h.getLongestDurationMovies))
		v1.POST("/new-movie", Handle(h.createNewMovie))
		v1.GET("/top-rated-movies", Handle(h.getTopRatedMovies))
		v1.GET("/genre-movies-with-subtotals", Handle(h.getGenreMoviesWithSubtotals))
		v1.POST("/update-runtime-minutes", Handle(h.updateRuntimeMinutes))
		// v1.POST("/update-runtime-minutes-dec", Handle(h.updateRuntimeMinutesDec))

	}

	return app.R
}

func (h MyHandler) index(c *gin.Context) error {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
	return nil
}

// func getLongestDurationMovies(c *gin.Context, app *Config) error {
// 	var movies []Movie
// 	err :=
// }
