package session

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"net/http"
	"time"

	"bookings.com/m/data"
)

// create a map to store tokens in.
// This is not the safest option, but for this proof of concept works well for verifying that a user has their own token being passed.
// var SessionTokens = make(map[string]int)

// GenerateSecureToken takes in a length (desired length of secure token) as an int and returns a string and an error.
// a slice of bytes of the given length is instantiated.
// rand.Read is used to generate random bytes, and the hexadecimal encoding of these random bytes is returned.
func GenerateSecureToken(length int) (string, error) {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return "", err
	} else {
		return hex.EncodeToString(b), nil
	}
}

// UserTokenAuthentication takes a token as a string and returns an int.
// This function is used to ensure that the user attempting a PUT http request is editing their own user object and not someone elses.
// The verification of the ID is performed in the relevant handler file on the ID that is returned from this function.
func UserTokenAuthentication(token string, db *sql.DB) int {
	// Query the db for the token, if present return the id, otherwise return -1

	rows, err := db.Query("SELECT user_id FROM sessiontokens WHERE token = $1;", token)
	if err != nil {
		return -1
	}

	for rows.Next() {
		var user data.User
		err := rows.Scan(&user.ID)
		if err != nil {
			return -1
		} else {
			return user.ID
		}
	}
	return -1

}

// StoreCookie takes an HTTP ResponseWriter and a token as a string as parameters.
// This function stores the token as an http.Cookie struct.
// The cookie is then set to the ResponseWriters header.
func StoreCookie(rw http.ResponseWriter, token string) {
	cookie := http.Cookie{
		Name:     "session_token",
		Value:    token,
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
	}

	http.SetCookie(rw, &cookie)
}

// RetrieveCookie takes an http Request as a parameter anf returns a string.
// This function attempts to retrieve the set cookie from the http request.
// The value of the cookie is returned if there is no error, otherwise an empty string is returned.
func RetrieveCookie(r *http.Request) string {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		return ""
	}
	return cookie.Value
}
