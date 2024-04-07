package main

import (
	"net/http"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	// Retrieve the profile from the session using SCS
	profile, ok := app.session.Get(r.Context(), "profile").(map[string]interface{})
	if !ok {
		// Handle the case where the profile is not set
		profile = map[string]interface{}{}
	}

	// Add any additional data that you want to pass to the template.
	data := map[string]interface{}{
		"Name":       profile["name"],
		"Email":      profile["email"],
		"ShowFooter": true,
		// Add any other data here.
	}
	// Pass the data to the template.
	RenderTemplate(w, "home", &data)
}

func (app *application) UserHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve the profile from the session using SCS
	profile, ok := app.session.Get(r.Context(), "profile").(map[string]interface{})
	if !ok {
		app.clientError(w, http.StatusUnauthorized)
		return
	}

	// Add any additional data that you want to pass to the template.
	data := map[string]interface{}{
		"Name":       profile["name"],
		"Email":      profile["email"],
		"ShowFooter": true,
		// Add any other data here.
	}
	// Pass the data to the template.
	RenderTemplate(w, "user", &data)
}

// handlers.go

func (app *application) githubIssuesHandler(w http.ResponseWriter, r *http.Request) {
	// Example data passed to the template, with ShowFooter set to false to hide the footer
	data := map[string]interface{}{
		"ShowFooter": false,
	}
	RenderTemplate(w, "githubIssues", data)
}
