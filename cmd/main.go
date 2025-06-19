package main

import (
	"awesomeProject/internal/apiServer/httpServer"
	"awesomeProject/internal/config"
	"awesomeProject/internal/repository/postgresql"
	"awesomeProject/internal/userservice"
	"awesomeProject/logger"
)

func main() {
	cfg := config.MustLoad()

	log := logger.New()

	db, _ := postgresql.NewPostgreSQL(cfg, log)

	userService := userservice.NewUserService(db, log)

	defer db.Close()

	server := httpServer.NewServer(userService, log)

	server.Run(cfg)
}
