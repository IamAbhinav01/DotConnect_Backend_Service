package app

import (
	"DotConnect/config"
	"DotConnect/controllers"
	"DotConnect/repositories"
	"DotConnect/router"
	"DotConnect/services"
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

func (app *Application) Run() error {

	err := config.ConnectToDB()
	if err != nil {
		return fmt.Errorf("Unable to establish connection to the database %v", err)
	}

	fmt.Println("Succesfully connected to the database")

	wsRepo := repositories.NewWorkspaceRepository(config.DotConnect)
	chRepo := repositories.NewChannelRepository(config.Channels)
	wsServ := services.NewWorkspaceService(wsRepo, chRepo)
	wsCntrl := controllers.NewWorkspaceController(wsServ)
	wsRouter := router.NewWorkspaceRouter(wsCntrl)

	chServ := services.NewChannelService(chRepo, wsRepo)
	chCntrl := controllers.NewChannelController(chServ)
	chRouter := router.NewChannelRouter(chCntrl)

	chiRouter := router.SetupRouter(router.AppRouters{WorkspaceRouter: wsRouter, ChannelRouter: chRouter})

	port := app.PORT
	if !strings.HasPrefix(port, ":") {
		port = ":" + port
	}

	server := http.Server{
		Addr:         port,
		Handler:      chiRouter,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	fmt.Println("Server started on port ", port)
	return server.ListenAndServe()
}
