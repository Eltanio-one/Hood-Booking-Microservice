package handlers

import (
	"log"
	"net/http"
	"reflect"

	"bookings.com/m/data"
	"bookings.com/m/database"
	"bookings.com/m/session"
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
// Before each request is handled, the session token is authenticated to ensure login has been performed.
func (b *Bookings) ServeHTTP(rw http.ResponseWriter, r *http.Request) {

	// Initialise database connection
	db, err := database.InitialiseConnection(b.l)
	if err != nil {
		b.l.Println("Database connection error", err)
		return
	}
	defer db.Close()

	if r.Method == http.MethodGet {
		token := session.RetrieveCookie(r)
		if token == "" {
			http.Error(rw, "Unable to retrieve cookie", http.StatusBadRequest)
			return
		}

		b.getBookings(rw, r)
		return
	}

	// TODO: Need to ensure the user can only create a booking with their ID.
	if r.Method == http.MethodPost {
		token := session.RetrieveCookie(r)
		if token == "" {
			http.Error(rw, "Unable to retrieve cookie", http.StatusBadRequest)
			return
		}

		b.addBooking(rw, r, token)
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
func (b *Bookings) addBooking(rw http.ResponseWriter, r *http.Request, token string) {
	b.l.Println("Handling POST request")

	book := &data.Booking{}

	err := book.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "Unable to Marshal JSON", http.StatusBadRequest)
	}

	// TODO: check for missing values
	if result := checkMissingValuesBooking(book); !result {
		http.Error(rw, "Please ensure there is no missing data entered", http.StatusBadRequest)
		return
	}

	// TODO: Ensure that ID value from `SessionTokens` token key matches the ID the user is trying to book for.
	idCheck := session.UserTokenAuthentication(token)
	if idCheck == -1 {
		http.Error(rw, "Error whilst trying to retrieve matching user ID using token", http.StatusBadRequest)
		return
	}

	// verify the hood exists in the hoodList
	if hoodCheck := checkHoodExists(book.HoodNumber); !hoodCheck {
		http.Error(rw, "That hood number does not exist", http.StatusBadRequest)
		return
	}

	// the following two functions can be reworked once the database is connected, can get user ID verification via session token.
	// and query the db with the user ID to get the relevant information to do booking checks, won't need position in list.
	// get the position of the user in the userlist (deprecated once database connection added here).
	pos := getUserListPos(idCheck)
	if pos == -1 {
		http.Error(rw, "User does not exist in the user list, consider registration", http.StatusBadRequest)
		return
	}

	// ensure that the user booking is the same as the user being booked for (also depracated once session tokens added).
	userName := data.UserList[pos].Name
	if userName != book.UserName {
		http.Error(rw, "Currently cannot book a hood for another user", http.StatusBadRequest)
		return
	}

	// verify hood and user are free at the booked time.
	for _, booking := range data.BookingList {
		if booking.BookingDate.Equal(book.BookingDate) {
			if booking.UserName == userName {
				b.l.Printf("You are already booked into hood %d at the requested time!", booking.HoodNumber)
				http.Error(rw, "Booking failed as previous booking exists at this time", http.StatusBadRequest)
				return
			}
			if booking.HoodNumber == book.HoodNumber {
				b.l.Printf("This hood is already booked at that time by %s!", booking.UserName)
				http.Error(rw, "Booking failed as previous booking exists at this time", http.StatusBadRequest)
				return
			}
		}
	}

	b.l.Printf("Booking: %#v", book)
	data.AddBooking(book)
}

// checkMissingValues ensures that the user has entered all required data for the booking.
// The function takes the previously created Booking struct as a pointer.
// The Booking struct is then dereferenced to allow the use of reflect.NumField() to calculate the number of fields within the struct.
// If the count of missing values is 1, then no data is missing and the function returns true to enable continuation of the request.
func checkMissingValuesBooking(b *data.Booking) bool {
	var count int

	vals := reflect.ValueOf(b).Elem()

	for i := 0; i < vals.NumField(); i++ {
		if fieldValue := vals.Field(i); fieldValue.IsZero() {
			count++
		}
	}
	return count == 1
}

// getUserListPos takes an id as an int and returns an int.
// This function is used to iterate over the UserList and return the position of the user with the matching ID in the userList.
// If no matching ID is found, then -1 is returned.
func getUserListPos(id int) int {
	for i, u := range data.UserList {
		if u.ID == id {
			return i
		}
	}
	return -1
}

// checkHoodExists takes a hood number as an int and returns a bool.
// This function is used to ensure that the hood number included in the users POST request exists.
func checkHoodExists(hoodNumber int) bool {
	for _, hood := range data.HoodList {
		if hood.Hood_Number == hoodNumber {
			return true
		}
	}
	return false
}
