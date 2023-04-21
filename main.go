package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mychatGPT/proxy"
	"net/http"
)

type ChatGPTRequest struct {
	Messages []Message `json:"messages"`
	Model    string    `json:"model"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatGPTResponse struct {
	Id      string `json:"id"`
	Object  string `json:"object"`
	Created int    `json:"created"`
	Model   string `json:"model"`
	Usage   struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
	Choices []struct {
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
		Index        int    `json:"index"`
	} `json:"choices"`
}

func main() {
	apiKey := "xxx"
	userMessage := "What are the benefits of using Go language?"

	messages := []Message{
		{
			Role:    "system",
			Content: "You are a helpful assistant that provides information about the Go programming language.",
		},
		{
			Role:    "user",
			Content: userMessage,
		},
	}

	chatGPTRequest := ChatGPTRequest{
		Messages: messages,
		Model:    "gpt-3.5-turbo",
	}

	queryChatGPT(apiKey, chatGPTRequest)
}

func queryChatGPT(apiKey string, chatGPTRequest ChatGPTRequest) {
	url := "https://api.openai.com/v1/chat/completions"

	client := &http.Client{}

	p, err := proxy.RoundRobinProxySwitcher("http://127.0.0.1:8888")
	if err != nil {
		fmt.Println("Error creating proxy:", err)
		return
	}
	transport := http.DefaultTransport.(*http.Transport)
	transport.Proxy = p
	client.Transport = transport

	requestBody, err := json.Marshal(chatGPTRequest)
	if err != nil {
		fmt.Println("Error marshalling request:", err)
		return
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	fmt.Println(string(body))

	var chatGPTResponse ChatGPTResponse
	err = json.Unmarshal(body, &chatGPTResponse)
	if err != nil {
		fmt.Println("Error unmarshalling response:", err)
		return
	}

	fmt.Println("Generated response:", chatGPTResponse.Choices[0].Message.Content)
}
