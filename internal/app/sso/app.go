package ssoapp

import (
	"log/slog"
	grpcapp "sso/internal/app/grpc"
	"sso/internal/config"
	"sso/internal/services/auth"
	"sso/internal/storage/postgresql"
	"time"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(
	log *slog.Logger,
	grpcPort int,
	cfg config.Db,
	tokenTTL time.Duration,
) *App {
	// data layer
	storage, err := postgresql.New(cfg)
	if err != nil {
		panic(err)
	}

	// service layer
	authService := auth.New(log, storage, storage, storage, tokenTTL)

	// transport layer
	grpcApp := grpcapp.New(log, authService, grpcPort)

	return &App{
		GRPCServer: grpcApp,
	}
}
