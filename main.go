package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Artikel struct {
	ArtikelID  int       `json:"artikelid"`
	Judul      string    `json:"judul"`
	Gambar     string    `json:"gambar"`
	Waktu      time.Time `json:"waktu"`
	KategoriID int       `json:"kategoriid"`
}

func main() {
	// Koneksi database
	connStr := "postgres://postgres.jfugaikxhuxsryzqpres:mtlfztox1987@aws-1-ap-southeast-1.pooler.supabase.com:5432/postgres"
	pool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		log.Fatalf("Gagal konek database: %v", err)
	}
	defer pool.Close()

	r := gin.Default()

	// Secret key JWT
	AccessToken := []byte("AccessToken")
	RefreshToken := []byte("RefreshToken")

	// Middleware untuk cek token
	Middleware := func(c *gin.Context) {
		Authoriz := c.GetHeader("Authorization")
		if Authoriz == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token tidak ada"})
			c.Abort()
			return
		}
		auth, err := jwt.Parse(Authoriz, func(token *jwt.Token) (interface{}, error) {
			return AccessToken, nil
		})
		if err != nil || !auth.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token tidak valid"})
			c.Abort()
			return

		}

	}

	// LOGIN
	r.POST("/login", func(c *gin.Context) {
		var req struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

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
		AccessClaims := jwt.MapClaims{
			"username": req.Username,
			"expired":  time.Now().Add(15 * time.Minute).Unix(),
		}
		access := jwt.NewWithClaims(jwt.SigningMethodHS256, AccessClaims)
		at, _ := access.SignedString(AccessToken)

		// Refresh token 7 hari
		RefreshClaims := jwt.MapClaims{
			"username": req.Username,
			"expired":  time.Now().Add(7 * 24 * time.Hour).Unix(),
		}
		refresh := jwt.NewWithClaims(jwt.SigningMethodHS256, RefreshClaims)
		rt, _ := refresh.SignedString(RefreshToken)

		c.JSON(http.StatusOK, gin.H{
			"access_token":  at,
			"refresh_token": rt,
		})
	})

	// REFRESH
	r.POST("/refresh", func(c *gin.Context) {
		var req struct {
			RefreshToken1 string `json:"refresh_token"`
		}

		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Request tidak valid"})
			return
		}

		token, err := jwt.Parse(req.RefreshToken1, func(token *jwt.Token) (interface{}, error) {
			return RefreshToken, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token tidak valid"})
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		username := claims["username"].(string)

		// Buat token baru
		AccessClaims := jwt.MapClaims{
			"username": username,
			"expired":  time.Now().Add(15 * time.Minute).Unix(),
		}
		accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, AccessClaims)
		at, _ := accessToken.SignedString(AccessToken)

		RefreshClaims := jwt.MapClaims{
			"username": username,
			"expired":  time.Now().Add(7 * 24 * time.Hour).Unix(),
		}
		refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, RefreshClaims)
		rt, _ := refreshToken.SignedString(RefreshToken)

		c.JSON(http.StatusOK, gin.H{
			"access_token":  at,
			"refresh_token": rt,
		})
	})
	// NEWS (protected)
	r.GET("/news", Middleware, func(c *gin.Context) {
		DaftarKategori := c.DefaultQuery("kategori", "")
		LimitHalaman := c.DefaultQuery("limit", "")
		SortWaktu := c.DefaultQuery("sort", "")
		var rows pgx.Rows

		if DaftarKategori != "" && LimitHalaman != "" {
			rows, err = pool.Query(context.Background(), fmt.Sprintf(`
				SELECT a.artikelid, a.judul, a.gambar, a.waktu, a.kategoriid
				FROM artikel a
				JOIN kategori k ON a.kategoriid = k.kategoriid
				WHERE k.kategorinama = $1
				ORDER BY a.waktu %s
				LIMIT $2`, SortWaktu), DaftarKategori, LimitHalaman)
		} else if DaftarKategori != "" {
			rows, err = pool.Query(context.Background(), fmt.Sprintf(`
				SELECT a.artikelid, a.judul, a.gambar, a.waktu, a.kategoriid
				FROM artikel a
				JOIN kategori k ON a.kategoriid = k.kategoriid
				WHERE k.kategorinama = $1
				ORDER BY a.waktu %s`, SortWaktu), DaftarKategori)
		} else if LimitHalaman != "" {
			rows, err = pool.Query(context.Background(), fmt.Sprintf(`
				SELECT a.artikelid, a.judul, a.gambar, a.waktu, a.kategoriid
				FROM artikel a
				ORDER BY a.waktu %s
				LIMIT $1`, SortWaktu), LimitHalaman)
		} else {
			rows, err = pool.Query(context.Background(), fmt.Sprintf(`
				SELECT a.artikelid, a.judul, a.gambar, a.waktu, a.kategoriid
				FROM artikel a
				ORDER BY a.waktu %s`, SortWaktu))
		}

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Query gagal"})
			return
		}
		defer rows.Close()

		var articles []Artikel
		for rows.Next() {
			var article Artikel
			if err := rows.Scan(&article.ArtikelID, &article.Judul, &article.Gambar, &article.Waktu, &article.KategoriID); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal proses data"})
				return
			}
			articles = append(articles, article)
		}

		c.JSON(http.StatusOK, articles)
	})

	r.Run(":8080")
}
