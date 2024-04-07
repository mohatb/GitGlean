package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type GitHubIssue struct {
	URL         string `json:"url"`
	CommentsURL string `json:"comments_url"`
	Title       string `json:"title"`
	Body        string `json:"body"`
	State       string `json:"state"`
	HTMLURL     string `json:"html_url"`
	Comments    int    `json:"comments"` // Add this line
	Reactions   struct {
		TotalCount int `json:"total_count"`
	} `json:"reactions"` // Add this line
}

type GitHubComment struct {
	Body string `json:"body"`
}

// Fetch GitHub issues with authentication
func FetchGitHubIssues(url string) ([]GitHubIssue, error) {
	accessToken := os.Getenv("GITHUB_ACCESS_TOKEN") // Retrieve the access token from environment variable
	if accessToken == "" {
		return nil, fmt.Errorf("GitHub access token not set")
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Set authorization header with the personal access token
	req.Header.Set("Authorization", "token "+accessToken)
	//req.Header.Add("Accept", "application/vnd.github.v3+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var issues []GitHubIssue
	if err := json.NewDecoder(resp.Body).Decode(&issues); err != nil {
		return nil, err
	}

	return issues, nil
}

func FetchIssueComments(commentsURL string) ([]GitHubComment, error) {
	accessToken := os.Getenv("GITHUB_ACCESS_TOKEN") // Retrieve the access token from environment variable
	if accessToken == "" {
		return nil, fmt.Errorf("GitHub access token not set")
	}

	req, err := http.NewRequest("GET", commentsURL, nil)
	if err != nil {
		return nil, err
	}

	// Set authorization header with the personal access token
	req.Header.Set("Authorization", "token "+accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var comments []GitHubComment
	if err := json.NewDecoder(resp.Body).Decode(&comments); err != nil {
		return nil, err
	}

	return comments, nil
}

// FetchGitHubIssue fetches details of a single GitHub issue.
func FetchGitHubIssue(issueURL string) (*GitHubIssue, error) {
	// Assuming you have the access token stored in an environment variable.
	accessToken := os.Getenv("GITHUB_ACCESS_TOKEN")
	if accessToken == "" {
		return nil, fmt.Errorf("GitHub access token not set")
	}

	// Create the HTTP request for the GitHub API.
	req, err := http.NewRequest("GET", issueURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "token "+accessToken)

	// Send the request.
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Decode the JSON response into the GitHubIssue struct.
	var issue GitHubIssue
	if err := json.NewDecoder(resp.Body).Decode(&issue); err != nil {
		return nil, err
	}

	return &issue, nil
}
