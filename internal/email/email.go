package email

import (
	"fmt"
	"log"
	"net/smtp"
	"os"
)

func SendVerificationEmail(toEmail, token string) error {
	from := os.Getenv("EMAIL_FROM")
	username := os.Getenv("EMAIL_USERNAME")
	password := os.Getenv("EMAIL_PASSWORD")
	smtpHost := os.Getenv("EMAIL_HOST")
	smtpPort := os.Getenv("EMAIL_PORT")

	if from == "" || username == "" || password == "" || smtpHost == "" || smtpPort == "" {
		return fmt.Errorf("email configuration is missing")
	}

	verificationLink := fmt.Sprintf("%s?token=%s", os.Getenv("EMAIL_VERIFICATION_URL"), token)
	message := []byte(fmt.Sprintf("Subject: Email Verification\n\nPlease verify your email by clicking the link: %s", verificationLink))

	auth := smtp.PlainAuth("", username, password, smtpHost)
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{toEmail}, message)
	if err != nil {
		log.Printf("Failed to send email: %v", err)
		return fmt.Errorf("failed to send email: %w", err)
	}

	log.Println("Verification email sent successfully")
	return nil
}
