package session

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"time"
)

var SessionCookies = make(map[string]int)

func GenerateSecureToken(length int) (string, error) {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return "", err
	} else {
		return hex.EncodeToString(b), nil
	}
}

func AuthenticateToken(rw http.ResponseWriter, r *http.Request, ID int) bool {
	token, err := RetrieveCookie(r)
	if err != nil {
		http.Error(rw, "Unable to retrieve cookie", http.StatusBadRequest)
		return false
	}
	returnID := userTokenAuthentication(token)
	if ID == returnID {
		return true
	} else {
		return false
	}
}

func userTokenAuthentication(token string) int {
	id, ok := SessionCookies[token]
	if ok {
		return id
	}
	return -1
}

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

func RetrieveCookie(r *http.Request) (string, error) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}
