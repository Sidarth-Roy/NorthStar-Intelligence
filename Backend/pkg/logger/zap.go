package logger

import (
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	log  *zap.Logger
	once sync.Once
)

func InitLogger() {
	once.Do(func() {
		encoderConfig := zap.NewProductionEncoderConfig()
		encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		encoderConfig.TimeKey = "timestamp"
		encoderConfig.LevelKey = "level"
		encoderConfig.MessageKey = "message"
		encoderConfig.CallerKey = "caller"      // Shows file:line
		encoderConfig.FunctionKey = "func"      // Shows function name

		core := zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			zapcore.AddSync(os.Stdout),
			zap.InfoLevel,
		)

		// AddCaller adds the file/line, AddStacktrace adds trace on Errors
		log = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel)).With(
			zap.String("service", "northstar-backend"),
			zap.String("environment", os.Getenv("APP_ENV")),
		)
	})
}

func Get() *zap.Logger {
	return log
}