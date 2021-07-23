package main

import (
	"encoding/json"
	"fmt"
	"github.com/TykTechnologies/tyk/ctx"
	"github.com/TykTechnologies/tyk/log"
	jwt "github.com/dgrijalva/jwt-go"
	"net/http"
)

const (
	authorizationHeader = "Authorization"
)

// User struct
type User struct {
	Username string `json:"username"`
	Name     string `json:"name"`
}

var logger = log.Get()

// ApiKeyToJwt creates a signed JWT from an APIKey
func ApiKeyToJwt(w http.ResponseWriter, r *http.Request) {
	logger.Info("ApiKeyToJwt main starting")

	// Lookup user details from developer metadata
	var user User
	session := ctx.GetSession(r)

	var tykUserFields = session.MetaData["tyk_user_fields"]
	userFieldsStr, err := json.Marshal(tykUserFields)
	logger.Info("tyk_user_fields= ", userFieldsStr)

	err = json.Unmarshal([]byte(userFieldsStr), &user)
	if len(user.Name) == 0 || len(user.Username) == 0 {
		writeError(w, "Developer Identity metadata not set")
		return
	}

	logger.Info("Developer Username: ", user.Username)
	logger.Info("Developer Name: ", user.Name)

	// Now create the JWT
	mySigningKey := []byte("my-256-bit-secret")

	type MyCustomClaims struct {
		Name string `json:"name"`
		jwt.StandardClaims
	}

	// Create the Claims
	claims := MyCustomClaims{
		user.Name,
		jwt.StandardClaims{
			Subject:   user.Username,
			ExpiresAt: 15000,
			Issuer:    "Tyk",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(mySigningKey)
	fmt.Printf("%v %v", signedToken, err)

	//Add JWT to Auth Header
	r.Header.Set(authorizationHeader, signedToken)
}

func writeError(w http.ResponseWriter, errorMessage string) {
	w.WriteHeader(http.StatusInternalServerError)
	jsonData, err := json.Marshal(errorMessage)
	if err != nil {
		logger.Info("Couldn't marshal")
		return
	}
	_, _ = w.Write(jsonData)
}


// Run on startup by Tyk when loaded.  Bootstrapping the service here
func init() {
	logger.Info("log ApiKeyToJwt init")
}

func main() {}
