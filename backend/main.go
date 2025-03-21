package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/KennethTrinh/go-vite-app/config"
	"github.com/KennethTrinh/go-vite-app/initializers"
	"github.com/KennethTrinh/go-vite-app/router"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	fiberLogger "github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/rs/zerolog/log"
)

func init() {
	err := config.LoadConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load environment variables!")
	}
	initializers.InitLogger()
	initializers.ConnectDB()
}

func main() {
	app := fiber.New()
	app.Use(fiberLogger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     config.Env.AllowedOrigins,
		AllowCredentials: true,
	}))

	router.SetupRoutes(app)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	go func() {
		<-sigCh
		log.Info().Msg("Shutting down server...")
		if err := app.Shutdown(); err != nil {
			log.Fatal().Err(err).Msg("Error shutting down server")
		}
	}()

	log.Info().Msg("Starting server on port " + config.Env.ServerPort)
	if err := app.Listen(":" + config.Env.ServerPort); err != nil && err != fiber.ErrGracefulTimeout {
		log.Fatal().Err(err).Msg("Failed to start server!")
	}
}
