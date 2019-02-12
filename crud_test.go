package main

import (
	"database/sql"
	"io/ioutil"
	"os"
	"testing"
)

// Creates a temporary database file
func MakeTestFile(t *testing.T) string {
	file, err := ioutil.TempFile("./", "test_db")
	if err != nil {
		t.Fatal(err)
	}
	file.Close()
	return file.Name()
}

// Creates a database file, a connection, all the needed database tables, foreign keys etc.
func CreateTestDb(t *testing.T) string {
	var err error
	testDbName := MakeTestFile(t)

	db, err = sql.Open("sqlite3", testDbName)
	if err != nil {
		t.Fatal("Failed to open database:", err)
	}

	file, err := os.Open("dbSource.sql")
	if err != nil {
		t.Fatal(err)
	}

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}

	if _, err = db.Exec(string(fileBytes)); err != nil {
		t.Fatal("Failed to create database tables:", err)
	}

	// Prepare all the sql statements for later use
	prepareDatabaseStatements()

	return testDbName
}

// Destroys the test database, closes the prepared statements etc
func DestroyTestDb(testDbName string) {
	defer os.Remove(testDbName)
	defer db.Close()
	defer closeDatabaseStatements()
}

// Tests the MeetUp.Read method too
func TestMeetUp_Create(t *testing.T) {
	testDbName := CreateTestDb(t)
	defer DestroyTestDb(testDbName)

	var tests = []MeetUp{
		{UserHash: "abc", AdminHash: "def", Description: "meetUp description"},
		{Id: -1, UserHash: "abc", AdminHash: "def", Description: "meetUp description"},
	}

	for _, meetUp := range tests {
		if err := meetUp.Create(); err != nil {
			t.Errorf("meetUp create failed: %s\n", err)
		} else if meetUp.Id <= 0 {
			t.Fatal("no id returned for inserted meetup row")
		} else {
			retMeetUp := MeetUp{}
			if err = retMeetUp.Read(meetUp.Id); err != nil {
				t.Fatalf("couldn't read row back from meetup table: %s\n", err)
			}

			if retMeetUp.Id != meetUp.Id || retMeetUp.UserHash != meetUp.UserHash || retMeetUp.AdminHash != meetUp.AdminHash || retMeetUp.Description != meetUp.Description {
				t.Errorf("returned row from DB was different to the one inserted. inserted: %+v, returned: %+v\n", meetUp, retMeetUp)
			}
		}
	}
}

func TestMeetUp_Update(t *testing.T) {
	testDbName := CreateTestDb(t)
	defer DestroyTestDb(testDbName)

	var meetUp = MeetUp{UserHash: "abc", AdminHash: "def", Description: "test description"}

	if err := meetUp.Create(); err != nil {
		t.Fatalf("create failed: %s\n", err)
	}

	meetUp.UserHash = "xyz"
	meetUp.AdminHash = "uvw"
	meetUp.Description = "rst"
	if err := meetUp.Update(); err != nil {
		t.Fatalf("update failed: %s\n", err)
	}

	retMeetUp := MeetUp{}
	if err := retMeetUp.Read(meetUp.Id); err != nil {
		t.Fatalf("couldn't read row back from meetup table: %s\n", err)
	}

	if retMeetUp.Id != meetUp.Id || retMeetUp.UserHash != "xyz" || retMeetUp.AdminHash != "uvw" || retMeetUp.Description != "rst" {
		t.Errorf("returned row from DB was different to the one updated. updated: %+v, returned: %+v\n", meetUp, retMeetUp)
	}
}

func TestMeetUp_Delete(t *testing.T) {
	testDbName := CreateTestDb(t)
	defer DestroyTestDb(testDbName)

	var meetUp = MeetUp{UserHash: "abc", AdminHash: "def", Description: "test description"}

	if err := meetUp.Create(); err != nil {
		t.Fatalf("create failed: %s\n", err)
	}

	retMeetUp := MeetUp{}
	if err := retMeetUp.Read(meetUp.Id); err != nil {
		t.Fatalf("couldn't read row back from meetup table: %s\n", err)
	}

	if err := meetUp.Delete(); err != nil {
		t.Fatalf("delete failed: %s\n", err)
	}

	retMeetUp = MeetUp{}
	if err := retMeetUp.Read(meetUp.Id); err.Error() != "no rows" {
		t.Fatalf("couldn't read row back from meetup table: %s\n", err)
	}
}

// Tests the Admin.Read method too
func TestAdmin_Create(t *testing.T) {
	testDbName := CreateTestDb(t)
	defer DestroyTestDb(testDbName)

	meetUp := MeetUp{UserHash: "abc", AdminHash: "def", Description: "meetUp description"}
	if err := meetUp.Create(); err != nil {
		t.Errorf("meetUp create failed: %s\n", err)
	}

	var tests = []Admin{
		{Email: "abc", Alerts: true},
		{Id: -1, Email: "def", Alerts: false},
	}

	for _, admin := range tests {
		admin.IdMeetUp = meetUp.Id

		if err := admin.Create(); err != nil {
			t.Errorf("admin create failed: %s\n", err)
		} else if admin.Id <= 0 {
			t.Fatal("no id returned for inserted admin row")
		} else {
			retAdmin := Admin{}
			if err = retAdmin.Read(admin.Id); err != nil {
				t.Fatalf("couldn't read row back from admin table: %s\n", err)
			}

			if retAdmin.Id != admin.Id || retAdmin.IdMeetUp != admin.IdMeetUp || retAdmin.Email != admin.Email || retAdmin.Alerts != admin.Alerts {
				t.Errorf("returned row from DB was different to the one inserted. inserted: %+v, returned: %+v\n", admin, retAdmin)
			}
		}
	}
}

func TestAdmin_Update(t *testing.T) {
	testDbName := CreateTestDb(t)
	defer DestroyTestDb(testDbName)

	var meetUp = MeetUp{UserHash: "abc", AdminHash: "def", Description: "test description"}

	if err := meetUp.Create(); err != nil {
		t.Fatalf("create failed: %s\n", err)
	}

	var admin = Admin{IdMeetUp: meetUp.Id, Email: "abc", Alerts: true}
	if err := admin.Create(); err != nil {
		t.Fatalf("create failed: %s\n", err)
	}

	admin.Email = "xyz"
	admin.Alerts = false

	if err := admin.Update(); err != nil {
		t.Fatalf("update failed: %s\n", err)
	}

	retAdmin := Admin{}
	if err := retAdmin.Read(admin.Id); err != nil {
		t.Fatalf("couldn't read row back from admin table: %s\n", err)
	}

	if retAdmin.Id != admin.Id || retAdmin.Email != "xyz" || retAdmin.Alerts != false || retAdmin.IdMeetUp != admin.IdMeetUp {
		t.Errorf("returned row from DB was different to the one updated. updated: %+v, returned: %+v\n", admin, retAdmin)
	}
}

func TestAdmin_Delete(t *testing.T) {
	testDbName := CreateTestDb(t)
	defer DestroyTestDb(testDbName)

	var meetUp = MeetUp{UserHash: "abc", AdminHash: "def", Description: "test description"}

	if err := meetUp.Create(); err != nil {
		t.Fatalf("create failed: %s\n", err)
	}

	var admin = Admin{IdMeetUp: meetUp.Id, Email: "abc", Alerts: true}
	if err := admin.Create(); err != nil {
		t.Fatalf("create failed: %s\n", err)
	}

	retAdmin := Admin{}
	if err := retAdmin.Read(admin.Id); err != nil {
		t.Fatalf("couldn't read row back from admin table: %s\n", err)
	}

	if err := admin.Delete(); err != nil {
		t.Fatalf("delete failed: %s\n", err)
	}

	retAdmin = Admin{}
	if err := retAdmin.Read(admin.Id); err.Error() != "no rows" {
		t.Fatalf("couldn't read row back from admin table: %s\n", err)
	}
}

// Tests the Date.Read method too
func TestDate_Create(t *testing.T) {
	testDbName := CreateTestDb(t)
	defer DestroyTestDb(testDbName)

	meetUp := MeetUp{UserHash: "abc", AdminHash: "def", Description: "meetUp description"}
	if err := meetUp.Create(); err != nil {
		t.Errorf("meetUp create failed: %s\n", err)
	}

	var tests = []Date{
		{Date: 1549537200000},
		{Id: -1, Date: 1549537200000},
	}

	for _, date := range tests {
		date.IdMeetUp = meetUp.Id

		if err := date.Create(); err != nil {
			t.Errorf("date create failed: %s\n", err)
		} else if date.Id <= 0 {
			t.Fatal("no id returned for inserted date row")
		} else {
			retDate := Date{}
			if err = retDate.Read(date.Id); err != nil {
				t.Fatalf("couldn't read row back from date table: %s\n", err)
			}

			if retDate.Id != date.Id || retDate.IdMeetUp != date.IdMeetUp || retDate.Date != date.Date {
				t.Errorf("returned row from DB was different to the one inserted. inserted: %+v, returned: %+v\n", date, retDate)
			}
		}
	}
}

func TestDate_Update(t *testing.T) {
	testDbName := CreateTestDb(t)
	defer DestroyTestDb(testDbName)

	var meetUp = MeetUp{UserHash: "abc", AdminHash: "def", Description: "test description"}
	if err := meetUp.Create(); err != nil {
		t.Fatalf("create failed: %s\n", err)
	}

	var date = Date{IdMeetUp: meetUp.Id, Date: 1549537200000}
	if err := date.Create(); err != nil {
		t.Fatalf("create failed: %s\n", err)
	}

	var newTimestamp int64 = 1549537200099
	date.Date = newTimestamp

	if err := date.Update(); err != nil {
		t.Fatalf("update failed: %s\n", err)
	}

	retDate := Date{}
	if err := retDate.Read(date.Id); err != nil {
		t.Fatalf("couldn't read row back from date table: %s\n", err)
	}

	if retDate.Id != date.Id || retDate.Date != newTimestamp || retDate.IdMeetUp != date.IdMeetUp {
		t.Errorf("returned row from DB was different to the one updated. updated: %+v, returned: %+v\n", date, retDate)
	}
}

func TestDate_Delete(t *testing.T) {
	testDbName := CreateTestDb(t)
	defer DestroyTestDb(testDbName)

	var meetUp = MeetUp{UserHash: "abc", AdminHash: "def", Description: "test description"}
	if err := meetUp.Create(); err != nil {
		t.Fatalf("create failed: %s\n", err)
	}

	var date = Date{IdMeetUp: meetUp.Id, Date: 1549537200000}
	if err := date.Create(); err != nil {
		t.Fatalf("create failed: %s\n", err)
	}

	retDate := Date{}
	if err := retDate.Read(date.Id); err != nil {
		t.Fatalf("couldn't read row back from date table: %s\n", err)
	}

	if err := date.Delete(); err != nil {
		t.Fatalf("delete failed: %s\n", err)
	}

	retDate = Date{}
	if err := retDate.Read(date.Id); err.Error() != "no rows" {
		t.Fatalf("couldn't read row back from date table: %s\n", err)
	}
}

// Tests the User.Read method too
func TestUser_Create(t *testing.T) {
	testDbName := CreateTestDb(t)
	defer DestroyTestDb(testDbName)

	meetUp := MeetUp{UserHash: "abc", AdminHash: "def", Description: "meetUp description"}
	if err := meetUp.Create(); err != nil {
		t.Errorf("meetUp create failed: %s\n", err)
	}

	var date = Date{IdMeetUp: meetUp.Id, Date: 1549537200000}
	if err := date.Create(); err != nil {
		t.Fatalf("create failed: %s\n", err)
	}

	var tests = []User{
		{Name: "bob", Available: true},
		{Id: -1, Name: "harry", Available: false},
	}

	for _, user := range tests {
		user.IdDate = date.Id

		if err := user.Create(); err != nil {
			t.Errorf("user create failed: %s\n", err)
		} else if user.Id <= 0 {
			t.Fatal("no id returned for inserted user row")
		} else {
			retUser := User{}
			if err = retUser.Read(user.Id); err != nil {
				t.Fatalf("couldn't read row back from user table: %s\n", err)
			}

			if retUser.Id != user.Id || retUser.IdDate != user.IdDate || retUser.Name != user.Name || retUser.Available != user.Available {
				t.Errorf("returned row from DB was different to the one inserted. inserted: %+v, returned: %+v\n", user, retUser)
			}
		}
	}
}

func TestUser_Update(t *testing.T) {
	testDbName := CreateTestDb(t)
	defer DestroyTestDb(testDbName)

	var meetUp = MeetUp{UserHash: "abc", AdminHash: "def", Description: "test description"}
	if err := meetUp.Create(); err != nil {
		t.Fatalf("create failed: %s\n", err)
	}

	var date = Date{IdMeetUp: meetUp.Id, Date: 1549537200000}
	if err := date.Create(); err != nil {
		t.Fatalf("create failed: %s\n", err)
	}

	var user = User{IdDate: date.Id, Name: "bob", Available: true}
	if err := user.Create(); err != nil {
		t.Fatalf("create failed: %s\n", err)
	}

	var newName string = "harry"
	var newAvailable bool = false
	user.Name = newName
	user.Available = newAvailable

	if err := user.Update(); err != nil {
		t.Fatalf("update failed: %s\n", err)
	}

	retUser := User{}
	if err := retUser.Read(user.Id); err != nil {
		t.Fatalf("couldn't read row back from user table: %s\n", err)
	}

	if retUser.Id != user.Id || retUser.Name != newName || retUser.IdDate != user.IdDate || retUser.Available != user.Available {
		t.Errorf("returned row from DB was different to the one updated. updated: %+v, returned: %+v\n", date, retUser)
	}
}

func TestUser_Delete(t *testing.T) {
	testDbName := CreateTestDb(t)
	defer DestroyTestDb(testDbName)

	var meetUp = MeetUp{UserHash: "abc", AdminHash: "def", Description: "test description"}
	if err := meetUp.Create(); err != nil {
		t.Fatalf("create failed: %s\n", err)
	}

	var date = Date{IdMeetUp: meetUp.Id, Date: 1549537200000}
	if err := date.Create(); err != nil {
		t.Fatalf("create failed: %s\n", err)
	}

	var user = User{IdDate: date.Id, Name: "bob", Available: true}
	if err := user.Create(); err != nil {
		t.Fatalf("create failed: %s\n", err)
	}

	retUser := User{}
	if err := retUser.Read(user.Id); err != nil {
		t.Fatalf("couldn't read row back from user table: %s\n", err)
	}

	if err := user.Delete(); err != nil {
		t.Fatalf("delete failed: %s\n", err)
	}

	retUser = User{}
	if err := retUser.Read(user.Id); err.Error() != "no rows" {
		t.Fatalf("couldn't read row back from date table: %s\n", err)
	}
}
