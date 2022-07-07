package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/rpeaugh/bookings/pkg/config"
	"github.com/rpeaugh/bookings/pkg/handlers"
	"github.com/rpeaugh/bookings/pkg/render"
)

const portNumber = ":8080"

// This must be outside the main function so it is available for the middleware.
var app config.AppConfig
var session *scs.SessionManager

// main is the main application function
func main() {
	// Change this to true when in production
	app.InProduction = false

	// Set up a session with a lifetime of 24 hours.  Stored in cookies.
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("Cannot create template cache")
	}

	app.TemplateCache = tc
	app.UseCache = app.InProduction // Should be true for production.

	repo := handlers.NewRepo(&app)

	handlers.NewHandlers(repo)

	render.NewTemplates(&app)

	fmt.Printf("Starting application on port %s \n", portNumber)

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)
}
