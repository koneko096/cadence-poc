package booking

import (
	"context"
	"log"
	"net/http"

	"go.temporal.io/sdk/client"
)

type Server struct {
	Port   string
	Client *client.Client
}

func (s *Server) handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/book", s.book)
	return mux
}

func (s *Server) Run() {
	httpServer := &http.Server{
		Addr:    ":" + s.Port,
		Handler: s.handler(),
	}

	err := httpServer.ListenAndServe()
	if err != nil {
		log.Fatalln("unable to start Server", err)
	}
}

func (s *Server) book(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		workflowOptions := client.StartWorkflowOptions{
			TaskQueue: "BOOKING_QUEUE",
		}

		param := BookingRequest{}

		c := *s.Client
		_, err := c.ExecuteWorkflow(context.Background(), workflowOptions, BookingWorkflow, &param)
		if err != nil {
			log.Printf("Workflow crashed\n%v", err)
			http.Error(w, "Internal server error.", http.StatusInternalServerError)
		}
	default:
		http.Error(w, "Page not found.", http.StatusNotFound)
	}
}
