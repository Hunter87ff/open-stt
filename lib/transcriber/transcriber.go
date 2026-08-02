package transcriber

import (
	"os/exec"
	"strings"
	"os"
	config "open-sst/lib/config"
)


func checkHost() bool {
	cmd := exec.Command("ping", "-c", "1", config.TranscribeEndpoint)
	err := cmd.Run()
	return err == nil
}


func parseResponse(response string) string {
	// Check if the response contains the <asr_text> tag
	if strings.Contains(response, "<asr_text>") {
		// Split the response by the <asr_text> tag and return the text after it
		parts := strings.Split(response, "<asr_text>")
		if len(parts) > 1 {
			return strings.TrimSpace(parts[1])
		}
	}

	// If the tag is not found, return the original response
	return strings.TrimSpace(response)
}


func checkPathExistsOrNot(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return true
}


func Transcribe(path string) (string, error) {
	audioExists := checkPathExistsOrNot(path)
	if !audioExists {
		return "", nil
	}

	hostStatus := checkHost()
	if !hostStatus {
		return "", nil
	}

	cmd := exec.Command("curl", "-X", "POST", config.TranscribeEndpoint, "-F", "file=@"+path, "-F", "model="+config.TranscribeModel)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	response := string(output)
	transcribedText := parseResponse(response)
	os.Remove(path) // Clean up the temporary audio file after transcription
	return transcribedText, nil
}




