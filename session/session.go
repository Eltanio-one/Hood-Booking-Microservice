package session

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"time"
)

var SessionTokens = make(map[string]int)

func GenerateSecureToken(length int) (string, error) {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return "", err
	} else {
		return hex.EncodeToString(b), nil
	}
}

func AuthenticateToken(rw http.ResponseWriter, r *http.Request) string {
	token, err := RetrieveCookie(r)
	if err != nil {
		return ""
	}
	return token
	// returnID := userTokenAuthentication(token)
	// if ID == returnID {
	// 	return true
	// } else {
	// 	return false
	// }
}

func userTokenAuthentication(token string) int {
	id, ok := SessionTokens[token]
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
