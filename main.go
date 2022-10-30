package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type Ticket struct {
	ID     int    `json:"id"`
	Owner  string `json:"owner"`
	Status string `json:"status"`
}

var tickets = []Ticket{

	{
		ID:     1,
		Owner:  "Faiyaj",
		Status: "approved",
	},

	{
		ID:     2,
		Owner:  "Zaman",
		Status: "pending",
	},
}

func home(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Welcome to our API")
}

func fetchAllTickets(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(tickets)
}

func fetchTicket(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	// mux.Vars returns all path parameters as a map
	id := mux.Vars(r)["id"]
	currentTicketId, _ := strconv.Atoi(id)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	for _, ticket := range tickets {
		if ticket.ID == currentTicketId {
			json.NewEncoder(w).Encode(ticket)
		}
	}

}

func createTicket(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Ticket Created")
}

func updateTicket(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Ticket Updated")
}

func deleteTicket(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Ticket Deleted")
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/", home)
	router.HandleFunc("/tickets", fetchAllTickets).Methods("GET")
	router.HandleFunc("/tickets/{id}", fetchTicket).Methods("GET")
	router.HandleFunc("/ticket", createTicket).Methods("POST")
	router.HandleFunc("/tikets/{id}", updateTicket).Methods("PATCH")
	router.HandleFunc("/tickets/{id}", deleteTicket).Methods("DELETE")

	runServer := &http.Server{
		Handler: router,
		Addr:    "127.0.0.1:8000",

		// Good practice to enforce timeouts for servers
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(runServer.ListenAndServe())

}
