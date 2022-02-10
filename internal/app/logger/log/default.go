package log

import "github.com/docplanner/helm-repo-updater/internal/app/logger"

var (
	// Logger is the global logger instance to be used before the container is initialized.
	Logger logger.Logger
)

func init() {
	Logger = logger.NewZeroLogger(logger.NewConsoleZeroLogger("info"))
}
