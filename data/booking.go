// Package data provides all structs and helper functions.
package data

import (
	"encoding/json"
	"io"
	"time"
)

// Booking is the struct that contains the fields defining a booking.
// This includes;
// the user ID of the user that booked the slot,
// the ID of the hood that was booked,
// the time and date of the booking.
type Booking struct {
	ID          int       `json:"id"`
	UserName    string    `json:"user_name"`
	HoodNumber  int       `json:"hood_number"`
	BookingDate time.Time `json:"booking_time"`
}

// BookingsList is a type defined to characterise an array of the Booking struct type variables.
// This is mainly used for defining the temporary booking list, and also in GET requests of bookings where the bookingList is queried.
type BookingsList []*Booking

// GetBookings returns the bookinglist above.
// This bookingList is to be used as a test for HTTP requests while the database is not linked.
func GetBookings() BookingsList {
	return BookingList
}

// FromJSON can be used on Booking type variables.
// It takes in an io.Reader parameter, and instantiates a decoder that writes to the io.Reader.
// Uses the json decoder to decode the data stored in the io.Reader and store this data in the booking object.
func (b *Booking) FromJSON(r io.Reader) error {
	dec := json.NewDecoder(r)
	return dec.Decode(b)
}

// ToJSON can be used on Bookings type variables.
// It takes in an io.Writer parameter, and instantiates an encoder that writes to the io.Writer.
// Uses the json encoder to encode the data stored in the Booking object to the io.Writer.
func (b *BookingsList) ToJSON(w io.Writer) error {
	enc := json.NewEncoder(w)
	return enc.Encode(b)
}

// AddBooking takes in a Booking struct, and is used to add the passed struct to the temporary bookingList (this will be deprecated once connected to a database).
// The function calls a secondary helper function, GetNextBookingId, see below for details.
func AddBooking(b *Booking) {
	b.ID = GetNextBookingID()
	BookingList = append(BookingList, b)
}

// GetNextBookingID is used to find the next numerical ID number and returns an integer of that value.
// Using the length of the bookingList, it finds the ID of the last added booking and returns that value plus 1.
func GetNextBookingID() int {
	lastBooking := BookingList[len(BookingList)-1]
	return lastBooking.ID + 1
}

// bookingList is a temporary list of bookings used for testing purposes, that will be deprecated once a database is incorporated into this project.
var BookingList = BookingsList{
	{
		ID:          1,
		UserName:    "dan",
		HoodNumber:  101,
		BookingDate: time.Date(2022, time.January, 15, 0, 0, 0, 0, time.UTC),
	},
	{
		ID:          2,
		UserName:    "warren",
		HoodNumber:  201,
		BookingDate: time.Date(2022, time.January, 16, 0, 0, 0, 0, time.UTC),
	},
}
