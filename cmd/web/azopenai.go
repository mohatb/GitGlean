package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type OpenAIRequest struct {
	Prompt           string  `json:"prompt"`
	MaxTokens        int     `json:"max_tokens"`
	Temperature      float64 `json:"temperature"`
	TopP             float64 `json:"top_p"`
	FrequencyPenalty float64 `json:"frequency_penalty"`
	PresencePenalty  float64 `json:"presence_penalty"`
}

type OpenAIResponse struct {
	Choices []struct {
		ContentFilterResults struct {
			Hate struct {
				Filtered bool
				Severity string
			}
			SelfHarm struct {
				Filtered bool
				Severity string
			}
			Sexual struct {
				Filtered bool
				Severity string
			}
		} `json:"content_filter_results"`
		FinishReason string      `json:"finish_reason"`
		Index        int         `json:"index"`
		Logprobs     interface{} `json:"logprobs"` // Using interface{} as it seems to be null or complex
		Message      struct {
			Content string `json:"content"`
			Role    string `json:"role"`
		} `json:"message"`
	} `json:"choices"`
}

func callOpenAISummarize(issue GitHubIssue, comments []GitHubComment) (string, error) {
	azureOpenAIEndpoint := os.Getenv("AZURE_OPENAI_ENDPOINT")
	apiKey := os.Getenv("AZURE_OPENAI_API_KEY")

	if azureOpenAIEndpoint == "" || apiKey == "" {
		log.Println("AZURE_OPENAI_ENDPOINT or AZURE_OPENAI_API_KEY is not set")
		return "", fmt.Errorf("AZURE_OPENAI_ENDPOINT or AZURE_OPENAI_API_KEY is not set")
	}

	// Concatenate all comments into a single string
	var commentsText string
	for _, comment := range comments {
		commentsText += comment.Body + "\n\n"
	}

	// Construct the payload for the OpenAI API call
	prompt := fmt.Sprintf("Title: %s\n\nIssue Body: %s\n\nComments: %s\n\nSummarize the key points of the discussion, including any proposed solutions or workarounds mentioned in the comments. Highlight the main concerns raised and summarize the conversation trajectory, noting any consensus or solutions that users found helpful.", issue.Title, issue.Body, commentsText)
	openAIReq := map[string]interface{}{
		"messages": []map[string]interface{}{
			{
				"role":    "system",
				"content": "You are an assistant tasked with summarizing the essence of GitHub issue discussions. Your goal is to analyze the given issue title, body, and user comments to provide a detailed summary. Focus on identifying main concerns, summarizing discussions around proposed solutions or workarounds, and noting any consensus or effective solutions mentioned by users. Do not create solutions or recommendations.",
			},
			{
				"role":    "user",
				"content": prompt,
			},
		},
	}

	payloadBytes, err := json.Marshal(openAIReq)
	if err != nil {
		log.Println("Error marshaling payload:", err)
		return "", err
	}

	log.Printf("Sending to OpenAI: %s\n", string(payloadBytes))

	// Create and send the request to OpenAI
	azureOpenAIURL := azureOpenAIEndpoint + "/openai/deployments/gpt4/chat/completions?api-version=2024-02-15-preview"
	req, err := http.NewRequest("POST", azureOpenAIURL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		log.Println("Error creating request to OpenAI:", err)
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("api-key", apiKey)

	// Log the request URL and headers
	log.Printf("Request URL: %s", req.URL.String())
	log.Printf("Request headers: %v", req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error sending request to OpenAI:", err)
		return "", err
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading response body from OpenAI:", err)
		return "", err
	}

	log.Printf("Raw response from OpenAI: %s", string(bodyBytes))

	var aiResp OpenAIResponse
	if err := json.Unmarshal(bodyBytes, &aiResp); err != nil {
		log.Println("Error unmarshaling response from OpenAI:", err)
		return "", err
	}

	if len(aiResp.Choices) > 0 {
		return aiResp.Choices[0].Message.Content, nil // Correctly accessing the Text field
	}
	return "", fmt.Errorf("OpenAI response did not contain a summary")
}
