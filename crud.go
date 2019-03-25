package main

import (
	"errors"
	"fmt"
	"github.com/lib/pq"
)

// CRUD, don't edit by hand.

/*
MeetUp CRUD
*/
func (m *MeetUp) Create() error {
	//`INSERT INTO meetup.meetup(userhash, adminhash, adminemail, sendalerts, dates, description) values($1,$2,$3,$4,$5,$6) RETURNING idmeetup`,
	err := preparedStmts["insertMeetup"].QueryRow(m.UserHash, m.AdminHash, m.AdminEmail, m.SendAlerts, pq.Array(m.Dates), m.Description).Scan(&m.Id)
	if err != nil {
		return err
	}

	return nil
}
func (m *MeetUp) Read(id int64) (retErr error) {
	rows, retErr := preparedStmts["selectMeetup"].Query(id)
	if retErr != nil {
		return
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			retErr = fmt.Errorf("%s unable to close rows %s", retErr, closeErr)
		}
	}()

	if rows.Next() {
		retErr = rows.Scan(&m.Id, &m.UserHash, &m.AdminHash, &m.AdminEmail, &m.SendAlerts, pq.Array(&m.Dates), &m.Description)
		if retErr != nil {
			return
		}
	} else {
		retErr = errors.New("no rows")
		return
	}

	return nil
}
func (m *MeetUp) Update() error {
	_, err := preparedStmts["updateMeetup"].Exec(m.AdminEmail, m.SendAlerts, pq.Array(m.Dates), m.Description, m.Id)
	if err != nil {
		return err
	}

	return nil
}
func (m *MeetUp) Delete() error {
	if _, err := preparedStmts["deleteMeetup"].Exec(m.Id); err != nil {
		return err
	}

	return nil
}

/*
User CRUD
*/
func (u *User) Create() error {
	err := preparedStmts["insertUser"].QueryRow(u.IdMeetUp, u.Name, pq.Array(u.Dates)).Scan(&u.Id)
	if err != nil {
		return err
	}

	return nil
}
func (u *User) Read(id int64) (retErr error) {
	rows, retErr := preparedStmts["selectUser"].Query(id)
	if retErr != nil {
		return
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			retErr = fmt.Errorf("%s unable to close rows %s", retErr, closeErr)
		}
	}()

	if rows.Next() {
		retErr = rows.Scan(&u.Id, &u.IdMeetUp, &u.Name, pq.Array(&u.Dates))
		if retErr != nil {
			return
		}
	} else {
		retErr = errors.New("no rows")
		return
	}

	return nil
}
func (u *User) Update() error {
	_, err := preparedStmts["updateUser"].Exec(u.Name, pq.Array(u.Dates), u.Id)
	if err != nil {
		return err
	}

	return nil
}
func (u *User) Delete() error {
	_, err := preparedStmts["deleteUser"].Exec(u.Id)
	if err != nil {
		return err
	}

	return nil
}
