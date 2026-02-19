package logger

import (
	"context"
	"os"

	"github.com/sirupsen/logrus"
)

type contextKey string

const (
	RequestIDKey contextKey = "request_id"
	UserIDKey    contextKey = "user_id"
)

var log = logrus.New()

func init() {
	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(logrus.InfoLevel)
}

func WithContext(ctx context.Context) *logrus.Entry {
	fields := logrus.Fields{}

	if reqID, ok := ctx.Value(RequestIDKey).(string); ok {
		fields["request_id"] = reqID
	}
	if userID, ok := ctx.Value(UserIDKey).(string); ok {
		fields["user_id"] = userID
	}

	return log.WithFields(fields)
}

func Info(ctx context.Context, msg string, fields map[string]interface{}) {
	WithContext(ctx).WithFields(fields).Info(msg)
}

func Error(ctx context.Context, msg string, fields map[string]interface{}) {
	WithContext(ctx).WithFields(fields).Error(msg)
}

func Warn(ctx context.Context, msg string, fields map[string]interface{}) {
	WithContext(ctx).WithFields(fields).Warn(msg)
}

func Debug(ctx context.Context, msg string, fields map[string]interface{}) {
	WithContext(ctx).WithFields(fields).Debug(msg)
}
