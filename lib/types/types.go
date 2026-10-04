package types

type Platforms struct {
	Win   string
	Linux string
	Mac   string
}

type Hotkeys struct {
	Record string
	Quit   string
}

type STTServer struct {
	Port   string
	Model  string
	Config *string
	Endpoint string
}

type Config struct {
	Hotkeys             Hotkeys
	STT                 STTServer
	RecordingSampleRate string
}
