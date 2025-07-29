package domain

import "time"

type User struct {
	ID        int64
	Email     string
	FirstName string
	LastName  string
	Password  Password
	Activated bool
	Version   int
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Password struct {
	Plaintext *string
	Hash      []byte
}
