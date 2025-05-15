package service

import (
	"context"
	"juggler/internal/config"
	"sync"
)

type dataProcessor interface {
	StartJuggling(ctx context.Context) error
	StopJuggling()
}

type jugglingService struct {
	config         *config.App
	ctx            context.Context
	cancel         context.CancelFunc
	wg             sync.WaitGroup
	availableBalls chan int64
	ballStates     map[int64]string
	mu             sync.Mutex
}

func New(config *config.App) dataProcessor {
	ctx, cancel := context.WithCancel(context.Background())
	availableBalls := make(chan int64, config.N)
	ballStates := make(map[int64]string)

	for i := int64(1); i <= config.N; i++ {
		availableBalls <- i
		ballStates[i] = "в руках"
	}

	return &jugglingService{
		config:         config,
		ctx:            ctx,
		cancel:         cancel,
		availableBalls: availableBalls,
		ballStates:     ballStates,
	}
}
