package logger

type Mode string

const (
	ModeProduction  Mode = "production"
	ModeDevelopment Mode = "development"
)

type LogLevel string

const (
	DebugLevel LogLevel = "debug"
	InfoLevel  LogLevel = "info"
	WarnLevel  LogLevel = "warn"
	ErrorLevel LogLevel = "error"
	FatalLevel LogLevel = "fatal"
)

type Config struct {
	Level        LogLevel
	EnableCaller bool
	EnableTrace  bool
	Mode         Mode
}

func DefaultConfig(mode Mode) *Config {
	return &Config{
		Level:        DebugLevel,
		EnableCaller: false,
		EnableTrace:  false,
		Mode:         mode,
	}
}
