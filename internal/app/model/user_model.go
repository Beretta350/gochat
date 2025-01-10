package model

// TODO: Put a validator for the user (required fields and some rules)

type User struct {
	Username string `json:"username"`
	Fullname string `json:"fullname"`
	Status   string `json:"status"`
}

func NewUser(username, fullname, status string) *User {
	return &User{Username: username, Fullname: fullname, Status: status}
}
