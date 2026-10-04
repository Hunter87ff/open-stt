// this module creates a thread to run the stt model

package config;

import (
	"os"
	"fmt"
	"os/exec"
	"open-stt/lib/logging"
)

var logger = logging.Logger

func checkLlamaServer() error {
	// Check if llama-server is installed
	_, err := exec.LookPath("llama-server")
	if err != nil {
		logger.Error(fmt.Sprintf("llama-server not found in PATH: %v", err))
		return err
	}
	return nil
}

/*
Implementation : Pending...

This function sets up the llama-server. It checks if the server is installed, and if not, it downloads and installs it. This function is called when the user first runs the program to ensure that the necessary dependencies are in place for speech-to-text functionality.
*/
func setupLlamaServer() error {
	// Setup llama-server. this method will be called when the user first runs the program. it will check if llama-server is installed, and if not, it will download and install it.
	return nil
}


/*
wraps llama-server inference in a goroutine to run in the background

-hf : model name
-c : context size
-np : number of threads
--flash-attn : enable flash attention
--n-gpu-layers : number of layers to run on GPU
--port : port to run the server on 
*/
func StartASRThread(){
	if err := checkLlamaServer(); err != nil {
		logger.Error(fmt.Sprintf("Error: %v", err))
		os.Exit(1)
	}

	go func(){
		logger.Info("Starting ASR thread...")
		exec.Command("llama-server", "-hf", STT.Model, "-c", "350", "-np", "1", "--flash-attn", "on", "--n-gpu-layers", "999", "--port", STT.Port).Run()
	}()
	logger.Info("ASR thread started.")

}