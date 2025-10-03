package logging

import (
	"io"
	"os"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"achievement-management/internal/config"
)

// Logger ログ機能のインターフェース
type Logger interface {
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Warn(args ...interface{})
	Warnf(format string, args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	WithField(key string, value interface{}) Logger
	WithFields(fields map[string]interface{}) Logger
}

// LogrusLogger logrusを使用したLogger実装
type LogrusLogger struct {
	logger *logrus.Logger
	entry  *logrus.Entry
}

// NewLogger 新しいLoggerを作成
func NewLogger(config *config.Config) (Logger, error) {
	logger := logrus.New()
	
	// ログレベルを設定
	level, err := logrus.ParseLevel(config.Logging.Level)
	if err != nil {
		return nil, err
	}
	logger.SetLevel(level)
	
	// ログフォーマットを設定
	switch strings.ToLower(config.Logging.Format) {
	case "json":
		logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: time.RFC3339,
		})
	case "text":
		logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: time.RFC3339,
		})
	default:
		logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: time.RFC3339,
		})
	}
	
	// 出力先を設定
	switch strings.ToLower(config.Logging.Output) {
	case "stdout":
		logger.SetOutput(os.Stdout)
	case "stderr":
		logger.SetOutput(os.Stderr)
	default:
		// ファイルパスとして扱う
		file, err := os.OpenFile(config.Logging.Output, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return nil, err
		}
		logger.SetOutput(file)
	}
	
	return &LogrusLogger{
		logger: logger,
		entry:  logrus.NewEntry(logger),
	}, nil
}

// NewLoggerWithOutput 出力先を指定してLoggerを作成
func NewLoggerWithOutput(config *config.Config, output io.Writer) Logger {
	logger := logrus.New()
	
	// ログレベルを設定
	level, _ := logrus.ParseLevel(config.Logging.Level)
	logger.SetLevel(level)
	
	// ログフォーマットを設定
	switch strings.ToLower(config.Logging.Format) {
	case "json":
		logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: time.RFC3339,
		})
	case "text":
		logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: time.RFC3339,
		})
	default:
		logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: time.RFC3339,
		})
	}
	
	logger.SetOutput(output)
	
	return &LogrusLogger{
		logger: logger,
		entry:  logrus.NewEntry(logger),
	}
}

// Debug デバッグレベルのログを出力
func (l *LogrusLogger) Debug(args ...interface{}) {
	l.entry.Debug(args...)
}

// Debugf フォーマット付きデバッグレベルのログを出力
func (l *LogrusLogger) Debugf(format string, args ...interface{}) {
	l.entry.Debugf(format, args...)
}

// Info 情報レベルのログを出力
func (l *LogrusLogger) Info(args ...interface{}) {
	l.entry.Info(args...)
}

// Infof フォーマット付き情報レベルのログを出力
func (l *LogrusLogger) Infof(format string, args ...interface{}) {
	l.entry.Infof(format, args...)
}

// Warn 警告レベルのログを出力
func (l *LogrusLogger) Warn(args ...interface{}) {
	l.entry.Warn(args...)
}

// Warnf フォーマット付き警告レベルのログを出力
func (l *LogrusLogger) Warnf(format string, args ...interface{}) {
	l.entry.Warnf(format, args...)
}

// Error エラーレベルのログを出力
func (l *LogrusLogger) Error(args ...interface{}) {
	l.entry.Error(args...)
}

// Errorf フォーマット付きエラーレベルのログを出力
func (l *LogrusLogger) Errorf(format string, args ...interface{}) {
	l.entry.Errorf(format, args...)
}

// WithField フィールド付きのLoggerを返す
func (l *LogrusLogger) WithField(key string, value interface{}) Logger {
	return &LogrusLogger{
		logger: l.logger,
		entry:  l.entry.WithField(key, value),
	}
}

// WithFields 複数フィールド付きのLoggerを返す
func (l *LogrusLogger) WithFields(fields map[string]interface{}) Logger {
	return &LogrusLogger{
		logger: l.logger,
		entry:  l.entry.WithFields(fields),
	}
}

// AccessLogger アクセスログ用のLogger
type AccessLogger struct {
	logger Logger
}

// NewAccessLogger アクセスログ用のLoggerを作成
func NewAccessLogger(config *config.Config) (*AccessLogger, error) {
	// アクセスログ用の設定を作成
	accessConfig := *config
	accessConfig.Logging.Level = "info"
	
	logger, err := NewLogger(&accessConfig)
	if err != nil {
		return nil, err
	}
	
	return &AccessLogger{
		logger: logger,
	}, nil
}

// LogRequest HTTPリクエストをログに記録
func (a *AccessLogger) LogRequest(method, path, remoteAddr string, statusCode int, duration time.Duration) {
	a.logger.WithFields(map[string]interface{}{
		"method":      method,
		"path":        path,
		"remote_addr": remoteAddr,
		"status_code": statusCode,
		"duration_ms": duration.Milliseconds(),
		"type":        "access",
	}).Info("HTTP request")
}

// ErrorLogger エラーログ用のLogger
type ErrorLogger struct {
	logger Logger
}

// NewErrorLogger エラーログ用のLoggerを作成
func NewErrorLogger(config *config.Config) (*ErrorLogger, error) {
	// エラーログ用の設定を作成
	errorConfig := *config
	errorConfig.Logging.Level = "error"
	
	logger, err := NewLogger(&errorConfig)
	if err != nil {
		return nil, err
	}
	
	return &ErrorLogger{
		logger: logger,
	}, nil
}

// LogError エラーをログに記録
func (e *ErrorLogger) LogError(operation, component string, err error, fields map[string]interface{}) {
	logFields := map[string]interface{}{
		"operation": operation,
		"component": component,
		"error":     err.Error(),
		"type":      "error",
	}
	
	// 追加フィールドをマージ
	for k, v := range fields {
		logFields[k] = v
	}
	
	e.logger.WithFields(logFields).Error("Operation failed")
}

// LogDatabaseError データベースエラーをログに記録
func (e *ErrorLogger) LogDatabaseError(operation, table string, err error) {
	e.LogError(operation, "database", err, map[string]interface{}{
		"table": table,
	})
}

// LogServiceError サービスエラーをログに記録
func (e *ErrorLogger) LogServiceError(service, operation string, err error) {
	e.LogError(operation, "service", err, map[string]interface{}{
		"service": service,
	})
}

// LogAPIError APIエラーをログに記録
func (e *ErrorLogger) LogAPIError(endpoint, method string, statusCode int, err error) {
	e.LogError("api_request", "handler", err, map[string]interface{}{
		"endpoint":    endpoint,
		"method":      method,
		"status_code": statusCode,
	})
}