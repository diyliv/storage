package models

import "time"

type User struct {
	Id                  int
	UserName            string
	UserEmail           string
	UserToken           string
	UserHashedPassword  string
	UserUpdatedPassword time.Time
	UserCreatedAt       time.Time
}
