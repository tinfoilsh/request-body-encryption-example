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

type chatRequest struct {
	Model    string `json:"model"`
	Messages []struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"messages"`
}

func main() {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Fatal("OPENAI_API_KEY environment variable is not set")
	}

	client := openai.NewClient(option.WithAPIKey(apiKey))

	http.HandleFunc("/api/chat", chatHandler(&client))

	log.Println("Server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func chatHandler(client *openai.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			return
		}
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var in chatRequest
		if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}
		if in.Model == "" {
			in.Model = "gpt-3.5-turbo"
		}

		msgs := make([]openai.ChatCompletionMessageParamUnion, 0, len(in.Messages))
		for _, msg := range in.Messages {
			switch msg.Role {
			case "user":
				msgs = append(msgs, openai.UserMessage(msg.Content))
			case "system":
				msgs = append(msgs, openai.SystemMessage(msg.Content))
			case "assistant":
				msgs = append(msgs, openai.AssistantMessage(msg.Content))
			case "":
				http.Error(w, "message role is required", http.StatusBadRequest)
				return
			default:
				http.Error(w, "unsupported message role: "+msg.Role, http.StatusBadRequest)
				return
			}
		}

		stream := client.Chat.Completions.NewStreaming(r.Context(), openai.ChatCompletionNewParams{
			Model:    in.Model,
			Messages: msgs,
		})

		wantStream := strings.Contains(r.Header.Get("Accept"), "text/event-stream")
		var flusher http.Flusher
		if wantStream {
			w.Header().Set("Content-Type", "text/event-stream")
			w.Header().Set("Cache-Control", "no-cache")
			w.Header().Set("Connection", "keep-alive")
			var ok bool
			flusher, ok = w.(http.Flusher)
			if !ok {
				http.Error(w, "streaming not supported", http.StatusInternalServerError)
				return
			}
		}

		var content strings.Builder
		for stream.Next() {
			chunk := stream.Current()
			if len(chunk.Choices) == 0 {
				continue
			}
			delta := chunk.Choices[0].Delta.Content
			if delta == "" {
				continue
			}
			content.WriteString(delta)
			if wantStream {
				payload := strings.ReplaceAll(delta, "\n", "\ndata: ")
				fmt.Fprintf(w, "data: %s\n\n", payload)
				flusher.Flush()
			}
		}
		if err := stream.Err(); err != nil {
			log.Printf("stream error: %v", err)
			http.Error(w, "failed to create chat completion", http.StatusInternalServerError)
			return
		}

		if wantStream {
			fmt.Fprint(w, "data: [DONE]\n\n")
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"choices": []map[string]any{
				{
					"message": map[string]string{
						"role":    "assistant",
						"content": content.String(),
					},
				},
			},
		})
	}
}
