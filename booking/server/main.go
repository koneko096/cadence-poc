package main

import (
	"cadence-poc/booking"
	"context"
	"encoding/json"
	"log"
	"net/http"

	"go.temporal.io/sdk/client"
)

type Server struct {
	Client *client.Client
}

func main() {
	c, err := client.NewClient(client.Options{})
	if err != nil {
		log.Fatalln("unable to create Temporal client\n", err)
	}
	defer c.Close()

	s := Server{&c}
	httpServer := &http.Server{
		Addr:    ":8090",
		Handler: s.handler(),
	}
	log.Printf("Listening on localhost:8090")

	err = httpServer.ListenAndServe()
	if err != nil {
		log.Fatalln("unable to start Server", err)
	}
}

func (s *Server) handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/book", s.book)
	return mux
}

func (s *Server) book(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		workflowOptions := client.StartWorkflowOptions{
			TaskQueue: "BOOKING_QUEUE",
		}

		var param booking.BookingRequest
		err := json.NewDecoder(r.Body).Decode(&param)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		c := *s.Client
		_, err = c.ExecuteWorkflow(context.Background(), workflowOptions, booking.BookingWorkflow, &param)
		if err != nil {
			log.Printf("Workflow crashed\n%v", err)
			http.Error(w, "Internal server error.", http.StatusInternalServerError)
		}
	default:
		http.Error(w, "Page not found.", http.StatusNotFound)
	}
}
