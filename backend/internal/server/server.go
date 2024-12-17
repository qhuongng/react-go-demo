package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"chi-mysql-boilerplate/internal/database"
)

type Server struct {
	port int
	db   database.Db
}

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	dbService, err := database.New()
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize database: %v", err))
	}

	NewServer := &Server{
		port: port,
		db:   dbService,
	}

	// declare server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
