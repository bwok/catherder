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
	result, err := preparedStmts["insertMeetup"].Exec(m.UserHash, m.AdminHash, m.Description)
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
		retErr = rows.Scan(&m.Id, &m.UserHash, &m.AdminHash, &m.Description)
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
	_, err := preparedStmts["updateMeetup"].Exec(m.UserHash, m.AdminHash, m.Description, m.Id)
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
Admin CRUD
*/
func (a *Admin) Create() error {
	result, err := preparedStmts["insertAdmin"].Exec(a.IdMeetUp, a.Email, a.Alerts)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	a.Id = id
	return nil
}
func (a *Admin) Read(id int64) (retErr error) {
	rows, retErr := preparedStmts["selectAdmin"].Query(id)
	if retErr != nil {
		return
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			retErr = fmt.Errorf("%s unable to close rows %s", retErr, closeErr)
		}
	}()

	if rows.Next() {
		retErr = rows.Scan(&a.Id, &a.IdMeetUp, &a.Email, &a.Alerts)
		if retErr != nil {
			return
		}
	} else {
		retErr = errors.New("no rows")
		return
	}

	return nil
}
func (a *Admin) Update() error {
	_, err := preparedStmts["updateAdmin"].Exec(a.Email, a.Alerts, a.Id)
	if err != nil {
		return err
	}

	return nil
}
func (a *Admin) Delete() error {
	_, err := preparedStmts["deleteAdmin"].Exec(a.Id)
	if err != nil {
		return err
	}

	return nil
}

/*
Date CRUD
*/

func (d *Date) Create() error {
	result, err := preparedStmts["insertDate"].Exec(d.IdMeetUp, d.Date)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	d.Id = id

	return nil
}
func (d *Date) Read(id int64) (retErr error) {
	rows, retErr := preparedStmts["selectDate"].Query(id)
	if retErr != nil {
		return
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			retErr = fmt.Errorf("%s unable to close rows %s", retErr, closeErr)
		}
	}()

	if rows.Next() {
		retErr = rows.Scan(&d.Id, &d.IdMeetUp, &d.Date)
		if retErr != nil {
			return
		}
	} else {
		retErr = errors.New("no rows")
		return
	}

	return nil
}
func (d *Date) Update() error {
	_, err := preparedStmts["updateDate"].Exec(d.Date, d.Id)
	if err != nil {
		return err
	}

	return nil
}
func (d *Date) Delete() error {
	_, err := preparedStmts["deleteDate"].Exec(d.Id)
	if err != nil {
		return err
	}

	return nil
}

/*
User CRUD
*/
func (u *User) Create() error {
	result, err := preparedStmts["insertUser"].Exec(u.IdDate, u.Name, u.Available)
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
		retErr = rows.Scan(&u.Id, &u.IdDate, &u.Name, &u.Available)
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
	_, err := preparedStmts["updateUser"].Exec(u.Name, u.Available, u.Id)
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
