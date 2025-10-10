package main

import (
	"goquery-example/api"
	"goquery-example/db"
	"log"
)

func main() {
	// Koneksi database
	err := db.StartConnection()
	if err != nil {
		log.Fatalf("Gagal konek database: %v\n", err)
	}
	defer db.Pool.Close()

	// Scheduler buat scraping (sehari 1x)

	// Host REST API
	api.RunAPI()
}
