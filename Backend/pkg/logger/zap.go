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
		var encoder zapcore.Encoder
		var stacktraceLevel zapcore.LevelEnabler
		
		env := os.Getenv("APP_ENV")

		if env == "production" {
			// --- PRODUCTION: Keep your original JSON logic ---
			encoderConfig := zap.NewProductionEncoderConfig()
			encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
			encoderConfig.TimeKey = "timestamp"
			encoderConfig.LevelKey = "level"
			encoderConfig.MessageKey = "message"
			encoderConfig.CallerKey = "caller"
			encoderConfig.FunctionKey = "func"
			
			encoder = zapcore.NewJSONEncoder(encoderConfig)
			stacktraceLevel = zapcore.ErrorLevel // Stacktrace on every error
		} else {
			// --- DEVELOPMENT: Optimized for readability ---
			encoderConfig := zap.NewDevelopmentEncoderConfig()
			encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder // Adds colors
			encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("15:04:05") // Just the time
			
			encoder = zapcore.NewConsoleEncoder(encoderConfig)
			stacktraceLevel = zapcore.DPanicLevel // No stacktraces for standard Errors
		}

		core := zapcore.NewCore(
			encoder,
			zapcore.AddSync(os.Stdout),
			zap.DebugLevel,
		)

		// Passing zapcore.LevelEnabler directly here
		log = zap.New(core, zap.AddCaller(), zap.AddStacktrace(stacktraceLevel)).With(
			zap.String("service", "northstar-backend"),
			zap.String("environment", env),
		)
	})
}

func Get() *zap.Logger {
	return log
}

// Sync flushes any buffered log entries. 
// Call this in your main.go: defer logger.Sync()
func Sync() {
	if log != nil {
		_ = log.Sync()
	}
}