package middleware

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func generateToken(length int) string {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	return base64.URLEncoding.EncodeToString(b)
}

var (
	AccessTokenPassword  = "YBtazYQBk1fKxqsDUYsYwZ6L1dUjpb8b8v9kz3fytTo=" // TODO: harusnya ini dibikin sesulit mungkin sih kayak password
	RefreshTokenPassword = "jmsZlCYzQiIXOo9W2Uoz0GdzyBgRbacx_Ip21XgY69k="
	AccessTokenType      = 0
	RefreshTokenType     = 1
)

func GenerateJWT(tokenType int, callback func(jwt.MapClaims)) (string, error) {
	if tokenType != AccessTokenType && tokenType != RefreshTokenType {
		return "", fmt.Errorf("invalid token type")
	}

	jwtClaims := jwt.MapClaims{}
	callback(jwtClaims)

	jwtClaims["gt"] = generateToken(32)

	if tokenType == AccessTokenType {
		jwtClaims["expired"] = time.Now().Add(15 * time.Minute).Unix()
		accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims)
		accessTokenString, err := accessToken.SignedString([]byte(AccessTokenPassword))
		if err != nil {
			log.Println("Error signing access token: ", err)
			return "", err
		}
		return accessTokenString, nil
	}

	jwtClaims["expired"] = time.Now().Add(7 * 24 * time.Hour).Unix()
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(RefreshTokenPassword))
	if err != nil {
		log.Println("Error signing refresh token: ", err)
		return "", err
	}
	return refreshTokenString, nil
}

// Middleware untuk cek token
func VerifyJWT(c *gin.Context) {
	authoriz := c.GetHeader("Authorization")
	if authoriz == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token tidak ada"})
		c.Abort()
		return
	}

	if authoriz[:7] != "Bearer " {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Format token salah"})
		c.Abort()
		return
	}

	accessToken := authoriz[7:]

	auth, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(AccessTokenPassword), nil
	})
	if err != nil || !auth.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token tidak valid"})
		c.Abort()
		return
	}
	c.Next()
}
