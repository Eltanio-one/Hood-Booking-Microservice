package handlers

import (
	"log"
	"net/http"

	"bookings.com/m/data"
	"bookings.com/m/session"
)

// Hoods struct is created to enable dependency injection of a logger.
type Hoods struct {
	l *log.Logger
}

// NewHoodHandler takes a logger object and returns a Hoods object.
// The logger passed will be assigned to the Hoods object logger field.
// This function is used in the main() function to return the Hoods handler that is required to pass to the created servemux.
func NewHoodHandler(l *log.Logger) *Hoods {
	return &Hoods{l}
}

// ServeHTTP is called on a Hoods object.
// It takes an http ResponseWriter and Request as parameters.
// This function deals with all HTTP request methods that are queried, so far GET and POST requests are handled.
// Before each request is handled, the session token is authenticated to ensure login has been performed.
func (h *Hoods) ServeHTTP(rw http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		token := session.RetrieveCookie(r)
		if token == "" {
			http.Error(rw, "Unable to retrieve cookie", http.StatusBadRequest)
			return
		}

		h.getHoods(rw, r)
		return
	}

	if r.Method == http.MethodPost {
		token := session.RetrieveCookie(r)
		if token == "" {
			http.Error(rw, "Unable to retrieve cookie", http.StatusBadRequest)
			return
		}

		h.addHood(rw, r)
		return
	}

	rw.WriteHeader(http.StatusMethodNotAllowed)
}

// getHoods is called on a Hoods object and takes an http ResponseWriter and Request as parameters.
// This function is responsible for handling GET requests for Hoods.
// It calls functions "GetHoods" and "ToJSON" from the booking data file to retrieve and encode the data to be presented to the user.
func (h *Hoods) getHoods(rw http.ResponseWriter, r *http.Request) {
	h.l.Println("Handling GET request for hoods")

	// retrieve hoodList
	hoodList := data.GetHoods()

	// encode data
	err := hoodList.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to Marshal JSON", http.StatusInternalServerError)
	}
}

// addHood is called on a Hoods struct object and takes an HTTP ResponseWriter and Request as parameters.
// This function is responsible for handling POST requests for Hoods.
// It calls the function "FromJSON" from the Hood data file to decode the data being passed by the user.
// The decoded data is then passed to the function AddHood from the Hood data file to add the file to the hoodList.
func (h *Hoods) addHood(rw http.ResponseWriter, r *http.Request) {
	h.l.Println("Handling POST request for hoods")

	hd := &data.Hood{}

	err := hd.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "Unable to Marshal JSON", http.StatusBadRequest)
	}

	h.l.Printf("Hood: %#v", hd)
	data.AddHood(hd)
}
