package main

import (
	"url_shortener/cmd/server"
	_ "github.com/joho/godotenv/autoload"
)

func main(){
	server.Start()
}