package app

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/shamank/edutour_auth_service/app/config"
	handler "github.com/shamank/edutour_auth_service/app/internal/controller/http"
	"github.com/shamank/edutour_auth_service/app/internal/repository"
	"github.com/shamank/edutour_auth_service/app/internal/service"
	"github.com/shamank/edutour_auth_service/app/pkg/auth"
	"github.com/shamank/edutour_auth_service/app/server"
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Run(configDir string) {

	cfg, err := config.Init(configDir)
	if err != nil {
		//logger....
		return
	}

	db, err := sqlx.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.User, cfg.Postgres.Password, cfg.Postgres.DBName, cfg.Postgres.SSLMode))
	if err != nil {
		logrus.Fatalf("error occurred in connecting to postgres: %s", err.Error())
	}

	repos := repository.NewRepository(db)
	services := service.NewServices(repos)

	tokenManager, err := auth.NewManager(cfg.AuthConfig.JWT.SignedKey)
	if err != nil {
		log.Fatalf("error occured generate tokenManager: %s", err.Error())
	}

	handlers := handler.NewHandler(services, tokenManager)

	srv := server.NewServer(cfg, handlers.InitAPI())

	go func() {
		if err := srv.Start(); err != nil {
			return
		}

	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	const timeout = 5 * time.Second

	ctx, shutdown := context.WithTimeout(context.Background(), timeout)
	defer shutdown()

	if err := srv.Stop(ctx); err != nil {
		return
	}
}
