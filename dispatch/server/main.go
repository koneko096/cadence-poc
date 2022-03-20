package main

import (
	"context"
	"log"
	"net/http"

	"go.temporal.io/sdk/client"
)

type Server struct {
	Client *client.Client
}

func main() {
	// Create the client object just once per process
	c, err := client.NewClient(client.Options{})
	if err != nil {
		log.Fatalln("unable to create Temporal client", err)
	}
	defer c.Close()

	s := Server{&c}
	httpServer := &http.Server{
		Addr:    ":8091",
		Handler: s.handler(),
	}
	log.Printf("Listening on localhost:8091")

	err = httpServer.ListenAndServe()
	if err != nil {
		log.Fatalln("unable to start Server", err)
	}
}

func (s *Server) handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/arrive", s.arrive)
	return mux
}

func (s *Server) arrive(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		c := *s.Client
		userID := r.URL.Query().Get("uid")
		err := c.SignalWorkflow(context.Background(), "dispatcher-"+userID, "", "DRIVER_ARRIVE", "asdf")
		if err != nil {
			log.Printf("Error signaling client:\n%v", err)
			http.Error(w, "Error.", http.StatusInternalServerError)
		}
	default:
		http.Error(w, "Page not found.", http.StatusNotFound)
	}
}
