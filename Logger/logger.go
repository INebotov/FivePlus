package Logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var null = &zap.Logger{}

func GetLogger(config zap.Config) (*zap.Logger, error) {
	config.EncoderConfig.FunctionKey = zapcore.OmitKey
	config.EncoderConfig.LineEnding = zapcore.DefaultLineEnding
	config.EncoderConfig.EncodeLevel = zapcore.LowercaseLevelEncoder
	config.EncoderConfig.EncodeTime = zapcore.EpochTimeEncoder
	config.EncoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
	config.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder

	loger, err := config.Build()
	if err != nil {
		return null, err
	}

	loger.Error("Hello!")
	return loger, nil
}
