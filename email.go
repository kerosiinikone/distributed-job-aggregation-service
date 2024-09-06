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
	d := NewDialer()
	m := NewMessage(*postings, job.EmailAddr)

	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}

func NewDialer() gomail.Dialer {
	return *gomail.NewDialer(
		"smtp.gmail.com", 
		587, 
		os.Getenv("EMAIL_ADDRESS"), 
		os.Getenv("EMAIL_PASSWORD"),
	)
}

func NewMessage(postings JobResults, addr string) *gomail.Message {
	postLinks := make([]string, len(postings.Results))
	for i, p := range postings.Results {
		postLinks[i] = p.Link
	}

	m := gomail.NewMessage()
	m.SetHeader("From", os.Getenv("EMAIL_ADDRESS"))
	m.SetHeader("To", addr)
	m.SetHeader("Subject", "Your job postings")
	m.SetBody("text/html", strings.Join(postLinks[:], "\n\n"))

	return m
}