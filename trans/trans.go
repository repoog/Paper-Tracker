package trans

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const URL = "http://localhost:11434"

type Request struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type Response struct {
	Model              string `json:"model"`
	CreatedAt          string `json:"created_at"`
	Response           string `json:"response"`
	Done               bool   `json:"done"`
	Context            []int  `json:"context"`
	TotalDuration      int64  `json:"total_duration"`
	LoadDuration       int64  `json:"load_duration"`
	PromptEvalCount    int    `json:"prompt_eval_count"`
	PromptEvalDuration int64  `json:"prompt_eval_duration"`
	EvalCount          int    `json:"eval_count"`
	EvalDuration       int64  `json:"eval_duration"`
}

func Trans(originalText string) (string, error) {
	llamaUrl := URL + "/api/generate"
	data := Request{
		Model:  "llama3.1",
		Prompt: "作为网络安全专家，请将以下文字翻译为中文，并仅输出翻译结果：" + originalText,
		Stream: false,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("[!] Error building JSON request")
	}

	resp, err := http.Post(llamaUrl, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("[!] Error sending request to llama: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("[!] Error reading response of llama: %v", err)
	}

	var response Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", fmt.Errorf("[!] Error unmarshalling response: %v", err)
	}

	return response.Response, nil
}
