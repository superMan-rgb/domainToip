package other

import "net/http"

type Fofa struct {
	email string
	key   string
	http  *http.Client
}

func NewClient(email, key string) *Fofa {
	return &Fofa{
		email: email,
		key:   key,
		http:  &http.Client{},
	}
}
