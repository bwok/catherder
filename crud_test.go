package main

import (
	"database/sql"
	"io/ioutil"
	"os"
	"reflect"
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

	if db, err = sql.Open("sqlite3", "file:"+testDbName+"?_foreign_keys=1"); err != nil {
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
		{UserHash: "abc", AdminHash: "def", AdminEmail: "testy@testy.test", SendAlerts: true, Dates: []int64{1550401200000, 1550487600000, 1550574000000, 1550660400000, 1550746800000, 1550833200000, 1550919600000, 1551006000000}, Description: "meetUp description"},
		{Id: -1, UserHash: "abc", AdminHash: "def", AdminEmail: "testy@testy.test", SendAlerts: true, Dates: []int64{}, Description: "meetUp description"},
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

			if retMeetUp.Id != meetUp.Id || compareMeetUpObjects(retMeetUp, meetUp) == false {
				t.Errorf("returned row from DB was different to the one inserted. inserted: %+v, returned: %+v\n", meetUp, retMeetUp)
			}

			if err := retMeetUp.Delete(); err != nil {
				t.Fatalf("delete failed: %s\n", err)
			}
		}
	}
}

func TestMeetUp_Update(t *testing.T) {
	testDbName := CreateTestDb(t)
	defer DestroyTestDb(testDbName)

	var meetUp = MeetUp{UserHash: "abc", AdminHash: "def", AdminEmail: "testy@testy.test", SendAlerts: true, Dates: []int64{1550401200000, 1550487600000, 1550574000000, 1550660400000, 1550746800000, 1550833200000, 1550919600000, 1551006000000}, Description: "meetUp description"}

	if err := meetUp.Create(); err != nil {
		t.Fatalf("create failed: %s\n", err)
	}

	meetUp.AdminEmail = "abc@def.hij"
	meetUp.SendAlerts = false
	meetUp.Dates = []int64{1550401200000}
	meetUp.Description = "rst"
	if err := meetUp.Update(); err != nil {
		t.Fatalf("update failed: %s\n", err)
	}

	retMeetUp := MeetUp{}
	if err := retMeetUp.Read(meetUp.Id); err != nil {
		t.Fatalf("couldn't read row back from meetup table: %s\n", err)
	}

	if retMeetUp.Id != meetUp.Id || compareMeetUpObjects(retMeetUp, meetUp) == false {
		t.Errorf("returned row from DB was different to the one updated. updated: %+v, returned: %+v\n", meetUp, retMeetUp)
	}

	if err := retMeetUp.Delete(); err != nil {
		t.Fatalf("delete failed: %s\n", err)
	}
}

func TestMeetUp_Delete(t *testing.T) {
	testDbName := CreateTestDb(t)
	defer DestroyTestDb(testDbName)

	var meetUp = MeetUp{UserHash: "abc", AdminHash: "def", AdminEmail: "testy@testy.test", SendAlerts: true, Dates: []int64{1550401200000, 1550487600000, 1550574000000, 1550660400000, 1550746800000, 1550833200000, 1550919600000, 1551006000000}, Description: "meetUp description"}

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

// Tests the User.Read method too
func TestUser_Create(t *testing.T) {
	testDbName := CreateTestDb(t)
	defer DestroyTestDb(testDbName)

	var meetUp = MeetUp{UserHash: "abc", AdminHash: "def", AdminEmail: "testy@testy.test", SendAlerts: true, Dates: []int64{1550401200000, 1550487600000, 1550574000000, 1550660400000, 1550746800000, 1550833200000, 1550919600000, 1551006000000}, Description: "meetUp description"}
	if err := meetUp.Create(); err != nil {
		t.Errorf("meetUp create failed: %s\n", err)
	}

	var tests = []User{
		{Name: "bob", Dates: []int64{1550401200000, 1550487600000, 1550574000000}},
		{Id: -1, Name: "harry", Dates: []int64{}},
	}

	for _, user := range tests {
		user.IdMeetUp = meetUp.Id

		if err := user.Create(); err != nil {
			t.Errorf("user create failed: %s\n", err)
		} else if user.Id <= 0 {
			t.Fatal("no id returned for inserted user row")
		} else {
			retUser := User{}
			if err = retUser.Read(user.Id); err != nil {
				t.Fatalf("couldn't read row back from user table: %s\n", err)
			}

			if retUser.Id != user.Id || retUser.IdMeetUp != user.IdMeetUp || retUser.Name != user.Name || reflect.DeepEqual(retUser.Dates, user.Dates) == false {
				t.Errorf("returned row from DB was different to the one inserted. inserted: %+v, returned: %+v\n", user, retUser)
			}

			if err := retUser.Delete(); err != nil {
				t.Fatalf("delete failed: %s\n", err)
			}
		}
	}

	if err := meetUp.Delete(); err != nil {
		t.Fatalf("delete failed: %s\n", err)
	}
}

func TestUser_Update(t *testing.T) {
	testDbName := CreateTestDb(t)
	defer DestroyTestDb(testDbName)

	var meetUp = MeetUp{UserHash: "abc", AdminHash: "def", AdminEmail: "testy@testy.test", SendAlerts: true, Dates: []int64{1550401200000, 1550487600000, 1550574000000, 1550660400000, 1550746800000, 1550833200000, 1550919600000, 1551006000000}, Description: "meetUp description"}
	if err := meetUp.Create(); err != nil {
		t.Fatalf("create failed: %s\n", err)
	}

	var user = User{IdMeetUp: meetUp.Id, Name: "bob", Dates: []int64{1550401200000, 1550487600000, 1550574000000}}
	if err := user.Create(); err != nil {
		t.Fatalf("create failed: %s\n", err)
	}

	var newName = "harry"
	var newAvailable = []int64{1550487600000}
	user.Name = newName
	user.Dates = newAvailable

	if err := user.Update(); err != nil {
		t.Fatalf("update failed: %s\n", err)
	}

	retUser := User{}
	if err := retUser.Read(user.Id); err != nil {
		t.Fatalf("couldn't read row back from user table: %s\n", err)
	}

	if retUser.Id != user.Id || retUser.Name != newName || retUser.IdMeetUp != user.IdMeetUp || reflect.DeepEqual(retUser.Dates, user.Dates) == false {
		t.Errorf("returned row from DB was different to the one updated. updated: %+v, returned: %+v\n", user, retUser)
	}

	if err := meetUp.Delete(); err != nil {
		t.Fatalf("delete failed: %s\n", err)
	}
}

func TestUser_Delete(t *testing.T) {
	testDbName := CreateTestDb(t)
	defer DestroyTestDb(testDbName)

	var meetUp = MeetUp{UserHash: "abc", AdminHash: "def", AdminEmail: "testy@testy.test", SendAlerts: true, Dates: []int64{1550401200000, 1550487600000, 1550574000000, 1550660400000, 1550746800000, 1550833200000, 1550919600000, 1551006000000}, Description: "meetUp description"}
	if err := meetUp.Create(); err != nil {
		t.Fatalf("create failed: %s\n", err)
	}

	var user = User{IdMeetUp: meetUp.Id, Name: "bob", Dates: []int64{1550401200000, 1550487600000, 1550574000000}}
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

	if err := meetUp.Delete(); err != nil {
		t.Fatalf("delete failed: %s\n", err)
	}
}
