package utils

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/smtp"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"

	"department-management/config"
	"department-management/models"
)

var cfg, _ = config.LoadConfig()

// Hash a given password using bcrypt
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// Compare a hashed password with a plain text password
func ComparePasswords(hashedPassword, plainPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
}

// Send a confirmation email with a given token
func SendConfirmationEmail(email, link string) error {
	from := cfg.EmailFrom
	password := cfg.EmailPassword
	to := email
	smtpHost := cfg.SMTPHost
	smtpPort := cfg.SMTPPort

	subject := "Confirm your account"
	body := fmt.Sprintf("Click the link to confirm your account: %s", link)

	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: " + subject + "\n\n" +
		body

	// Create a custom TLS config to skip verification
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true, // Disable certificate verification (dangerous for production)
		ServerName:         smtpHost,
	}

	// Dial the connection
	conn, err := tls.Dial("tcp", smtpHost+":"+smtpPort, tlsConfig)
	if err != nil {
		log.Printf("Failed to dial SMTP server: %v", err)
		return err
	}
	defer conn.Close()

	// Send the email
	c, err := smtp.NewClient(conn, smtpHost)
	if err != nil {
		log.Printf("Failed to create SMTP client: %v", err)
		return err
	}

	// Authenticate
	auth := smtp.PlainAuth("", from, password, smtpHost)
	err = c.Auth(auth)
	if err != nil {
		log.Printf("Failed to authenticate SMTP: %v", err)
		return err
	}

	// Set the sender and recipient
	err = c.Mail(from)
	if err != nil {
		log.Printf("Failed to set sender: %v", err)
		return err
	}

	err = c.Rcpt(to)
	if err != nil {
		log.Printf("Failed to set recipient: %v", err)
		return err
	}

	// Send the message
	wc, err := c.Data()
	if err != nil {
		log.Printf("Failed to get data writer: %v", err)
		return err
	}

	_, err = wc.Write([]byte(msg))
	if err != nil {
		log.Printf("Failed to send message: %v", err)
		return err
	}

	err = wc.Close()
	if err != nil {
		log.Printf("Failed to close writer: %v", err)
		return err
	}

	c.Quit()

	log.Printf("Confirmation email sent to: %s", to)
	return nil
}

// Retrieve the JWT secret from the configuration
func GetJWTSecret() string {
	return cfg.JWTSecret
}

// Generate a JWT token for a given user
func GenerateJWT(userID, email, role string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &models.Claims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.JWTSecret))
}
