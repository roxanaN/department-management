package models

import "github.com/golang-jwt/jwt/v4"

type Department struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Flags int    `json:"flags"`
}

type User struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	Password      string `json:"password"`
	IsPasswordSet bool   `json:"is_password_set"`
	Activated     bool   `json:"activated"`
	Role          string `json:"role"`
	Department    string `json:"department"`
}

type ShortUser struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

type Role struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type CompleteDepartment struct {
	Department Department  `json:"department"`
	Users      []ShortUser `json:"users"`
}

type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}
