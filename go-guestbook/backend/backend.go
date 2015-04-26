package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/JamesClonk/vcap"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
)

var (
	redisServiceInstance = "redis-go-guestbook"
	env                  *vcap.VCAP
)

func init() {
	var err error
	env, err = vcap.New()
	if err != nil {
		fmt.Printf("ERROR: %v", err)
	}
}

func main() {
	backendRegistration()

	r := render.New()
	router := mux.NewRouter()

	router.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "Welcome to the Cloud Foundry go-guestbook sample app backend!")
	})

	router.HandleFunc("/entries", func(w http.ResponseWriter, req *http.Request) {
		entries, err := getEntries()
		if err != nil {
			fmt.Fprintf(w, "ERROR getEntries: %v", err)
			return
		}
		r.JSON(w, http.StatusOK, entries)
	}).Methods("GET")

	router.HandleFunc("/entry", func(w http.ResponseWriter, req *http.Request) {
		text := req.URL.Query().Get("text")
		if text != "" {
			insertEntry(text)
		}
	}).Methods("GET")

	router.HandleFunc("/entry", func(w http.ResponseWriter, req *http.Request) {
		err := req.ParseForm()
		if err != nil {
			fmt.Fprintf(w, "ERROR ParseForm: %v", err)
			return
		}

		text := req.FormValue("text")
		if text != "" {
			insertEntry(text)
		}
	}).Methods("POST")

	n := negroni.Classic()
	n.UseHandler(router)
	n.Run(fmt.Sprintf(":%d", env.Port))
}

func backendRegistration() {
	ticker := time.NewTicker(time.Second * 10)
	go func() {
		for range ticker.C {
			registerBackend()
		}
	}()
}
