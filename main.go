package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	a "github.com/l3njo/dropnote-web/app"
)

var (
	app a.App
	port string
	signals     chan os.Signal
)

func cleanup() {
	log.Println("Shutting down server.")
}

func handle(e error) {
	if e != nil {
		log.Println(e)
	}
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
	handle(e)

	port = os.Getenv("PORT")
}

func main() {
	app.Init()
	app.Run(port)
}
