package routes

import (
	"log"
	"net/http"
	"src/conf"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type Handler func(*gin.Context) error

func Handle(h Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := h(c)

		if err != nil {
			cerr := errors.Cause(err)

			if cerr == conf.ErrNotFound {
				c.JSON(http.StatusNotFound, gin.H{
					"message": "Not Found",
					"success": false,
				})
				c.AbortWithStatus(404)
				return
			}

			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "internal server error",
				"success": false,
			})
			c.Status(500)
			return
		}
	}
}
