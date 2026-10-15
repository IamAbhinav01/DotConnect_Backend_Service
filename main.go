package main

import (
	"DotConnect/app"
	"DotConnect/config"
	"fmt"
)

func main() {
	config.Load()
	application := app.NewApplication()
	err:=application.Run()
	if err != nil{
		fmt.Println("Error occured : ",err)
	}
}