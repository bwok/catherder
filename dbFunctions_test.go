package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/mail"
	"testing"
)

// Tests most of the other custom sql functions in dbFunctions as a side effect.
func TestMeetUp_CreateMeetUp(t *testing.T) {
	testDbName := CreateTestDb(t)
	defer DestroyTestDb(testDbName)

	var meetUpObj = MeetUp{
		UserHash:  "8d9d7c59eec27a7aee55536582e45afb18f072c282edd22474a0db0676d74299",
		AdminHash: "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
		Dates: Dates{
			{Date: 1549537200000, Users: Users{
				{Name: "user1", Available: true},
				{Name: "user2", Available: false},
				{Name: "user3", Available: false},
			}},
			{Date: 1549623600000, Users: Users{
				{Name: "user1", Available: false},
				{Name: "user2", Available: false},
				{Name: "user3", Available: false},
			}},
			{Date: 1549710000000, Users: Users{
				{Name: "user1", Available: true},
				{Name: "user2", Available: true},
				{Name: "user3", Available: true},
			}},
			{Date: 1549796400000, Users: Users{
				{Name: "user1", Available: true},
				{Name: "user2", Available: false},
				{Name: "user3", Available: true},
			}},
			{Date: 1549882800000, Users: Users{
				{Name: "user1", Available: false},
				{Name: "user2", Available: true},
				{Name: "user3", Available: true},
			}},
		},
		Admin:       Admin{Email: "testy@testy.test", Alerts: true},
		Description: "ljkas;ldfjk;asldkjf",
	}

	if err := meetUpObj.CreateMeetUp(); err != nil {
		t.Fatalf("CreateMeetUp() failed: %s\n", err)
	}

	var retMeetUp = MeetUp{}
	if err := retMeetUp.GetByUserHash(meetUpObj.UserHash); err != nil {
		t.Fatalf("GetByUserHash() failed: %s\n", err)
	}

	if compareMeetUpObjects(meetUpObj, retMeetUp) == false {
		t.Fatalf("MeetUp objects were different.")
	}

	retMeetUp = MeetUp{}
	if err := retMeetUp.GetByAdminHash(meetUpObj.AdminHash); err != nil {
		t.Fatalf("GetByUserHash() failed: %s\n", err)
	}

	if compareMeetUpObjects(meetUpObj, retMeetUp) == false {
		t.Fatalf("MeetUp objects were different.")
	}
}

func TestMeetUp_UpdateMeetUpDeleteDates(t *testing.T) {
	testDbName := CreateTestDb(t)
	defer DestroyTestDb(testDbName)
	var dbMeetUpObj MeetUp

	var meetUpObj = MeetUp{
		UserHash:  "8d9d7c59eec27a7aee55536582e45afb18f072c282edd22474a0db0676d74299",
		AdminHash: "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
		Dates: Dates{
			{Date: 1549537200000, Users: Users{
				{Name: "user1", Available: true},
				{Name: "user2", Available: false},
				{Name: "user3", Available: false},
			}},
			{Date: 1549623600000, Users: Users{
				{Name: "user1", Available: false},
				{Name: "user2", Available: false},
				{Name: "user3", Available: false},
			}},
			{Date: 1549710000000, Users: Users{
				{Name: "user1", Available: true},
				{Name: "user2", Available: true},
				{Name: "user3", Available: true},
			}},
			{Date: 1549796400000, Users: Users{
				{Name: "user1", Available: true},
				{Name: "user2", Available: false},
				{Name: "user3", Available: true},
			}},
			{Date: 1549882800000, Users: Users{
				{Name: "user1", Available: false},
				{Name: "user2", Available: true},
				{Name: "user3", Available: true},
			}},
		},
		Admin:       Admin{Email: "testy@testy.test", Alerts: true},
		Description: "ljkas;ldfjk;asldkjf",
	}

	// Save meetup object and children to database, retrieve again aftward so we have all required IDs
	if err := meetUpObj.CreateMeetUp(); err != nil {
		t.Fatalf("CreateMeetUp() failed: %s\n", err)
	}
	if err := dbMeetUpObj.GetByUserHash(meetUpObj.UserHash); err != nil {
		t.Fatalf("GetByUserHash() failed: %s\n", err)
	}

	// Modify the meetup object
	var newMeetUpObj = MeetUp{
		UserHash:  meetUpObj.UserHash,
		AdminHash: meetUpObj.AdminHash,
		Dates: Dates{
			{Date: 1549537200000, Users: Users{
				{Name: "user1", Available: true},
				{Name: "user2", Available: false},
				{Name: "user3", Available: false},
			}},
		},
		Admin:       Admin{Email: "klm@nop.qrs", Alerts: false},
		Description: "hij",
	}

	// Update the database, and compare the updated receiver object the new meetup object
	if err := dbMeetUpObj.UpdateMeetUpDeleteDates(&newMeetUpObj); err != nil {
		t.Fatalf("GetByUserHash() failed: %s\n", err)
	}
	if compareMeetUpObjects(dbMeetUpObj, newMeetUpObj) == false {
		t.Errorf("MeetUp objects were different.\nupdated: %+v\nnew: %+v\n", meetUpObj, newMeetUpObj)
	}

	// Compare the database version with the updated receiver object.
	// Get saved meetup object and children with IDs etc.
	dbMeetUpObj = MeetUp{}
	if err := dbMeetUpObj.GetByUserHash(meetUpObj.UserHash); err != nil {
		t.Fatalf("GetByUserHash() failed: %s\n", err)
	}
	if compareMeetUpObjects(dbMeetUpObj, newMeetUpObj) == false {
		t.Errorf("MeetUp objects were different.\ndatabase: %+v\nnew: %+v\n", dbMeetUpObj, newMeetUpObj)
	}
}

func TestMeetUp_DeleteByAdminHash(t *testing.T) {
	testDbName := CreateTestDb(t)
	defer DestroyTestDb(testDbName)

	var meetUpObj = MeetUp{
		UserHash:  "8d9d7c59eec27a7aee55536582e45afb18f072c282edd22474a0db0676d74299",
		AdminHash: "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
		Dates: Dates{
			{Date: 1549537200000, Users: Users{
				{Name: "user1", Available: true},
				{Name: "user2", Available: false},
				{Name: "user3", Available: false},
			}},
			{Date: 1549623600000, Users: Users{
				{Name: "user1", Available: false},
				{Name: "user2", Available: false},
				{Name: "user3", Available: false},
			}},
			{Date: 1549710000000, Users: Users{
				{Name: "user1", Available: true},
				{Name: "user2", Available: true},
				{Name: "user3", Available: true},
			}},
			{Date: 1549796400000, Users: Users{
				{Name: "user1", Available: true},
				{Name: "user2", Available: false},
				{Name: "user3", Available: true},
			}},
			{Date: 1549882800000, Users: Users{
				{Name: "user1", Available: false},
				{Name: "user2", Available: true},
				{Name: "user3", Available: true},
			}},
		},
		Admin:       Admin{Email: "testy@testy.test", Alerts: true},
		Description: "ljkas;ldfjk;asldkjf",
	}

	if err := meetUpObj.CreateMeetUp(); err != nil {
		t.Fatalf("CreateMeetUp() failed: %s\n", err)
	}

	if err := meetUpObj.DeleteByAdminHash(meetUpObj.AdminHash); err != nil {
		t.Fatalf("DeleteByAdminHash() failed: %s\n", err)
		return
	}

	if err := meetUpObj.GetByUserHash(meetUpObj.UserHash); err.Error() != "no rows matching the userhash" {
		t.Fatalf("GetByUserHash() failed: %s\n", err)
	}
}

func TestMeetUp_GetByUserHash(t *testing.T) {
	testDbName := CreateTestDb(t)
	defer DestroyTestDb(testDbName)

	var meetUpObj = MeetUp{
		UserHash:  "8d9d7c59eec27a7aee55536582e45afb18f072c282edd22474a0db0676d74299",
		AdminHash: "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
		Dates: Dates{
			{Date: 1549537200000, Users: Users{
				{Name: "user1", Available: true},
				{Name: "user2", Available: false},
				{Name: "user3", Available: false},
			}},
			{Date: 1549623600000, Users: Users{
				{Name: "user1", Available: false},
				{Name: "user2", Available: false},
				{Name: "user3", Available: false},
			}},
			{Date: 1549710000000, Users: Users{
				{Name: "user1", Available: true},
				{Name: "user2", Available: true},
				{Name: "user3", Available: true},
			}},
			{Date: 1549796400000, Users: Users{
				{Name: "user1", Available: true},
				{Name: "user2", Available: false},
				{Name: "user3", Available: true},
			}},
			{Date: 1549882800000, Users: Users{
				{Name: "user1", Available: false},
				{Name: "user2", Available: true},
				{Name: "user3", Available: true},
			}},
		},
		Admin:       Admin{Email: "testy@testy.test", Alerts: true},
		Description: "ljkas;ldfjk;asldkjf",
	}

	if err := meetUpObj.CreateMeetUp(); err != nil {
		t.Fatalf("CreateMeetUp() failed: %s\n", err)
	}

	var dbMeetUp = MeetUp{}
	if err := dbMeetUp.GetByUserHash(meetUpObj.UserHash); err != nil && err.Error() != "no rows matching the userhash" {
		t.Errorf("GetByUserHash() failed: %s\n", err)
	}

	if err := validateMeetUpObject(dbMeetUp); err != nil {
		t.Errorf("MeetUp validation failed: %s\n", err)
	}

	if compareMeetUpObjects(meetUpObj, dbMeetUp) == false {
		t.Errorf("MeetUp objects were different.")
	}
}

func TestMeetUp_GetByAdminHash(t *testing.T) {
	testDbName := CreateTestDb(t)
	defer DestroyTestDb(testDbName)

	var meetUpObj = MeetUp{
		UserHash:  "8d9d7c59eec27a7aee55536582e45afb18f072c282edd22474a0db0676d74299",
		AdminHash: "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
		Dates: Dates{
			{Date: 1549537200000, Users: Users{
				{Name: "user1", Available: true},
				{Name: "user2", Available: false},
				{Name: "user3", Available: false},
			}},
			{Date: 1549623600000, Users: Users{
				{Name: "user1", Available: false},
				{Name: "user2", Available: false},
				{Name: "user3", Available: false},
			}},
			{Date: 1549710000000, Users: Users{
				{Name: "user1", Available: true},
				{Name: "user2", Available: true},
				{Name: "user3", Available: true},
			}},
			{Date: 1549796400000, Users: Users{
				{Name: "user1", Available: true},
				{Name: "user2", Available: false},
				{Name: "user3", Available: true},
			}},
			{Date: 1549882800000, Users: Users{
				{Name: "user1", Available: false},
				{Name: "user2", Available: true},
				{Name: "user3", Available: true},
			}},
		},
		Admin:       Admin{Email: "testy@testy.test", Alerts: true},
		Description: "ljkas;ldfjk;asldkjf",
	}

	if err := meetUpObj.CreateMeetUp(); err != nil {
		t.Fatalf("CreateMeetUp() failed: %s\n", err)
	}

	var dbMeetUp = MeetUp{}
	if err := dbMeetUp.GetByAdminHash(meetUpObj.AdminHash); err != nil && err.Error() != "no rows matching the adminhash" {
		t.Errorf("GetByAdminHash() failed: %s\n", err)
	}

	if err := validateMeetUpObject(dbMeetUp); err != nil {
		t.Errorf("MeetUp validation failed: %s\n", err)
	}

	if compareMeetUpObjects(meetUpObj, dbMeetUp) == false {
		t.Errorf("MeetUp objects were different.")
	}
}

func TestAdmin_GetByMeetUpId(t *testing.T) {
	// TODO test
}

func TestDates_GetAllByMeetUpId(t *testing.T) {
	// TODO test
}

func TestUsers_GetAllByDateId(t *testing.T) {
	// TODO test
}

func TestUsers_CreateUsers(t *testing.T) {
	// TODO test
}


func TestMeetUp_MarshalJSON(t *testing.T) {
	var meetUpObj = MeetUp{
		UserHash:  "8d9d7c59eec27a7aee55536582e45afb18f072c282edd22474a0db0676d74299",
		AdminHash: "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
		Dates: Dates{
			{Date: 1549537200000, Users: Users{
				{Name: "user1", Available: true},
			}},
		},
		Admin:       Admin{Email: "testy@testy.test", Alerts: true},
		Description: "ljkas;ldfjk;asldkjf",
	}

	if js, err := json.MarshalIndent(&meetUpObj, "", "	"); err != nil {
		t.Errorf("MeetUp couldn't be marshalled. error: %s", err)
	} else {
		fmt.Printf("%s\n", js)
	}

}

/*
Helper functions to compare some of the properties of various database objects.
The IDs don't get compared as one of the passed objects usually doesn't have any
*/
func compareMeetUpObjects(obj1, obj2 MeetUp) bool {
	if obj1.UserHash != obj2.UserHash || obj1.AdminHash != obj2.AdminHash || obj1.Description != obj2.Description {
		return false
	}
	if compareAdminObjects(obj1.Admin, obj2.Admin) == false || compareDatesObject(obj1.Dates, obj2.Dates) == false {
		return false
	}
	return true
}
func compareAdminObjects(obj1, obj2 Admin) bool {
	if obj1.Email != obj2.Email || obj1.Alerts != obj2.Alerts {
		return false
	}
	return true
}
func compareDatesObject(obj1, obj2 Dates) bool {
	if len(obj1) != len(obj2) {
		return false
	}
	for i := 0; i < len(obj1); i++ {
		if compareDateObjects(obj1[i], obj2[i]) == false {
			return false
		}
	}
	return true
}
func compareDateObjects(obj1, obj2 Date) bool {
	if obj1.Date != obj2.Date || compareUsersObject(obj1.Users, obj2.Users) == false {
		return false
	}
	return true
}
func compareUsersObject(obj1, obj2 Users) bool {
	if len(obj1) != len(obj2) {
		return false
	}
	for i := 0; i < len(obj1); i++ {
		if compareUserObjects(obj1[i], obj2[i]) == false {
			return false
		}
	}
	return true
}
func compareUserObjects(obj1, obj2 User) bool {
	if obj1.Name != obj2.Name || obj1.Available != obj2.Available {
		return false
	}
	return true
}

/*
Helper functions to validate that a meetup object and its children have realistic IDs, and that the IDs match.
*/
func validateMeetUpObject(meetUp MeetUp) error {
	if meetUp.Id <= 0 {
		return errors.New("invalid meetup id")
	} else if validateHash(meetUp.UserHash) != nil {
		return errors.New("invalid user hash")
	} else if validateHash(meetUp.AdminHash) != nil {
		return errors.New("invalid admin hash")
	} else if err := validateAdminObject(meetUp.Id, meetUp.Admin); err != nil {
		return err
	} else if err := validateDatesObject(meetUp.Id, meetUp.Dates); err != nil {
		return err
	} else {
		return nil
	}
}
func validateAdminObject(meetUpId int64, admin Admin) error {
	// Admin.Alerts not checked because it can only be true/false anyway
	if admin.Id <= 0 {
		return errors.New("invalid Admin id")
	} else if admin.IdMeetUp != meetUpId {
		return errors.New("meetup id and Admin.IdMeetUp do not match")
	} else if _, err := mail.ParseAddress(admin.Email); err != nil {
		return fmt.Errorf("invalid email address %q. error: %s", admin.Email, err)
	} else {
		return nil
	}
}
func validateDatesObject(meetUpId int64, dates Dates) error {
	for i := 0; i < len(dates); i++ {
		if err := validateDateObject(meetUpId, dates[i]); err != nil {
			return err
		}
	}
	return nil
}
func validateDateObject(meetUpId int64, date Date) error {
	if date.Id <= 0 {
		return errors.New("invalid Date id")
	} else if date.IdMeetUp != meetUpId {
		return errors.New("meetup id and Date id do not match")
	} else if date.Date <= 0 {
		return errors.New("date timestamp <= 0")
	} else if err := validateUsersObject(date.Id, date.Users); err != nil {
		return err
	} else {
		return nil
	}
}
func validateUsersObject(dateId int64, users Users) error {
	for i := 0; i < len(users); i++ {
		if err := validateUserObject(dateId, users[i]); err != nil {
			return err
		}
	}
	return nil
}
func validateUserObject(dateId int64, user User) error {
	// User.Available not checked because it can only be true/false anyway
	if user.Id <= 0 {
		return errors.New("invalid User id")
	} else if user.IdDate != dateId {
		return errors.New("date id and user id do not match")
	} else if user.Name == "" {
		return errors.New("user name is an empty string")
	} else {
		return nil
	}
}
