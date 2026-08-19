package helpers

import (
	"backend-test/internal/models"
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type ClaimToken struct {
	UserID uint            `json:"userId"`
	Role   models.UserRole `json:"role"`
	jwt.RegisteredClaims
}

func GenerateToken(_ context.Context, user *models.User, now time.Time) (string, error) {
	secret := GetEnv("APP_SECRET", "")
	if len(secret) < 32 {
		return "", errors.New("APP_SECRET must contain at least 32 characters")
	}
	claims := ClaimToken{UserID: user.ID, Role: user.Role, RegisteredClaims: jwt.RegisteredClaims{Issuer: GetEnv("APP_NAME", "maintenance-request-log"), Subject: user.Email, IssuedAt: jwt.NewNumericDate(now), ExpiresAt: jwt.NewNumericDate(now.Add(8 * time.Hour))}}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret))
}

func ValidateToken(token string) (*ClaimToken, error) {
	secret := GetEnv("APP_SECRET", "")
	parsed, err := jwt.ParseWithClaims(token, &ClaimToken{}, func(t *jwt.Token) (interface{}, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, errors.New("invalid signing method")
		}
		return []byte(secret), nil
	}, jwt.WithIssuer(GetEnv("APP_NAME", "maintenance-request-log")), jwt.WithExpirationRequired())
	if err != nil {
		return nil, err
	}
	claims, ok := parsed.Claims.(*ClaimToken)
	if !ok || !parsed.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}
