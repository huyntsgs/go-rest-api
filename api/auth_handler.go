package api

import (
	"log"
	"net/http"

	"os"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/gin-gonic/gin"
)

func Authorize() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenKey := []byte(os.Getenv("TOKEN_KEY"))
		var keyFunc jwt.Keyfunc = func(t *jwt.Token) (interface{}, error) { return tokenKey, nil }
		// Get token from request
		token, err := request.ParseFromRequest(c.Request, request.OAuth2Extractor, keyFunc)

		if err != nil {
			switch err.(type) {
			case *jwt.ValidationError: // JWT validation error
				vErr := err.(*jwt.ValidationError)
				switch vErr.Errors {
				case jwt.ValidationErrorExpired: //JWT expired
					log.Println("[Authorize]- Token Expired, get a new one")
				default:
					log.Printf("[Authorize]- ValidationError error: %+v\n", err)
				}
			default:
				log.Printf("[Authorize]- Token parse error: %v\n", err)
			}
			log.Println("[Authorize]- Token Expired, get a new one")
			GinAbort(c, http.StatusUnauthorized, TOKEN_EXPIRED, "")
			return
		}
		if token.Valid {
			// Set claimInfo to conext for using in backward router
			claims := token.Claims.(jwt.MapClaims)
			c.Set("claimInfo", claims)
			c.Next()
		} else {
			log.Printf("[Authorize]- Token is invalid")
			GinAbort(c, http.StatusUnauthorized, TOKEN_INVALID, "")
		}
	}
}
