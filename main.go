package main

import (
	"log"
	"marathon-postgresql/config"
	"marathon-postgresql/server"
	"os"

	_ "github.com/lib/pq"
)

func getConfigFileName() string {
	env := os.Getenv("ENV")

	if env != "" {
		return "runners-" + env
	}

	return "runners"
}

func main() {
	log.Println("Starting Runners App")
	log.Println("Initializing configuration")
	config := config.InitConfig(getConfigFileName())
	log.Println("Initializing Database")
	dbHandler := server.InitDatabase(config)
	log.Println("Initializing Database")
	go server.InitPrometheus()
	log.Println("Initializing HTTP Server")
	httpServer := server.InitHttpServer(config, dbHandler)
	httpServer.Start()
}
