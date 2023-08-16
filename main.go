package main

import (
	"log"
	"marathon-postgresql/config"
	"marathon-postgresql/server"

	_ "github.com/lib/pq"
)

func main() {
	log.Println("Starting Runners App")
	log.Println("Initializing configuration")
	config := config.InitConfig("runners.toml")
	log.Println("Initializing Database")
	dbHandler := server.InitDatabase(config)
	log.Println("Initializing HTTP Server")
	httpServer := server.InitHttpServer(config, dbHandler)
	httpServer.Start()
}
