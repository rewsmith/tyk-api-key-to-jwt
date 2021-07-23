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
	apikey = "APIKey"
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
	//var tykUserFields = session.MetaData["tyk_user_fields"].(string)
	var tykUserFields = session.MetaData["tyk_user_fields"]
	userFieldsStr, _ := json.Marshal(tykUserFields)
	json.Unmarshal([]byte(userFieldsStr), &user)

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

func returnNoAuth(w http.ResponseWriter, errorMessage string) {
	jsonData, err := json.Marshal(errorMessage)
	if err != nil {
		logger.Info("Couldn't marshal")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, _ = w.Write(jsonData)
}


// Run on startup by Tyk when loaded.  Bootstrapping the service here
func init() {
	logger.Info("log ApiKeyToJwt init")
}

func main() {}
