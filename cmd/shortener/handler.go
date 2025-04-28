package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
)

var (
	urlStore = make(map[string]string)
)

func shortHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	body := make([]byte, r.ContentLength)
	if _, err := r.Body.Read(body); err != nil && err.Error() != "EOF" {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}

	originalURL := string(body)
	shortKey := generateShortKey()
	urlStore[shortKey] = originalURL
	shortenedURL := fmt.Sprintf("http://localhost:8080/%s", shortKey)

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(201)

	w.Write([]byte(shortenedURL))
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	originalURL, exists := urlStore[id]
	if !exists {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Location", originalURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func generateShortKey() string {
	b := make([]byte, 6)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return base64.URLEncoding.EncodeToString(b)[:8]
}
