package main

import (
	"fmt"
	"log"
	"net/smtp"
)


// Sends an email to the given email address on meetup creation.
func sendCreationEmail(meetUp MeetUp, serverHostname string) {
	subject := "A meetup was created."

	// TODO URL encode
	messageBody := "A meetup has been created.\r\n\r\n" +
		"Use this link to edit the meetup: https://" + serverHostname + "/edit?id=" + meetUp.AdminHash + "\r\n\r\n" +
		"Share this link with the meetup participants: https://" + serverHostname + "/view?id=" + meetUp.UserHash + "\r\n"

	go sendMail(meetUp.AdminEmail, subject, messageBody)
}

// Sends an email to the given email address when a user is added or updated.
func sendUserChangedEmail(user User, emailAddress, userHash, serverHostname string, userAdded bool) {
	var subject string

	if userAdded == true {
		subject = fmt.Sprintf("New user %q was added to the meetup.", user.Name)
	} else {
		subject = fmt.Sprintf("User %q changed their meetup availability.", user.Name)
	}

	// TODO URL encode
	messageBody := fmt.Sprintf("View the meetup here: https://%s/view?id=%s\r\n", serverHostname, userHash)

	go sendMail(emailAddress, subject, messageBody)
}

// Given an address, subject and mail body, sends an email to the address.
func sendMail(emailAddress, subject, messageBody string) {
	// TODO supply these values via a config file
	username := ""
	password := ""
	host := ""
	port := "587"

	from := username
	to := []string{emailAddress}
	msg := []byte("To: " + emailAddress + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		messageBody + "\r\n")

	auth := smtp.PlainAuth("", username, password, host)
	err := smtp.SendMail(host + ":" + port, auth, from, to, msg)
	if err != nil {
		log.Printf("sending notification email failed (%s)", err)
	}
}
