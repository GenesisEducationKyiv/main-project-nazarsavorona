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
		return "Debug"
	case Info:
		return "Info"
	case Error:
		return "Error"
	default:
		return "Unknown"
	}
}
