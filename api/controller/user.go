package controller

import (
	"goquery-example/api/middleware"
	"goquery-example/api/schema"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func StartUserRouter(engine *gin.Engine) {
	r := engine.Group("/users")
	r.GET("/login", loginUserHandler)
	r.POST("/refresh", refreshTokenHandler)
}

// REFRESH
func refreshTokenHandler(c *gin.Context) {
	var req schema.RefreshTokenRequest

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Request tidak valid"})
		return
	}

	token, err := jwt.Parse(req.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		return req.RefreshToken, nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token tidak valid"})
		return
	}

	claims := token.Claims.(jwt.MapClaims)
	username := claims["username"].(string)

	// Buat token baru
	accessToken, err := middleware.GenerateJWT(middleware.AccessTokenType, func(a jwt.MapClaims) {
		a["username"] = username
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error when generating access token"})
		return
	}

	resp := schema.RefreshTokenResponse{
		AccessToken: accessToken,
	}
	c.JSON(http.StatusOK, resp)
}

func loginUserHandler(c *gin.Context) {
	var req schema.LoginUserRequest

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Request tidak valid"})
		return
	}

	// Dummy user: admin/password1
	if req.Username != "admin" || req.Password != "password1" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Username/password salah"})
		return
	}

	// Access token 15 menit
	accessToken, err := middleware.GenerateJWT(middleware.AccessTokenType, func(a jwt.MapClaims) {
		a["username"] = req.Username
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error when generating access token"})
		return
	}

	// Refresh token 7 hari
	refreshToken, err := middleware.GenerateJWT(middleware.RefreshTokenType, func(a jwt.MapClaims) {
		a["username"] = req.Username
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error when generating refresh token"})
		return
	}

	resp := schema.LoginUserResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	c.JSON(http.StatusOK, resp)
}
