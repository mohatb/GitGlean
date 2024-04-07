// cmd/web/apiHandlers.go

package main

import (
	"encoding/json"
	"net/http"
)

// listGitHubIssuesHandler handles requests to list GitHub issues for a specific repository.
func (app *application) listGitHubIssuesHandler(w http.ResponseWriter, r *http.Request) {
	// Here, we're hardcoding the URL for simplicity, but you could modify this to accept dynamic input
	issues, err := FetchGitHubIssues("https://api.github.com/repos/Azure/AKS/issues?per_page=100")
	if err != nil {
		app.errorLog.Printf("Error fetching GitHub issues: %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Setting the header to application/json for the response
	w.Header().Set("Content-Type", "application/json")
	// Encoding the issues into JSON and sending the response
	err = json.NewEncoder(w).Encode(issues)
	if err != nil {
		app.errorLog.Printf("Error encoding JSON response: %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (app *application) analyzeGitHubIssueHandler(w http.ResponseWriter, r *http.Request) {
	issueNumber := r.URL.Query().Get("issue_number")
	if issueNumber == "" {
		http.Error(w, "Issue number not specified", http.StatusBadRequest)
		return
	}

	issue, err := FetchGitHubIssue("https://api.github.com/repos/Azure/AKS/issues/" + issueNumber)
	if err != nil {
		app.errorLog.Printf("Error fetching GitHub issue: %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	comments, err := FetchIssueComments(issue.CommentsURL)
	if err != nil {
		app.errorLog.Printf("Error fetching GitHub comments: %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	summary, err := callOpenAISummarize(*issue, comments) // Make sure this function is correctly implemented to handle a single issue and its comments.
	if err != nil {
		app.errorLog.Printf("Error summarizing with OpenAI: %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"summary": summary})
}
