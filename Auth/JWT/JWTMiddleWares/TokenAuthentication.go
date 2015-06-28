package JWTMiddleWares

import (
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/francoishill/go-simple-token-auth/Auth/JWT/JWTBackend"
	"net/http"
)

func RequireTokenAuthentication(rw http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	authBackend := JWTBackend.CurrentGlobalAuthenticationBackend

	token, err := jwt.ParseFromRequest(req, authBackend.GetPublicKey)
	//TODO: We should probably check for specific errors like 'jwt.ValidationErrorExpired'

	if err == nil && token.Valid && !authBackend.IsInBlacklist(req.Header.Get("Authorization")) {
		next(rw, req)
	} else {
		rw.WriteHeader(http.StatusUnauthorized)
	}
}
