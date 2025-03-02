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

type CartRequest struct {
	Cart      []string `json:"cart"`
	ProductID string   `json:"productId,omitempty"`
}

var products = []Product{
	{ID: "1", Name: "Laptop", Price: 999.99},
	{ID: "2", Name: "Mouse", Price: 24.99},
}

func main() {
	http.HandleFunc("/products", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(products)
	})

	http.HandleFunc("/v2/cart/add", func(w http.ResponseWriter, r *http.Request) {
		var req CartRequest
		json.NewDecoder(r.Body).Decode(&req)
		req.Cart = append(req.Cart, req.ProductID)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(req)
	})

	http.HandleFunc("/v2/cart/show", func(w http.ResponseWriter, r *http.Request) {
		var req CartRequest
		json.NewDecoder(r.Body).Decode(&req)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(req)
	})

	log.Printf("Stateless server starting on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
