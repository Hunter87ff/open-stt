package recorder

import (
	"os"
	"fmt"
	"sync"
	"time"
	"os/exec"
	config "open-sst/lib/config"
)


type Recorder struct {
	mu        sync.Mutex
	cmd       *exec.Cmd
	audioPath string
}


func checkDependencies() error {
	// Check if ffmpeg is installed
	_, err := exec.LookPath("ffmpeg")
	if err != nil {
		return fmt.Errorf("ffmpeg is not installed or not found in PATH. Please install ffmpeg to use this application")
	}

	return nil
}


func createTempFile() (string, error) {
	tmpFile, err := os.CreateTemp("", "open-sst-*.wav")
	if err != nil {
		return "", err
	}
	if err := tmpFile.Close(); err != nil {
		_ = os.Remove(tmpFile.Name())
		return "", err
	}
	return tmpFile.Name(), nil
}


func startAudioCapture(outputPath string) (*exec.Cmd, error) {

	err := checkDependencies()
	if err != nil {
		return nil, err
	}

    // Try ffmpeg first (cross-platform)
    cmd := exec.Command(
        "ffmpeg",
        "-hide_banner",
        "-loglevel", "error",
        "-f", "auto",  // Try auto-detect input (works on Windows/macOS/Linux)
        "-i", "default",
        "-ac", "1",
        "-ar", config.RecordingSampleRate,
        "-y",
        outputPath,
    )
    if err := cmd.Start(); err != nil {
        return nil, fmt.Errorf("ffmpeg failed: %v. Ensure ffmpeg is installed.", err)
    }
    return cmd, nil
}


func (r *Recorder) Start() error {
	r.mu.Lock()
	filePath, err := createTempFile()
	if err != nil {
		return err
	}

	cmd, err := startAudioCapture(filePath)
	if err != nil {
		_ = os.Remove(filePath)
		return err
	}

	r.cmd = cmd
	r.audioPath = filePath
	r.mu.Unlock()

	return nil
}


func (r *Recorder) Stop() string {
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