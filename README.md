# open-stt

Push-to-talk speech-to-text. Hold **F6** to record, release to transcribe and paste the result directly at the cursor. **F7** quits.

## Features

- **Push-to-Talk Transcription**: Hold hotkey to record, release to transcribe and auto-paste at the active cursor.
- **Automatic ASR Backend Management**: Spawns and manages `llama-server` automatically in the background using Hugging Face GGUF models.
- **Cross-Platform Typing/Pasting**:
  - **Linux**: X11 (`xdotool`) and Wayland (`wtype`).
  - **Windows**: PowerShell SendKeys.
  - **macOS**: AppleScript (`osascript`).
- **Clipboard Sync**: Transcribed text is also copied directly to your clipboard.
- **Color-Coded Leveled Logging**: Structured logging with `DEBUG`, `INFO`, `WARN`, `ERROR`, and `FATAL` levels.

## Requirements

- **Go** (1.20+)
- **ffmpeg** (in `PATH` for microphone audio capture)
- **llama-server** (in `PATH`, from [llama.cpp](https://github.com/ggerganov/llama.cpp))
- **Paste/Type Helper**:
  - **Linux**: `xdotool` (X11) or `wtype` (Wayland)
  - **Windows**: PowerShell (built-in)
  - **macOS**: AppleScript (built-in)

## Setup

```sh
go mod tidy
go build -o open-stt .
```

## Run

```sh
./open-stt
```

1. On launch, `open-stt` verifies required dependencies (`ffmpeg`, `llama-server`) and boots the ASR thread (`llama-server`) with the configured model on port `8080`.
2. Hold **F6** to speak. Release to stop recording.
3. The audio is transcribed and typed at the cursor position (and saved to clipboard).
4. Press **F7** to quit the application.

## Configuration

Settings are defined in [`lib/config/config.go`](file:///home/hunter87/Desktop/projects/side-projects/open-stt/lib/config/config.go):

| Setting | Default | Description |
|---|---|---|
| `Hotkeys.Record` | `"f6"` | Hotkey to hold while speaking |
| `Hotkeys.Quit` | `"f7"` | Hotkey to exit the application |
| `STT.Port` | `"8080"` | Port for the local `llama-server` ASR service |
| `STT.Model` | `"hunter87/Qwen3-ASR-0.6B-GGUF"` | Hugging Face model repository / GGUF model |
| `STT.Endpoint` | `"http://localhost:8080/v1/audio/transcriptions"` | OpenAI-compatible transcription API endpoint |
| `RecordingSampleRate` | `"16000"` | Audio recording sample rate in Hz |

## Notes

- **Microphone**: Requires PulseAudio or PipeWire on Linux for default microphone capture via `ffmpeg`.
- **Server Startup**: `llama-server` runs with GPU acceleration enabled by default (`--n-gpu-layers 999` and `--flash-attn on`).
