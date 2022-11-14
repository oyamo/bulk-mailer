package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"strings"
)

type Recipient struct {
	Name  string
	Email string
}

func (r *Recipient) String() string {
	return r.Email
}

func main() {
	var emailFile string
	var htmlFile string
	var envFile string
	var subject string
	var emailColumn int
	var nameColumd int
	var recipients []Recipient
	var htmlTemplate string

	// Parse command line arguments
	flag.StringVar(&emailFile, "csv", "", "CSV file containing email addresses")
	flag.StringVar(&htmlFile, "body", "", "HTML file")
	flag.StringVar(&subject, "subject", "", "Email subject")
	flag.StringVar(&envFile, "env", ".env", "EMAIL PASSWORD Environment file")

	flag.Parse()

	if emailFile == "" || htmlFile == "" || envFile == "" || subject == "" {
		flag.Usage()
		os.Exit(1)
	}

	// Read email file
	email, err := os.Open(emailFile)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error: %v (email file) \n", err)
		os.Exit(1)
	}

	// Read HTML file
	html, err := os.Open(htmlFile)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error: %v (html file) \n", err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(html)
	for scanner.Scan() {
		htmlTemplate = strings.Join([]string{htmlTemplate, scanner.Text()}, "")
	}

	// Read environment file
	err = godotenv.Load(envFile)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error: %v (env file) \n", err)
		os.Exit(1)
	}

	// Read files to strings
	emails, err := csv.NewReader(email).ReadAll()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error: %v (email file) \n", err)
		os.Exit(1)
	}

	if len(emails) != 2 {
		_, _ = fmt.Fprintf(os.Stderr, "Error: (email file) Must contain 2 columns \n")
		os.Exit(1)
	}

	// Test if first column contains email addresses
	// Loop from fileRder[0][:4] to test only first 4 rows
	for _, email := range emails[0][:4] {
		if strings.Contains(email, "@") {
			emailColumn = 0
			nameColumd = 1
			break
		}
	}

	// Test if second column contains email addresses
	// Loop from fileRder[1][:4] to test only first 4 rows
	for _, email := range emails[1][:4] {
		if strings.Contains(email, "@") {
			emailColumn = 1
			nameColumd = 0
			break
		}
	}

	// Create recipients
	for index, _ := range emails[0] {
		if index == 0 {
			continue
		}
		recipients = append(recipients, Recipient{
			Name:  emails[nameColumd][index],
			Email: emails[emailColumn][index],
		})
	}

	// Create mailer
	mailer := NewMailer(os.Getenv("EMAIL"), os.Getenv("PASSWORD"), "smtp.gmail.com")

	// Send emails
	err = mailer.SendMail(recipients, subject, htmlTemplate)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error: %v (email file) \n", err)
		os.Exit(1)
	}

	_, _ = fmt.Fprintf(os.Stdout, "Emails sent successfully \n")
}
