package entity

import "github.com/google/uuid"

type User struct {
    ID       uuid.UUID `gorm:"type:uuid;primaryKey"`
    Name     string
    Email    string `gorm:"unique"`
    Password string
    Role     string
}