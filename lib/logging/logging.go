package logging

import (
	"open-stt/lib/constants"
	"open-stt/lib/types"
)

type _Logger struct {
	Level int
	Levels  types.LoggerLevels
	Prefixes types.LogPrefixes
}

var Logger  = &_Logger{
	Level: constants.LogLevels.INFO,
	Levels: constants.LogLevels,
	Prefixes: constants.LogPrefixes,
}

/*
private log method that chesk the log level and prints the message with the appropriate prefix
*/
func (l *_Logger) log(level int, message string){
	if level < l.Level {
		return
	}

	var prefix string
	switch level {
		case l.Levels.DEBUG:
			prefix = l.Prefixes.Debug
		case l.Levels.INFO:
			prefix = l.Prefixes.Info
		case l.Levels.WARN:
			prefix = l.Prefixes.Warn
		case l.Levels.ERROR:
			prefix = l.Prefixes.Error
		case l.Levels.FATAL:
			prefix = l.Prefixes.Fatal
		default:
			prefix = ""
	}

	println(prefix + " " + message)
}

func (l *_Logger) Debug(message string){
	l.log(l.Levels.DEBUG, message)
}

func (l *_Logger) Info(message string){
	l.log(l.Levels.INFO, message)
}

func (l *_Logger) Warn(message string){
	l.log(l.Levels.WARN, message)
}

func (l *_Logger) Error(message string){
	l.log(l.Levels.ERROR, message)
}

func (l *_Logger) Fatal(message string){
	l.log(l.Levels.FATAL, message)
}