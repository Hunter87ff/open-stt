/*
This package handles keyboard action and corosponding actions. It is responsible for starting and stopping the recording, as well as sending the recorded audio to the inference server for transcription.
*/

package controller

import (
	"os"
	"fmt"
	"sync"
	"open-stt/lib/config"
	"open-stt/lib/recorder"
	"github.com/robotn/gohook"
	"open-stt/lib/transcriber"
	"golang.design/x/clipboard"
)

// Controller manages the recording and transcription workflow
type Controller struct {
	recorder *recorder.Recorder
	fp       string
	mu       sync.Mutex
	text     string
}

func (c *Controller) pasteText() error {
	if c.text == "" {
		return nil
	}

	err := clipboard.Write(clipboard.FmtText, []byte(c.text))
	if err != nil {
		return fmt.Errorf("failed to write to clipboard: %w", err)
	}
	return nil
}

func (c *Controller) handleRecording(audioPath string) {
	if audioPath == "" {
		return
	}

	transcript, err := transcriber.Transcribe(audioPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "transcribe error: %v\n", err)
		_ = os.Remove(audioPath)
		return
	}

	c.mu.Lock()
	c.text = transcript
	c.fp = audioPath
	c.mu.Unlock()

	if err := c.pasteText(); err != nil {
		fmt.Fprintf(os.Stderr, "paste error: %v\n", err)
	}

	_ = os.Remove(audioPath)
}

func Listen() {
	var rec recorder.Recorder
	var ctrl Controller

	hook.Register(hook.KeyDown, []string{config.RecordHotkey}, func(e hook.Event) {

		if err := rec.Start(); err != nil {
			fmt.Fprintf(os.Stderr, "record start error: %v\n", err)
		}
	})

	hook.Register(hook.KeyUp, []string{config.RecordHotkey}, func(e hook.Event) {

		audioPath := rec.Stop()
		go ctrl.handleRecording(audioPath)
	})

	hook.Register(hook.KeyDown, []string{config.QuitHotkey}, func(e hook.Event) {
		os.Exit(0)
	})

	s := hook.Start()
	if s == nil {
		fmt.Fprintln(os.Stderr, "failed to start global hook")
		return
	}

	<-hook.Process(s)

}
