package main

import (
	"encoding/gob"
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/joho/godotenv"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	session  *scs.SessionManager // Add this line

}

func init() {
	// Register the map[string]interface{} type with gob
	// to ensure it can be stored in the session.
	gob.Register(map[string]interface{}{})
}

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Print("Error loading .env file")
	}

	addr := flag.String("addr", ":3000", "HTTP network address")

	sessionManager := scs.New()
	sessionManager.Lifetime = 24 * time.Hour
	sessionManager.Cookie.Persist = true
	sessionManager.Cookie.SameSite = http.SameSiteLaxMode
	sessionManager.Cookie.Secure = false // Set to true in production

	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
		session:  sessionManager,
	}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		// Call the new app.routes() method to get the servemux containing our routes.
		Handler: app.routes(),
	}

	ip := app.getIP()
	if ip == "" {
		app.errorLog.Fatal("Unable to retrieve IP address")
	}

	infoLog.Printf("Starting server on http://%s%s", ip, *addr)
	err = srv.ListenAndServe()
	app.errorLog.Fatal(err)
}
