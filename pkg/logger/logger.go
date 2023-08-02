package logger

type Logger interface {
	Log(level Level, message string)
}
