package service

import (
	"context"
	"juggler/internal/config"
	"sync"
)

type dataProcessor interface {
	StartJuggling(ctx context.Context) error
	StopJuggling(cancel context.CancelFunc)
}

type jugglingService struct {
	config         *config.App
	wg             sync.WaitGroup
	availableBalls chan int64
	ballStates     map[int64]string
	mu             sync.Mutex
}

func New(config *config.App) dataProcessor {
	availableBalls := make(chan int64, config.N)
	ballStates := make(map[int64]string)

	for i := int64(1); i <= config.N; i++ {
		availableBalls <- i
		ballStates[i] = "в руках"
	}

	return &jugglingService{
		config:         config,
		availableBalls: availableBalls,
		ballStates:     ballStates,
	}
}
