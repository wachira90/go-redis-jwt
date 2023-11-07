package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"context"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
)

var redisClient *redis.Client

func main() {
	// Initialize the Redis client
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Replace with your Redis server address
		Password: "",               // No password by default
		DB:       0,                // Default DB
	})

	// Initialize the router
	r := mux.NewRouter()

	// Define CRUD endpoints
	r.HandleFunc("/create", createHandler).Methods("POST")
	r.HandleFunc("/read/{key}", readHandler).Methods("GET")
	r.HandleFunc("/update/{key}", updateHandler).Methods("PUT")
	r.HandleFunc("/delete/{key}", deleteHandler).Methods("DELETE")

	fmt.Println("Server is running on :8080")
	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}

func createHandler(w http.ResponseWriter, r *http.Request) {
	key := mux.Vars(r)["key"]
	value := r.FormValue("value")

	err := redisClient.Set(context.Background(), key, value, 0).Err()
	if err != nil {
		http.Error(w, "Error creating key-value pair", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Key-value pair created successfully"))
}

func readHandler(w http.ResponseWriter, r *http.Request) {
	key := mux.Vars(r)["key"]

	value, err := redisClient.Get(context.Background(), key).Result()
	if err == redis.Nil {
		http.Error(w, "Key not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Error reading key-value pair", http.StatusInternalServerError)
		return
	}

	data := map[string]string{
		"key":   key,
		"value": value,
	}

	responseJSON, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJSON)
}

func updateHandler(w http.ResponseWriter, r *http.Request) {
	key := mux.Vars(r)["key"]
	value := r.FormValue("value")

	err := redisClient.Set(context.Background(), key, value, 0).Err()
	if err != nil {
		http.Error(w, "Error updating key-value pair", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Key-value pair updated successfully"))
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	key := mux.Vars(r)["key"]

	err := redisClient.Del(context.Background(), key).Err()
	if err != nil {
		http.Error(w, "Error deleting key-value pair", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Key-value pair deleted successfully"))
}
