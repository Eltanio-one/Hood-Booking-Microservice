package data

import (
	"encoding/json"
	"io"
)

// Hood struct created with necessary information to identify each hood.
type Hood struct {
	ID          int    `json:"id"`
	Hood_Number int    `json:"hood_number"`
	Room        string `json:"room"`
}

// HoodsList is a type defined to characterise an array of the User struct type variables.
// This is mainly used for defining the temporary hoodlist, and also in GET/PUT requests of the current registered hoods.
type HoodsList []*Hood

// GetHoods returns the hoodlist above.
// This hoodList is to be used as a test for HTTP requests while the database is not linked.
func GetHoods() Hoods {
	return hoodList
}

// FromJSON can be used on Hood struct objects.
// It takes in an io.Writer parameter, and instantiates a decoder that writes to the io.Writer, before returning an error (this shou ld be nil if all has worked).
// Uses the decoder to decode the data read from the io.Reader and store it in the Hood object.
func (h *Hood) FromJSON(r io.Reader) error {
	dec := json.NewDecoder(r)
	return dec.Decode(h)
}

// ToJSON can be used on Hoods type variables.
// It takes in an io.Writer parameter, and instantiates an encoder that writes to the io.Writer.
// Uses the encoder to encode the data to the io.Writer which will then display the data in JSON format to the user.
func (h *Hoods) ToJSON(w io.Writer) error {
	enc := json.NewEncoder(w)
	return enc.Encode(h)
}

// AddHood takes a Hood struct object as a parameter.
// This function is used to collect the next available hood ID and assign this to the passed Hood object, before appending this hood object to the hoodList.
func AddHood(h *Hood) {
	h.ID = GetNextHoodID()
	hoodList = append(hoodList, h)
}

// GetNextHoodID returns the next available ID as an integer.
// Using the length of the hoodList, it finds the ID of the last added booking and returns that value plus 1.
func GetNextHoodID() int {
	lastHood := hoodList[len(hoodList)-1]
	return lastHood.ID + 1
}

// func UpdateHood(id int, h *Hood) error {
// 	_, pos, err := findUser(id)
// 	if err != nil {
// 		return err
// 	}

// 	h.ID = id
// 	hoodList[pos] = h
// 	return err
// }

// // create structured error
// var ErrHoodNotFound = fmt.Errorf("Hood Not Found")

// func findHood(id int) (*Hood, int, error) {
// 	for i, h := range hoodList {
// 		if h.ID == id {
// 			return h, i, nil
// 		}
// 	}
// 	return nil, 0, ErrHoodNotFound
// }

// hoodList is a temporary list of users used for testing purposes, that will be deprecated once a database is incorporated into this project.
var hoodList = HoodsList{
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
