package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/k8s-2025-pschoeppner/ctf/pkg/flagset"
	"github.com/k8s-2025-pschoeppner/ctf/pkg/k8s"
	"github.com/k8s-2025-pschoeppner/ctf/pkg/types"
)

const (
	envPodName      = "POD_NAME"
	envPodNamespace = "POD_NAMESPACE"
)

var (
	flagName = flag.String("flag", "", "The name of the flag to capture")
	server   = flag.String("server", "http://ctf-server:8080", "The server to capture the flag on")
)

func getEnvVar(name string) (string, error) {
	value := os.Getenv(name)
	if value == "" {
		return "", fmt.Errorf("env var %s not set", name)
	}
	return value, nil
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	flag.Parse()

	podName, err := getEnvVar(envPodName)
	if err != nil {
		logger.Error("read env var", slog.String("env", envPodName), slog.String("err", err.Error()))
		os.Exit(1)
	}
	podNamespace, err := getEnvVar(envPodNamespace)
	if err != nil {
		logger.Error("read env var", slog.String("env", envPodNamespace), slog.String("err", err.Error()))
		os.Exit(1)
	}

	if *flagName == "" {
		logger.Error("flag not set")
		os.Exit(1)
	}

	client, err := k8s.InClusterConfig()
	if err != nil {
		logger.Error("initialise in-cluster kubernetes client", slog.String("err", err.Error()))
		os.Exit(1)
	}

	flagSet := flagset.NewFlagSet(client, logger)

	flag, ok := flagSet[*flagName]
	if !ok {
		logger.Error("flag not found", slog.String("flag", *flagName), "flag", flagSet)
		os.Exit(1)
	}

	ctx := context.Background()

	clientID := uuid.New().String()
	req := types.Request{
		ID:           clientID,
		PodName:      podName,
		PodNamespace: podNamespace,
		Args:         make(map[string]string),
	}

	for _, fulfiller := range flag.Fulfillers {
		err = fulfiller(ctx, req, client)
		if err != nil {
			logger.Error("fulfilling condition", slog.String("flag", *flagName), slog.String("err", err.Error()))
			os.Exit(1)
		}
	}

	body, err := req.ToJSON()
	if err != nil {
		logger.Error("marshalling request", slog.String("flag", *flagName), slog.String("err", err.Error()))
		os.Exit(1)
	}

	resp, err := http.Post(fmt.Sprintf("%s/%s", *server, *flagName), "application/json", bytes.NewReader(body))
	if err != nil {
		logger.Error("sending request", slog.String("flag", *flagName), slog.String("err", err.Error()))
		os.Exit(1)
	}
	defer resp.Body.Close()

	flagValue, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error("reading response", slog.String("flag", *flagName), slog.String("err", err.Error()))
		os.Exit(1)
	}

	if resp.StatusCode != http.StatusOK {
		logger.Error("flag capture failed", slog.String("flag", *flagName), slog.String("status", resp.Status), slog.String("body", string(flagValue)))
		os.Exit(1)
	}

	logger.Info("flag captured 🎉🎉🎉", slog.String("flag", *flagName), slog.String("value", string(flagValue)))
}
