package main

import (
	"DotConnect/app"
	"fmt"
)

func main() {
	application := app.NewApplication()
	err:=application.Run()
	if err != nil{
		fmt.Println("Error occured : ",err)
	}
}