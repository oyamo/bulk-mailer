package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"io"
	"os"
	"strings"
	"time"
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
	var emailColumn int = -1
	var nameColumn int = -1
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

	b, err := io.ReadAll(html)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error: %v (html file) \n", err)
		os.Exit(1)
	}

	htmlTemplate = string(b)
	fmt.Println(htmlTemplate)

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

	if len(emails[0]) != 2 {
		_, _ = fmt.Fprintf(os.Stderr, "Error: (email file) Must contain 2 columns \n")
		os.Exit(1)
	}

	// Column
	if strings.Contains(strings.ToLower(emails[0][0]), "email") {
		emailColumn = 0
	} else if strings.Contains(strings.ToLower(emails[0][1]), "email") {
		emailColumn = 1
	}

	if strings.Contains(strings.ToLower(emails[0][0]), "name") {
		nameColumn = 0
	} else if strings.Contains(strings.ToLower(emails[0][1]), "name") {
		nameColumn = 1
	}

	if emailColumn == -1 {
		_, _ = fmt.Fprintf(os.Stderr, "error: (email file) Must contain a column with the word 'email' in it\n")
		os.Exit(1)
	}

	if nameColumn == -1 {
		_, _ = fmt.Fprintf(os.Stderr, "Error: (email file) Must contain a column with the word 'name' in it\n")
		os.Exit(1)
	}

	// Create recipients
	for index, _ := range emails[1:] {
		recipients = append(recipients, Recipient{
			Name:  emails[index+1][nameColumn],
			Email: emails[index+1][emailColumn],
		})
	}

	fmt.Println("[bulk-mailer@v1.0.0] https://github.com/oyamo/bulk-mailer")

	// Create mailer
	mailer := NewMailer(os.Getenv("EMAIL"), os.Getenv("PASSWORD"), "smtp.gmail.com")

	start := time.Now()

	// Send emails
	err = mailer.SendMail(recipients, subject, htmlTemplate)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	_, _ = fmt.Fprintf(os.Stdout, "%d emails of %d bytes sent successfully in %d Milliseconds\n", len(recipients), len(htmlTemplate), time.Since(start).Milliseconds())
}
