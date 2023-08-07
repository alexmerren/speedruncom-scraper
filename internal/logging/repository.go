package logger

type Logger interface {
	Debug(msg interface{})
	Info(msg interface{})
	Error(msg interface{})
	Fatal(msg interface{})
	With(key string, value interface{}) Logger
	WithError(err error) Logger
	Sync() error
}
