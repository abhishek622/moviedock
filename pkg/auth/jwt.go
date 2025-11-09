package auth

import (
        "errors"
        "os"
        "time"

        "github.com/golang-jwt/jwt/v5"
)

var (
        ErrInvalidToken = errors.New("invalid token")
)

func getJWTSecret() []byte {
        secret := os.Getenv("SESSION_SECRET")
        if secret == "" {
                secret = "default-dev-secret-change-in-production"
        }
        return []byte(secret)
}

type Claims struct {
        UserID   string `json:"user_id"`
        Username string `json:"username"`
        jwt.RegisteredClaims
}

func GenerateToken(userID, username string) (string, error) {
        claims := Claims{
                UserID:   userID,
                Username: username,
                RegisteredClaims: jwt.RegisteredClaims{
                        ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
                        IssuedAt:  jwt.NewNumericDate(time.Now()),
                },
        }

        token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
        return token.SignedString(getJWTSecret())
}

func ValidateToken(tokenString string) (*Claims, error) {
        token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
                if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                        return nil, ErrInvalidToken
                }
                return getJWTSecret(), nil
        })

        if err != nil {
                return nil, err
        }

        if claims, ok := token.Claims.(*Claims); ok && token.Valid {
                return claims, nil
        }

        return nil, ErrInvalidToken
}

func GenerateSystemToken() (string, error) {
        claims := Claims{
                UserID:   "system",
                Username: "system",
                RegisteredClaims: jwt.RegisteredClaims{
                        ExpiresAt: jwt.NewNumericDate(time.Now().Add(365 * 24 * time.Hour)),
                        IssuedAt:  jwt.NewNumericDate(time.Now()),
                },
        }

        token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
        return token.SignedString(getJWTSecret())
}
