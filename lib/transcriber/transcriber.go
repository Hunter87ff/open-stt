package transcriber

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	config "open-sst/lib/config"
	"os"
	"path/filepath"
	"strings"
)

var client = &http.Client{}

type transcriptionResponse struct {
	Text string `json:"text"`
}

func stripASRText(response string) string {
	if strings.Contains(response, "<asr_text>") {
		parts := strings.Split(response, "<asr_text>")
		if len(parts) > 1 {
			return strings.TrimSpace(parts[1])
		}
	}

	return strings.TrimSpace(response)
}

func parseResponse(response string) string {
	var parsed transcriptionResponse
	if err := json.Unmarshal([]byte(response), &parsed); err == nil && parsed.Text != "" {
		return stripASRText(parsed.Text)
	}

	return stripASRText(response)
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

	file := path
	f, err := os.Open(file)
	if err != nil {
		return "", err
	}
	defer f.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filepath.Base(file))
	if err != nil {
		return "", err
	}
	_, err = io.Copy(part, f)
	if err != nil {
		return "", err
	}
	writer.WriteField("model", config.TranscribeModel)
	err = writer.Close()
	if err != nil {
		return "", err
	} 

	req, err := http.NewRequest(http.MethodPost, config.TranscribeEndpoint, body)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	output, err := io.ReadAll(res.Body)

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return "", fmt.Errorf("transcription failed: %s: %s", res.Status, strings.TrimSpace(string(output)))
	}
	
	response := string(output)
	transcribedText := parseResponse(response)
	os.Remove(path) // Clean up the temporary audio file after transcription
	return transcribedText, nil
}
