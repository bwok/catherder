package main

import (
	"testing"
)

/*
Helper functions to compare the properties of various database objects
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
		t.Fatalf("CreateMeetUp() failed %s\n", err)
	}

	var retMeetUp = MeetUp{}
	if err := retMeetUp.GetByUserHash(meetUpObj.UserHash); err != nil {
		t.Fatalf("GetByUserHash() failed %s\n", err)
	}

	if compareMeetUpObjects(meetUpObj, retMeetUp) == false {
		t.Fatalf("MeetUp objects were different.")
	}

	retMeetUp = MeetUp{}
	if err := retMeetUp.GetByAdminHash(meetUpObj.AdminHash); err != nil {
		t.Fatalf("GetByUserHash() failed %s\n", err)
	}

	if compareMeetUpObjects(meetUpObj, retMeetUp) == false {
		t.Fatalf("MeetUp objects were different.")
	}
}
