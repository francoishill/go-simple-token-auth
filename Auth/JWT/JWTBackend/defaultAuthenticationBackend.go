package JWTBackend

import (
	jwt "github.com/dgrijalva/jwt-go"
	"io/ioutil"
	"time"
)

type authenticateMethodSignature func(email, username, password string) (authenticated bool, UUID string)
type getUUIDFromEmailAndUserNameMethodSignature func(email, username string) (gotUUID bool, UUID string)

func CreateDefaultAuthenticationBackend(privateKeyPath, publicKeyPath string, jwtExpirationDelta int, authenticateMethod authenticateMethodSignature, getUUIDFromEmailAndUserNameMethod getUUIDFromEmailAndUserNameMethodSignature, keyValueStoreProvider IKeyValueStoreProvider) IAuthenticationBackend {
	if CurrentGlobalAuthenticationBackend == nil {
		CurrentGlobalAuthenticationBackend = &defaultAuthenticationBackend{
			privateKey:                        getPrivateKeyFromFile(privateKeyPath),
			publicKey:                         getPublicKeyFromFile(publicKeyPath),
			jwtExpirationDelta:                jwtExpirationDelta,
			authenticateMethod:                authenticateMethod,
			getUUIDFromEmailAndUserNameMethod: getUUIDFromEmailAndUserNameMethod,
			keyValueStoreProvider:             keyValueStoreProvider,
		}
	}

	return CurrentGlobalAuthenticationBackend
}

type defaultAuthenticationBackend struct {
	privateKey                        []byte
	publicKey                         []byte
	jwtExpirationDelta                int
	authenticateMethod                authenticateMethodSignature
	getUUIDFromEmailAndUserNameMethod getUUIDFromEmailAndUserNameMethodSignature
	keyValueStoreProvider             IKeyValueStoreProvider
}

const (
	tokenDuration = 72
	expireOffset  = 3600
)

func (this *defaultAuthenticationBackend) GetPublicKey(token *jwt.Token) (interface{}, error) {
	return this.publicKey, nil
}

func (this *defaultAuthenticationBackend) GenerateToken(userUUID string) (string, error) {
	token := jwt.New(jwt.GetSigningMethod("RS256"))
	token.Claims["exp"] = time.Now().Add(time.Hour * time.Duration(this.jwtExpirationDelta)).Unix()
	token.Claims["iat"] = time.Now().Unix()
	token.Claims["sub"] = userUUID
	tokenString, err := token.SignedString(this.privateKey)
	if err != nil {
		panic(err)
		return "", err
	}
	return tokenString, nil
}

func (this *defaultAuthenticationBackend) Authenticate(email, username, password string) (authenticated bool, UUID string) {
	return this.authenticateMethod(email, username, password)
}

func (this *defaultAuthenticationBackend) GetUUIDFromEmailAndUserName(email, username string) (gotUUID bool, UUID string) {
	return this.getUUIDFromEmailAndUserNameMethod(email, username)
}

func (this *defaultAuthenticationBackend) getTokenRemainingValidity(timestamp interface{}) int {
	if validity, ok := timestamp.(float64); ok {
		tm := time.Unix(int64(validity), 0)
		remainer := tm.Sub(time.Now())
		if remainer > 0 {
			return int(remainer.Seconds() + expireOffset)
		}
	}
	return expireOffset
}

func (this *defaultAuthenticationBackend) Logout(tokenString string, token *jwt.Token) error {
	storeConnection := this.keyValueStoreProvider.GetNewConnection()
	return storeConnection.SetValue(tokenString, tokenString, this.getTokenRemainingValidity(token.Claims["exp"]))
}

func (this *defaultAuthenticationBackend) IsInBlacklist(token string) bool {
	storeConnection := this.keyValueStoreProvider.GetNewConnection()
	tokenInStore, _ := storeConnection.GetValue(token)

	return tokenInStore != nil
}

func getPrivateKeyFromFile(privateKeyPath string) []byte {
	privateKey, err := ioutil.ReadFile(privateKeyPath)
	if err != nil {
		panic(err)
	}

	return privateKey
}

func getPublicKeyFromFile(publicKeyPath string) []byte {
	publicKey, err := ioutil.ReadFile(publicKeyPath)
	if err != nil {
		panic(err)
	}

	return publicKey
}
