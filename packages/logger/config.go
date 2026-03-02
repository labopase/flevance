package logger

type Mode string

const (
	ModeProduction  Mode = "production"
	ModeStaging     Mode = "staging"
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
	Level        LogLevel `mapstructure:"level"`
	EnableCaller bool     `mapstructure:"enable_caller"`
	EnableTrace  bool     `mapstructure:"enable_trace"`
	Environment  Mode     `mapstructure:"environment"`
}

func DefaultConfig() *Config {
	return &Config{
		Level:        DebugLevel,
		EnableCaller: false,
		EnableTrace:  false,
		Environment:  ModeProduction,
	}
}
