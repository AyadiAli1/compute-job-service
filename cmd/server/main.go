// Command server starts the compute job gRPC service.
package main

import (
	"log/slog"
	"os"
)

// version is injected at build time via -ldflags "-X main.version=...".
var version = "dev"

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	logger.Info("compute-job-service starting", "version", version)
}
