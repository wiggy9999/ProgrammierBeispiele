package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type Product struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

var products = []Product{
	{ID: "1", Name: "Laptop", Price: 999.99},
	{ID: "2", Name: "Mouse", Price: 24.99},
}

var cart []string // In-Memory-Zustand

func main() {
	http.HandleFunc("/products", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(products)
	})

	http.HandleFunc("/cart/add", func(w http.ResponseWriter, r *http.Request) {
		productID := r.URL.Query().Get("productId")
		cart = append(cart, productID)
		w.Write([]byte("Product added to cart\n"))
	})

	http.HandleFunc("/cart/show", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(cart)
	})

	log.Printf("Stateful server starting on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
