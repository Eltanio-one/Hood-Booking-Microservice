package data

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
)

// User struct created with necessary information to identify each user.
type User struct {
	ID                  int    `json:"id"`
	Name                string `json:"name"`
	Hash                string `json:"hash"`
	Email               string `json:"email"`
	Emergency_Telephone int    `json:"emergency_telephone"`
	Research_Group      string `json:"research_group"`
}

// UsersList is a type defined to characterise an array of the User struct type variables.
// This is mainly used for defining the temporary userlist, and also in GET/PUT requests of the current registered users.
type UsersList []*User

var UserList UsersList

// GetUsers returns the temporary userlist.
// This userList is to be used as a test for HTTP requests while a database is not incorporated into the project.
func GetUsers(db *sql.DB) UsersList {
	// TODO: query the db, convert sql rows to an array, return
	rows, err := db.Query("SELECT id, username, email, emergency_telephone, research_group FROM users;")
	if err != nil {
		return nil
	}

	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Emergency_Telephone, &user.Research_Group)
		if err != nil {
			return nil
		}
		UserList = append(UserList, &user)
	}
	return UserList
}

// FromJSON can be used on User type objects.
// It takes in an io.Reader parameter, and instantiates a decoder that writes to the io.Reader.
// Uses the json decoder to decode the data stored in the io.Reader and store this data in the User object.
func (u *User) FromJSON(r io.Reader) error {
	dec := json.NewDecoder(r)
	return dec.Decode(u)
}

// ToJSON can be used on UsersList type objects.
// It takes in an io.Writer parameter, and instantiates an encoder that writes to the io.Writer.
// Uses the json encoder to encode the data stored in the User object to the io.Writer.
func (u *UsersList) ToJSON(w io.Writer) error {
	enc := json.NewEncoder(w)
	return enc.Encode(u)
}

// AddUser takes a User struct object as a parameter.
// It calls the function GetNextUserID to collect the next available user ID and assign it to the passed User struct object.
// The user object is then inserted into the database.
func AddUser(u *User, db *sql.DB) error {
	u.ID = GetNextUserID(db)
	if u.ID == -1 {
		return ErrDBQueryError
	}

	// connect to database and user insert to add user object.
	_, err := db.Exec("INSERT INTO users (id, username, passhash, email, emergency_telephone, research_group) VALUES ($1, $2, $3, $4, $5, $6)", u.ID, u.Name, u.Hash, u.Email, u.Emergency_Telephone, u.Research_Group)
	if err != nil {
		return err
	}
	return nil
}

// GetNextUserID is used to find the next numerical ID number and returns an integer of that value.
// Using the length of the UserList, it finds the ID of the last added booking and returns that value plus 1.
func GetNextUserID(db *sql.DB) int {
	// TODO: query the length of the users rows
	var maxID int
	// db query
	rows, err := db.Query("SELECT MAX(id) FROM users")
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

// UpdateUser takes a user ID and a User struct object as parameters and returns an error.
// This function finds the position of the user in the UserList based on their ID, and assigns the passed id to the User struct object.
// The User object at the position located during the function is then overwritten by the passed User parameter.
// TODO: currently all user info is overwritten apart from ID, so need to ensure that if no data was supplied in the PUT that the field is not changed.
func UpdateUser(rw http.ResponseWriter, id int, u *User) error {
	matchedUser, pos, err := findUser(id)
	if err != nil {
		return err
	}

	// ensure user can only change their own data
	if u.ID != id {
		http.Error(rw, "Cannot edit another users data", http.StatusBadRequest)
		return err
	}

	// ensure hash can't be changed
	if u.Hash != "" {
		http.Error(rw, "Unable to edit password hash", http.StatusBadRequest)
		return err
	}

	// go through the data from the request, and any empty fields replace with the data currently stored before updating the userList.
	replaceEmptyFields(matchedUser, u)

	UserList[pos] = u
	return err
}

// findUser takes a user ID as a parameter and returns the corresponding User object, the position of this user in the UserList, and an error.
// If no corresponding ID is found, the function returns the structured ErrUserNotFound alongside nil values for the other return values.
func findUser(id int) (*User, int, error) {
	for i, u := range UserList {
		if u.ID == id {
			return u, i, nil
		}
	}
	return nil, 0, ErrUserNotFound
}

// replaceEmptyFields takes two pointers to User struct objects.
// One of these is pulled from the userList and contains the current data.
// One contains data that has been passed by a user in a PUT request.
// This function iterates over the fields of a User struct and ensures that any missing data the user does not want to change does not get overwritten.
// This is done by replacing any empty fields with the data that is currently stored for that user.
func replaceEmptyFields(stored *User, update *User) {
	vals := reflect.ValueOf(update).Elem()

	for i := 0; i < vals.NumField(); i++ {
		if field := vals.Field(i); field.IsZero() {
			currentField := reflect.ValueOf(stored).Elem().Field(i)
			field.Set(currentField)
		}
	}
}

// create structured error
var ErrUserNotFound = fmt.Errorf("user not found")
var ErrDBQueryError = fmt.Errorf("error querying database")
