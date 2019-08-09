package main

import (
	"encoding/gob"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	a "github.com/l3njo/dropnote-web/app"
	c "github.com/l3njo/dropnote-web/controllers"
)

var (
	app     a.Application
	port    string
	signals chan os.Signal
)

func cleanup() {
	log.Println("Shutting down server.")
}

func init() {
	signals = make(chan os.Signal)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-signals
		cleanup()
		os.Exit(1)
	}()

	e := godotenv.Load()
	c.Handle(e)

	port = os.Getenv("PORT")
	gob.Register(&c.SessionData{})
}

func main() {
	app.Init()
	app.Run(port)
}
