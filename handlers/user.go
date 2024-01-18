package handlers

import (
	"log"
	"net/http"
	"regexp"
	"strconv"

	"bookings.com/m/data"
	"bookings.com/m/session"
)

// create a Users struct to enable addition of a logger.
type Users struct {
	l *log.Logger
}

// NewBookingHandler takes a logger object and returns a Bookings object.
// The logger passed will be assigned to the Bookings object logger field.
// This function is used in the main() function to return the Bookings handler that is required to pass to the created servemux.
func NewUserHandler(l *log.Logger) *Users {
	return &Users{l}
}

// ServeHTTP is called on a Users object.
// It takes an http ResponseWriter and Request as parameters.
// This function deals with GET and PUT HTTP request methods that are queried, as POST methods are covered in the registration handler.
// The session cookie that is generated and stored at login is retrieved here to authenticate the user before returning data client-side.
func (u *Users) ServeHTTP(rw http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		token := session.AuthenticateToken(rw, r)
		if token == "" {
			http.Error(rw, "Unable to retrieve cookie", http.StatusBadRequest)
			return
		}
		u.getUsers(rw, r)
		return
	}

	if r.Method == http.MethodPut {
		u.l.Println("Handling PUT Request for user")

		// authenticate cookie
		token := session.AuthenticateToken(rw, r)
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

		if err != nil {
			u.l.Println("Invalid URI unable to convert ID to integer")
			http.Error(rw, "Invalid URI", http.StatusBadRequest)
			return
		}

		u.updateUsers(id, rw, r)

	}

	rw.WriteHeader(http.StatusMethodNotAllowed)
}

func (u *Users) getUsers(rw http.ResponseWriter, r *http.Request) {
	u.l.Println("Handling GET request for users")

	userList := data.GetUsers()

	err := userList.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to Marshal JSON", http.StatusInternalServerError)
	}
}

func (u *Users) updateUsers(id int, rw http.ResponseWriter, r *http.Request) {
	u.l.Println("Handle PUT request for user")

	ur := &data.User{}

	err := ur.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "Unable to unmarshal JSON", http.StatusBadRequest)
	}

	err = data.UpdateUser(id, ur)
	if err == data.ErrUserNotFound {
		http.Error(rw, "User not found", http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(rw, "User not found", http.StatusInternalServerError)
		return
	}
}
