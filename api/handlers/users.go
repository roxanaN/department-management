package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"

	"department-management/db"
	"department-management/models"
)

// Create new user
func AddUser(w http.ResponseWriter, r *http.Request) {
	var user models.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid requestBody", http.StatusBadRequest)
		return
	}

	// Fetch role id for existing role name
	roleID, err := db.GetRoleIDByName(user.Role)
	if err != nil {
		log.Printf("Error getting role ID: %v", err)
		http.Error(w, "Error getting role ID", http.StatusInternalServerError)
		return
	}

	// Fetch department id for existing department name
	departmentID, err := db.GetDepartmentIDByName(user.Department)
	if err != nil {
		log.Printf("Error getting department ID: %v", err)
		http.Error(w, "Error getting department ID", http.StatusInternalServerError)
		return
	}

	// Add user to database
	err = db.AddUser(uuid.New().String(), user.Email, roleID, departmentID)
	if err != nil {
		log.Printf("Error adding user with department: %v", err)
		http.Error(w, "Error adding user with department", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User added successfully"))
}

// Delete user by user id
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		UserID string `json:"user_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		log.Printf(": %v", err)
		http.Error(w, "Invalid requestBody. Provide user_id", http.StatusBadRequest)
		return
	}

	// Update database
	err := db.DeleteUser(requestBody.UserID)
	if err != nil {
		log.Printf("Error deleting user: %v", err)
		http.Error(w, "Error deleting user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User deleted successfully"))
}

// Fetch all existing users
func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := db.GetAllUsers()
	if err != nil {
		log.Printf("Error getting all users: %v", err)
		http.Error(w, "Error getting all users", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// Change user role
func ChangeUserRole(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		UserID  string `json:"user_id"`
		NewRole string `json:"new_role"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Invalid requestBody", http.StatusBadRequest)
		return
	}

	// Fetch role id by role name
	roleID, err := db.GetRoleIDByName(requestBody.NewRole)
	if err != nil {
		log.Printf("Error getting role ID: %v", err)
		http.Error(w, "Error getting role ID", http.StatusInternalServerError)
		return
	}

	// Update database
	err = db.ChangeUserRole(requestBody.UserID, roleID)
	if err != nil {
		log.Printf("Error changing user role: %v", err)
		http.Error(w, "Error changing user role", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User role changed successfully"))
}

// Migrate user to another department
func ChangeUserDepartment(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		UserID        string `json:"user_id"`
		NewDepartment string `json:"new_department"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Invalid requestBody", http.StatusBadRequest)
		return
	}

	// Fetch department id by department name
	departmentID, err := db.GetDepartmentIDByName(requestBody.NewDepartment)
	if err != nil {
		log.Printf("Error getting department ID: %v", err)
		http.Error(w, "Error getting department ID", http.StatusInternalServerError)
		return
	}

	// Update database
	err = db.ChangeUserDepartment(requestBody.UserID, departmentID)
	if err != nil {
		log.Printf("Error changing user department: %v", err)
		http.Error(w, "Error changing user department", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User department changed successfully"))
}

// Fetch user details for specified user id
func GetUsersByDepartmentID(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		DepartmentID string `json:"department_id"`
	}
	json.NewDecoder(r.Body).Decode(&requestBody)

	if requestBody.DepartmentID == "" {
		log.Printf("Error: Missing department_id")
		http.Error(w, "Missing department_id", http.StatusBadRequest)
		return
	}

	// Update database
	users, err := db.GetUsersByDepartment(requestBody.DepartmentID)
	if err != nil {
		log.Printf("Error getting users: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// Fetch all departments and associated users under the specified department
func GetCompleteHierarchy(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		Name string `json:"department_name"`
	}
	json.NewDecoder(r.Body).Decode(&requestBody)

	if requestBody.Name == "" {
		log.Printf("Error: Missing department_name")
		http.Error(w, "Missing department_name", http.StatusBadRequest)
		return
	}

	// Fetch departments hierarchy
	hierarchy, err := db.GetHierarchy(requestBody.Name)
	if err != nil {
		log.Printf("Error getting hierarchy: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	completeHierarchy := []models.CompleteDepartment{}

	// Iterate through all sub-departments to get users from each department
	for _, dept := range hierarchy {
		users, err := db.GetUsersByDepartment(dept.ID)
		if err != nil {
			log.Printf("Error getting users for department %s: %v", dept.ID, err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Populate the department-users structure
		completeDept := models.CompleteDepartment{
			Department: dept,
			Users:      users,
		}

		// Add object to hierarchy
		completeHierarchy = append(completeHierarchy, completeDept)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(completeHierarchy)
}
