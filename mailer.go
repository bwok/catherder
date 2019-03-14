package main

import (
	"fmt"
	"log"
	"net/smtp"
	"net/url"
)

// Sends an email to the given email address on meetup creation.
func sendCreationEmail(meetUp MeetUp, serverHostname string) {
	subject := "A meetup was created."
	messageBody := "A meetup has been created.\r\n\r\n" +
		"Use this link to edit the meetup: https://" + serverHostname + "/edit?id=" + url.QueryEscape(meetUp.AdminHash) + "\r\n\r\n" +
		"Share this link with the meetup participants: https://" + serverHostname + "/view?id=" + url.QueryEscape(meetUp.UserHash) + "\r\n"

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

	messageBody := fmt.Sprintf("View the meetup here: https://%s/view?id=%s\r\n", serverHostname, url.QueryEscape(userHash))

	go sendMail(emailAddress, subject, messageBody)
}

// Given an address, subject and mail body, sends an email to the address.
func sendMail(emailAddress, subject, messageBody string) {
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
	err := smtp.SendMail(host+":"+port, auth, from, to, msg)
	if err != nil {
		log.Printf("sending notification email failed (%s)", err)
	}
}
