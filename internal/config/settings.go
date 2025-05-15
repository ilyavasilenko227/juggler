package config

type App struct {
	T int64 `env:"JUGGLER_DURATION_MINUTS"`
	N int64 `env:"JUGGLER_BALL_COUNT"`
}
