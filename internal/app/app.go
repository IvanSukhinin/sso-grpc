package app

import (
	"log/slog"
	"os"
	"os/signal"
	ssoapp "sso/internal/app/sso"
	"sso/internal/config"
	"syscall"
)

const envDev = "dev"

type App struct {
	sso *ssoapp.App
	cfg *config.Config
	log *slog.Logger
}

func New() *App {
	cfg := config.MustLoad()
	log := setupLogger(cfg.Env)
	return &App{
		sso: ssoapp.New(log, cfg.GRPC.Port, cfg.Db, cfg.TokenTTL),
		cfg: cfg,
		log: log,
	}
}

func (a *App) Run() {
	a.log.Info("start", slog.Any("config", a.cfg))

	go a.sso.GRPCServer.MustRun()

	// gracefull shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	sign := <-stop

	a.sso.GRPCServer.Stop()
	a.log.Info("app stopped by signal " + sign.String())
}

func setupLogger(env string) *slog.Logger {
	level := slog.LevelInfo
	if env == envDev {
		level = slog.LevelDebug
	}
	return slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: level}),
	)
}
