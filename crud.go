package main

import (
	"errors"
	"fmt"
)

// CRUD, don't edit by hand.

/*
MeetUp CRUD
*/
func (m *MeetUp) Create() error {
	result, err := preparedStmts["insertMeetup"].Exec(m.UserHash, m.AdminHash, m.AdminEmail, m.SendAlerts, m.Dates, m.Description)
	if err != nil {
		return err
	}

	m.Id, err = result.LastInsertId()
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
		retErr = rows.Scan(&m.Id, &m.UserHash, &m.AdminHash, &m.AdminEmail, &m.SendAlerts, &m.Dates, &m.Description)
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
	_, err := preparedStmts["updateMeetup"].Exec(m.AdminEmail, m.SendAlerts, m.Dates, m.Description, m.Id)
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
	result, err := preparedStmts["insertUser"].Exec(u.IdMeetUp, u.Name, u.Dates)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	u.Id = id
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
		retErr = rows.Scan(&u.Id, &u.IdMeetUp, &u.Name, &u.Dates)
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
	_, err := preparedStmts["updateUser"].Exec(u.Name, u.Dates, u.Id)
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
