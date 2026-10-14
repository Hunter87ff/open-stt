# open-stt

Push-to-talk speech-to-text. Hold **F6** to record, release to transcribe with local AI and automatically type the result directly at your cursor. Press **F7** to quit.

---

## Features

- **Push-to-Talk Transcription**: Hold a hotkey to record audio, release to transcribe and auto-type at the active cursor.
- **Automatic ASR Backend Management**: Spawns and manages `llama-server` automatically in the background using Hugging Face GGUF models.
- **Cross-Platform Typing & Pasting**:
  - **Linux**: X11 (`xdotool`) and Wayland (`wtype`).
  - **Windows**: PowerShell SendKeys.
  - **macOS**: AppleScript (`osascript`).
- **Clipboard Sync**: Transcribed text is simultaneously copied directly to your system clipboard.
- **Local & Private**: Completely private inference once the model weights are downloaded.
- **Color-Coded Logging**: Leveled logging (`DEBUG`, `INFO`, `WARN`, `ERROR`, `FATAL`) for quick monitoring and debugging.

---

## Requirements & Prerequisites

`open-stt` requires **Go**, **llama-server** (from `llama.cpp`), **ffmpeg**, and platform-specific typing utilities.

### 1. Install Go (1.20+)

- **Ubuntu / Debian**:
  ```sh
  sudo apt update
  sudo apt install -y golang-go
  ```
  *(Or install the latest official release from [go.dev/dl](https://go.dev/dl/).)*
- **Arch Linux**:
  ```sh
  sudo pacman -S go
  ```
- **macOS** (Homebrew):
  ```sh
  brew install go
  ```
- **Windows**:
  ```powershell
  winget install GoLang.Go
  ```

Verify:
```sh
go version
```

---

### 2. Install `llama.cpp` (`llama-server`)

`open-stt` relies on `llama-server` in your system `PATH` to run local ASR inference.

#### Option A: Download Pre-built Binaries (Easiest)

1. Go to the [llama.cpp Releases page](https://github.com/ggerganov/llama.cpp/releases).
2. Download the pre-built archive matching your OS and hardware:
   - **Linux (x64 CPU)**: `llama-bXXXX-bin-ubuntu-x64.zip`
   - **Linux (CUDA)**: `llama-bXXXX-bin-ubuntu-x64-cuda-cuXX.X.zip`
   - **macOS (Metal)**: `llama-bXXXX-bin-macos-arm64.zip`
   - **Windows**: `llama-bXXXX-bin-win-cuda-cuXX.X-x64.zip` or CPU zip
3. Extract the archive, locate the `llama-server` (or `llama-server.exe`) binary, and move it to a folder in your `PATH` (e.g. `/usr/local/bin` on Linux/macOS):
   ```sh
   sudo mv llama-server /usr/local/bin/
   sudo chmod +x /usr/local/bin/llama-server
   ```

#### Option B: Build From Source with CMake

If you want hardware optimization (e.g., NVIDIA CUDA or Vulkan acceleration):

```sh
# Clone llama.cpp repository
git clone https://github.com/ggerganov/llama.cpp
cd llama.cpp

# For CPU build:
cmake -B build
cmake --build build --config Release -j --target llama-server

# OR for NVIDIA GPU (CUDA) build:
cmake -B build -DGGML_CUDA=ON
cmake --build build --config Release -j --target llama-server

# OR for Vulkan build:
cmake -B build -DGGML_VULKAN=ON
cmake --build build --config Release -j --target llama-server

# Copy llama-server to PATH
sudo cp build/bin/llama-server /usr/local/bin/
```

#### Option C: Homebrew (macOS)

```sh
brew install llama.cpp
```

Verify `llama-server` is in your `PATH`:
```sh
llama-server --help
```

---

### 3. Install `ffmpeg` & Audio Backends

`ffmpeg` is used to capture audio from your microphone.

- **Ubuntu / Debian**:
  ```sh
  sudo apt install -y ffmpeg
  ```
  *Note:* Ensure PulseAudio or PipeWire (`pipewire-pulse`) is running, as `open-stt` uses the default recording device via Pulse.
- **Arch Linux**:
  ```sh
  sudo pacman -S ffmpeg pipewire-pulse
  ```
- **macOS**:
  ```sh
  brew install ffmpeg
  ```
- **Windows**:
  ```powershell
  winget install Gyan.FFmpeg
  ```

---

### 4. Install Typing & CGO Dependencies

- **Linux (X11)**:
  ```sh
  sudo apt install -y xdotool libx11-dev libxtst-dev gcc
  ```
- **Linux (Wayland)**:
  ```sh
  sudo apt install -y wtype
  ```
- **macOS**:
  - Uses built-in AppleScript (`osascript`).
  - Grant your terminal or IDE accessibility access under **System Settings > Privacy & Security > Accessibility** so global hotkeys and typing simulation can operate.
- **Windows**:
  - Uses built-in PowerShell SendKeys.

---

## Build & Setup

1. Clone this repository:
   ```sh
   git clone https://github.com/hunter87/open-stt.git
   cd open-stt
   ```

2. Download Go module dependencies:
   ```sh
   go mod tidy
   ```

3. Build the binary:
   ```sh
   go build -o open-stt .
   ```

---

## Usage

Run the compiled executable:

```sh
./open-stt
```

### Controls

| Action | Key | Description |
|---|---|---|
| **Record** | Hold **F6** | Hold down while speaking into your microphone |
| **Transcribe & Paste** | Release **F6** | Releases recording; audio is transcribed, pasted at the cursor, and copied to clipboard |
| **Quit** | Press **F7** | Gracefully shuts down `open-stt` |

> [!NOTE]
> On the first launch, `llama-server` downloads the model weights (`hunter87/Qwen3-ASR-0.6B-GGUF`) from Hugging Face into `~/.cache/llama.cpp/`. This initial download may take a minute depending on your internet connection.

---

## Configuration

Configurations are located in [`lib/config/config.go`](file:///home/hunter87/Desktop/projects/side-projects/open-stt/lib/config/config.go) and [`lib/config/stt.go`](file:///home/hunter87/Desktop/projects/side-projects/open-stt/lib/config/stt.go):

| Setting | Default Value | Description |
|---|---|---|
| `Hotkeys.Record` | `"f6"` | Global hotkey to hold while speaking |
| `Hotkeys.Quit` | `"f7"` | Global hotkey to exit application |
| `STT.Port` | `"8080"` | Port for the local `llama-server` instance |
| `STT.Model` | `"hunter87/Qwen3-ASR-0.6B-GGUF"` | Hugging Face repository or path for GGUF model |
| `STT.Endpoint` | `"http://localhost:8080/v1/audio/transcriptions"` | OpenAI-compatible transcription endpoint |
| `RecordingSampleRate` | `"16000"` | Audio capture sample rate in Hz |

### Customizing the Server Startup

In [`lib/config/stt.go`](file:///home/hunter87/Desktop/projects/side-projects/open-stt/lib/config/stt.go), `llama-server` is launched with GPU offloading and flash attention enabled:
```go
exec.Command("llama-server", "-hf", STT.Model, "-c", "350", "-np", "1", "--flash-attn", "on", "--n-gpu-layers", "999", "--port", STT.Port)
```
If your machine does not have a dedicated GPU or has limited VRAM, you can adjust `--n-gpu-layers` or lower context `-c`.

---

## Troubleshooting

- **`llama-server not found in PATH`**:
  Make sure `llama-server` is executable and placed in `/usr/local/bin` (or another folder in `$PATH`). Run `which llama-server` to confirm.
- **No audio recorded / Empty transcripts**:
  Verify your microphone is functioning and set as the default input device in PulseAudio or PipeWire. You can test audio recording with:
  ```sh
  ffmpeg -f pulse -i default -t 3 test.wav
  ```
- **Text does not paste automatically**:
  - On Linux X11, ensure `xdotool` is installed.
  - On Linux Wayland, ensure `wtype` is installed.
  - The transcribed text is always saved to your system clipboard regardless of whether auto-pasting succeeds, so you can manually press `Ctrl+V` (or `Cmd+V`).
- **Global hotkeys not registering on macOS**:
  Grant accessibility permissions to your terminal in **System Settings > Privacy & Security > Accessibility**.
