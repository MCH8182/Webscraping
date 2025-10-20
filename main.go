package main

import (
	"goquery-example/api"
	"goquery-example/db"
	"goquery-example/scraper"
	"log"

	"github.com/jasonlvhit/gocron"
)

func main() {
	// Koneksi database
	err := db.StartConnection()
	if err != nil {
		log.Fatalf("Gagal konek database: %v\n", err)
	}
	defer db.Pool.Close()

	// Scheduler buat scraping (sehari 1x)
	gocron.Every(1).Day().At("08:00").Do(scraper.ScrapeNews)
	go gocron.Start()
	// Host REST API
	api.RunAPI()
	select {}
}
