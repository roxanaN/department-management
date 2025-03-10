package db

import (
	"database/sql"
	"department-management/config"
	"department-management/models"
	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

// Initializes the database connection.
func InitDB(cfg *config.Config) (*sql.DB, error) {
	dsn := cfg.DBUser + ":" + cfg.DBPassword + "@tcp(" + cfg.DBHost + ":" + cfg.DBPort + ")/" + cfg.DBName
	database, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	DB = database

	return database, nil
}

// Functions for calling stored procedures
func CreateDepartment(deptID, deptName string) error {
	_, err := DB.Exec("CALL CreateDepartment(?, ?)", deptID, deptName)
	return err
}

func UpdateDepartment(deptID, deptName string, deptFlags int) error {
	_, err := DB.Exec("CALL UpdateDepartment(?, ?, ?)", deptID, deptName, deptFlags)
	return err
}

func DeleteDepartment(deptID string) error {
	_, err := DB.Exec("CALL DeleteDepartment(?)", deptID)
	return err
}

func GetAllDepartments() ([]models.Department, error) {
	rows, err := DB.Query("CALL GetAllDepartments()")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var departments []models.Department
	for rows.Next() {
		var department models.Department
		if err := rows.Scan(&department.ID, &department.Name, &department.Flags); err != nil {
			return nil, err
		}
		departments = append(departments, department)
	}
	return departments, nil
}

func GetHierarchy(deptName string) ([]models.Department, error) {
	rows, err := DB.Query("CALL GetHierarchy(?)", deptName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var hierarchy []models.Department
	for rows.Next() {
		var dept models.Department
		if err := rows.Scan(&dept.ID, &dept.Name, &dept.Flags); err != nil {
			return nil, err
		}
		hierarchy = append(hierarchy, dept)
	}
	return hierarchy, nil
}

func CreateUser(user *models.User) (string, error) {
	var userID string
	err := DB.QueryRow("CALL CreateUser(?, ?, ?)", user.Email, user.Password, user.Role).Scan(&userID)
	if err != nil {
		return "", err
	}
	return userID, nil
}

func GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := DB.QueryRow("CALL GetUserByEmail(?)", email).Scan(&user.ID, &user.Email, &user.Password, &user.IsPasswordSet, &user.Activated, &user.Role, &user.Department)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func SetUserPassword(email, hashedPassword string) error {
	_, err := DB.Exec("CALL SetUserPassword(?, ?)", email, hashedPassword)
	return err
}

func ActivateUserByEmail(email string) error {
	_, err := DB.Exec("CALL ActivateUser(?)", email)
	return err
}

func UpdateUserPasswordStatus(email string, status bool) error {
	_, err := DB.Exec("CALL UpdateUserPasswordStatus(?, ?)", email, status)
	return err
}

func GetRoleIDByName(roleName string) (string, error) {
	var roleID string
	query := "CALL GetRoleIDByName(?)"
	err := DB.QueryRow(query, roleName).Scan(&roleID)
	if err != nil {
		return "", err
	}
	return roleID, nil
}

func GetDepartmentIDByName(departmentName string) (string, error) {
	var departmentID string
	query := "CALL GetDepartmentIDByName(?)"
	err := DB.QueryRow(query, departmentName).Scan(&departmentID)
	if err != nil {
		return "", err
	}
	return departmentID, nil
}

func AddDepartmentToHierarchy(parentID, childID string) error {
	query := "CALL AddDepartmentToHierarchy(?, ?)"
	_, err := DB.Exec(query, parentID, childID)
	return err
}

func AddUser(userID, userEmail, roleID, departmentID string) error {
	query := "CALL AddUser(?, ?, ?, ?)"
	_, err := DB.Exec(query, userID, userEmail, roleID, departmentID)
	return err
}

func ChangeUserDepartment(userID, newDepartmentID string) error {
	query := "CALL ChangeUserDepartment(?, ?)"
	_, err := DB.Exec(query, userID, newDepartmentID)
	return err
}

func ChangeUserRole(userID, roleID string) error {
	query := "CALL ChangeUserRole(?, ?)"
	_, err := DB.Exec(query, userID, roleID)
	return err
}

func GetAllUsers() ([]models.User, error) {
	query := "CALL GetAllUsers()"
	rows, err := DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.Email, &user.Role, &user.Department)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func DeleteUser(userID string) error {
	query := "CALL DeleteUser(?)"
	_, err := DB.Exec(query, userID)
	return err
}

func GetUsersByDepartment(departmentID string) ([]models.ShortUser, error) {
	query := "CALL GetUsersByDepartment(?)"
	rows, err := DB.Query(query, departmentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.ShortUser
	for rows.Next() {
		var user models.ShortUser
		if err := rows.Scan(&user.ID, &user.Email, &user.Role); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}
