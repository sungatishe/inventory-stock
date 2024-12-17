package app

import (
	db "auth-service/db/config"
	"auth-service/internal/repository"
	"auth-service/internal/server"
	"auth-service/internal/service"
	"os"
)

func Run() {
	port := os.Getenv("PORT")

	db.InitDb()

	authRepo := repository.NewUserRepository(db.DB)
	authService := service.NewAuthService(authRepo)

	authServer := server.NewAuthServer(authService)

	authServer.Run(port)
}
