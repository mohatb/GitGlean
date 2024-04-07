package main

import (
	"context"
	"encoding/base64"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"

	oidc "github.com/coreos/go-oidc"
	"golang.org/x/oauth2"
)

type Authenticator struct {
	Provider *oidc.Provider
	Config   oauth2.Config
	Ctx      context.Context
}

func (app *application) NewAuthenticator() (*Authenticator, error) {
	ctx := context.Background()
	provider, err := oidc.NewProvider(ctx, "https://login.microsoftonline.com/"+os.Getenv("tenantid")+"/v2.0")
	if err != nil {
		app.infoLog.Printf("failed to get provider: %v", err)
		return nil, err
	}

	conf := oauth2.Config{
		ClientID:     os.Getenv("clientid"),
		ClientSecret: os.Getenv("clientsecret"),
		RedirectURL:  os.Getenv("callbackurl"),
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "profile"},
	}

	app.infoLog.Print(conf)

	return &Authenticator{
		Provider: provider,
		Config:   conf,
		Ctx:      ctx,
	}, nil
}

func (app *application) CallbackHandler(w http.ResponseWriter, r *http.Request) {
	receivedState := r.URL.Query().Get("state")
	expectedState, ok := app.session.Get(r.Context(), "oauthState").(string)
	if !ok || receivedState != expectedState {
		app.infoLog.Println("Invalid state parameter")
		http.Error(w, "Invalid state parameter", http.StatusBadRequest)
		return
	}

	authenticator, err := app.NewAuthenticator()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	token, err := authenticator.Config.Exchange(context.TODO(), r.URL.Query().Get("code"))
	if err != nil {
		app.errorLog.Printf("no token found: %v", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		http.Error(w, "No id_token field in oauth2 token.", http.StatusInternalServerError)
		return
	}

	oidcConfig := &oidc.Config{
		ClientID: os.Getenv("clientid"),
	}

	idToken, err := authenticator.Provider.Verifier(oidcConfig).Verify(context.TODO(), rawIDToken)
	if err != nil {
		http.Error(w, "Failed to verify ID Token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var profile map[string]interface{}
	if err := idToken.Claims(&profile); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	app.session.Put(r.Context(), "id_token", rawIDToken)
	app.session.Put(r.Context(), "access_token", token.AccessToken)
	app.session.Put(r.Context(), "profile", profile)

	log.Print(profile)

	// Redirect to logged-in page
	http.Redirect(w, r, "/user", http.StatusSeeOther)
}

func (app *application) LoginHandler(w http.ResponseWriter, r *http.Request) {
	// Generate random state
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	state := base64.StdEncoding.EncodeToString(b)

	app.infoLog.Printf("Generated state: %s", state)

	app.session.Put(r.Context(), "oauthState", state)

	authenticator, err := app.NewAuthenticator()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, authenticator.Config.AuthCodeURL(state), http.StatusFound)
}

func (app *application) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// Clear all values stored in the session
	app.session.Remove(r.Context(), "id_token")
	app.session.Remove(r.Context(), "access_token")
	app.session.Remove(r.Context(), "profile")
	app.session.Remove(r.Context(), "oauthState")

	// Create the logout URL with the post_logout_redirect_uri parameter
	logoutUrl, err := url.Parse("https://login.microsoftonline.com/" + os.Getenv("tenantid") + "/oauth2/logout?client_id=" + os.Getenv("clientid") + "&post_logout_redirect_uri=" + url.QueryEscape("http://localhost:3000"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Redirect the user to the logout URL
	http.Redirect(w, r, logoutUrl.String(), http.StatusTemporaryRedirect)
}

func (app *application) IsAuthenticated(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	profile := app.session.Get(r.Context(), "profile")
	if profile == nil {
		app.errorLog.Println("User is not authenticated")
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	app.infoLog.Println("User is authenticated:", profile)
	next(w, r)
}
