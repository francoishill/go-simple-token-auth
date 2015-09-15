package JWTMiddleWares

import (
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/francoishill/go-simple-token-auth/Auth/JWT/JWTBackend"
)

type iTokenHandlers interface {
	HandleTokenExtractedSuccessfully(http.ResponseWriter, *http.Request, *jwt.Token)
}

type JWTMiddleWares struct {
	TokenHandlers iTokenHandlers
}

func (j *JWTMiddleWares) RequireTokenAuthentication(rw http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	authBackend := JWTBackend.CurrentGlobalAuthenticationBackend

	token, err := jwt.ParseFromRequest(req, authBackend.GetPublicKey)
	//TODO: We should probably check for specific errors like 'jwt.ValidationErrorExpired'

	if err == nil && token.Valid && !authBackend.IsInBlacklist(req.Header.Get("Authorization")) {
		if j.TokenHandlers != nil {
			j.TokenHandlers.HandleTokenExtractedSuccessfully(rw, req, token)
		}
		next(rw, req)
	} else {
		http.Error(rw, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
	}
}

func New(tokenHandlers iTokenHandlers) *JWTMiddleWares {
	return &JWTMiddleWares{
		tokenHandlers,
	}
}
