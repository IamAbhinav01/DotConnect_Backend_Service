package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func Load() error {
	return godotenv.Load()
}

func init(){
	Load()
}


func GetString(key string) string{
	value,ok:=os.LookupEnv(key)
	if !ok{
		return "Error occured while fetching the env variables"
	}

	return value
}

func GetInt(key string) int{
	value,ok:=os.LookupEnv(key)
	if !ok{
		return 0
	}
	valueInt,err:=strconv.Atoi(value)
	if err!=nil{
		return 0
	}

	return valueInt
}