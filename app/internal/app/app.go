package app

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/patrickmn/go-cache"
	"github.com/shamank/eduTour-backend/app/config"
	handler "github.com/shamank/eduTour-backend/app/internal/delivery/http"
	"github.com/shamank/eduTour-backend/app/internal/repository"
	"github.com/shamank/eduTour-backend/app/internal/service"
	"github.com/shamank/eduTour-backend/app/pkg/auth"
	"github.com/shamank/eduTour-backend/app/pkg/hash"
	"github.com/shamank/eduTour-backend/app/server"
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
		logrus.Fatalf("error occured init config: %s", err.Error())
	}

	db, err := sqlx.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.User, cfg.Postgres.Password, cfg.Postgres.DBName, cfg.Postgres.SSLMode))
	if err != nil {
		logrus.Fatalf("error occurred in connecting to postgres: %s", err.Error())
	}

	repos := repository.NewRepository(db)

	memcache := cache.New(5*time.Minute, 10*time.Minute)

	JWTConfig := cfg.AuthConfig.JWT

	tokenManager, err := auth.NewManager(JWTConfig.SignedKey, JWTConfig.AccessTokenTTL, JWTConfig.RefreshTokenTTL)
	if err != nil {
		log.Fatalf("error occured generate tokenManager: %s", err.Error())
	}
	hasher := hash.NewSHA1Hasher(cfg.AuthConfig.PasswordSalt)

	services := service.NewServices(repos, memcache, hasher, tokenManager)

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
