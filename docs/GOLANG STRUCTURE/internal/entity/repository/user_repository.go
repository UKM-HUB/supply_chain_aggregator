package repository

import (
    "project/internal/entity"

    "gorm.io/gorm"
)

type UserRepository struct {
    DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
    return &UserRepository{DB: db}
}

func (r *UserRepository) Create(user *entity.User) error {
    return r.DB.Create(user).Error
}

func (r *UserRepository) FindByEmail(email string) (*entity.User, error) {
    var user entity.User

    err := r.DB.Where("email = ?", email).First(&user).Error

    return &user, err
}