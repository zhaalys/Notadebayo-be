package user

import (
	"errors"
	"os"
	"strings"
	"time"

	"tasklybe/pkg/db"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrJWTSecretMissing   = errors.New("jwt secret missing")
)

func RegisterUser(input RegisterUserRequest) (*User, error) {
	// Check if email already exists
	var existing User
	err := db.DB.First(&existing, "email = ?", input.Email).Error
	if err == nil {
		return nil, ErrEmailAlreadyExists
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// Hash password
	hashed, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create user
	u := User{
		Email:    input.Email,
		Password: string(hashed),
		Name:     strings.TrimSpace(input.Name),
	}

	if err := db.DB.Create(&u).Error; err != nil {
		return nil, err
	}

	return &u, nil
}

func LoginUser(input LoginUserRequest) (*User, string, error) {
	// Check if user exists
	var u User
	if err := db.DB.First(&u, "email = ?", input.Email).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, "", ErrInvalidCredentials
		}
		return nil, "", err
	}

	// Compare password
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(input.Password)); err != nil {
		return nil, "", ErrInvalidCredentials
	}

	// Generate API key
	apiKey, err := generateAPIKey(&u)
	if err != nil {
		return nil, "", err
	}

	return &u, apiKey, nil
}

type authClaims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

func generateAPIKey(u *User) (string, error) {
	secret := strings.TrimSpace(os.Getenv("JWT_SECRET"))
	if secret == "" {
		return "", ErrJWTSecretMissing
	}

	now := time.Now()
	claims := authClaims{
		Email: u.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   u.ID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(24 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
