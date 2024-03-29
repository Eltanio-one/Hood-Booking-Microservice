package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"reflect"

	"bookings.com/m/data"
	"bookings.com/m/database"
	"golang.org/x/crypto/bcrypt"
)

// Registers is a instantiated as a struct to allow dependency injection of a logger.
type Registers struct {
	l *log.Logger
}

// NewRegisterHandler takes a logger as a parameter and returns a Registers struct that is assigned that logger.
// This function is used to generate the registration handler used in our main function to handle any requests on the registration url path specified in the mux.Handle function.
// This function is used in the main() function to allow user registration.
func NewRegisterHandler(l *log.Logger) *Registers {
	return &Registers{l}
}

// ServeHTTP is called on Registers objects and takes an http ResponseWriter and Request as parameters.
// For registration, only POST http requests are handled. An http.StatusMethodNotAllowed is passed to the ResponseWriter if any other request types are performed.
// A database connection is initialised and the closure of the connection deferred until the function is returned (a request has been made).
func (reg *Registers) ServeHTTP(rw http.ResponseWriter, r *http.Request) {

	// Initialise database connection
	db, err := database.InitialiseConnection(reg.l)
	if err != nil {
		reg.l.Println("Database connection error", err)
		return
	}
	defer db.Close()

	if r.Method == http.MethodPost {
		reg.register(rw, r, db)
		return
	}

	rw.WriteHeader(http.StatusMethodNotAllowed)
}

// register is called on Registers handler struct objects and takes an http ResponseWriter, http Request and the sql database connection as parameters.
// This function is used to authenticate the data that is provided by the user during registration and then add the user to the userList to enable login.
// In order, the data provided in the request body is decoded into a newly instantiated User object.
// checkMissingValues ensure the user has no left any required field empty.
// checkExistingUser ensures that a user with the same name doesn't already exist (due to the small number of people who will be using the API, using the name as an identifier is permissible).
// hashPass hashes the users provided password using the bcrypt package.
// The user is then added to the users table of the database using the data.AddUser function.
func (reg *Registers) register(rw http.ResponseWriter, r *http.Request, db *sql.DB) {
	reg.l.Println("Registering new user...")

	// attempt to decode request body into new user struct
	usr := &data.User{}
	if err := json.NewDecoder(r.Body).Decode(usr); err != nil {
		http.Error(rw, "Unable to unmarshal JSON", http.StatusBadRequest)
		return
	}

	// check for missing data in registration request
	if result := checkMissingValues(usr); !result {
		http.Error(rw, "Please ensure there is no missing data entered", http.StatusBadRequest)
		return
	}

	// need to check user of same name doesn't exist
	nameCheck, err := checkExistingUser(usr, db)
	if nameCheck {
		http.Error(rw, "A user with that name already exists, ensure you don't already have an account", http.StatusBadRequest)
		return
	}
	if err != nil {
		reg.l.Println(err)
		return
	}

	// hash password
	if hash, err := hashPass(usr.Hash); err != nil {
		http.Error(rw, "Hashing of password failed", http.StatusBadRequest)
		return
	} else {
		usr.Hash = hash
	}

	// add new user to the UserList
	if err := data.AddUser(usr, db); err != nil {
		reg.l.Println(err)
		http.Error(rw, "Error adding user to database", http.StatusInternalServerError)
		return
	}

	reg.l.Println("Registration complete!")
}

// checkMissingValues ensures that the user has entered all required data for registration.
// The function takes the previously created User struct as a pointer.
// The user struct is then dereferenced to allow the use of reflect.NumField() to calculate the number of fields within the struct.
// The user should have supplied 5/6 fields, with the hash being blank as this will be generated later in the registration request.
// If the count of missing values is 1, then no data is missing and the function returns true to enable continuation of the request.
func checkMissingValues(u *data.User) bool {
	var count int

	vals := reflect.ValueOf(u).Elem()

	for i := 0; i < vals.NumField(); i++ {
		if fieldValue := vals.Field(i); fieldValue.IsZero() {
			count++
		}
	}
	return count == 1
}

// checkExistingUser takes a User struct object and a sql DB connection as parameters, and returns a bool and an error.
// queries the users table of the database to make sure that a use wth the same name provided at registration doesn't exist.
// As there are a small and limited number of employees who would use this microservice, using name as an identifier of duplicate registrations is fair.
// The function gathers all usernames from the users table of the database.
// Iterating over the rows returned by the query, the username of each row is compared to the name of the user passed as a parameter.
// If the passed name matches any user, then true is returned which will halt the registration.
func checkExistingUser(u *data.User, db *sql.DB) (bool, error) {

	// query db for all usernames
	userRows, err := db.Query("SELECT username FROM users;")
	if err != nil {
		return false, err
	}

	for userRows.Next() {
		var user data.User
		err := userRows.Scan(&user.Name)
		if err != nil {
			return false, err
		}
		if u.Name == user.Name {
			return true, nil
		}
	}
	return false, nil
}

// hashPass takes in the user's inputted password, and returns a hashed version of the password and an error.
// This function uses bcrypt to hash the password.
func hashPass(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}
