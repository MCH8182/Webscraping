package middleware

var (
	AccessToken      = []byte("AccessToken") // TODO: harusnya ini dibikin sesulit mungkin sih kayak password
	RefreshToken     = []byte("RefreshToken")
	AccessTokenType  = 0
	RefreshTokenType = 1
)

func GenerateJWT(tokenType int, callback func(jwt.MapClaims)) (string, error) {
	if tokenType != AccessTokenType && tokenType != RefreshTokenType {
		return "", fmt.Errorf("invalid token type")
	}

	jwtClaims := jwt.MapClaims{}
	callback(jwtClaims)

	if tokenType == AccessTokenType {
		jwtClaims["expired"] = time.Now().Add(15 * time.Minute).Unix()
		accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, middleware.AccessToken)
		accessTokenString, err := accessToken.SignedString(middleware.AccessToken)
		if err != nil {
			return "", err
		}
		return accessTokenString, nil
	}

	jwtClaims["expired"] = time.Now().Add(7 * 24 * time.Hour).Unix()
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, middleware.RefreshToken)
	refreshTokenString, err := refreshToken.SignedString(middleware.RefreshToken)
	if err != nil {
		return "", err
	}
	return refreshTokenString, nil
}

// Middleware untuk cek token
func VerifyJWT(c *gin.Context) {
	Authoriz := c.GetHeader("Authorization")
	if Authoriz == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token tidak ada"})
		c.Abort()
		return
	}
	// TODO: harusnya cek juga formatnya "Bearer {access_token}"
	auth, err := jwt.Parse(Authoriz, func(token *jwt.Token) (interface{}, error) {
		return AccessToken, nil
	})
	if err != nil || !auth.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token tidak valid"})
		c.Abort()
		return
	}
	c.Next()
}
