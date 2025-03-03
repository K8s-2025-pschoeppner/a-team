package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/k8s-2025-pschoeppner/ctf/pkg/flagset"
	"github.com/k8s-2025-pschoeppner/ctf/pkg/k8s"
)

func main() {
	// init new slog logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	// init new k8s client
	client, err := k8s.InClusterConfig()
	if err != nil {
		logger.Error("initialise in-cluster kubernetes client", slog.String("err", err.Error()))
		os.Exit(1)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	router := http.NewServeMux()

	// TODO
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Handler for /
		w.WriteHeader(http.StatusNotImplemented)
		w.Write([]byte("Hello, world!"))
	})

	flagSet := flagset.NewFlagSet(client, logger)
	cfg, err := loadConfig("/etc/config/config.json")
	if err != nil {
		logger.Error("load config", slog.String("err", err.Error()))
		os.Exit(1)
	}
	if err := newFlagSetFromConfig(cfg, flagSet); err != nil {
		logger.Error("populate flag set from config", slog.String("err", err.Error()))
		os.Exit(1)
	}

	for name, flag := range flagSet {
		router.HandleFunc("/"+name, flag.Handler(ctx))
	}

	serverLogger := slog.NewLogLogger(logger.Handler(), slog.LevelError)
	server := &http.Server{
		Addr:     ":8080",
		Handler:  router,
		ErrorLog: serverLogger,
	}
	go func() {
		logger.InfoContext(ctx, "server started", slog.String("addr", server.Addr))
		if err := server.ListenAndServe(); err != nil {
			logger.ErrorContext(ctx, "server error", slog.String("err", err.Error()))
			os.Exit(1)
		}
		logger.InfoContext(ctx, "server stopped")
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	gracefulCtx, gracefulCancel := context.WithTimeout(ctx, 10*time.Second)
	defer gracefulCancel()

	if err := server.Shutdown(gracefulCtx); err != nil {
		logger.ErrorContext(ctx, "graceful server shutdown", slog.String("err", err.Error()))
		os.Exit(1)
	}
	logger.InfoContext(ctx, "server shutdown")
}
