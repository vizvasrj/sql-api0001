package conf

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func (app *Config) Listen_for_shutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println("Waiting for shutdown")
	<-quit
	// app.quit <- true
	fmt.Println("I am sutting")
	app.Db.Close()
	os.Exit(0)
}
