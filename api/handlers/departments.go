package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"

	"department-management/db"
	"department-management/models"
)

// Add a new department
func CreateDepartment(w http.ResponseWriter, r *http.Request) {
	// If the new department is top level, it will not have a parent department.
	// Otherwise, it will be attached under an existing department
	var dept struct {
		models.Department
		ParentName string `json:"parent_name"`
	}

	// Extract department details from body
	if err := json.NewDecoder(r.Body).Decode(&dept); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Generate UUID for the new department
	deptID := uuid.New().String()

	// Add the new department to database
	err := db.CreateDepartment(deptID, dept.Name)
	if err != nil {
		log.Printf("Error creating department: %v", err)
		http.Error(w, "Error creating department", http.StatusInternalServerError)
		return
	}

	// Check if it's a top level department
	if dept.ParentName != "" {
		parentID, err := db.GetDepartmentIDByName(dept.ParentName)
		if err != nil {
			log.Printf("Error getting parent department ID: %v", err)
			http.Error(w, "Error getting parent department ID", http.StatusInternalServerError)
			return
		}

		// Add relationship to database
		err = db.AddDepartmentToHierarchy(parentID, deptID)
		if err != nil {
			log.Printf("Error adding department to hierarchy: %v", err)
			http.Error(w, "Error adding department to hierarchy", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Department added successfully"))
}

// Update department fields
func UpdateDepartment(w http.ResponseWriter, r *http.Request) {
	// All the fields required
	var dept models.Department
	json.NewDecoder(r.Body).Decode(&dept)

	// Add changes to database
	err := db.UpdateDepartment(dept.ID, dept.Name, dept.Flags)
	if err != nil {
		log.Printf("Error updating department: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Department updated successfully"))

}

// Update flags field
func DeleteDepartment(w http.ResponseWriter, r *http.Request) {
	// Department id required1
	var requestBody struct {
		ID string `json:"id"`
	}
	json.NewDecoder(r.Body).Decode(&requestBody)

	if requestBody.ID == "" {
		log.Printf("Error: Missing department ID")
		http.Error(w, "Missing department ID", http.StatusBadRequest)
		return
	}

	// Update database
	err := db.DeleteDepartment(requestBody.ID)
	if err != nil {
		log.Printf("Error deleting department: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Department marked as deleted successfully"))
}

// Get all existing departments
func GetAllDepartments(w http.ResponseWriter, r *http.Request) {
	departments, err := db.GetAllDepartments()
	if err != nil {
		log.Printf("Error getting departments: %v", err)
		http.Error(w, "Error getting departments", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(departments)
}

// Get sub-departments for specified department name
func GetHierarchy(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		Name string `json:"name"`
	}
	json.NewDecoder(r.Body).Decode(&requestBody)

	if requestBody.Name == "" {
		log.Printf("Error: Missing department name")
		http.Error(w, "Missing department name", http.StatusBadRequest)
		return
	}

	// Fetch departments
	hierarchy, err := db.GetHierarchy(requestBody.Name)
	if err != nil {
		log.Printf("Error getting hierarchy: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(hierarchy)
}
