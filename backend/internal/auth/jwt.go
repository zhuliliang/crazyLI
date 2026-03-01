package auth

import (
    "time"

    "github.com/golang-jwt/jwt/v5"
    "github.com/google/uuid"
)

// Claims represents session claims embedded in cookies.
type Claims struct {
    UserID uuid.UUID `json:"user_id"`
    Email  string    `json:"email"`
    Name   string    `json:"name"`
    jwt.RegisteredClaims
}

// CreateSessionToken returns a signed JWT string.
func CreateSessionToken(secret string, userID uuid.UUID, email, name string) (string, error) {
    claims := Claims{
        UserID: userID,
        Email:  email,
        Name:   name,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
        },
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(secret))
}

// ParseSession verifies the JWT and returns claims.
func ParseSession(secret, tokenStr string) (Claims, error) {
    token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        return []byte(secret), nil
    })
    if err != nil {
        return Claims{}, err
    }
    if claims, ok := token.Claims.(*Claims); ok && token.Valid {
        return *claims, nil
    }
    return Claims{}, jwt.ErrTokenInvalidClaims
}
