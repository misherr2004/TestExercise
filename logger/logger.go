package logger

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"os"
)

type LoggerConfig struct {
	Level       string   `yaml:"level"`
	OutputPaths []string `yaml:"output_paths"`
}

func InitLogger(cfg *LoggerConfig) (*zap.Logger, error) {
	const op = "InitLogger"

	var level zapcore.Level
	if err := level.UnmarshalText([]byte(cfg.Level)); err != nil {
		log.Printf("error with unmarshaling in %s: %v", op, err)
		return nil, fmt.Errorf("функция: %s, ошибка: %w", op, err)
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	file, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("error with opening file in %s: %v", op, err)
		return nil, fmt.Errorf("функция %s: %w", op, err)
	}

	fileWriteSyncer := zapcore.AddSync(file)
	consoleWriteSyncer := zapcore.AddSync(os.Stdout)

	core := zapcore.NewTee(
		zapcore.NewCore(
			zapcore.NewConsoleEncoder(encoderConfig),
			fileWriteSyncer,
			level,
		),
		zapcore.NewCore(
			zapcore.NewConsoleEncoder(encoderConfig),
			consoleWriteSyncer,
			level,
		),
	)

	return zap.New(core, zap.AddCaller()), nil
}

func New() *zap.Logger {
	const op = "New"
	cfg := &LoggerConfig{
		Level:       "info",
		OutputPaths: []string{"logs.txt", "stdout"},
	}
	logger, err := InitLogger(cfg)
	if err != nil {
		log.Fatalf("error witn init logger in %s: %v", op, err)
	}
	return logger
}
