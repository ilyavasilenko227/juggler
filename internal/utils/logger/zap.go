package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Zap *zap.SugaredLogger

func Init() {
	config := zap.NewDevelopmentEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	consoleEncoder := zapcore.NewConsoleEncoder(config)
	core := zapcore.NewCore(consoleEncoder, zapcore.AddSync(zapcore.Lock(os.Stdout)), zap.DebugLevel)
	Zap = zap.New(core).Sugar()

}
