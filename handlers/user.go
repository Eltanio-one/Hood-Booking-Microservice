package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"regexp"
	"strconv"

	"bookings.com/m/data"
	"bookings.com/m/database"
	"bookings.com/m/session"
)

// create a Users struct to enable addition of a logger and to be used as a handler (http.handler) struct.
type Users struct {
	l *log.Logger
}

// NewBookingHandler takes a logger object and returns a User object.
// The logger passed will be assigned to the Users object logger field.
// This function is used in the main() function to return the Users handler that is required to pass to the created servemux and handle relevant http requests on the passed url path.
func NewUserHandler(l *log.Logger) *Users {
	return &Users{l}
}

// ServeHTTP is called on a Users object.
// It takes an http ResponseWriter and Request as parameters.
// This function deals with GET and PUT HTTP request methods that are queried, as POST methods are covered in the registration handler.
// The session cookie that is generated and stored at login is retrieved here to authenticate the user before returning data client-side.
func (u *Users) ServeHTTP(rw http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		token := session.RetrieveCookie(r)
		if token == "" {
			http.Error(rw, "Unable to retrieve cookie", http.StatusBadRequest)
			return
		}

		// Initialise database connection
		db, err := database.InitialiseConnection(u.l)
		if err != nil {
			u.l.Println("Database connection error", err)
			return
		}
		defer db.Close()

		u.getUsers(rw, r, db)
		return
	}

	if r.Method == http.MethodPut {
		u.l.Println("Handling PUT Request for user")

		// authenticate cookie
		token := session.RetrieveCookie(r)
		if token == "" {
			http.Error(rw, "Unable to retrieve cookie", http.StatusBadRequest)
			return
		}

		// expect the ID in the URI
		reg := regexp.MustCompile(`/([0-9]+)`)
		g := reg.FindAllStringSubmatch(r.URL.Path, -1)

		// ensure only one ID has been returned
		if len(g) != 1 {
			u.l.Println("Invalid URI more than one ID")
			http.Error(rw, "Invalid URI", http.StatusBadRequest)
			return
		}
		// ensure only one capture group returned
		if len(g[0]) != 2 {
			u.l.Println("Invalid URI more than one capture group")
			http.Error(rw, "Invalid URI", http.StatusBadRequest)
			return
		}
		// collect the ID
		idString := g[0][1]
		// convert to int
		id, err := strconv.Atoi(idString)

		// ensure the user can only update their own User object.
		idCheck := session.UserTokenAuthentication(token)
		if idCheck == -1 {
			http.Error(rw, "Error whilst trying to retrieve matching user ID using token", http.StatusBadRequest)
			return
		} else if idCheck != id {
			http.Error(rw, "Permission Denied, User IDs do not match", http.StatusBadRequest)
			return
		}

		if err != nil {
			u.l.Println("Invalid URI unable to convert ID to integer")
			http.Error(rw, "Invalid URI", http.StatusBadRequest)
			return
		}

		u.updateUsers(id, rw, r)

	}

	rw.WriteHeader(http.StatusMethodNotAllowed)
}

// getUsers is called on a Users tye object and takes an HTTP ResponseWriter and Request as parameters.
// This function is involved in handling GET requests of all users.
// Session cookies are authenticated before this function is executed in the ServeHTTP function.
// The full userList is pulled using GetUsers and the data from this list is encoded into this struct.
func (u *Users) getUsers(rw http.ResponseWriter, r *http.Request, db *sql.DB) {
	u.l.Println("Handling GET request for users")

	userList := data.GetUsers(db)

	err := userList.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to Marshal JSON", http.StatusInternalServerError)
	}
}

// updateUsers is called on Users type objects and takes the ID of the user to be updated as an int, and an HTTP ResponseWriter and Request as parameters.
// This function is involved in handling PUT requests for users,
// The data stored in the request body is decoded into a newly instantiated User struct object.
// Using the UpdateUser function, the User with the corresponding ID is updated.
func (u *Users) updateUsers(id int, rw http.ResponseWriter, r *http.Request) {
	u.l.Println("Handle PUT request for user")

	ur := &data.User{}

	// decode supplied data
	err := ur.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "Unable to unmarshal JSON", http.StatusBadRequest)
		return
	}

	err = data.UpdateUser(rw, id, ur)
	if err == data.ErrUserNotFound {
		http.Error(rw, "User not found", http.StatusBadRequest)
		return
	}

	u.l.Println("Update complete!")

}
