package service

import (
	"context"
	"juggler/internal/utils/logger"
	"math/rand"
	"time"
)

func (s *jugglingService) StartJuggling(ctx context.Context) error {
	jugglingDuration := time.Duration(s.config.T) * time.Minute
	timer := time.NewTimer(jugglingDuration)
	defer timer.Stop()

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	logger.Zap.Infof("Жонглирование начато, мячи %d, минуты %d", s.config.N, s.config.T)
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-timer.C:
			logger.Zap.Info("Время жонглирования истекло. Новые мячи не подбрасываются.")
			return nil
		case <-ticker.C:
			id := <-s.availableBalls
			s.printBallStates()

			s.mu.Lock()
			s.ballStates[id] = "в полёте"
			s.mu.Unlock()

			s.wg.Add(1)
			go func(ballID int64) {
				defer s.wg.Done()
				s.runBall(ctx, ballID)
			}(id)
		}
	}
}

func (s *jugglingService) runBall(ctx context.Context, id int64) {
	flightTime := time.Duration(rand.Intn(6)+5) * time.Second
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	logger.Zap.Infof("Мяч #%d: начал лететь", id)

	for i := 0; i < int(flightTime.Seconds()); i++ {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			logger.Zap.Infof("Мяч #%d: %d/%d сек", id, i+1, int(flightTime.Seconds()))
		}
	}

	logger.Zap.Infof("Мяч #%d: упал", id)

	s.mu.Lock()
	s.ballStates[id] = "в руках"
	s.mu.Unlock()

	s.availableBalls <- id
}

func (s *jugglingService) printBallStates() {
	s.mu.Lock()
	defer s.mu.Unlock()

	inAir := 0
	inHand := 0
	for _, state := range s.ballStates {
		switch state {
		case "в полёте":
			inAir++
		case "в руках":
			inHand++
		}
	}

	logger.Zap.Infow("Текущие состояния мячей",
		"в_полёте", inAir,
		"в_руках", inHand,
		"статусы", s.ballStates,
	)
}

func (s *jugglingService) StopJuggling(cancel context.CancelFunc) {
	s.wg.Wait()
	cancel()
}
