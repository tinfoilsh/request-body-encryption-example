package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

const (
	upstreamURL   = "https://ehbp.inf6.tinfoil.sh/v1/chat/completions"
	allowHeaders  = "Accept, Authorization, Content-Type, Ehbp-Client-Public-Key, Ehbp-Encapsulated-Key"
	exposeHeaders = "Ehbp-Encapsulated-Key, Ehbp-Client-Public-Key, Ehbp-Fallback"
)

var (
	tinfoilRequestHeaders  = []string{"Ehbp-Client-Public-Key", "Ehbp-Encapsulated-Key"}
	tinfoilResponseHeaders = []string{"Ehbp-Encapsulated-Key", "Ehbp-Client-Public-Key", "Ehbp-Fallback"}
)

func main() {
	http.HandleFunc("/v1/chat/completions", proxyHandler)

	log.Println("proxy listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func proxyHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", allowHeaders)
	w.Header().Set("Access-Control-Expose-Headers", exposeHeaders)

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	req, err := http.NewRequestWithContext(r.Context(), http.MethodPost, upstreamURL, r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	if accept := r.Header.Get("Accept"); accept != "" {
		req.Header.Set("Accept", accept)
	}

	apiKey := os.Getenv("TINFOIL_API_KEY")
	if apiKey == "" {
		http.Error(w, "TINFOIL_API_KEY not set", http.StatusInternalServerError)
		return
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)

	copyHeaders(req.Header, r.Header, tinfoilRequestHeaders...)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	copyHeaders(w.Header(), resp.Header, tinfoilResponseHeaders...)
	if ct := resp.Header.Get("Content-Type"); ct != "" {
		w.Header().Set("Content-Type", ct)
	}
	if te := resp.Header.Get("Transfer-Encoding"); te != "" {
		w.Header().Set("Transfer-Encoding", te)
		w.Header().Del("Content-Length")
	}

	w.WriteHeader(resp.StatusCode)

	if flusher, ok := w.(http.Flusher); ok {
		fw := flushWriter{ResponseWriter: w, Flusher: flusher}
		if _, copyErr := io.Copy(&fw, resp.Body); copyErr != nil {
			log.Printf("stream copy failed: %v", copyErr)
		}
		return
	}

	if _, copyErr := io.Copy(w, resp.Body); copyErr != nil {
		log.Printf("response copy failed: %v", copyErr)
	}
}

type flushWriter struct {
	http.ResponseWriter
	http.Flusher
}

func (fw *flushWriter) Write(p []byte) (int, error) {
	n, err := fw.ResponseWriter.Write(p)
	if fw.Flusher != nil {
		fw.Flush()
	}
	return n, err
}

func copyHeaders(dst, src http.Header, keys ...string) {
	for _, key := range keys {
		if value := src.Get(key); value != "" {
			dst.Set(key, value)
		}
	}
}
