package utils

import (
	"fmt"
	"os"

	"github.com/resend/resend-go/v3"
)

func SendForgotPasswordEmail(email string, token string) error {
	apiKey := os.Getenv("RESEND_API_KEY")
	client := resend.NewClient(apiKey)

	params := &resend.SendEmailRequest{
		From:    os.Getenv("EMAIL_SENDER"),
		To:      []string{email},
		Subject: "Your Forgot Password email",
		Html:    fmt.Sprintf("<p>Please click the link below to reset your password:<br/> http://localhost:4000/api/auth/reset-password?token=%s</p>", token),
	}

	sent, err := client.Emails.Send(params)
	if err != nil {
		fmt.Println("Error sending email", err)
		return err
	}

	fmt.Println("Sent id:", sent.Id)
	return nil
}

func SendVerificationEmail(email string, token string) error {
	apiKey := os.Getenv("RESEND_API_KEY")
	client := resend.NewClient(apiKey)

	params := &resend.SendEmailRequest{
		From:    os.Getenv("EMAIL_SENDER"),
		To:      []string{email},
		Subject: "Your email verification",
		Html:    fmt.Sprintf("<p>Please click the link below to verify your email:<br/> http://localhost:4000/api/auth/verify-email?token=%s</p>", token),
	}

	sent, err := client.Emails.Send(params)
	if err != nil {
		fmt.Println("Error sending email", err)
		return err
	}

	fmt.Println("Sent id:", sent.Id)
	return nil
}
