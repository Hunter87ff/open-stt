package recorder

import (
	"os"
	"fmt"
	"sync"
	"time"
	"os/exec"
	"open-stt/lib/config"
	"open-stt/lib/logging"
)


var logger = logging.Logger

type Recorder struct {
	mu        sync.Mutex
	cond      *sync.Cond
	starting  bool
	recording bool
	cmd       *exec.Cmd
	audioPath string
}

func (r *Recorder) ensureCond() {
	if r.cond == nil {
		r.cond = sync.NewCond(&r.mu)
	}
}

func checkDependencies() error {
	// Check if ffmpeg is installed
	_, err := exec.LookPath("ffmpeg")
	if err != nil {
		logger.Error("ffmpeg is not installed or not found in PATH. Please install ffmpeg to use this application")
		return err
	}

	return nil
}

func createTempFile() (string, error) {
	tmpFile, err := os.CreateTemp("", "open-stt-*.wav")
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
		"-f", "pulse",
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
	r.ensureCond()
	if r.starting || r.recording {
		r.mu.Unlock()
		return nil
	}
	r.starting = true
	filePath, err := createTempFile()
	if err != nil {
		r.starting = false
		r.cond.Broadcast()
		r.mu.Unlock()
		return err
	}

	logger.Info("Recording started...")

	cmd, err := startAudioCapture(filePath)
	if err != nil {
		_ = os.Remove(filePath)
		r.starting = false
		r.cond.Broadcast()
		r.mu.Unlock()
		return err
	}

	r.cmd = cmd
	r.audioPath = filePath
	r.starting = false
	r.recording = true
	r.cond.Broadcast()
	r.mu.Unlock()

	return nil
}

func (r *Recorder) Stop() string {
	r.mu.Lock()
	r.ensureCond()
	for r.starting {
		r.cond.Wait()
	}
	if !r.recording {
		r.mu.Unlock()
		return ""
	}
	r.recording = false
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
	logger.Info("Recording stopped.")

	return audioPath
}
