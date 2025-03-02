// main.go
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

func main() {
    http.HandleFunc("/products", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(products)
    })

    log.Printf("Server starting on port 8080...")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatal(err)
    }
}
