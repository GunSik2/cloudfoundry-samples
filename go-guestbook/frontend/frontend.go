package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"time"

	"github.com/JamesClonk/vcap"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
)

var (
	redisServiceInstance = "redis-go-guestbook"
	env                  *vcap.VCAP
)

type Entry struct {
	Timestamp time.Time
	Text      string
}

func init() {
	var err error
	env, err = vcap.New()
	if err != nil {
		fmt.Printf("ERROR: %v", err)
	}
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "Welcome to the Cloud Foundry go-guestbook sample app frontend!\n")

		counter, err := getHitCounter()
		if err != nil {
			fmt.Fprintf(w, "ERROR getHitCounter: %v", err)
			return
		}
		fmt.Fprintf(w, "This is page hit #%v", counter)
	})

	router.HandleFunc("/backends", func(w http.ResponseWriter, req *http.Request) {
		backends, err := discoverBackends()
		if err != nil {
			fmt.Fprintf(w, "ERROR discoverBackends: %v", err)
			return
		}
		for idx, backend := range backends {
			fmt.Fprintf(w, "go-guestbook-backend #%v: %v\n", idx, backend)
		}
	})

	router.HandleFunc("/entries", func(w http.ResponseWriter, req *http.Request) {
		backend, err := getBackendURL()
		if err != nil {
			fmt.Fprintf(w, "ERROR getBackendURL: %v", err)
			return
		}

		resp, err := http.Get(backend + "/entries")
		if err != nil {
			fmt.Fprintf(w, "ERROR GET entries: %v", err)
			return
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Fprintf(w, "ERROR GET entries: %v", err)
			return
		}
		fmt.Fprintf(w, "%s", string(body))
	})

	router.HandleFunc("/entry", func(w http.ResponseWriter, req *http.Request) {
		text := req.URL.Query().Get("text")
		if text != "" {
			backend, err := getBackendURL()
			if err != nil {
				fmt.Fprintf(w, "ERROR getBackendURL: %v", err)
				return
			}

			_, err = http.PostForm(backend+"/entry", url.Values{"text": {text}})
			if err != nil {
				fmt.Fprintf(w, "ERROR POST entry: %v", err)
				return
			}
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
			backend, err := getBackendURL()
			if err != nil {
				fmt.Fprintf(w, "ERROR getBackendURL: %v", err)
				return
			}

			_, err = http.PostForm(backend+"/entry", url.Values{"text": {text}})
			if err != nil {
				fmt.Fprintf(w, "ERROR POST entry: %v", err)
				return
			}
		}
	}).Methods("POST")

	n := negroni.Classic()
	n.Use(&HitCounter{})
	n.UseHandler(router)
	n.Run(fmt.Sprintf(":%d", env.Port))
}

func getBackendURL() (string, error) {
	// discovering and choosing backends should probably be cached a bit,
	// instead of actually looking it up on redis every time
	rand.Seed(time.Now().UTC().UnixNano())

	backends, err := discoverBackends()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("http://%s", backends[rand.Intn(len(backends))]), nil
}
