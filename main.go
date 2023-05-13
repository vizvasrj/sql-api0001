package main

import (
	"log"
	"os"
	"src/conf"
	"src/database"
	"src/routes"
	"src/seeder"

	"github.com/gin-gonic/gin"
)

func main() {
	// time.Sleep(500 * time.Second)
	router := gin.Default()
	db, err := database.Get_database()
	if err != nil {
		log.Fatalln("error in opening postgres connection", err)
		// log.Println("Error conecting database")
		// time.Sleep(600 * time.Second)
	}
	quit := make(chan bool)
	var infoLog *log.Logger = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	var errorLog *log.Logger = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Llongfile)

	app := conf.Config{
		R:     router,
		Db:    db,
		Quit:  quit,
		Error: errorLog,
		Info:  infoLog,
		E:     conf.AddCallerInfo,
	}
	go app.Listen_for_shutdown()
	// err = seed_data(&app)
	// if err != nil {
	// 	app.Error.Fatal("something wrong with seed_data", err)
	// }
	routes.Setup_router(&app)
	err = seeder.Seeder_data(&app)
	if err != nil {
		log.Fatal(err)
	}
	if err != nil {
		log.Fatal(err)
	}
	router.Run()
	// <-app.quit
}
