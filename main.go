package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	hook "github.com/robotn/gohook"
)

const (
	recordHotkey        = "f6"
	transcribeEndpoint  = "http://127.0.0.1:8080/v1/audio/transcriptions"
	transcribeModel     = "qwen3-asr-0.6b"
	recordingSampleRate = "16000"
)

type recorder struct {
	mu        sync.Mutex
	cmd       *exec.Cmd
	audioPath string
}


func (r *recorder) start() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.cmd != nil {
		return nil
	}

	tmpFile, err := os.CreateTemp("", "open-sst-*.wav")
	if err != nil {
		return err
	}
	if err := tmpFile.Close(); err != nil {
		_ = os.Remove(tmpFile.Name())
		return err
	}

	cmd, err := startAudioCapture(tmpFile.Name())
	if err != nil {
		_ = os.Remove(tmpFile.Name())
		return err
	}

	r.cmd = cmd
	r.audioPath = tmpFile.Name()
	return nil
}

func (r *recorder) stop() string {
	r.mu.Lock()
	cmd := r.cmd
	audioPath := r.audioPath
	r.cmd = nil
	r.audioPath = ""
	r.mu.Unlock()

	if cmd == nil {
		return ""
	}

	_ = cmd.Process.Signal(os.Interrupt)
	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()

	select {
	case <-time.After(3 * time.Second):
		_ = cmd.Process.Kill()
		<-done
	case <-done:
	}

	return audioPath
}

func startAudioCapture(outputPath string) (*exec.Cmd, error) {
	if _, err := exec.LookPath("ffmpeg"); err == nil {
		cmd := exec.Command(
			"ffmpeg",
			"-hide_banner",
			"-loglevel", "error",
			"-f", "pulse",
			"-i", "default",
			"-ac", "1",
			"-ar", recordingSampleRate,
			"-y",
			outputPath,
		)
		if err := cmd.Start(); err != nil {
			return nil, err
		}
		return cmd, nil
	}

	if _, err := exec.LookPath("arecord"); err == nil {
		cmd := exec.Command(
			"arecord",
			"-f", "S16_LE",
			"-c", "1",
			"-r", recordingSampleRate,
			"-t", "wav",
			outputPath,
		)
		if err := cmd.Start(); err != nil {
			return nil, err
		}
		return cmd, nil
	}

	return nil, fmt.Errorf("need ffmpeg or arecord for microphone capture")
}

func transcribe(audioPath string) (string, error) {
	file, err := os.Open(audioPath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	part, err := writer.CreateFormFile("file", filepath.Base(audioPath))
	if err != nil {
		return "", err
	}
	if _, err := io.Copy(part, file); err != nil {
		return "", err
	}
	if err := writer.WriteField("model", transcribeModel); err != nil {
		return "", err
	}
	if err := writer.Close(); err != nil {
		return "", err
	}

	req, err := http.NewRequest(http.MethodPost, transcribeEndpoint, &body)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{Timeout: 5 * time.Minute}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("transcription failed: %s: %s", resp.Status, strings.TrimSpace(string(responseBody)))
	}

	var payload struct {
		Text string `json:"text"`
	}
	if err := json.Unmarshal(responseBody, &payload); err == nil && payload.Text != "" {
		return strings.TrimSpace(payload.Text), nil
	}

	return strings.TrimSpace(string(responseBody)), nil
}

func pasteText(text string) error {
	if text == "" {
		return nil
	}

	if _, err := exec.LookPath("xdotool"); err == nil {
		return exec.Command("xdotool", "type", "--clearmodifiers", "--delay", "0", text).Run()
	}

	if _, err := exec.LookPath("wtype"); err == nil {
		return exec.Command("wtype", text).Run()
	}

	return fmt.Errorf("need xdotool or wtype to paste text")
}

func handleRecording(audioPath string) {
	if audioPath == "" {
		return
	}

	transcript, err := transcribe(audioPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "transcribe error: %v\n", err)
		_ = os.Remove(audioPath)
		return
	}

	if err := pasteText(transcript); err != nil {
		fmt.Fprintf(os.Stderr, "paste error: %v\n", err)
	}

	_ = os.Remove(audioPath)
}

func main() {
	var rec recorder

	hook.Register(hook.KeyDown, []string{recordHotkey}, func(e hook.Event) {
		if err := rec.start(); err != nil {
			fmt.Fprintf(os.Stderr, "record start error: %v\n", err)
		}
	})

	hook.Register(hook.KeyUp, []string{recordHotkey}, func(e hook.Event) {
		audioPath := rec.stop()
		go handleRecording(audioPath)
	})

	hook.Register(hook.KeyDown, []string{"ctrl", "shift", "q"}, func(e hook.Event) {
		os.Exit(0)
	})

	s := hook.Start()
	if s == nil {
		fmt.Fprintln(os.Stderr, "failed to start global hook")
		return
	}

	<-hook.Process(s)
}
