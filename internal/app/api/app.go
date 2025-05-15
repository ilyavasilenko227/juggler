package api

import (
	"context"
	"fmt"
	"juggler/internal/config"
	"juggler/internal/service"
	"juggler/internal/utils/logger"

	"go.uber.org/zap"
)

func init() {
	logger.Init()
}

func Run() {
	appCfg, err := config.Config()
	if err != nil {
		panic(fmt.Errorf("can't configure application: %w", err))
	}

	processor := service.New(&appCfg)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	jugglingDone := make(chan struct{})

	go func() {
		if err := processor.StartJuggling(ctx); err != nil {
			logger.Zap.Error("Juggling error", zap.Error(err))
		}
		close(jugglingDone)
	}()

	<-jugglingDone
	processor.StopJuggling(cancel)
}
