package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

type Product struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

func main() {
	// Verbindung zur Datenbank
	db, _ := sql.Open("postgres", "postgres://postgres:password@db:5432/shopdb?sslmode=disable")

	// API-Endpunkt für Produkte
	http.HandleFunc("/products", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Produkte aus der Datenbank lesen
		rows, _ := db.Query("SELECT id, name, price FROM products")

		var products []Product
		for rows.Next() {
			var p Product
			rows.Scan(&p.ID, &p.Name, &p.Price)
			products = append(products, p)
		}

		json.NewEncoder(w).Encode(products)
	})

	log.Printf("Server starting on port 8080...")
	http.ListenAndServe(":8080", nil)
}
