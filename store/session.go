package store

import (
	"net/http"
	"time"

	"github.com/google/uuid"
)

type Values map[string]interface{}

type session struct {
	Values  Values
	expires time.Time
}

var sessionMap = make(map[string]*session)

func Session(name string, w http.ResponseWriter, r *http.Request) (Values, error) {
	c, err := r.Cookie(name)
	var sessionToken string
	if err != nil {
		if err == http.ErrNoCookie {
			sessionToken = uuid.New().String()
		} else {
			return nil, err
		}
	} else {
		sessionToken = c.Value
	}

	expiresAt := time.Now().Add(120 * time.Second)
	http.SetCookie(w, &http.Cookie{
		Name:    name,
		Value:   sessionToken,
		Expires: expiresAt,
	})

	if s, ok := sessionMap[sessionToken]; ok {
		s.expires = expiresAt
		return s.Values, nil
	} else {
		s = &session{Values: make(map[string]interface{}), expires: expiresAt}
		sessionMap[sessionToken] = s
		return s.Values, nil
	}
}

func (v Values) Sync() error {
	return nil
}
