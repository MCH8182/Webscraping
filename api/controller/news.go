package controller

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx"
)

// TODO: startNewsRouter
// TODO: handler'an
// NEWS (protected)

func StartNewsRouter(engine *gin.Engine) {
	engine.GET("/news", GetNews)
}

func GetNews(c *gin.Context) {
	DaftarKategori := c.DefaultQuery("kategori", "")
	LimitHalaman := c.DefaultQuery("limit", "")
	SortWaktu := c.DefaultQuery("sort", "")
	var rows pgx.Rows
	var err error

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

}
