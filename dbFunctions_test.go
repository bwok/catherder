package main

import (
	"testing"
)

/*
	var meetUpObj = MeetUp{
		Id:        1,
		UserHash:  "8d9d7c59eec27a7aee55536582e45afb18f072c282edd22474a0db0676d74299",
		AdminHash: "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
		Dates: Dates{
			{Id: 1, IdMeetUp: 1, Date: 1549537200000, Users: Users{
				{Id: 1, IdDate: 1, Name: "user1", Available: true},
				{Id: 2, IdDate: 1, Name: "user2", Available: false},
				{Id: 3, IdDate: 1, Name: "user3", Available: false},
			}},
			{Id: 2, IdMeetUp: 1, Date: 1549623600000, Users: Users{
				{Id: 4, IdDate: 2, Name: "user1", Available: false},
				{Id: 5, IdDate: 2, Name: "user2", Available: false},
				{Id: 6, IdDate: 2, Name: "user3", Available: false},
			}},
			{Id: 3, IdMeetUp: 1, Date: 1549710000000, Users: Users{
				{Id: 7, IdDate: 3, Name: "user1", Available: true},
				{Id: 8, IdDate: 3, Name: "user2", Available: true},
				{Id: 9, IdDate: 3, Name: "user3", Available: true},
			}},
			{Id: 4, IdMeetUp: 1, Date: 1549796400000, Users: Users{
				{Id: 10, IdDate: 4, Name: "user1", Available: true},
				{Id: 11, IdDate: 4, Name: "user2", Available: false},
				{Id: 12, IdDate: 4, Name: "user3", Available: true},
			}},
			{Id: 5, IdMeetUp: 1, Date: 1549882800000, Users: Users{
				{Id: 3, IdDate: 5, Name: "user1", Available: false},
				{Id: 14, IdDate: 5, Name: "user2", Available: true},
				{Id: 15, IdDate: 5, Name: "user3", Available: true},
			}},
		},
		Admin:       Admin{Id: 1, IdMeetUp: 1, Email: "testy@testy.test", Alerts: true},
		Description: "ljkas;ldfjk;asldkjf",
	}
*/

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
}
