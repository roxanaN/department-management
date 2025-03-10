package api

import (
	"department-management/api/handlers"
	"department-management/api/middleware"
	"github.com/gorilla/mux"
)

func RegisterRoutes() *mux.Router {
	router := mux.NewRouter()

	// Public routes
	router.HandleFunc("/register", handlers.Register).Methods("POST")
	router.HandleFunc("/confirm", handlers.ConfirmEmail).Methods("GET")
	router.HandleFunc("/login", handlers.Login).Methods("POST")

	// Protected routes
	protected := router.PathPrefix("/api").Subrouter()
	protected.Use(middleware.AuthMiddleware)

	// Departments routes
	protected.HandleFunc("/departments", handlers.CreateDepartment).Methods("POST")
	protected.HandleFunc("/departments", handlers.UpdateDepartment).Methods("PUT")
	protected.HandleFunc("/departments", handlers.GetAllDepartments).Methods("GET")
	protected.HandleFunc("/departments/delete", handlers.DeleteDepartment).Methods("PUT")
	protected.HandleFunc("/departments/hierarchy", handlers.GetHierarchy).Methods("GET")

	// User routes
	protected.HandleFunc("/users", handlers.AddUser).Methods("POST")
	protected.HandleFunc("/users", handlers.GetAllUsers).Methods("GET")
	protected.HandleFunc("/users", handlers.DeleteUser).Methods("DELETE")
	protected.HandleFunc("/users/department", handlers.ChangeUserDepartment).Methods("PUT")
	protected.HandleFunc("/users/role", handlers.ChangeUserRole).Methods("PUT")
	protected.HandleFunc("/users/departments", handlers.GetUsersByDepartmentID).Methods("GET")
	protected.HandleFunc("/users/hierarchy", handlers.GetCompleteHierarchy).Methods("GET")

	return router
}
