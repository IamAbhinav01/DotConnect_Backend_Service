package app

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

type Application struct {
	PORT string
}

func NewApplication() *Application {
	return &Application{
		PORT: os.Getenv("PORT"),
	}
}

func (app *Application) Run() error{
	
	server:= http.Server{
		Addr: ":4000",
		Handler:  http.DefaultServeMux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	
	fmt.Println("Server started on port : ",4000)
	return server.ListenAndServe()
}