package jwtMiddleware

import jwt "github.com/golang-jwt/jwt/v5"

type Claims struct {
	UserID   int64  `json:"uid"`
	Username string `json:"username"`
	Email    string `json:"email"`
	jwt.RegisteredClaims
}
