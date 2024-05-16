package interceptor

import (
	"fmt"
	"goworkshop/model"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var secretKey = "87654321"

func JwtSign(payload model.User) string {

	atClaims := jwt.MapClaims{}

	// Payload begin
	atClaims["id"] = payload.ID
	atClaims["username"] = payload.Username
	atClaims["level"] = payload.Level
	atClaims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	// Payload end
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, _ := at.SignedString([]byte(secretKey))
	return token

}

func JwtVerify(c *gin.Context) {
	tokenString := strings.Split(c.Request.Header["Authorization"][0], " ")[1]
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method : %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})

	if claim, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		c.Set("jwt_username", claim["username"])
		c.Next()
	} else {
		c.JSON(http.StatusOK, gin.H{"result": "nok", "message": "invalid token", "error": err})
		c.Abort()
	}
}
