package middleware

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"ta_microservice_peminjaman/app/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

var SECRET = []byte(os.Getenv("SECRET.KEY"))

func DecodeToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", t.Header["alg"])
		}
		return SECRET, nil
	})

	claims, oke := token.Claims.(jwt.MapClaims)
	if oke && token.Valid {
		fmt.Println("token_type and exp :", claims["token_type"], claims["exp"])
	} else {
		log.Println("this not valid")
		log.Println(err)
		return nil, err
	}
	return claims, nil
}

func Jwt() gin.HandlerFunc {
	return func(c *gin.Context) {
		res := models.JsonResponse{Success: true}

		if len(strings.Split(c.Request.Header.Get("Authorization"), " ")) != 2 {
			errorMsg := "invalid token"
			res.Success = false
			res.Error = &errorMsg
			c.JSON(http.StatusUnauthorized, res)
			c.Abort()
			return
		}

		accessToken := strings.Split(c.Request.Header.Get("Authorization"), " ")[1]

		claims, err := DecodeToken(accessToken)
		if err != nil {
			errMsg := err.Error()
			res.Success = false
			res.Error = &errMsg
			c.JSON(401, res)
			c.Abort()
			return
		}

		_, found := claims["token_type"]
		if !found {
			errMsg := err.Error()
			res.Success = false
			res.Error = &errMsg
			c.JSON(http.StatusUnauthorized, res)
			c.Abort()
			return
		}

		if claims["token_type"] != "access_token" {
			errMsg := err.Error()
			res.Success = false
			res.Error = &errMsg
			c.JSON(401, res)
			c.Abort()
			return
		}

		c.Next()
	}
}
