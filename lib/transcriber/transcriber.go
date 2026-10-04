package transcriber

import (
	"io"
	"os"
	"fmt"
	"bytes"
	"strings"
	"net/http"
	"path/filepath"
	"encoding/json"
	"mime/multipart"
	"open-stt/lib/config"
)

var client = &http.Client{}

type transcriptionResponse struct {
	Text string `json:"text"`
}

/*
This function checks if the provided path exists or not. It uses os.Stat to get the file information and checks for os.IsNotExist error to determine if the path does not exist. It returns true if the path exists, and false otherwise.
*/
func stripASRText(response string) string {
	if strings.Contains(response, "<asr_text>") {
		parts := strings.Split(response, "<asr_text>")
		if len(parts) > 1 {
			return strings.TrimSpace(parts[1])
		}
	}

	return strings.TrimSpace(response)
}

/*
This function parses the response from the transcription service. It attempts to unmarshal the JSON response into a transcriptionResponse struct. If successful and the Text field is not empty, it returns the stripped text. If unmarshalling fails or the Text field is empty, it returns the stripped original response.
*/
func parseResponse(response string) string {
	var parsed transcriptionResponse
	if err := json.Unmarshal([]byte(response), &parsed); err == nil && parsed.Text != "" {
		return stripASRText(parsed.Text)
	}

	return stripASRText(response)
}


/*
This function checks if the provided path exists or not. It uses os.Stat to get the file information and checks for os.IsNotExist error to determine if the path does not exist. It returns true if the path exists, and false otherwise.
*/
func checkPathExistsOrNot(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

/*
This function transcribes the audio file at the given path using the configured STT model. It first checks if the audio file exists, then creates a multipart/form-data request to send the audio file to the STT endpoint. It reads the response and parses it to extract the transcribed text. If successful, it returns the transcribed text; otherwise, it returns an error.
*/
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
	writer.WriteField("model", config.STT.Model)
	err = writer.Close()
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(http.MethodPost, config.STT.Endpoint, body)
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
