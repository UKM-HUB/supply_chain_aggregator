package usecase

import (
    "project/internal/entity"
    "project/internal/helper"
    "project/internal/repository"

    "github.com/google/uuid"
)

type AuthUsecase struct {
    Repo *repository.UserRepository
}

func NewAuthUsecase(repo *repository.UserRepository) *AuthUsecase {
    return &AuthUsecase{Repo: repo}
}

func (u *AuthUsecase) Register(user *entity.User) error {
    user.ID = uuid.New()

    return u.Repo.Create(user)
}

func (u *AuthUsecase) Login(email string) (string, error) {
    user, err := u.Repo.FindByEmail(email)
    if err != nil {
        return "", err
    }

    return helper.GenerateJWT(user.ID.String())
}