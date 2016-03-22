package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/JamesClonk/vcap"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
)

const (
	mongoDbServiceInstance = "mongodb-docker-guestbook"
)

var env *vcap.VCAP

func init() {
	// parse cloudfoundry VCAP_* env data
	var err error
	env, err = vcap.New()
	if err != nil {
		log.Fatalf("ERROR: %v\n", err)
	}
}

func main() {
	// create render instance
	r := render.New(render.Options{
		IndentJSON: true,
	})

	// setup routes
	router := mux.NewRouter()
	router.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		r.JSON(w, http.StatusOK, "Welcome to the Cloud Foundry docker-guestbook sample app backend!")
	})
	router.HandleFunc("/entries", getEntriesHandler(r)).Methods("GET")
	router.HandleFunc("/entry", postEntryHandler(r)).Methods("POST")

	// setup negroni
	n := negroni.Classic()
	n.UseHandler(router)
	n.Run(fmt.Sprintf(":%d", env.Port))
}

func getEntriesHandler(r *render.Render) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		entries, err := getEntries()
		if err != nil {
			r.JSON(w, http.StatusInternalServerError, err)
			return
		}
		r.JSON(w, http.StatusOK, entries)
	}
}

func postEntryHandler(r *render.Render) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		err := req.ParseForm()
		if err != nil {
			r.JSON(w, http.StatusInternalServerError, err)
			return
		}

		text := req.FormValue("text")
		if text == "" {
			r.JSON(w, http.StatusExpectationFailed, "No text provided")
		}

		if err := insertEntry(text); err != nil {
			r.JSON(w, http.StatusInternalServerError, err)
			return
		}
		r.JSON(w, http.StatusCreated, nil)
	}
}
