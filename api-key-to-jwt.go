package main

import (
	"encoding/json"
	"fmt"
	"github.com/TykTechnologies/tyk/ctx"
	jwt "github.com/dgrijalva/jwt-go"
	"log"
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


// ApiKeyToJwt creates a signed JWT from an APIKey
func ApiKeyToJwt(w http.ResponseWriter, r *http.Request) {
	log.Println("ApiKeyToJwt main starting")

	// Lookup user details from developer metadata
	var user User
	session := ctx.GetSession(r)
	user.Username = session.MetaData["tyk_developer_username"].(string)
	user.Name = session.MetaData["tyk_developer_name"].(string)
	log.Println("Developer Username: ", user.Username)
	log.Println("Developer Name: ", user.Name)

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
		log.Println("Couldn't marshal")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, _ = w.Write(jsonData)
}


// Run on startup by Tyk when loaded.  Bootstrapping the service here
func init() {
	log.Println("log ApiKeyToJwt init")
}

func main() {}
