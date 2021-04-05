package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Shodocan/UserService/internal/configs"
	"github.com/Shodocan/UserService/internal/database"
	"github.com/Shodocan/UserService/internal/database/mongo"
	"github.com/Shodocan/UserService/internal/database/mongoredis"
	"github.com/Shodocan/UserService/internal/web"
	"github.com/joho/godotenv"
)

// @title User Service
// @version 1.0
// @description This is a sample swagger for User Service
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email wdcasonatto@gmail.com
// @host localhost:8080
// @BasePath /api/v1
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {
	log.SetFlags(log.Llongfile | log.LstdFlags)

	// load .env file from current directory into env variables
	err := godotenv.Load(".env")
	if err != nil {
		log.Printf("error reading .env file: %v", err)
	}

	// load env vars into struct
	config, err := configs.GetEnvConfig()
	if err != nil {
		log.Fatal(err)
		return
	}

	logger := configs.NewLog()

	//migration

	err = mongo.Migrate(config)
	if err != nil {
		panic(err)
	}

	// instantiate a postgres database instance
	var database database.MongoDB
	if config.RedisDBActive == "true" {
		database, err = mongoredis.NewDB(config, logger)
		if err != nil {
			log.Printf("database err %s", err)
			return
		}
	} else {
		database, err = mongo.NewDB(config, logger)
		if err != nil {
			log.Printf("database err %s", err)
			return
		}
	}

	// // create the fiber server
	server := web.Router(config, database)

	// listen on port and serve
	port := os.Getenv("PORT")
	if port == "" {
		// Default to port 8080
		port = "8080"
	}

	listenShutdown(func() error {
		log.Println(database.Disconnect())
		return server.Shutdown()
	})
	if err := server.Listen(fmt.Sprintf(":%v", port)); err != nil {
		log.Println(err)
	}
}

func listenShutdown(shutdownFunc func() error) {
	go func() {
		shutdown := make(chan os.Signal, 1)
		signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)
		sig := <-shutdown
		log.Printf("main : %+v : Start shutdown", sig)
		panic(shutdownFunc())
	}()
}
