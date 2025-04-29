package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
)

// Product repräsentiert ein Produkt im Shop
type Product struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

// Simulierte Produktdatenbank
var products = map[int]Product{
	1: {ID: 1, Name: "Laptop", Price: 999.99},
	2: {ID: 2, Name: "Smartphone", Price: 499.99},
	3: {ID: 3, Name: "Headphones", Price: 79.99},
}

var ctx = context.Background()
var rdb *redis.Client

func main() {
	// Redis-Client initialisieren
	rdb = redis.NewClient(&redis.Options{
		Addr:     "redis:6379", // Redis-Server im Docker-Netzwerk
		Password: "",           // kein Passwort
		DB:       0,            // Standard-DB
	})

	// Überprüfen, ob Redis erreichbar ist
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Konnte keine Verbindung zu Redis herstellen: %v", err)
	}

	// Router einrichten
	r := mux.NewRouter()
	r.HandleFunc("/api/products/{id}", getProduct).Methods("GET")
	r.HandleFunc("/api/cache/product/{id}", cacheProduct).Methods("POST")
	r.HandleFunc("/api/products/{id}", updateProduct).Methods("PUT")

	// Server starten
	fmt.Println("Server läuft auf Port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

// getProduct gibt ein Produkt zurück, prüft zuerst im Cache
func getProduct(w http.ResponseWriter, r *http.Request) {
	// ID aus URL-Parameter extrahieren
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Ungültige Produkt-ID", http.StatusBadRequest)
		return
	}

	// Versuchen, das Produkt aus dem Cache zu holen
	cacheKey := fmt.Sprintf("product:%d", id)
	cachedProduct, err := rdb.Get(ctx, cacheKey).Result()

	if err == nil {
		// Cache-Treffer: Produkt sofort zurückgeben
		fmt.Printf("Cache HIT für Produkt %d\n", id)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(cachedProduct))
		return
	}

	// Cache-Fehltreffer: Aus der "Datenbank" holen (mit simulierter Verzögerung)
	fmt.Printf("Cache MISS für Produkt %d\n", id)
	product, exists := products[id]
	if !exists {
		http.Error(w, "Produkt nicht gefunden", http.StatusNotFound)
		return
	}

	// Datenbanklatenz simulieren (2 Sekunden)
	time.Sleep(2 * time.Second)

	// Produkt als JSON zurückgeben
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}

// cacheProduct speichert ein Produkt im Cache
func cacheProduct(w http.ResponseWriter, r *http.Request) {
	// ID aus URL-Parameter extrahieren
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Ungültige Produkt-ID", http.StatusBadRequest)
		return
	}

	// Prüfen, ob das Produkt existiert
	product, exists := products[id]
	if !exists {
		http.Error(w, "Produkt nicht gefunden", http.StatusNotFound)
		return
	}

	// Produkt als JSON konvertieren
	productJSON, err := json.Marshal(product)
	if err != nil {
		http.Error(w, "Fehler bei der JSON-Konvertierung", http.StatusInternalServerError)
		return
	}

	// Im Redis-Cache speichern
	cacheKey := fmt.Sprintf("product:%d", id)
	err = rdb.Set(ctx, cacheKey, productJSON, 0).Err() // 0 = kein Ablaufdatum
	if err != nil {
		http.Error(w, "Fehler beim Caching", http.StatusInternalServerError)
		return
	}

	// Erfolgsmeldung zurückgeben
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": fmt.Sprintf("Produkt %d wurde im Cache gespeichert", id),
	})
}

// updateProduct aktualisiert ein Produkt und invalidiert den Cache
func updateProduct(w http.ResponseWriter, r *http.Request) {
	// ID aus URL-Parameter extrahieren
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Ungültige Produkt-ID", http.StatusBadRequest)
		return
	}

	// Prüfen, ob das Produkt existiert
	_, exists := products[id]
	if !exists {
		http.Error(w, "Produkt nicht gefunden", http.StatusNotFound)
		return
	}

	// Anfragekörper entschlüsseln
	var updatedProduct Product
	err = json.NewDecoder(r.Body).Decode(&updatedProduct)
	if err != nil {
		http.Error(w, "Ungültiger Anfragekörper", http.StatusBadRequest)
		return
	}

	// ID beibehalten
	updatedProduct.ID = id

	// Produkt in der "Datenbank" aktualisieren
	products[id] = updatedProduct

	// Cache für dieses Produkt invalidieren
	cacheKey := fmt.Sprintf("product:%d", id)
	err = rdb.Del(ctx, cacheKey).Err()
	if err != nil {
		http.Error(w, "Fehler beim Invalidieren des Cache", http.StatusInternalServerError)
		return
	}

	// Erfolgsmeldung zurückgeben
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": fmt.Sprintf("Produkt %d aktualisiert und Cache invalidiert", id),
	})
}
