package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

// extra database functions that don't live in crud,
// database struct definitions, custom sql etc.

type User struct {
	Id       int64
	IdMeetUp int64
	Name     string  `json:"name"`
	Dates    []int64 `json:"dates"` // dates the user is available for. This is a UNIX timestamp in milliseconds, as per ecma script defines it. "The number of milliseconds between 1 January 1970 00:00:00 UTC and the given date."
}
type Users []User
type MeetUp struct {
	Id          int64
	UserHash    string  `json:"userhash"`
	AdminHash   string  `json:"adminhash"`
	AdminEmail  string  `json:"adminemail"`
	SendAlerts  bool    `json:"sendalerts"`
	Dates       []int64 `json:"dates"` // This is a UNIX timestamp in milliseconds, as per ecma script defines it. "The number of milliseconds between 1 January 1970 00:00:00 UTC and the given date."
	Description string  `json:"description"`
	Users       Users   `json:"users"`
}

// Prepared statements that functions can use.
// They get closed at the termination of the program in closeDatabaseStatements()
var preparedStmts = make(map[string]*sql.Stmt)

// prepares all the required statements for later use.
// log.Fatal on error
func prepareDatabaseStatements() {
	// A map of sql statements that get prepared in prepareDatabaseStatements()
	var prepStmtInit = map[string]string{
		"insertMeetup":            `INSERT INTO meetup(userhash, adminhash, adminemail, sendalerts, dates, description) values(?,?,?,?,?,?)`,
		"selectMeetup":            `SELECT idmeetup, userhash, adminhash, adminemail, sendalerts, dates, description FROM meetup WHERE idmeetup = ?`,
		"updateMeetup":            `UPDATE meetup SET adminemail = ?, sendalerts = ?, dates = ?, description = ? WHERE idmeetup = ?`,
		"deleteMeetup":            `DELETE from meetup WHERE idmeetup = ?`,
		"selectMeetupByUserhash":  `SELECT idmeetup, userhash, adminhash, adminemail, sendalerts, dates, description FROM meetup WHERE userhash = ?`,
		"selectMeetupByAdminhash": `SELECT idmeetup, userhash, adminhash, adminemail, sendalerts, dates, description FROM meetup WHERE adminhash = ?`,
		"deleteMeetupByAdminhash": `DELETE from meetup WHERE adminhash = ?`,

		"insertUser":            `INSERT INTO "user"(idmeetup, name, dates) values(?,?,?)`,
		"selectUser":            `SELECT iduser,idmeetup,name,dates FROM "user" WHERE iduser = ?`,
		"updateUser":            `UPDATE "user" SET name = ?, dates = ? WHERE iduser = ?`,
		"deleteUser":            `DELETE from "user" WHERE iduser = ?`,
		"selectUsersByMeetUpid": `SELECT * FROM "user" WHERE idmeetup = ?`,
	}

	for key, val := range prepStmtInit {
		stmt, err := db.Prepare(val)

		if err != nil {
			log.Fatalf("prepareDatabaseStatements failed on key %q, val %q, error (%s)", key, val, err)
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

// MarshalJSON Set json output format and fields
func (m *MeetUp) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		UserHash    string  `json:"userhash"`
		AdminHash   string  `json:"adminhash"`
		AdminEmail  string  `json:"adminemail"`
		SendAlerts  bool    `json:"sendalerts"`
		Dates       []int64 `json:"dates"`
		Description string  `json:"description"`
		Users       Users   `json:"users"`
	}{
		m.UserHash,
		m.AdminHash,
		m.AdminEmail,
		m.SendAlerts,
		m.Dates,
		m.Description,
		m.Users,
	})
}

// MarshalJSON Set json output format and fields
func (u *User) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Name  string  `json:"name"`
		Dates []int64 `json:"dates"`
	}{
		u.Name,
		u.Dates,
	})
}

// DeleteByAdminHash Deletes a meetup by its admin hash. Deletes get cascaded to the other tables.
func (m *MeetUp) DeleteByAdminHash(adminHash string) error {
	if _, err := preparedStmts["deleteMeetupByAdminhash"].Exec(adminHash); err != nil {
		return err
	}

	return nil
}

// GetByUserHash Selects a MeetUp row by the user hash
// Also gets all sub objects of the MeetUp row from the date, admin and user tables.
func (m *MeetUp) GetByUserHash(userHash string) (retErr error) {
	rows, retErr := preparedStmts["selectMeetupByUserhash"].Query(userHash)
	if retErr != nil {
		return
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			retErr = fmt.Errorf("%s unable to close rows %s", retErr, closeErr)
		}
	}()

	if rows.Next() {
		var datesBlob []byte
		retErr = rows.Scan(&m.Id, &m.UserHash, &m.AdminHash, &m.AdminEmail, &m.SendAlerts, &datesBlob, &m.Description)
		if retErr != nil {
			return
		}
		m.Dates = convertBlobToDates(datesBlob)
	} else {
		retErr = errors.New("no rows matching the userhash")
		return
	}

	// Read all users with meetup id
	retErr = m.Users.GetAllByMeetUpId(m.Id)
	if retErr != nil {
		return
	}

	return nil
}

// GetByAdminHash Selects a MeetUp row by the admin hash
// Also gets all sub objects of the MeetUp row from the date, admin and user tables.
func (m *MeetUp) GetByAdminHash(adminHash string) (retErr error) {
	rows, retErr := preparedStmts["selectMeetupByAdminhash"].Query(adminHash)
	if retErr != nil {
		return
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			retErr = fmt.Errorf("%s unable to close rows %s", retErr, closeErr)
		}
	}()

	if rows.Next() {
		var datesBlob []byte
		retErr = rows.Scan(&m.Id, &m.UserHash, &m.AdminHash, &m.AdminEmail, &m.SendAlerts, &datesBlob, &m.Description)
		if retErr != nil {
			return
		}
		m.Dates = convertBlobToDates(datesBlob)
	} else {
		retErr = errors.New("no rows matching the adminhash")
		return
	}

	// Read all users with meetupid
	retErr = m.Users.GetAllByMeetUpId(m.Id)
	if retErr != nil {
		return
	}

	return nil
}

// GetAllByMeetUpId Selects all User rows with meetup id
func (u *Users) GetAllByMeetUpId(idMeetUp int64) (retErr error) {
	rows, retErr := preparedStmts["selectUsersByMeetUpid"].Query(idMeetUp)
	if retErr != nil {
		return
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			retErr = fmt.Errorf("%s unable to close rows %s", retErr, closeErr)
		}
	}()

	*u = make(Users, 0)
	for rows.Next() {
		var user = User{}
		var datesBlob []byte
		retErr = rows.Scan(&user.Id, &user.IdMeetUp, &user.Name, &datesBlob)
		if retErr != nil {
			return
		}
		user.Dates = convertBlobToDates(datesBlob)
		*u = append(*u, user)
	}

	return nil
}
