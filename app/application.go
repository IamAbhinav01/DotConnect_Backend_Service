package app

import (
	"DotConnect/config"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type Application struct {
	PORT string
}

func NewApplication() *Application {
	return &Application{
		PORT: config.GetString("PORT"),
	}
}

func (app *Application) Run() error{
	

	err := config.ConnectToDB()
	if err != nil{
		return  fmt.Errorf("Unable to establish connection to the database %v",err)
	}

	fmt.Println("Succesfully connected to the database")

	port := app.PORT
	if(!strings.HasPrefix(port,":")){
		port = ":"+port
	}



	server:= http.Server{
		Addr: port,
		Handler:  http.DefaultServeMux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	
	fmt.Println("Server started on port ",port)
	return server.ListenAndServe()
}