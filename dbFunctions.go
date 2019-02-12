package main

import (
	"database/sql"
	"errors"
	"fmt"
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

// Creates a meetup and all children in one transaction. Does not update the receiver.
func (m *MeetUp) CreateMeetUp() error {
	insTx, err := db.Begin()
	if err != nil {
		return err
	}

	// Do MeetUp Insert
	resultMeetUp, err := insTx.Stmt(preparedStmts["insertMeetup"]).Exec(m.UserHash, m.AdminHash, m.Description)
	if err != nil {
		if rollBkErr := insTx.Rollback(); rollBkErr != nil {
			return fmt.Errorf("%s: %s", err, rollBkErr)
		}
		return err
	}
	idMeetUp, err := resultMeetUp.LastInsertId()
	if err != nil {
		if rollBkErr := insTx.Rollback(); rollBkErr != nil {
			return fmt.Errorf("%s: %s", err, rollBkErr)
		}
		return err
	}

	// Do Admin Insert
	_, err = insTx.Stmt(preparedStmts["insertAdmin"]).Exec(idMeetUp, m.Admin.Email, m.Admin.Alerts)
	if err != nil {
		if rollBkErr := insTx.Rollback(); rollBkErr != nil {
			return fmt.Errorf("%s: %s", err, rollBkErr)
		}
		return err
	}

	for _, date := range m.Dates {
		// Do Dates Insert
		resultDate, err := insTx.Stmt(preparedStmts["insertDate"]).Exec(idMeetUp, date.Date)
		if err != nil {
			if rollBkErr := insTx.Rollback(); rollBkErr != nil {
				return fmt.Errorf("%s: %s", err, rollBkErr)
			}
			return err
		}
		idDate, err := resultDate.LastInsertId()
		if err != nil {
			if rollBkErr := insTx.Rollback(); rollBkErr != nil {
				return fmt.Errorf("%s: %s", err, rollBkErr)
			}
			return err
		}
		// Do Users Insert
		for _, user := range date.Users {
			_, err := insTx.Stmt(preparedStmts["insertUser"]).Exec(idDate, user.Name, user.Available)
			if err != nil {
				if rollBkErr := insTx.Rollback(); rollBkErr != nil {
					return fmt.Errorf("%s: %s", err, rollBkErr)
				}
				return err
			}
		}
	}

	err = insTx.Commit()
	if err != nil {
		return err
	}
	return nil
}

// Selects a MeetUp row by the user hash
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
		retErr = rows.Scan(&m.Id, &m.UserHash, &m.AdminHash, &m.Description)
		if retErr != nil {
			return
		}
	} else {
		retErr = errors.New("no rows matching the userhash")
		return
	}

	// Read admin
	retErr = m.Admin.GetByMeetUpId(m.Id)
	if retErr != nil {
		return
	}

	// Read all dates with dateid
	retErr = m.Dates.GetAllByMeetUpId(m.Id)
	if retErr != nil {
		return
	}

	return nil
}

// Selects a MeetUp row by the admin hash
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
		retErr = rows.Scan(&m.Id, &m.UserHash, &m.AdminHash, &m.Description)
		if retErr != nil {
			return
		}
	} else {
		retErr = errors.New("no rows matching the adminhash")
		return
	}

	// Read admin
	retErr = m.Admin.GetByMeetUpId(m.Id)
	if retErr != nil {
		return
	}

	// Read all dates with dateid
	retErr = m.Dates.GetAllByMeetUpId(m.Id)
	if retErr != nil {
		return
	}

	return nil
}

// Selects an Admin row by a meetup id
func (a *Admin) GetByMeetUpId(idMeetup int64) (retErr error) {
	rows, retErr := preparedStmts["selectAdminByMeetupid"].Query(idMeetup)
	if retErr != nil {
		return
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			retErr = fmt.Errorf("%s unable to close rows %s", retErr, closeErr)
		}
	}()

	for rows.Next() {
		retErr = rows.Scan(&a.Id, &a.IdMeetUp, &a.Email, &a.Alerts)
		if retErr != nil {
			return
		}
	}

	return nil
}

// Selects all Date rows by a meetup id
// Also gets all User sub objects for each Date object.
func (d *Dates) GetAllByMeetUpId(idMeetup int64) (retErr error) {
	rows, retErr := preparedStmts["selectDatesByMeetupid"].Query(idMeetup)
	if retErr != nil {
		return
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			retErr = fmt.Errorf("%s unable to close rows %s", retErr, closeErr)
		}
	}()

	for rows.Next() {
		var date = Date{}
		retErr = rows.Scan(&date.Id, &date.IdMeetUp, &date.Date)
		if retErr != nil {
			return
		}

		// Read all users with dateid
		retErr = date.Users.GetAllByDateId(date.Id)
		if retErr != nil {
			return
		}
		*d = append(*d, date)
	}

	return nil
}

// Selects all User rows with date id
func (u *Users) GetAllByDateId(idDate int64) (retErr error) {
	rows, retErr := preparedStmts["selectUsersByDateid"].Query(idDate)
	if retErr != nil {
		return
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			retErr = fmt.Errorf("%s unable to close rows %s", retErr, closeErr)
		}
	}()

	for rows.Next() {
		var user = User{}
		retErr = rows.Scan(&user.Id, &user.IdDate, &user.Name, &user.Available)
		if retErr != nil {
			return
		}
		*u = append(*u, user)
	}

	return nil
}

// Creates all users in one transaction. Does not update the receiver.
func (u *Users) CreateUsers() error {
	insTx, err := db.Begin()
	if err != nil {
		return err
	}

	// Do Users Insert
	for _, user := range *u {
		_, err := insTx.Stmt(preparedStmts["insertUser"]).Exec(user.IdDate, user.Name, user.Available)
		if err != nil {
			if rollErr := insTx.Rollback(); rollErr != nil {
				return rollErr
			}
			return err
		}
	}

	err = insTx.Commit()
	if err != nil {
		return err
	}
	return nil
}