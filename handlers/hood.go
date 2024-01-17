package handlers

import (
	"log"
	"net/http"

	"bookings.com/m/data"
)

// create a Hoods struct to enable addition of a logger.
type Hoods struct {
	l *log.Logger
}

// NewBookingHandler takes a logger object and returns a Bookings object.
// The logger passed will be assigned to the Bookings object logger field.
// This function is used in the main() function to return the Bookings handler that is required to pass to the created servemux.
func NewHoodHandler(l *log.Logger) *Hoods {
	return &Hoods{l}
}

// ServeHTTP is called on a Bookings object.
// It takes an http ResponseWriter and Request as parameters.
// This function deals with all HTTP request methods that are queried.
func (h *Hoods) ServeHTTP(rw http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		h.getHoods(rw, r)
		return
	}

	if r.Method == http.MethodPost {
		h.addHood(rw, r)
		return
	}

	rw.WriteHeader(http.StatusMethodNotAllowed)
}

func (h *Hoods) getHoods(rw http.ResponseWriter, r *http.Request) {
	h.l.Println("Handling GET request for hoods")

	hoodList := data.GetHoods()

	err := hoodList.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to Marshal JSON", http.StatusInternalServerError)
	}
}

func (h *Hoods) addHood(rw http.ResponseWriter, r *http.Request) {
	h.l.Println("Handling POST request for hoods")

	hd := &data.Hood{}

	err := hd.FromJSON(r.Body)
	if err != nil {
		// http.Error(rw, "Unable to Marshal JSON FUCK OFF", http.StatusBadRequest)
	}

	h.l.Printf("Hood: %#v", hd)
	data.AddHood(hd)
}

func (h *Hoods) updateUsers(id int, rw http.ResponseWriter, r *http.Request) {
	h.l.Println("Handle PUT request for hood")

	hd := &data.Hood{}

	err := hd.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "Unable to unmarshal JSON", http.StatusBadRequest)
	}

	err = data.UpdateHood(id, hd)
	if err == data.ErrHoodNotFound {
		http.Error(rw, "Hood not found", http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(rw, "Hood not found", http.StatusInternalServerError)
		return
	}
}
