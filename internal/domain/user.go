package domain

import "time"

type Role string

const (
    RoleUser  Role = "user"
    RoleAdmin Role = "admin"
)

type User struct {
    ID           string
    Username     string
    PasswordHash string    // never expose in JSON responses
    Role         Role
    CreatedAt    time.Time
}