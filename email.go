package main

import (
	"os"
	"strings"

	"gopkg.in/gomail.v2"
)

type Mailer interface {
	SendEmail(JobResults, JobRequest) error
}

type EmailService struct{}

func (e *EmailService) SendEmail(postings *JobResults, job *JobRequest) error {
	d := newDialer()
	m := newMessage(*postings, job.EmailAddr)

	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}

func newDialer() gomail.Dialer {
	return *gomail.NewDialer(
		"smtp.gmail.com", 
		587, 
		os.Getenv("EMAIL_ADDRESS"), 
		os.Getenv("EMAIL_PASSWORD"),
	)
}

func newMessage(postings JobResults, addr string) *gomail.Message {
	postLinks := make([]string, len(postings))
	for i, p := range postings {
		postLinks[i] = p.Link
	}

	m := gomail.NewMessage()
	m.SetHeader("From", os.Getenv("EMAIL_ADDRESS"))
	m.SetHeader("To", addr)
	m.SetHeader("Subject", "Your job postings")
	m.SetBody("text/html", strings.Join(postLinks[:], "\n\n"))

	return m
}