package data

import (
	"encoding/json"
	"fmt"
	"io"
)

type Users []*User

type User struct {
	ID                  int    `json:"id"`
	Name                string `json:"name"`
	Hash                string `json:"hash"`
	Email               string `json:"email"`
	Emergency_Telephone int    `json:"emergency_telephone"`
	Research_Group      string `json:"research_group"`
}

var UserList = Users{
	{
		ID:                  1,
		Name:                "Dan Haver",
		Hash:                "test123",
		Email:               "dan.haver@ANRI.net",
		Emergency_Telephone: 07712345677,
		Research_Group:      "Immunotherapy",
	},
	{
		ID:                  2,
		Name:                "Warren Patterson",
		Hash:                "test321",
		Email:               "warren.patterson@ANRI.net",
		Emergency_Telephone: 07727654323,
		Research_Group:      "Immunogenetics",
	},
}

// GetBookings returns the bookinglist above.
// This bookingList is to be used as a test for HTTP requests while the database is not linked.
func GetUsers() Users {
	return UserList
}

// FromJSON can be used on Bookings type variables.
// It takes in an io.Writer parameter, and instantiates a decoder that writes to the io.Writer.
// Uses the decoder to dencode the Bookings type the function is called on.
func (u *User) FromJSON(r io.Reader) error {
	dec := json.NewDecoder(r)
	return dec.Decode(u)
}

// ToJSON can be used on Bookings type variables.
// It takes in an io.Writer parameter, and instantiates an encoder that writes to the io.Writer.
// Uses the encoder to encode the Bookings type the function is called on.
func (u *Users) ToJSON(w io.Writer) error {
	enc := json.NewEncoder(w)
	return enc.Encode(u)
}

func AddUser(u *User) {
	u.ID = GetNextUserID()
	UserList = append(UserList, u)
}

func GetNextUserID() int {
	lastUser := UserList[len(bookingList)-1]
	return lastUser.ID + 1
}

func UpdateUser(id int, u *User) error {
	_, pos, err := findUser(id)
	if err != nil {
		return err
	}

	u.ID = id
	UserList[pos] = u
	return err
}

// create structured error
var ErrUserNotFound = fmt.Errorf("User Not Found")

func findUser(id int) (*User, int, error) {
	for i, u := range UserList {
		if u.ID == id {
			return u, i, nil
		}
	}
	return nil, 0, ErrUserNotFound
}
