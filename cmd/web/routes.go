package main

import (
	"net/http"

	"github.com/codegangsta/negroni"
)

func (app *application) routes() *negroni.Negroni {
	mux := http.NewServeMux()

	// Define your routes
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/login", app.LoginHandler)
	mux.HandleFunc("/logout", app.LogoutHandler)
	mux.HandleFunc("/callback", app.CallbackHandler)
	mux.HandleFunc("/user", app.UserHandler)
	mux.HandleFunc("/github-issues", app.githubIssuesHandler)
	mux.HandleFunc("/api/list-issues", app.listGitHubIssuesHandler)
	mux.HandleFunc("/api/analyze-issue", app.analyzeGitHubIssueHandler)

	// Create a Negroni instance which will hold all middleware and the mux.
	n := negroni.New()

	// First middleware: SCS LoadAndSave. It ensures session data is loaded/saved automatically.
	n.Use(negroni.HandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		ctx, err := app.session.Load(r.Context(), r.Header.Get("X-Session-Token"))
		if err != nil {
			app.serverError(w, err)
			return
		}

		next(w, r.WithContext(ctx))
	}))

	// Add your Negroni middleware for authentication before specific handlers if needed.
	// For example, you could add a middleware that checks if a user is authenticated
	// and only then calls the next handler in the chain.
	// This is where you would use app.IsAuthenticated, but note that integrating it directly
	// as middleware requires it to fit the negroni.HandlerFunc signature.
	// n.Use(negroni.HandlerFunc(app.IsAuthenticated))

	// The main router (mux) is added as the final handler.
	n.UseHandler(mux)

	return n
}
