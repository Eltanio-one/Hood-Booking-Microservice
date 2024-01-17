package handlers

import (
	"log"
	"net/http"

	"bookings.com/m/data"
)

// create a Bookings struct to enable addition of a logger.
type Bookings struct {
	l *log.Logger
}

// NewBookingHandler takes a logger object and returns a Bookings object.
// The logger passed will be assigned to the Bookings object logger field.
// This function is used in the main() function to return the Bookings handler that is required to pass to the created servemux.
func NewBookingHandler(l *log.Logger) *Bookings {
	return &Bookings{l}
}

// ServeHTTP is called on a Bookings object.
// It takes an http ResponseWriter and Request as parameters.
// This function deals with all HTTP request methods that are queried.
func (b *Bookings) ServeHTTP(rw http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		b.getBookings(rw, r)
		return
	}

	if r.Method == http.MethodPost {
		b.addBooking(rw, r)
		return
	}

	rw.WriteHeader(http.StatusMethodNotAllowed)
}

func (b *Bookings) getBookings(rw http.ResponseWriter, r *http.Request) {
	b.l.Println("Handling GET request")

	bookingList := data.GetBookings()

	err := bookingList.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to Marshal JSON", http.StatusInternalServerError)
	}
}

func (b *Bookings) addBooking(rw http.ResponseWriter, r *http.Request) {
	b.l.Println("Handling POST request")

	book := &data.Booking{}

	err := book.FromJSON(r.Body)
	if err != nil {
		// http.Error(rw, "Unable to Marshal JSON FUCK OFF", http.StatusBadRequest)
	}

	b.l.Printf("Booking: %#v", book)
	data.AddBooking(book)
}
