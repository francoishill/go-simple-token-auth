package JWTServices

import (
	"encoding/json"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/francoishill/go-simple-token-auth/Auth/JWT/JWTBackend"
	"net/http"
)

type tokenAuthentication struct {
	Token string `json:"token" form:"token"`
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func getResponseBytesFromToken(token string) []byte {
	response, err := json.Marshal(tokenAuthentication{token})
	checkError(err)
	return response
}

func Login(email, username, password string) (int, []byte) {
	authBackend := JWTBackend.CurrentGlobalAuthenticationBackend

	if authenticated, UUID := authBackend.Authenticate(email, username, password); authenticated {
		token, err := authBackend.GenerateToken(UUID)
		if err != nil {
			return http.StatusInternalServerError, []byte("")
		} else {
			return http.StatusOK, getResponseBytesFromToken(token)
		}
	}

	return http.StatusUnauthorized, []byte("")
}

func RefreshToken(email, username string) []byte {
	authBackend := JWTBackend.CurrentGlobalAuthenticationBackend

	if gotUUID, UUID := authBackend.GetUUIDFromEmailAndUserName(email, username); gotUUID {
		token, err := authBackend.GenerateToken(UUID)
		checkError(err)

		return getResponseBytesFromToken(token)
	} else {
		panic("User does not exist")
	}
}

func Logout(req *http.Request) error {
	authBackend := JWTBackend.CurrentGlobalAuthenticationBackend
	tokenRequest, err := jwt.ParseFromRequest(req, authBackend.GetPublicKey)
	if err != nil {
		return err
	}

	tokenString := req.Header.Get("Authorization")
	err = authBackend.Logout(tokenString, tokenRequest)
	if err != nil {
		return err
	}

	return nil
}
