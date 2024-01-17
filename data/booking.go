// Package data provides all structs and helper functions.
package data

import (
	"encoding/json"
	"io"
	"time"
)

type Bookings []*Booking

type Booking struct {
	ID             int       `json:"id"`
	User_ID        int       `json:"user_id"`
	Booked_Hood_ID int       `json:"booked_hood_id"`
	Booking_Date   time.Time `json:"booking_time"`
}

// GetBookings returns the bookinglist above.
// This bookingList is to be used as a test for HTTP requests while the database is not linked.
func GetBookings() Bookings {
	return bookingList
}

// FromJSON can be used on Bookings type variables.
// It takes in an io.Writer parameter, and instantiates a decoder that writes to the io.Writer.
// Uses the decoder to dencode the Bookings type the function is called on.
func (b *Booking) FromJSON(r io.Reader) error {
	dec := json.NewDecoder(r)
	return dec.Decode(b)
}

// ToJSON can be used on Bookings type variables.
// It takes in an io.Writer parameter, and instantiates an encoder that writes to the io.Writer.
// Uses the encoder to encode the Bookings type the function is called on.
func (b *Bookings) ToJSON(w io.Writer) error {
	enc := json.NewEncoder(w)
	return enc.Encode(b)
}

func AddBooking(b *Booking) {
	b.ID = GetNextBookingID()
	bookingList = append(bookingList, b)
}

func GetNextBookingID() int {
	lastBooking := bookingList[len(bookingList)-1]
	return lastBooking.ID + 1
}

var bookingList = []*Booking{
	{
		ID:             1,
		User_ID:        1,
		Booked_Hood_ID: 1,
		Booking_Date:   time.Now(),
	},
	{
		ID:             2,
		User_ID:        2,
		Booked_Hood_ID: 2,
		Booking_Date:   time.Now(),
	},
}
