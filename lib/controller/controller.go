/*
This package handles keyboard action and corosponding actions. It is responsible for starting and stopping the recording, as well as sending the recorded audio to the inference server for transcription.
*/

package controller

import (
	"fmt"
	"open-stt/lib/config"
	"open-stt/lib/logging"
	"open-stt/lib/recorder"
	"open-stt/lib/transcriber"
	"open-stt/lib/xpaste"
	"os"
	"sync"

	hook "github.com/robotn/gohook"
	clipboard "golang.design/x/clipboard"
)

var logger = logging.Logger

// Controller manages the recording and transcription workflow
type Controller struct {
	recorder *recorder.Recorder
	fp       string
	mu       sync.Mutex
	text     string
}

func (c *Controller) pasteText() error {
	if c.text == "" {
		logger.Debug("No speech detected (empty transcript).")
		return nil
	}
	logger.Info(fmt.Sprintf("Transcribed Text: %s", c.text))

	// Copy to clipboard
	clipboard.Write(clipboard.FmtText, []byte(c.text))

	// do a paste at cursor
	// Paste/type text at cursor
	if err := xpaste.PasteText(c.text); err != nil {
		logger.Error(fmt.Sprintf("failed to paste text: %s", err))
		return err
	}

	return nil
}

func (c *Controller) handleRecording(audioPath string) {
	if audioPath == "" {
		return
	}

	transcript, err := transcriber.Transcribe(audioPath)
	if err != nil {
		logger.Error(fmt.Sprintf("transcribe error: %v", err))
		_ = os.Remove(audioPath)
		return
	}

	c.mu.Lock()
	c.text = transcript
	c.fp = audioPath
	c.mu.Unlock()

	if err := c.pasteText(); err != nil {
		logger.Error(fmt.Sprintf("paste error: %v", err))
	}

	_ = os.Remove(audioPath)
}

func Listen() {
	var rec recorder.Recorder
	var ctrl Controller

	hook.Register(hook.KeyDown, []string{config.Hotkeys.Record}, func(e hook.Event) {

		if err := rec.Start(); err != nil {
			logger.Error(fmt.Sprintf("record start error: %v", err))
		}
	})

	hook.Register(hook.KeyUp, []string{config.Hotkeys.Record}, func(e hook.Event) {

		audioPath := rec.Stop()
		go ctrl.handleRecording(audioPath)
	})

	hook.Register(hook.KeyDown, []string{config.Hotkeys.Quit}, func(e hook.Event) {
		os.Exit(0)
	})

	s := hook.Start()
	if s == nil {
		logger.Error("failed to start global hook")
		return
	}

	<-hook.Process(s)

}
