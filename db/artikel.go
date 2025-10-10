package db

import "time"

type Artikel struct {
	ArtikelID  int       `json:"artikel_id"`
	Judul      string    `json:"judul"`
	Gambar     string    `json:"gambar"`
	Waktu      time.Time `json:"waktu"`
	KategoriID int       `json:"kategori_id"`
	NamaNegara string    `json:"-"`
}
