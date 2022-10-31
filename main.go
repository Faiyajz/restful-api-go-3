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

// add some tickets to our application by using a slice of tickets as our database
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
	fmt.Fprintf(w, "Welcome to our API.\n status %v\n", http.StatusOK)
}

func fetchAllTickets(w http.ResponseWriter, r *http.Request) {

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tickets)
}

func fetchTicket(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	// mux.Vars returns all path parameters as a map
	id := mux.Vars(r)["id"]                            //get id from request
	currentTicketId, _ := strconv.Atoi(id)             //convert id string type into int type id
	w.Header().Set("Content-Type", "application/json") //Set the headers and the response
	w.WriteHeader(http.StatusOK)

	//iterate all the tickets and match with the requested ticket id.
	for _, ticket := range tickets {

		//if matched then show the ticket
		if ticket.ID == currentTicketId {
			json.NewEncoder(w).Encode(ticket)
		}
	}

}

func createTicket(w http.ResponseWriter, r *http.Request) {
	var newTicket Ticket //create an instance of Ticket struct

	//read data from our requests by passing the body of our http request e.g. json.NewDecoder(r.Body)
	//Call .Decode() passing it a pointer to our newTicket Struct which is an instance of Ticket Struct
	//which allows it to match the json to the appropriate properties of the struct
	if err := json.NewDecoder(r.Body).Decode(&newTicket); err != nil {

		//send an internal server error
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Error parsing JSON request")

		log.Fatal(err)
	}

	tickets = append(tickets, newTicket)               //add new new ticket in the tickets slice
	w.Header().Set("Content-Type", "application/json") //Set the headers and the response
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(newTicket) //ticket created

}

func updateTicket(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]                //request id
	currentTicketId, _ := strconv.Atoi(id) //convert string to int

	var updatedTicket Ticket //create an instance of the Ticket struct

	w.Header().Set("Content-Type", "application/json") //Set the headers and the response

	//read data from our requests by passing the body of our http request e.g. json.NewDecoder(r.Body)
	//Call .Decode() passing it a pointer to our newTicket Struct which is an instance of Ticket Struct
	//which allows it to match the json to the appropriate properties of the struct
	if err := json.NewDecoder(r.Body).Decode(&updatedTicket); err != nil {
		//send an internal server error

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error parsing JSON request")

		log.Fatal(err)
	}

	//iterate all the tickets and match with the requested ticket id.
	for index, ticket := range tickets {

		//if the request id matched, update the ticket info
		if ticket.ID == currentTicketId {

			ticket.Owner = updatedTicket.Owner
			ticket.Status = updatedTicket.Status

			tickets[index] = ticket
			json.NewEncoder(w).Encode(ticket)
		}

	}

}

func deleteTicket(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]                //request id
	currentTicketId, _ := strconv.Atoi(id) //convert string to int

	var tmpTickets []Ticket //create an instance of Ticket struct

	//iterate over the tickets
	for index, ticket := range tickets {
		
		//if id matched then dlete the ticket
		if ticket.ID == currentTicketId {

			tmpTickets = append(tmpTickets[:1], tickets[index+1:]...)
			//So with tmpTickets[:1] we are telling our append function to use our tmpTickets slice up to, but not including,
			//the current index as the base of our new slice. By passing tickets[index+1:]... as the second argument,
			//we are appending each element in the original slice starting at the position after the current item.
			//The ... tells Go that even though we are passing append a slice, we want it to treat each element as a separate argument.
			//If we were to omit the ... it would throw an error because it would try to append the entire slice as a single argument
			//and the type wouldn’t match.

			fmt.Fprintf(w, "The ticket with ID %v has been deleted. \n", currentTicketId)
		}
	}

	tickets = tmpTickets
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
