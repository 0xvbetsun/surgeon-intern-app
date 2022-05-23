package middleware

import (
	"encoding/json"
	"fmt"
	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/form3tech-oss/jwt-go"
	"github.com/friendsofgo/errors"
	"net/http"
	"os"
)

type Jwks struct {
	Keys []JSONWebKeys `json:"keys"`
}

type JSONWebKeys struct {
	Kty string   `json:"kty"`
	Kid string   `json:"kid"`
	Use string   `json:"use"`
	N   string   `json:"n"`
	E   string   `json:"e"`
	X5c []string `json:"x5c"`
}

var jwtMiddleware = jwtmiddleware.New(jwtmiddleware.Options{
	ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
		// Verify 'aud' claim
		requiredAud := os.Getenv("AUTH0_AUDIENCE")
		aud, ok := token.Claims.(jwt.MapClaims)["aud"].([]interface{})
		if ok {
			// this means aud is a slice of strings
			s := make([]string, len(aud)) // https://stackoverflow.com/questions/44027826/convert-interface-to-string-in-golang
			for i, v := range aud {
				s[i] = fmt.Sprint(v)
			}
			token.Claims.(jwt.MapClaims)["aud"] = s
		}

		checkAud := token.Claims.(jwt.MapClaims).VerifyAudience(requiredAud, false)
		if !checkAud {
			return token, errors.New("Invalid audience.")
		}
		// Verify 'iss' claim
		iss := os.Getenv("AUTH0_DOMAIN")
		checkIss := token.Claims.(jwt.MapClaims).VerifyIssuer(iss, false)
		if !checkIss {
			return token, errors.New("Invalid issuer.")
		}

		cert, err := getPemCert(token)
		if err != nil {
			panic(err.Error())
		}

		result, _ := jwt.ParseRSAPublicKeyFromPEM([]byte(cert))
		return result, nil
	},
	SigningMethod: jwt.SigningMethodRS256,
})

func externalUserIDFromClaims(claims jwt.Claims) string {
	return claims.(jwt.MapClaims)["sub"].(string)
}

func getPemCert(token *jwt.Token) (string, error) {
	cert := ""
	resp, err := http.Get(os.Getenv("AUTH0_DOMAIN") + ".well-known/jwks.json")

	if err != nil {
		return cert, err
	}
	defer resp.Body.Close()

	var jwks = Jwks{}
	err = json.NewDecoder(resp.Body).Decode(&jwks)

	if err != nil {
		return cert, err
	}

	for k, _ := range jwks.Keys {
		if token.Header["kid"] == jwks.Keys[k].Kid {
			cert = "-----BEGIN CERTIFICATE-----\n" + jwks.Keys[k].X5c[0] + "\n-----END CERTIFICATE-----"
		}
	}

	if cert == "" {
		err := errors.New("Unable to find appropriate key.")
		return cert, err
	}

	return cert, nil
}
