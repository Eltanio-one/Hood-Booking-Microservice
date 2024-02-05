package data

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"

	"bookings.com/m/database"
)

// Hoods is a type defined to characterise an array of the User struct type variables.
// This is mainly used for defining the temporary hoodlist, and also in GET/PUT requests of the current registered hoods.
type HoodsList []*Hood

// Hood struct created with necessary information to identify each hood.
type Hood struct {
	ID          int    `json:"id"`
	Hood_Number int    `json:"hood_number"`
	Room        string `json:"room"`
}

var HoodList HoodsList

// GetHoods returns the hoodlist above.
// This hoodList is to be used as a test for HTTP requests while the database is not linked.
func GetHoods(db *sql.DB) HoodsList {

	rows, err := db.Query("SELECT id, hood_number, room FROM hoods;")
	if err != nil {
		return nil
	}

	for rows.Next() {
		var hood Hood
		err := rows.Scan(&hood.ID, &hood.Hood_Number, &hood.Room)
		if err != nil {
			return nil
		}
		HoodList = append(HoodList, &hood)
	}
	return HoodList
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
// Uses the json encoder to encode the data stored in the Hood object to the io.Writer.
func (h *HoodsList) ToJSON(w io.Writer) error {
	enc := json.NewEncoder(w)
	return enc.Encode(h)
}

// AddHood takes a Hood struct object as a parameter.
// This function is used to collect the next available hood ID and assign this to the passed Hood object, before appending this hood object to the hoodList.
func AddHood(h *Hood, db *sql.DB) error {
	h.ID = GetNextHoodID(db)
	if h.ID == -1 {
		return database.ErrDBQueryError
	}

	// Add user object to database.
	_, err := db.Exec("INSERT INTO hoods (id, hood_number, room) VALUES ($1, $2, $3)", h.ID, h.Hood_Number, h.Room)
	if err != nil {
		return err
	}
	return nil
}

// GetNextHoodID returns the next available ID as an integer.
// Using the length of the hoodList, it finds the ID of the last added booking and returns that value plus 1.
func GetNextHoodID(db *sql.DB) int {
	var maxID int
	// db query
	rows, err := db.Query("SELECT MAX(id) FROM hoods;")
	if err != nil {
		return -1
	}

	// check if no rows found
	if !rows.Next() {
		if err = rows.Err(); err != nil {
			if err == sql.ErrNoRows {
				return 1
			} else {
				fmt.Println(err)
				return -1
			}
		}
	}
	rows.Scan(&maxID)
	return maxID + 1
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

// create structured error
var ErrHoodNotFound = fmt.Errorf("Hood Not Found")

// func findHood(id int) (*Hood, int, error) {
// 	for i, h := range hoodList {
// 		if h.ID == id {
// 			return h, i, nil
// 		}
// 	}
// 	return nil, 0, ErrHoodNotFound
// }

// hoodList is a temporary list of users used for testing purposes, that will be deprecated once a database is incorporated into this project.
// var HoodList = Hoods{
// 	{
// 		ID:          1,
// 		Hood_Number: 101,
// 		Room:        "AN201",
// 	},
// 	{
// 		ID:          2,
// 		Hood_Number: 102,
// 		Room:        "AN201",
// 	},
// 	{
// 		ID:          3,
// 		Hood_Number: 103,
// 		Room:        "AN202",
// 	},
// 	{
// 		ID:          4,
// 		Hood_Number: 104,
// 		Room:        "AN202",
// 	},
// }
