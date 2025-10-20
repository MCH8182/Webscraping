package controller

import (
	"fmt"
	"goquery-example/api/middleware"
	"goquery-example/db"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

func StartNewsRouter(engine *gin.Engine) {
	// TODO tambah groups
	r := engine.Group("/news")
	r.GET("", middleware.VerifyJWT, getNews)
	// engine.GET("/news", getNews)
}

func getNews(c *gin.Context) {
	DaftarKategori := c.DefaultQuery("kategori", "")
	LimitHalaman := c.DefaultQuery("limit", "")
	SortWaktu := c.DefaultQuery("sort", "")
	var rows pgx.Rows
	var err error

	if DaftarKategori != "" && LimitHalaman != "" {
		rows, err = db.Pool.Query(c.Request.Context(), fmt.Sprintf(`
				SELECT a.artikelid, a.judul, a.gambar, a.waktu, a.kategoriid
				FROM artikel a
				JOIN kategori k ON a.kategoriid = k.kategoriid
				WHERE k.kategorinama = $1
				ORDER BY a.waktu %s
				LIMIT $2`, SortWaktu), DaftarKategori, LimitHalaman)
	} else if DaftarKategori != "" {
		rows, err = db.Pool.Query(c.Request.Context(), fmt.Sprintf(`
				SELECT a.artikelid, a.judul, a.gambar, a.waktu, a.kategoriid
				FROM artikel a
				JOIN kategori k ON a.kategoriid = k.kategoriid
				WHERE k.kategorinama = $1
				ORDER BY a.waktu %s`, SortWaktu), DaftarKategori)
	} else if LimitHalaman != "" {
		rows, err = db.Pool.Query(c.Request.Context(), fmt.Sprintf(`
				SELECT a.artikelid, a.judul, a.gambar, a.waktu, a.kategoriid
				FROM artikel a
				ORDER BY a.waktu %s
				LIMIT $1`, SortWaktu), LimitHalaman)
	} else {
		rows, err = db.Pool.Query(c.Request.Context(), fmt.Sprintf(`
				SELECT a.artikelid, a.judul, a.gambar, a.waktu, a.kategoriid
				FROM artikel a
				ORDER BY a.waktu %s`, SortWaktu))
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Query gagal"})
		return
	}
	defer rows.Close()

	var articles []db.Artikel
	for rows.Next() {
		var article db.Artikel
		if err := rows.Scan(&article.ArtikelID, &article.Judul, &article.Gambar, &article.Waktu, &article.KategoriID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal proses data"})
			return
		}
		articles = append(articles, article)
	}

	c.JSON(http.StatusOK, articles)

}
