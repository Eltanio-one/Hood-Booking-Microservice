package handlers

import (
	"log"
	"net/http"

	"bookings.com/m/data"
)

// Bookings struct is created to enable dependency injection of a logger.
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

// getBookings can be called on a Bookings object and takes an http ResponseWriter and Request as parameters.
// This function is responsible for handling GET requests for bookings.
// It calls functions "GetBookings" and "ToJSON" from the booking data file to retrieve and encode the data to be presented to the user.
func (b *Bookings) getBookings(rw http.ResponseWriter, r *http.Request) {
	b.l.Println("Handling GET request")

	// retrieve bookingList
	bookingList := data.GetBookings()

	// encode data
	err := bookingList.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to Marshal JSON", http.StatusInternalServerError)
	}
}

// addBooking can be called on a Bookings object and takes an http ResponseWriter and Request as parameters.
// This function is responsible for handling POST requests for bookings.
// It calls the function "FromJSON" from the booking data file to decode the data being passed by the user.
// The decoded data is then passed to the function AddBooking from the booking data file to add the file to the bookingList.
func (b *Bookings) addBooking(rw http.ResponseWriter, r *http.Request) {
	b.l.Println("Handling POST request")

	book := &data.Booking{}

	err := book.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "Unable to Marshal JSON", http.StatusBadRequest)
	}

	b.l.Printf("Booking: %#v", book)
	data.AddBooking(book)
}
