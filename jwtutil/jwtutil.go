package jwtutil

import (
	"fmt"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

var (
	jwtSigningMethod = jwt.SigningMethodHS256
	jwtSecret        = "JWT_SECRET"
	expDuration      = time.Hour * 24
)

func SetJwtSigningMethod(signingMethod *jwt.SigningMethodHMAC) {
	jwtSigningMethod = signingMethod
}
func SetJwtSecret(secret string) {
	jwtSecret = secret
}
func SetExpDuration(d time.Duration) {
	expDuration = d
}

func NewTokenWithSecret(m map[string]interface{}, jwtSecret string) (string, error) {
	claims := jwt.MapClaims{
		"nbf": time.Now().Unix(),
		"exp": time.Now().Add(expDuration).Unix(),
	}
	for k, v := range m {
		claims[k] = v
	}
	return jwt.NewWithClaims(jwtSigningMethod, claims).SignedString([]byte(jwtSecret))
}
func NewToken(m map[string]interface{}) (string, error) {
	return NewTokenWithSecret(m, jwtSecret)
}

func Renew(token string) (string, error) {
	claim, err := Extract(token)
	if err != nil {
		return "", err
	}
	claim["nbf"] = time.Now().Unix()
	claim["exp"] = time.Now().Add(expDuration).Unix()
	return jwt.NewWithClaims(jwtSigningMethod, claim).SignedString([]byte(jwtSecret))
}

func EditPayload(token string, m map[string]string) (string, error) {
	claimInfo, err := Extract(token)
	if err != nil {
		return "", err
	}

	for k, v := range m {
		claimInfo[k] = v
	}

	return jwt.NewWithClaims(jwtSigningMethod, claimInfo).SignedString([]byte(jwtSecret))
}
func Extract(token string) (jwt.MapClaims, error) {
	return ExtractWithSecret(token, jwtSecret)
}
func ExtractWithSecret(token, jwtSecret string) (jwt.MapClaims, error) {
	if token == "" {
		return nil, fmt.Errorf("Required authorization token not found")
	}

	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) { return []byte(jwtSecret), nil })
	if err != nil {
		return nil, fmt.Errorf("Error parsing token: %v", err)
	}

	if jwtSigningMethod != nil && jwtSigningMethod.Alg() != parsedToken.Header["alg"] {
		return nil, fmt.Errorf("Expected %s signing method but token specified %s",
			jwtSigningMethod.Alg(),
			parsedToken.Header["alg"])
	}

	if !parsedToken.Valid {
		return nil, fmt.Errorf("Token is invalid")
	}

	claimInfo, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token")
	}
	return claimInfo, nil
}
