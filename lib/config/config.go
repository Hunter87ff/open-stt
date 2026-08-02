/*this package contains application configurations*/
package config

var (
	RecordHotkey        = "f6"
	QuitHotkey          = "f7"
	TranscribeEndpoint  = "http://127.0.0.1:8080/v1/audio/transcriptions"
	TranscribeModel     = "qwen3-asr-0.6b"
	RecordingSampleRate = "16000"
)
