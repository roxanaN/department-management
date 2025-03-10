package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"regexp"

	"github.com/golang-jwt/jwt/v4"

	"department-management/db"
	"department-management/models"
	"department-management/utils"
)

// Authentication route
func Login(w http.ResponseWriter, r *http.Request) {
	// Extract email and password from body
	var user models.User
	json.NewDecoder(r.Body).Decode(&user)
	log.Printf("Login attempt for email: %s", user.Email)

	// Check if the user exists
	storedUser, err := db.GetUserByEmail(user.Email)
	if err != nil {
		log.Printf("Error fetching user with email %s: %v", user.Email, err)
		http.Error(w, "Invalid email or password. Please check your credentials and try again.", http.StatusUnauthorized)
		return
	}

	// Check if the account is activated
	if !storedUser.Activated {
		log.Printf("Register first: %s", storedUser.Email)
		http.Error(w, "Password not set. Please register first!", http.StatusUnauthorized)
		return
	}

	// Check if the entered password is correct
	if err := utils.ComparePasswords(storedUser.Password, user.Password); err != nil {
		log.Printf("Invalid password for user: %s. Error: %v", storedUser.Email, err)
		http.Error(w, "Invalid email or password. Please check your credentials and try again.", http.StatusUnauthorized)
		return
	}

	// Generate JWT token
	tokenString, err := utils.GenerateJWT(storedUser.ID, storedUser.Email, storedUser.Role)
	if err != nil {
		log.Printf("Error generating token for user: %s", storedUser.Email)
		http.Error(w, "Error generating token. Please try again later.", http.StatusInternalServerError)
		return
	}
	log.Printf("Token generated for user: %s. Token: %s", storedUser.Email, tokenString)

	// Add the token to the response
	response := map[string]string{
		"token": tokenString,
	}
	json.NewEncoder(w).Encode(response)
	log.Printf("Login successful for user: %s", storedUser.Email)
}

// Function that sets the password at the first login in the platform
func Register(w http.ResponseWriter, r *http.Request) {
	// Extract email and password from body
	var user models.User
	json.NewDecoder(r.Body).Decode(&user)
	log.Printf("Register attempt for email: %s", user.Email)

	// Check if user exists in the database
	// If the user does not exist in the database, the person is not part of the company
	storedUser, err := db.GetUserByEmail(user.Email)
	if err != nil {
		log.Printf("Error for user with email %s: %v", user.Email, err)
		http.Error(w, "User not found. Please contact your company administrator.", http.StatusNotFound)
		return
	}

	// Check if the password has already been set
	if storedUser.IsPasswordSet && storedUser.Activated {
		log.Printf("Error for user with email %s: %v", user.Email, err)
		http.Error(w, "User is already registered. Please go to the  page.", http.StatusConflict)
		return
	}

	// Validate password
	err = validatePassword(user.Password)
	if err != nil {
		log.Printf("Invalid password for email: %s", user.Email)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Hash the password
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		log.Printf("Error hashing password for email: %s", user.Email)
		http.Error(w, "Error hashing password. Please try again later.", http.StatusInternalServerError)
		return
	}
	user.Password = string(hashedPassword)

	// Save the password in the database
	err = db.SetUserPassword(user.Email, user.Password)
	if err != nil {
		log.Printf("Error for user with email %s: %v", user.Email, err)
		http.Error(w, "Error setting password. Please try again later.", http.StatusInternalServerError)
		return
	}

	// Generate confirmation token
	tokenString, err := utils.GenerateJWT(storedUser.ID, storedUser.Email, storedUser.Role)
	if err != nil {
		log.Printf("Error generating confirmation token for user: %s", user.Email)
		http.Error(w, "Error generating confirmation token. Please try again later.", http.StatusInternalServerError)
		return
	}

	confirmLink := fmt.Sprintf("http://localhost:8080/confirm?token=%s", tokenString)

	// Send confirmation email
	err = utils.SendConfirmationEmail(user.Email, confirmLink)
	if err != nil {
		log.Printf("Error ending email for user with email %s: %v", user.Email, err)
		http.Error(w, "Error sending email. Please try again later.", http.StatusInternalServerError)
		return
	}
	log.Printf("Confirmation email sent to: %s", user.Email)

	w.WriteHeader(http.StatusCreated)
	log.Printf("Registration successful for user: %s", user.Email)
}

// Manage the link accessed to confirm the account
// A confirmation link will be sent by email
// If the route is accessed, the account is successfully activated
func ConfirmEmail(w http.ResponseWriter, r *http.Request) {
	// Check the confirmation token
	tokenString := r.URL.Query().Get("token")
	if tokenString == "" {
		log.Println("No confirmation token provided")
		http.Error(w, "No confirmation token provided", http.StatusBadRequest)
		return
	}

	// Extract claims
	claims := &models.Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(utils.GetJWTSecret()), nil
	})

	// Token validation
	if err != nil || !token.Valid {
		log.Printf("Invalid or expired token: %s", tokenString)
		http.Error(w, "Invalid or expired token", http.StatusBadRequest)
		return
	}

	// Extract user details
	userID := claims.UserID
	email := claims.Email

	// Activate account for the user with the specified email
	err = db.ActivateUserByEmail(email)
	if err != nil {
		log.Printf("Error for user with email %s: %v", email, err)
		http.Error(w, "Error activating user. Please try again later.", http.StatusInternalServerError)
		return
	}
	log.Printf("User activated for email: %s", email)

	// Update IsPasswordSet only after activation
	// If the account has not been activated, we allow adding a new password with each registration attempt.
	err = db.UpdateUserPasswordStatus(email, true)
	if err != nil {
		log.Printf("Error for user with email %s: %v", email, err)
		http.Error(w, "Error updating user status. Please try again later.", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	log.Printf("Account confirmed successfully for user ID: %s", userID)
}

// Check if the password meets the requirements
func validatePassword(password string) error {
	// Check if the password is at least 5 characters long
	if len(password) < 5 {
		return errors.New("password must be at least 5 characters long")
	}

	// Regular expression to check for at least 3 digits
	re := regexp.MustCompile(`[0-9]`)
	matches := re.FindAllString(password, -1)

	// Check if there are at least 3 digits
	if len(matches) < 3 {
		return errors.New("password must contain at least 3 digits")
	}

	return nil
}
