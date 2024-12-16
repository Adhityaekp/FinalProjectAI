package service

import (
	"a21hc3NpZ25tZW50/model"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type AIService struct {
	Client HTTPClient
}

func (s *AIService) AnalyzeData(table map[string][]string, query, token string) (string, error) {
	if len(table) == 0 {
		return "", errors.New("table cannot be empty")
	}

	url := "https://api-inference.huggingface.co/models/google/tapas-large-finetuned-wtq"
	reqBody, _ := json.Marshal(map[string]interface{}{
		"inputs": map[string]interface{}{
			"query": query,
			"table": table,
		},
	})
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.Client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var response map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", err
	}

	if cells, ok := response["cells"].([]interface{}); ok && len(cells) > 0 {
		if result, ok := cells[0].(string); ok {
			return result, nil
		}
	}

	return "", errors.New("invalid response from AI model")
}

func (s *AIService) ChatWithAI(context, query, token string) (model.ChatResponse, error) {
	url := "https://api-inference.huggingface.co/models/microsoft/Phi-3.5-mini-instruct/v1/chat/completions"

	reqBody, _ := json.Marshal(map[string]interface{}{
		"model": "microsoft/Phi-3.5-mini-instruct",
		"messages": []map[string]string{
			{
				"role":    "user",
				"content": query,
			},
		},
		"max_tokens": 500,
		"stream":     false,
	})

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.Client.Do(req)
	if err != nil {
		return model.ChatResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return model.ChatResponse{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return model.ChatResponse{}, err
	}

	fmt.Println(string(body))

	var response struct {
		Choices []struct {
			Message struct {
				Role    string `json:"role"`
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return model.ChatResponse{}, err
	}

	if len(response.Choices) > 0 {
		return model.ChatResponse{
			GeneratedText: response.Choices[0].Message.Content,
		}, nil
	}

	return model.ChatResponse{}, fmt.Errorf("no response choices found")
}
