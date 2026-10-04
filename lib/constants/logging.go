package constants

import (
	"open-stt/lib/types"
)

var LogColors = types.Colors{
	White:     "\x1b[37m",
	Reset:     "\x1b[0m",
	Bright:    "\x1b[1m",
	Dim:       "\x1b[2m",
	Italic:    "\x1b[3m",
	Underline: "\x1b[4m",
	Inverse:   "\x1b[7m",
	Hidden:    "\x1b[8m",
	Red:       "\x1b[31m",
	Green:     "\x1b[32m",
	Yellow:    "\x1b[33m",
	Blue:      "\x1b[34m",
	Magenta:   "\x1b[35m",
	Cyan:      "\x1b[36m",
}

var LogLevels = types.LoggerLevels{
	DEBUG: 0,
	INFO:  1,
	WARN:  2,
	ERROR: 3,
	FATAL: 4,
}

var LogPrefixes = types.LogPrefixes{
	Debug: "\033[1;34m[DEBUG]\033[0m",
	Info:  "\033[1;32m[INFO]\033[0m",
	Warn:  "\033[1;33m[WARN]\033[0m",
	Error: "\033[1;31m[ERROR]\033[0m",
	Fatal: "\033[1;31m[FATAL]\033[0m",
}
