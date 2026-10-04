/*this package contains application configurations*/
package config

import (
	"fmt"
	"open-stt/lib/types"
)

var (
	Hotkeys = types.Hotkeys{
		Record: "f6",
		Quit:   "f7",
	}

	STT = types.STTServer{
		Port:     "8080",
		Model:    "hunter87/Qwen3-ASR-0.6B-GGUF",
		Endpoint: fmt.Sprintf("http://localhost:%s/v1/audio/transcriptions", "8080"),
		Config:   nil,
	}
	RecordingSampleRate = "16000"
)
