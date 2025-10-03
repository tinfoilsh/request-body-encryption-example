package main

import (
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/v1/chat/completions", func(w http.ResponseWriter, r *http.Request) {
		const (
			allowHeaders  = "Content-Type, Ehbp-Client-Public-Key, Ehbp-Encapsulated-Key"
			exposeHeaders = "Ehbp-Encapsulated-Key, Ehbp-Client-Public-Key, Ehbp-Fallback"
		)

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", allowHeaders)
		w.Header().Set("Access-Control-Expose-Headers", exposeHeaders)

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// Just add API key and proxy everything
		req, _ := http.NewRequest("POST", "https://ehbp.inf6.tinfoil.sh/v1/chat/completions", r.Body)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+os.Getenv("TINFOIL_API_KEY"))

		// Forward EHBP headers using canonical casing so Go preserves them
		if clientPubKey := r.Header.Get("Ehbp-Client-Public-Key"); clientPubKey != "" {
			req.Header.Set("Ehbp-Client-Public-Key", clientPubKey)
		}
		if encapsulatedKey := r.Header.Get("Ehbp-Encapsulated-Key"); encapsulatedKey != "" {
			req.Header.Set("Ehbp-Encapsulated-Key", encapsulatedKey)
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		defer resp.Body.Close()

		for _, header := range []string{
			"Ehbp-Encapsulated-Key",
			"Ehbp-Client-Public-Key",
			"Ehbp-Fallback",
		} {
			if value := resp.Header.Get(header); value != "" {
				w.Header().Set(header, value)
			}
		}

		w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
		w.WriteHeader(resp.StatusCode)
		io.Copy(w, resp.Body)
	})

	log.Println("Server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
