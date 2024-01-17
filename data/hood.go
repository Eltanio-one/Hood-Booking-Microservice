package data

import (
	"encoding/json"
	"fmt"
	"io"
)

type Hoods []*Hood

type Hood struct {
	ID          int    `json:"id"`
	Hood_Number int    `json:"hood_number"`
	Room        string `json:"room"`
}

// GetBookings returns the bookinglist above.
// This bookingList is to be used as a test for HTTP requests while the database is not linked.
func GetHoods() Hoods {
	return hoodList
}

// FromJSON can be used on Bookings type variables.
// It takes in an io.Writer parameter, and instantiates a decoder that writes to the io.Writer.
// Uses the decoder to dencode the Bookings type the function is called on.
func (h *Hood) FromJSON(r io.Reader) error {
	dec := json.NewDecoder(r)
	return dec.Decode(h)
}

// ToJSON can be used on Bookings type variables.
// It takes in an io.Writer parameter, and instantiates an encoder that writes to the io.Writer.
// Uses the encoder to encode the Bookings type the function is called on.
func (h *Hoods) ToJSON(w io.Writer) error {
	enc := json.NewEncoder(w)
	return enc.Encode(h)
}

func AddHood(h *Hood) {
	h.ID = GetNextHoodID()
	hoodList = append(hoodList, h)
}

func GetNextHoodID() int {
	lastHood := hoodList[len(hoodList)-1]
	return lastHood.ID + 1
}

func UpdateHood(id int, h *Hood) error {
	_, pos, err := findUser(id)
	if err != nil {
		return err
	}

	h.ID = id
	hoodList[pos] = h
	return err
}

// create structured error
var ErrHoodNotFound = fmt.Errorf("Hood Not Found")

func findHood(id int) (*Hood, int, error) {
	for i, h := range hoodList {
		if h.ID == id {
			return h, i, nil
		}
	}
	return nil, 0, ErrHoodNotFound
}

var hoodList = []*Hood{
	{
		ID:          1,
		Hood_Number: 101,
		Room:        "AN201",
	},
	{
		ID:          2,
		Hood_Number: 102,
		Room:        "AN201",
	},
	{
		ID:          3,
		Hood_Number: 103,
		Room:        "AN202",
	},
	{
		ID:          4,
		Hood_Number: 104,
		Room:        "AN202",
	},
}
