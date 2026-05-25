package usecase

import (
	"context"
	"errors"
	"b2b-aggregator/internal/auth/entity"
	"b2b-aggregator/internal/auth/repository"
	"b2b-aggregator/pkg/jwt"
	"golang.org/x/crypto/bcrypt"
)

type AuthUsecase struct {
	repo *repository.UserRepository
}

func NewAuthUsecase(repo *repository.UserRepository) *AuthUsecase {
	return &AuthUsecase{repo: repo}
}

func (u *AuthUsecase) Register(ctx context.Context, username, email, password, role string, lat, lng float64) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &entity.User{
		Username:  username,
		Email:     email,
		Password:  string(hashedPassword),
		Role:      role,
		Latitude:  lat,
		Longitude: lng,
	}

	return u.repo.CreateUser(ctx, user)
}

func (u *AuthUsecase) Login(ctx context.Context, email, password string) (string, error) {
	user, err := u.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return "", errors.New("email tidak ditemukan")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("password salah")
	}

	token, err := jwt.GenerateToken(user.ID, user.Role)
	if err != nil {
		return "", err
	}

	return token, nil
}
