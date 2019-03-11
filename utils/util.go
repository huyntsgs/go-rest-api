package utils

import (
	"fmt"
	"log"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

var KEY = []byte("urefmdaffafefea")

// GenToken encrypts values data with encryptKey key and
// add expire time expTime in minutes.
func GenToken(values map[string]interface{}, encryptedKey []byte, expTime int) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claim := token.Claims.(jwt.MapClaims)

	for k, v := range values {
		claim[k] = v
		claim["exp"] = time.Now().Add(time.Duration(expTime) * time.Minute).Unix()
	}

	tokenString, err := token.SignedString(encryptedKey)
	if err != nil {
		log.Printf("error SignedString %v", err)
		return "", err
	}

	return tokenString, nil
}

// DecodeToken decrypts tokenString with provided encryptedKey
func DecodeToken(tokenString string, encryptedKey []byte) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Error in DecodeToken")
		}
		return encryptedKey, nil
	})

	return token, err
}
