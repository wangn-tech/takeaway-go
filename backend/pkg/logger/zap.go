package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"takeaway-go/internal/app/config"
)

// Log 是一个全局的 zap Logger 实例
var Log *zap.Logger

// InitLogger 根据配置文件初始化 zap Logger
func InitLogger() {
	var core zapcore.Core
	var level zapcore.Level

	// 从配置中获取日志级别
	if err := level.UnmarshalText([]byte(config.AppConf.Log.Level)); err != nil {
		level = zapcore.InfoLevel // 解析失败则默认为 info
	}

	// 配置 Encoder
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder // 统一时间格式

	// 根据配置设置日志格式
	var encoder zapcore.Encoder
	if config.AppConf.Log.Format == "console" {
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder // 在 console 模式下使用带颜色的日志级别
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	}

	core = zapcore.NewCore(encoder, zapcore.Lock(os.Stdout), level)

	// 根据配置决定是否添加调用者信息和堆栈跟踪
	var options []zap.Option
	if config.AppConf.Log.ShowLine {
		options = append(options, zap.AddCaller())
	}
	if config.AppConf.Log.Stacktrace {
		options = append(options, zap.AddStacktrace(zap.ErrorLevel))
	}

	// 构造日志
	Log = zap.New(core, options...)
	zap.ReplaceGlobals(Log) // 替换 zap 的全局 Logger
}
