package utils

import (
	"log"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/joho/godotenv"

	jwt "github.com/dgrijalva/jwt-go"
)

//var key = []byte("urefmdaffafefea")

func TestToken(t *testing.T) {
	data := map[string]interface{}{
		"userid":   11,
		"username": "huyntsgs",
	}

	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	key := []byte(os.Getenv("TOKEN_KEY"))

	tokenString, err := GenToken(data, key, 1)
	if err != nil {
		t.Fail()
	}

	log.Println("tokenString: ", tokenString)
	token, err := DecodeToken(tokenString, key)
	if !token.Valid {
		log.Println("!token.Valid ")
		t.Fail()
	}

	claims := token.Claims.(jwt.MapClaims)

	if int(claims["userid"].(float64)) != 11 {
		log.Println("userid is not correct", reflect.TypeOf(claims["userid"]))
		t.Fail()
	}

	if claims["username"] != "huyntsgs" {
		log.Println("username is not correct")
		t.Fail()
	}
	log.Printf("exp %v\n", claims["exp"])

	//  After token timeline 1 minutes, token must expire
	time.Sleep(time.Minute * 2)
	exp := int64(claims["exp"].(float64))
	if exp < time.Now().Unix() {
		log.Println("not exp ")
		log.Printf("now %v\n", time.Now().Unix())
		t.Fail()
	}
}
