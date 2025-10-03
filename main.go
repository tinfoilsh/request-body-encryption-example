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
			allowHeaders  = "Accept, Authorization, Content-Type, Ehbp-Client-Public-Key, Ehbp-Encapsulated-Key"
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

		if ct := resp.Header.Get("Content-Type"); ct != "" {
			w.Header().Set("Content-Type", ct)
		}
		if te := resp.Header.Get("Transfer-Encoding"); te != "" {
			w.Header().Set("Transfer-Encoding", te)
			w.Header().Del("Content-Length")
		}

		w.WriteHeader(resp.StatusCode)

		if flusher, ok := w.(http.Flusher); ok {
			buf := make([]byte, 32*1024)
			for {
				n, readErr := resp.Body.Read(buf)
				if n > 0 {
					if _, writeErr := w.Write(buf[:n]); writeErr != nil {
						log.Printf("proxy write error: %v", writeErr)
						break
					}
					flusher.Flush()
				}
				if readErr != nil {
					if readErr != io.EOF {
						log.Printf("proxy read error: %v", readErr)
					}
					break
				}
			}
		} else {
			io.Copy(w, resp.Body)
		}
	})

	log.Println("Server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
