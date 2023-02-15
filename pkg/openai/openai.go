package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

var (
	IsAiEnabled bool
	Key         string
	Model       string
	ShellType   string
)

func OpenAI(apiKey string, model string, prompt string, cmdtype string) (map[string][]interface{}, error) {
	// Erstellen Sie eine neue Anfrage an die OpenAI API
	req, err := http.NewRequest("POST", "https://api.openai.com/v1/engines/text-davinci-003/completions", nil)
	if err != nil {
		fmt.Printf("Error creating API request: %s\n", err)
	}

	// Setzen Sie den API-Schlüssel und das Modell
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	// Stellen Sie die Anfrage an das Modell
	reqBody := map[string]interface{}{
		"prompt":            prompt + ". show a " + cmdtype + " example. Please do not write explanations. Please just a suggestion as" + cmdtype + " code.",
		"max_tokens":        100,
		"temperature":       0,
		"top_p":             1.0,
		"frequency_penalty": 0.2,
		"presence_penalty":  0.0,
	}
	jsonReq, err := json.Marshal(reqBody)
	if err != nil {
		fmt.Printf("Error encoding request body: %s\n", err)
	}
	req.Body = ioutil.NopCloser(bytes.NewBuffer(jsonReq))

	// Senden Sie die Anfrage an die API
	httpClient := &http.Client{}
	res, err := httpClient.Do(req.WithContext(context.Background()))
	if err != nil {
		fmt.Printf("Error sending API request: %s\n", err)
	}
	defer res.Body.Close()

	// Verarbeiten Sie die Antwort des Modells
	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("Error reading API response: %s\n", err)
	}
	var response map[string][]interface{}
	json.Unmarshal(resBody, &response)
	// if err := json.Unmarshal(resBody, &response); err != nil {
	// 	fmt.Printf("Error decoding API response: %s\n", err)
	// }
	return response, err
}

func PrintAiResponse(response map[string][]interface{}) string {
	out := response["choices"][0].(map[string]interface{})
	return string(strings.ReplaceAll(strings.TrimSpace(out["text"].(string)), "`", ""))
}
