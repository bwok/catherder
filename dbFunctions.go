package main

// extra database functions that don't live in crud,
// database struct definitions, custom sql etc.

type User struct {
	Id        int64
	IdDate    int64
	Name      string `json:"name"`
	Available bool   `json:"available"`
}
type Users []User
type Date struct {
	Id       int64
	IdMeetUp int64
	Date     int64 `json:"date"`
	Users    Users `json:"users"`
}
type Dates []Date
type Admin struct {
	Id       int64
	IdMeetUp int64
	Email    string `json:"email"`
	Alerts   bool   `json:"alerts"`
}
type MeetUp struct {
	Id          int64
	UserHash    string
	AdminHash   string
	Dates       Dates  `json:"dates"`
	Admin       Admin  `json:"admin"`
	Description string `json:"description"`
}
