package JWTBackend

import (
	jwt "github.com/dgrijalva/jwt-go"
)

var CurrentGlobalAuthenticationBackend IAuthenticationBackend

type IAuthenticationBackend interface {
	GenerateToken(userUUID string) (string, error)
	Authenticate(email, username, password string) (authenticated bool, UUID string)
	GetUUIDFromEmailAndUserName(email, username string) (gotUUID bool, UUID string)
	Logout(tokenString string, token *jwt.Token) error
	IsInBlacklist(token string) bool
	GetPublicKey(token *jwt.Token) (interface{}, error)
}
