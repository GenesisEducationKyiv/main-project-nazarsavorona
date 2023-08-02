package logger

type Level int

const (
	Debug Level = iota
	Info
	Error
)

func (l Level) String() string {
	switch l {
	case Debug:
		return "debug"
	case Info:
		return "info"
	case Error:
		return "error"
	default:
		return "unknown"
	}
}
