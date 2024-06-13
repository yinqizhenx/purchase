package logx

import (
	"fmt"
	"os"
	"path"
	"time"

	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/log"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"purchase/pkg/utils"
)

var (
	// defaultLogFolder        = "../../logs"
	defaultMaxAge     int64 = 30
	defaultRotateTime int64 = 24
	logPathEnv              = "LOG_FILE_PATH"
	logName                 = "%Y%m%d.log"
	logLinkName             = "log"
)

var _ log.Logger = (*ZapLogger)(nil)

type ZapLogger struct {
	logger *zap.Logger
}

func (zl *ZapLogger) Log(level log.Level, kvs ...interface{}) error {
	kvLen := len(kvs)
	// if kvLen == 0 || kvLen%2 != 0 {
	// 	zl.logger.Warn(fmt.Sprint("kvs must start with a message, and the rest appear in pairs: ", kvs))
	// 	return nil
	// }

	data := make([]zap.Field, 0, kvLen/2)
	if kvLen > 1 {
		for i := 0; i < kvLen; i += 2 {
			data = append(data, zap.Any(fmt.Sprint(kvs[i]), kvs[i+1]))
		}
	}

	msg := ""
	switch level {
	case log.LevelDebug:
		zl.logger.Debug(msg, data...)
	case log.LevelInfo:
		zl.logger.Info(msg, data...)
	case log.LevelWarn:
		zl.logger.Warn(msg, data...)
	case log.LevelError:
		zl.logger.Error(msg, data...)
	case log.LevelFatal:
		zl.logger.Fatal(msg, data...)
	}
	return nil
}

func (zl *ZapLogger) Sync() error {
	return zl.logger.Sync()
}

func (zl *ZapLogger) Close() error {
	return zl.Sync()
}

func newZapLog(cfg config.Config) *zap.Logger {
	encoder := buildEncoder()
	writer := buildWriter(cfg)
	core := zapcore.NewCore(encoder, writer, zapcore.InfoLevel)
	logger := zap.New(core)
	return logger
}

func buildEncoder() zapcore.Encoder {
	encConf := zap.NewProductionEncoderConfig()
	encConf.EncodeTime = zapcore.ISO8601TimeEncoder
	encConf.MessageKey = zapcore.OmitKey // 移除zap默认输出的msg
	return zapcore.NewJSONEncoder(encConf)
}

func buildWriter(cfg config.Config) zapcore.WriteSyncer {
	logFolder := GetLogFile()
	if f := os.Getenv(logPathEnv); f != "" {
		logFolder = f
	} else {
		log.Warnf("warning: no logger path found in env, use default value: %v \n", logFolder)
	}
	logPath := path.Join(logFolder, logName)
	linkName := path.Join(logFolder, logLinkName)

	maxAge, err := cfg.Value("logger.maxAge").Int()
	if err != nil || maxAge == 0 {
		log.Warnf("warning: no maxAge found in cfg, use default value: %v \n", defaultMaxAge)
		maxAge = defaultMaxAge
	}

	rotateTime, err := cfg.Value("logger.rotateTime").Int()
	if err != nil || rotateTime == 0 {
		log.Warnf("warning: no rotateTime found in cfg, use default value: %v \n", defaultRotateTime)
		rotateTime = defaultRotateTime
	}

	hook, err := rotatelogs.New(
		logPath,
		rotatelogs.WithLinkName(linkName),
		rotatelogs.WithMaxAge(time.Hour*24*time.Duration(maxAge)),
		rotatelogs.WithRotationTime(time.Hour*time.Duration(rotateTime)),
	)
	return zapcore.AddSync(hook)
}

func GetLogFile() string {
	root, err := utils.GetRootPath()
	if err != nil {
		panic(err)
	}
	return root + "/logs"
}
