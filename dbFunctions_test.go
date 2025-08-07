package main

import (
	"encoding/json"
	"reflect"
	"testing"
)

import (
	"errors"
)

func TestMeetUp_DeleteByAdminHash(t *testing.T) {
	testDbName := CreateTestDb(t)
	defer DestroyTestDb(testDbName)

	var meetUpObj = MeetUp{
		UserHash:    "8d9d7c59eec27a7aee55536582e45afb18f072c282edd22474a0db0676d74299",
		AdminHash:   "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
		Dates:       []int64{1550401200000, 1550487600000, 1550574000000, 1550660400000, 1550746800000, 1550833200000, 1550919600000, 1551006000000},
		Users:       Users{},
		Description: "ljkas;ldfjk;asldkjf",
	}

	if err := meetUpObj.Create(); err != nil {
		t.Fatalf("MeetUp.Create() failed: %s\n", err)
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
		UserHash:  "6cf51863dcbd352c9da7fc0670a34a7173056413214ac0e22f0effd1015a7fa907399f9107fd48689d3800043ff12bbd33a24a433a6ede783bda3423c9820278",
		AdminHash: "a39f823a49a4fbdfc2906a4baf0dce97a216d8b5bf6b0ab83a31d04c2d84ae619a6e017368a434ecb7b09b54015d22455062ac199ec48aa5b1c0dea830c3ecb6",
		Dates:     []int64{1550401200000, 1550487600000, 1550574000000, 1550660400000, 1550746800000, 1550833200000, 1550919600000, 1551006000000},
		Users: Users{
			{Name: "user1", Dates: []int64{1550401200000, 1550487600000, 1550574000000}},
			{Name: "user2", Dates: []int64{1550401200000, 1550574000000}},
			{Name: "user3", Dates: []int64{1550574000000}},
		},
		Description: "ljkas;ldfjk;asldkjf",
	}

	if err := meetUpObj.Create(); err != nil {
		t.Fatalf("MeetUp.Create() failed: %s\n", err)
	}
	for _, user := range meetUpObj.Users {
		user.IdMeetUp = meetUpObj.Id
		if err := user.Create(); err != nil {
			t.Fatalf("User.Create() failed: %s\n", err)
		}
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

	if err := meetUpObj.Delete(); err != nil {
		t.Fatalf("delete failed: %s\n", err)
	}
}

func TestMeetUp_GetByAdminHash(t *testing.T) {
	testDbName := CreateTestDb(t)
	defer DestroyTestDb(testDbName)

	var meetUpObj = MeetUp{
		UserHash:  "6cf51863dcbd352c9da7fc0670a34a7173056413214ac0e22f0effd1015a7fa907399f9107fd48689d3800043ff12bbd33a24a433a6ede783bda3423c9820278",
		AdminHash: "a39f823a49a4fbdfc2906a4baf0dce97a216d8b5bf6b0ab83a31d04c2d84ae619a6e017368a434ecb7b09b54015d22455062ac199ec48aa5b1c0dea830c3ecb6",
		Dates:     []int64{1550401200000, 1550487600000, 1550574000000, 1550660400000, 1550746800000, 1550833200000, 1550919600000, 1551006000000},
		Users: Users{
			{Name: "user1", Dates: []int64{1550401200000, 1550487600000, 1550574000000}},
			{Name: "user2", Dates: []int64{1550401200000, 1550574000000}},
			{Name: "user3", Dates: []int64{1550574000000}},
		},
		Description: "ljkas;ldfjk;asldkjf",
	}

	if err := meetUpObj.Create(); err != nil {
		t.Fatalf("Create() failed: %s\n", err)
	}
	for _, user := range meetUpObj.Users {
		user.IdMeetUp = meetUpObj.Id
		if err := user.Create(); err != nil {
			t.Fatalf("User.Create() failed: %s\n", err)
		}
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

	if err := meetUpObj.Delete(); err != nil {
		t.Fatalf("delete failed: %s\n", err)
	}
}

func TestMeetUp_MarshalJSON(t *testing.T) {
	var meetUpObj = MeetUp{
		UserHash:    "6cf51863dcbd352c9da7fc0670a34a7173056413214ac0e22f0effd1015a7fa907399f9107fd48689d3800043ff12bbd33a24a433a6ede783bda3423c9820278",
		AdminHash:   "a39f823a49a4fbdfc2906a4baf0dce97a216d8b5bf6b0ab83a31d04c2d84ae619a6e017368a434ecb7b09b54015d22455062ac199ec48aa5b1c0dea830c3ecb6",
		Dates:       []int64{1550401200000, 1550487600000, 1550574000000, 1550660400000, 1550746800000, 1550833200000, 1550919600000, 1551006000000},
		Users:       Users{},
		Description: "ljkas;ldfjk;asldkjf",
	}

	if _, err := json.MarshalIndent(&meetUpObj, "", "	"); err != nil {
		t.Errorf("MeetUp couldn't be marshalled. error: %s", err)
	} else {
		//fmt.Printf("%s\n", js)
	}

}

// Helper functions to compare some of the properties of various database objects.
// The IDs don't get compared as one of the passed objects usually doesn't have any
func compareMeetUpObjects(obj1, obj2 MeetUp) bool {
	if obj1.UserHash != obj2.UserHash || obj1.AdminHash != obj2.AdminHash || reflect.DeepEqual(obj1.Dates, obj2.Dates) == false || obj1.Description != obj2.Description {
		return false
	}
	if compareUsersObject(obj1.Users, obj2.Users) == false {
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
	if obj1.Name != obj2.Name || reflect.DeepEqual(obj1.Dates, obj2.Dates) == false {
		return false
	}
	return true
}

// Helper functions to validate that a meetup object and its children have realistic IDs, and that the IDs match.
func validateMeetUpObject(meetUp MeetUp) error {
	if meetUp.Id <= 0 {
		return errors.New("invalid meetup id")
	} else if validateHash(meetUp.UserHash) != nil {
		return errors.New("invalid user hash")
	} else if validateHash(meetUp.AdminHash) != nil {
		return errors.New("invalid admin hash")
	} else if err := validateUsersObject(meetUp.Id, meetUp.Users); err != nil {
		return err
	} else {
		return nil
	}
}
func validateUsersObject(meetUpId int64, users Users) error {
	for i := 0; i < len(users); i++ {
		if err := validateUserObject(meetUpId, users[i]); err != nil {
			return err
		}
	}
	return nil
}
func validateUserObject(meetUpId int64, user User) error {
	if user.Id <= 0 {
		return errors.New("invalid User id")
	} else if user.IdMeetUp != meetUpId {
		return errors.New("meetup id and user_idmeetup do not match")
	} else if user.Name == "" {
		return errors.New("user name is an empty string")
	} else {
		return nil
	}
}
