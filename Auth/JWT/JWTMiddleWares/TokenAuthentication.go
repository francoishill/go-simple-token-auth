package JWTMiddleWares

import (
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/francoishill/go-simple-token-auth/Auth/JWT/JWTBackend"
	"net/http"
)

var TokenHandlers iTokenHandlers = nil

type iTokenHandlers interface {
	HandleTokenExtractedSuccessfully(http.ResponseWriter, *http.Request, *jwt.Token)
}

func RequireTokenAuthentication(rw http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	authBackend := JWTBackend.CurrentGlobalAuthenticationBackend

	token, err := jwt.ParseFromRequest(req, authBackend.GetPublicKey)
	//TODO: We should probably check for specific errors like 'jwt.ValidationErrorExpired'

	if err == nil && token.Valid && !authBackend.IsInBlacklist(req.Header.Get("Authorization")) {
		if TokenHandlers != nil {
			TokenHandlers.HandleTokenExtractedSuccessfully(rw, req, token)
		}
		next(rw, req)
	} else {
		rw.WriteHeader(http.StatusUnauthorized)
	}
}
