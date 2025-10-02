package main

import (
	"io"
	"log"
	"net/http"
	"os"

	"github.com/tinfoilsh/encrypted-http-body-protocol/identity"
)

const (
	defaultEnclaveURL = "https://ehbp.inf6.tinfoil.sh/v1/chat/completions"
	keysEndpoint      = "/.well-known/hpke-keys"
)

func main() {
	serverIdentity, err := identity.FromFile("server_identity.json")
	if err != nil {
		log.Fatalf("failed to load HPKE identity: %v", err)
	}

	http.HandleFunc(keysEndpoint, func(w http.ResponseWriter, r *http.Request) {
		addCORSHeaders(w.Header())
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		if r.Method != http.MethodGet {
			http.NotFound(w, r)
			return
		}

		serverIdentity.ConfigHandler(w, r)
	})

	enclaveURL := os.Getenv("ENCLAVE_CHAT_URL")
	if enclaveURL == "" {
		enclaveURL = defaultEnclaveURL
	}

	apiKey := os.Getenv("TINFOIL_API_KEY")
	if apiKey == "" {
		apiKey = os.Getenv("OPENAI_API_KEY")
	}

	chatHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req, err := http.NewRequestWithContext(r.Context(), http.MethodPost, enclaveURL, r.Body)
		if err != nil {
			http.Error(w, "failed to create upstream request", http.StatusInternalServerError)
			return
		}

		req.Header.Set("Content-Type", "application/json")
		if apiKey != "" {
			req.Header.Set("Authorization", "Bearer "+apiKey)
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resp.StatusCode)
		if _, copyErr := io.Copy(w, resp.Body); copyErr != nil {
			log.Printf("write response: %v", copyErr)
		}
	})

	securedChat := serverIdentity.Middleware(false)(chatHandler)

	http.HandleFunc("/v1/chat/completions", func(w http.ResponseWriter, r *http.Request) {
		addCORSHeaders(w.Header())
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		securedChat.ServeHTTP(w, r)
	})

	log.Println("Server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func addCORSHeaders(h http.Header) {
	h.Set("Access-Control-Allow-Origin", "*")
	h.Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	h.Set(
		"Access-Control-Allow-Headers",
		"Content-Type, Ehbp-Client-Public-Key, Ehbp-Encapsulated-Key",
	)
	h.Set("Access-Control-Expose-Headers", "Ehbp-Encapsulated-Key")
}
