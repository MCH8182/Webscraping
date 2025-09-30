package controller

func StartUserRouter(engine *gin.Engine) {
	r := engine.Group("/users")
	r.GET("/login", loginUserHandler)
	// TODO: refresh token endpoint
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
