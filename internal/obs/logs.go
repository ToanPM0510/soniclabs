package obs

import "go.uber.org/zap"

func NewLogger() (*zap.Logger, func()) {
	logger, _ := zap.NewProduction()
	cleanup := func() {
		_ = logger.Sync()
	}
	return logger, cleanup
}
