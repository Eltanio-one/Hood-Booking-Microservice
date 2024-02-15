package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"bookings.com/m/data"
	"bookings.com/m/database"
	"bookings.com/m/session"
	"golang.org/x/crypto/bcrypt"
)

// create a Hoods struct to enable dependency injection of a logger.
type Logins struct {
	l *log.Logger
}

// NewBookingHandler takes a logger object and returns a Logins object.
// The logger passed will be assigned to the Logins object logger field.
// This function is used in the main() function to return the Bookings handler that is required to pass to the created servemux to handle http requests at the specified url path.
func NewLoginHandler(l *log.Logger) *Logins {
	return &Logins{l}
}

// ServeHTTP is called on a Logins struct.
// It takes an http ResponseWriter and Request as parameters.
// This function deals with all HTTP request methods that are queried.
// For logins, only POST requests are permitted and handled.
func (l *Logins) ServeHTTP(rw http.ResponseWriter, r *http.Request) {

	// Initialise database connection
	db, err := database.InitialiseConnection(l.l)
	if err != nil {
		l.l.Println("Database connection error", err)
		return
	}
	defer db.Close()

	if r.Method == http.MethodPost {
		l.login(rw, r, db)
		return
	}

	rw.WriteHeader(http.StatusMethodNotAllowed)
}

// login is called on a Logins struct, and takes an http ResponseWriter and request as arguments.
// This function handles logging in a user.
// It calls various helper functions to authenticate provided data.
// Once the user is authenticated, a session cookie is created and stored
func (l *Logins) login(rw http.ResponseWriter, r *http.Request, db *sql.DB) {
	l.l.Println("Logging in...")

	usr := &data.User{}

	// attempt to decode http request into usr
	err := json.NewDecoder(r.Body).Decode(usr)
	if err != nil {
		http.Error(rw, "Unable to unmarshal JSON", http.StatusBadRequest)
		return
	}

	// Ensure the username exists within the database
	matchedUser := checkUsername(usr.Name, db)
	if matchedUser == nil {
		http.Error(rw, "Invalid username", http.StatusBadRequest)
		return
	}
	if err = comparePasswords(matchedUser.Hash, usr.Hash); err != nil {
		http.Error(rw, "Incorrect password", http.StatusBadRequest)
		return
	}

	// generate secure token and store in session map
	token, err := session.GenerateSecureToken(32)
	if err != nil {
		http.Error(rw, "Failed to generate secure token", http.StatusBadRequest)
		return
	}

	// clear sessiontokens table (should be done in logout, not yet implemented) to avoid pkey error
	_, err = db.Exec("DELETE FROM sessiontokens WHERE user_id = $1;", matchedUser.ID)
	if err != nil {
		http.Error(rw, "Failed to clear previously stored token", http.StatusBadRequest)
		return
	}

	// Add session token to the database, not the local storage
	_, err = db.Exec("INSERT INTO sessiontokens (token, user_id) VALUES ($1, $2);", token, matchedUser.ID)
	if err != nil {
		http.Error(rw, "Failed to store secure token in database", http.StatusBadRequest)
		return
	}

	// Store cookie in Postman
	session.StoreCookie(rw, token)
	l.l.Println("User successfully logged in, welcome", matchedUser.Name)
}

// checekUsername takes a string containing the name provided during login.
// If the name matches a name that is already stored in the userList, then this user is returned.
// This allows checking of the hashed password with the password provided during the login.
func checkUsername(name string, db *sql.DB) *data.User {
	rows, err := db.Query("SELECT id, username, passhash FROM users;")
	if err != nil {
		return nil
	}

	for rows.Next() {
		var user data.User
		err := rows.Scan(&user.ID, &user.Name, &user.Hash)
		if err != nil {
			return nil
		}
		if user.Name == name {
			return &user
		}
	}
	return nil
}

// comparePasswords takes in the hashedPassword that was generated during registration and assigned to the User struct in the userList for this specific user.
// It also takes the inputPassword that was provided at the login request, and returns an error.
// Using bcrypt, it compares the hashedPassword and inputPassword, and will return an error if these do not match.
func comparePasswords(hashedPassword, inputPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(inputPassword))
}
