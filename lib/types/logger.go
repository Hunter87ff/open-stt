package types

type Colors struct {
	Bright    string
	Dim       string
	Italic    string
	Underline string
	Inverse   string
	Hidden    string
	Reset     string
	Black     string
	Red       string
	Green     string
	Yellow    string
	Blue      string
	Magenta   string
	Cyan      string
	White     string
}

type LoggerLevels struct {
	DEBUG int
	INFO  int
	WARN  int
	ERROR int
	FATAL int
}

type LogPrefixes struct {
	Debug string
	Info  string
	Warn  string
	Error string
	Fatal string
}
