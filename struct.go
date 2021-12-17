package main

import "time"

type PostData struct {
	Title   string `json:"title"`
	Comment string `json:"body"`
}

type lockData struct {
	LockReason string `json:"lock_reason"`
}

type Data struct {
	Number   int
	HTMLURL  string `json:"html_url"`
	Title    string
	State    string
	User     *User
	CreateAt time.Time `json:"created_at"`
	Body     string
}

type User struct {
	LoginName string `json:"login"`
	HTMLURL   string `json:"html_url"`
}
