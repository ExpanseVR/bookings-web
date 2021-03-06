package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ExpanseVR/bookings/internal/config"
	"github.com/ExpanseVR/bookings/internal/handlers"
	"github.com/ExpanseVR/bookings/internal/models"
	"github.com/ExpanseVR/bookings/internal/render"
	"github.com/alexedwards/scs/v2"
)

var app config.AppConfig

const portNumber = ":8080"

var session *scs.SessionManager

func main() {
	gob.Register(models.Reservation{})
	//change this to true when in production
	app.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
	}

	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app)
	handlers.NewHandler(repo)
	render.NewTemplates(&app)

	fmt.Println(fmt.Sprintf("Listening on port %s", portNumber))

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)
}
