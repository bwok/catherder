package main

import (
	"database/sql"
	"log"
)

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

var preparedStmts = make(map[string]*sql.Stmt) // Prepared statements that functions can use.

// A map of sql statements that get prepared in prepareDatabaseStatements()
// They get closed at the termination of the program in closeDatabaseStatements()
var prepStmtInit = map[string]string{
	"insertMeetup":            "INSERT INTO meetup(userhash, adminhash, description) values(?,?,?)",
	"selectMeetup":            "SELECT idmeetup, userhash, adminhash, description FROM meetup WHERE idmeetup= ?",
	"updateMeetup":            "UPDATE meetup SET userhash = ?, adminhash = ?, description = ? WHERE idmeetup = ?",
	"deleteMeetup":            "DELETE from meetup WHERE idmeetup = ?",
	"selectMeetupByUserhash":  "SELECT idmeetup, userhash, adminhash, description FROM meetup WHERE userhash= ?",
	"selectMeetupByAdminhash": "SELECT idmeetup, userhash, adminhash, description FROM meetup WHERE adminhash= ?",

	"insertAdmin":           "INSERT INTO admin(meetup_idmeetup, email, alerts) values(?,?,?)",
	"selectAdmin":           "SELECT idadmin, meetup_idmeetup, email, alerts FROM admin WHERE idadmin= ?",
	"updateAdmin":           "UPDATE admin SET email = ?, alerts = ? WHERE idadmin = ?",
	"deleteAdmin":           "DELETE from admin WHERE idadmin = ?",
	"selectAdminByMeetupid": "SELECT idadmin, meetup_idmeetup, email, alerts FROM admin WHERE meetup_idmeetup= ?",

	"insertDate":            "INSERT INTO date(meetup_idmeetup, date) values(?,?)",
	"selectDate":            "SELECT iddate,meetup_idmeetup,date FROM date WHERE iddate= ?",
	"updateDate":            "UPDATE date SET date = ? WHERE iddate = ?",
	"deleteDate":            "DELETE from date WHERE iddate = ?",
	"selectDatesByMeetupid": "SELECT iddate, meetup_idmeetup, date FROM date WHERE meetup_idmeetup= ?",

	"insertUser":          "INSERT INTO user(date_iddate, name, available) values(?,?,?)",
	"selectUser":          "SELECT iduser,date_iddate,name,available FROM user WHERE iduser= ?",
	"updateUser":          "UPDATE user SET name = ?, available = ? WHERE iduser = ?",
	"deleteUser":          "DELETE from user WHERE iduser = ?",
	"selectUsersByDateid": "SELECT iduser,date_iddate,name,available FROM user WHERE date_iddate= ?",
}

// prepares all the required statements for later use.
// log.Fatal on error
func prepareDatabaseStatements() {
	for key, val := range prepStmtInit {
		stmt, err := db.Prepare(val)

		if err != nil {
			log.Fatal(err)
		}
		preparedStmts[key] = stmt
	}
}

// Closes all prepared statements, logs any errors.
// Called by defer in main() on program termination.
func closeDatabaseStatements() {
	for key := range preparedStmts {
		if err := preparedStmts[key].Close(); err != nil {
			log.Println(err)
		}
	}
}
