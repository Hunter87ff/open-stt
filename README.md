# open-stt

Push-to-talk speech-to-text. Hold **F6** to record, release to transcribe and paste the result at the cursor. **F7** quits.

## Requirements

- Go, `ffmpeg`, `curl`
- Paste helper: `xdotool` (X11) or `wtype` (Wayland)
- `llama-server` with the ASR model (the transcription backend)

## Setup

```sh
go mod tidy
go build -o open-stt .
```

## Run

Terminal 1 — start the ASR server (serves OpenAI-compatible endpoint on `:8080`):

```sh
llama-server -hf hunter87/Qwen3-ASR-0.6B-GGUF -c 350 -np 1 --flash-attn on --n-gpu-layers 999 --no-mmap
```

Terminal 2 — run the app:

```sh
./open-stt
```

Hold `F6`, speak, release. The transcript is typed at the cursor. `F7` to quit.

## Configuration

`lib/config/config.go`:

| Var | Default |
|-----|---------|
| `RecordHotkey` | `f6` |
| `QuitHotkey` | `f7` |
| `TranscribeEndpoint` | `http://127.0.0.1:8080/v1/audio/transcriptions` |
| `TranscribeModel` | `qwen3-asr-0.6b` |
| `RecordingSampleRate` | `16000` |

## Notes

- Requires a Linux desktop with PulseAudio/PipeWire for mic capture.
- The ASR server must be running, or transcription will fail (error is printed to stderr).
