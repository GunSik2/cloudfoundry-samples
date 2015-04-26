package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/JamesClonk/vcap"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
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

	router := mux.NewRouter()
	router.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "Welcome to the Cloud Foundry go-guestbook sample app backend!")
	})

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
