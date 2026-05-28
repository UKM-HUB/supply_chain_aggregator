package usecase

import (
	"context"
	"errors"
	"strings"
	"time"

	"supply-chain-aggregator/services/auth-service/internal/entity"
	jwtmanager "supply-chain-aggregator/services/auth-service/internal/jwt"
	"supply-chain-aggregator/services/auth-service/internal/repository"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidCredential = errors.New("invalid email or password")

type RegisterInput struct {
	Name      string
	Email     string
	Phone     string
	Password  string
	Role      string
	Latitude  float64
	Longitude float64
}

type LoginOutput struct {
	Token string
	User  *entity.User
}

type AuthUsecase struct {
	userRepo   repository.UserRepository
	jwtManager *jwtmanager.Manager
}

func NewAuthUsecase(userRepo repository.UserRepository, jwtManager *jwtmanager.Manager) *AuthUsecase {
	return &AuthUsecase{
		userRepo:   userRepo,
		jwtManager: jwtManager,
	}
}

func (u *AuthUsecase) Register(ctx context.Context, input RegisterInput) (*entity.User, error) {
	role := strings.TrimSpace(input.Role)
	if role == "" {
		role = "SME"
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	user := &entity.User{
		ID:           uuid.NewString(),
		Name:         input.Name,
		Email:        input.Email,
		Phone:        input.Phone,
		PasswordHash: string(hashedPassword),
		Role:         role,
		Latitude:     input.Latitude,
		Longitude:    input.Longitude,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := u.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (u *AuthUsecase) Login(ctx context.Context, email, password string) (*LoginOutput, error) {
	user, err := u.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, ErrInvalidCredential
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, ErrInvalidCredential
	}

	token, err := u.jwtManager.GenerateToken(user.ID, user.Role)
	if err != nil {
		return nil, err
	}

	return &LoginOutput{
		Token: token,
		User:  user,
	}, nil
}

func (u *AuthUsecase) ValidateToken(token string) (*jwtmanager.CustomClaims, error) {
	return u.jwtManager.ValidateToken(token)
}
