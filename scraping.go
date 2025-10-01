package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Kategori struct {
	Nama     string
	Link     string
	Selector string
}

type Artikel struct {
	Judul      string
	Gambar     string
	Waktu      time.Time
	NamaNegara string
}

func main() {
	// Konek ke Database
	connStr := "postgres://postgres.jfugaikxhuxsryzqpres:mtlfztox1987@aws-1-ap-southeast-1.pooler.supabase.com:5432/postgres"
	pool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer pool.Close()

	//Isi dari Kategori Struct dalam bentuk slice
	Isikategori := []Kategori{
		{"Asia", "https://www.theguardian.com/world/asia", "#container-asia > ul > li"},
		{"America", "https://www.theguardian.com/us-news", "#container-us-latest-news > ul > li"},
		{"Europe", "https://www.theguardian.com/world/europe-news", "#container-latest-news > ul > li"},
		{"Australia", "https://www.theguardian.com/australia-news", "#container-australia-news > ul > li"},
	}

	//Looping dari Kategori agar bisa loop terhadap 4 negara
	for _, isinegara := range Isikategori {
		fmt.Printf("Ambil data dari artikel %s \n", isinegara.Nama)

		//make HTTP request
		res, err := http.Get(isinegara.Link)
		if err != nil {
			log.Fatalf("Gagal melakukan koneksi HTTP %s : %v\n", isinegara.Link, err)
		}
		defer res.Body.Close()

		//Cek Respon ke HTTP
		if res.StatusCode != 200 {
			log.Fatalf("Status code error: %d : %s\n", res.StatusCode, isinegara.Link)
		}

		//parse HTML
		doc, err := goquery.NewDocumentFromReader(res.Body)
		if err != nil {
			log.Fatalf("Gagal parse koneksi HTTP %s : %v\n", isinegara.Link, err)
		}

		var isian []Artikel

		//find data from HTML
		doc.Find(isinegara.Selector).Each(func(i int, s *goquery.Selection) {
			// ambil title dari tiap artikel
			title := s.Find("h3 > span").Text()

			// ambil gambar dari tiap artikel
			img := s.Find("picture > img")
			image, _ := img.Attr("src")

			// ambil datetime dari tiap artikel
			var clockfinal time.Time
			datetime := s.Find("footer > span > gu-island > time")
			clock, _ := datetime.Attr("datetime")
			if len(clock) > 0 {
				clockfinal, _ = time.Parse(time.RFC3339, clock)
			} else {
				clockfinal = time.Now()
			}

			//Simpan data ke Artikel struct
			Kontenartikel := Artikel{
				Judul:      title,
				Gambar:     image,
				Waktu:      clockfinal,
				NamaNegara: isinegara.Nama,
			}

			isian = append(isian, Kontenartikel)
		})

		//Lopping isian untuk simpen ke database
		for _, isi := range isian {
			// ambil kategoriid untuk ke database
			var iduntukkategori int
			err := pool.QueryRow(context.Background(), "SELECT kategoriid FROM kategori WHERE kategorinama = $1", isi.NamaNegara).Scan(&iduntukkategori)
			if err != nil {
				log.Fatalf("Gagal query kategori id %s : %v\n", isi.NamaNegara, err)
			}

			//Masukkan data ke database
			_, err = pool.Exec(context.Background(),
				"INSERT INTO artikel (judul, waktu, gambar, kategoriid) VALUES ($1, $2, $3, $4)",
				isi.Judul, isi.Waktu, isi.Gambar, iduntukkategori)
			if err != nil {
				log.Fatalf("Gagal menambahkan judul artikel %s: %v\n", isi.Judul, err)

			}
			// Print berhasil menyimpan artikel
			fmt.Printf("Berhasil menyimpan artikel %s\n", isi.Judul)

		}

	}

}
