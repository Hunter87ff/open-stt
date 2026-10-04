// this module creates a thread to run the stt model

package config;

import (
	"os"
	"fmt"
	"os/exec"
)


func checkLlamaServer() error {
	// Check if llama-server is installed
	_, err := exec.LookPath("llama-server")
	if err != nil {
		return fmt.Errorf("llama-server not found in PATH: %v", err)
	}
	return nil
}



// wraps llama-server -hf hunter87/Qwen3-ASR-0.6B-GGUF -c 350 -np 1 --flash-attn on --n-gpu-layers 999 
func StartASRThread(){
	if err := checkLlamaServer(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Starting ASR thread...")
	go func(){
		exec.Command("llama-server", "-hf", "hunter87/Qwen3-ASR-0.6B-GGUF", "-c", "350", "-np", "1", "--flash-attn", "on", "--n-gpu-layers", "999").Run()
	}()
	fmt.Println("ASR thread started.")

}