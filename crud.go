package main

import (
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

func (m *MeetUp) Create() error {
	datesBlob := convertDatesToBlob(m.Dates)

	result, err := preparedStmts["insertMeetup"].Exec(m.UserHash, m.AdminHash, m.AdminEmail, m.SendAlerts, datesBlob, m.Description)
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
		var datesBlob []byte
		retErr = rows.Scan(&m.Id, &m.UserHash, &m.AdminHash, &m.AdminEmail, &m.SendAlerts, &datesBlob, &m.Description)
		if retErr != nil {
			return
		}
		m.Dates = convertBlobToDates(datesBlob)
	} else {
		retErr = errors.New("no rows")
		return
	}

	return nil
}
func (m *MeetUp) Update() error {
	datesBlob := convertDatesToBlob(m.Dates)
	_, err := preparedStmts["updateMeetup"].Exec(m.AdminEmail, m.SendAlerts, datesBlob, m.Description, m.Id)
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

func (u *User) Create() error {
	datesBlob := convertDatesToBlob(u.Dates)

	result, err := preparedStmts["insertUser"].Exec(u.IdMeetUp, u.Name, datesBlob)
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
		var datesBlob []byte
		retErr = rows.Scan(&u.Id, &u.IdMeetUp, &u.Name, &datesBlob)
		if retErr != nil {
			return
		}
		u.Dates = convertBlobToDates(datesBlob)
	} else {
		retErr = errors.New("no rows")
		return
	}

	return nil
}
func (u *User) Update() error {
	datesBlob := convertDatesToBlob(u.Dates)

	_, err := preparedStmts["updateUser"].Exec(u.Name, datesBlob, u.Id)
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
