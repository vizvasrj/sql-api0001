package conf

import (
	"database/sql"
	"log"

	"github.com/gin-gonic/gin"
)

type Config struct {
	Db    *sql.DB
	R     *gin.Engine
	Quit  chan bool
	Error *log.Logger
	Info  *log.Logger
	E     func(error) error
	C     *gin.Context
}
