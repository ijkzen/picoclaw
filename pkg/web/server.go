package web

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/sipeed/picoclaw/pkg/agent"
	"github.com/sipeed/picoclaw/pkg/bus"
	"github.com/sipeed/picoclaw/pkg/config"
)

type Server struct {
	server    *http.Server
	cfg       *config.Config
	agentLoop *agent.AgentLoop
	msgBus    *bus.MessageBus
}

func NewServer(host string, port int, cfg *config.Config, agentLoop *agent.AgentLoop, msgBus *bus.MessageBus) *Server {
	mux := http.NewServeMux()

	s := &Server{
		cfg:       cfg,
		agentLoop: agentLoop,
		msgBus:    msgBus,
	}

	mux.HandleFunc("/api/config", s.handleConfig)
	mux.HandleFunc("/api/models", s.handleModels)
	mux.HandleFunc("/api/models/default", s.handleDefaultModel)
	mux.HandleFunc("/api/chat", s.handleChat)
	mux.HandleFunc("/api/chat/stream", s.handleChatStream)
	mux.HandleFunc("/api/status", s.handleStatus)

	staticFS, err := fs.Sub(DistFS, "dist/browser")
	if err != nil {
		staticFS = nil
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api/") {
			http.NotFound(w, r)
			return
		}

		if staticFS != nil {
			filePath := strings.TrimPrefix(r.URL.Path, "/")
			if filePath == "" {
				filePath = "index.html"
			}

			file, err := staticFS.Open(filePath)
			if err != nil {
				// SPA routing: return index.html for any non-file path
				file, err = staticFS.Open("index.html")
				if err != nil {
					http.NotFound(w, r)
					return
				}
				defer file.Close()

				// Always return text/html for SPA routes
				w.Header().Set("Content-Type", "text/html")
				io.Copy(w, file)
				return
			}
			defer file.Close()

			// Set content type based on file extension for actual files
			ext := path.Ext(filePath)
			w.Header().Set("Content-Type", getContentType(ext))
			io.Copy(w, file)
			return
		}

		http.Redirect(w, r, "/api/status", http.StatusTemporaryRedirect)
	})

	addr := fmt.Sprintf("%s:%d", host, port+1)
	s.server = &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	return s
}

func (s *Server) Start() error {
	return s.server.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func (s *Server) handleConfig(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method == http.MethodGet {
		json.NewEncoder(w).Encode(s.cfg)
		return
	}

	if r.Method == http.MethodPost {
		var newConfig config.Config
		if err := json.NewDecoder(r.Body).Decode(&newConfig); err != nil {
			http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusBadRequest)
			return
		}
		homeDir, _ := os.UserHomeDir()
		configPath := path.Join(homeDir, ".picoclaw", "config.json")
		if err := config.SaveConfig(configPath, &newConfig); err != nil {
			http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusInternalServerError)
			return
		}
		*s.cfg = newConfig
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
		return
	}

	http.Error(w, `{"error": "Method not allowed"}`, http.StatusMethodNotAllowed)
}

func (s *Server) handleModels(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method == http.MethodGet {
		json.NewEncoder(w).Encode(s.cfg.ModelList)
		return
	}
	http.Error(w, `{"error": "Method not allowed"}`, http.StatusMethodNotAllowed)
}

func (s *Server) handleDefaultModel(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method != http.MethodPost {
		http.Error(w, `{"error": "Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		ModelName string `json:"model_name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusBadRequest)
		return
	}

	s.cfg.Agents.Defaults.ModelName = req.ModelName
	s.saveConfig(w)
}

func (s *Server) handleChat(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method != http.MethodPost {
		http.Error(w, `{"error": "Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Content    string `json:"content"`
		SessionKey string `json:"session_key"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusBadRequest)
		return
	}

	if req.SessionKey == "" {
		req.SessionKey = "web:default"
	}

	ctx := context.Background()
	response, err := s.agentLoop.ProcessDirectWithChannel(ctx, req.Content, req.SessionKey, "web", "default")
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"response": response})
}

func (s *Server) handleChatStream(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	content := r.URL.Query().Get("content")
	sessionKey := r.URL.Query().Get("session_key")
	if sessionKey == "" {
		sessionKey = "web:default"
	}

	ctx := context.Background()
	response, err := s.agentLoop.ProcessDirectWithChannel(ctx, content, sessionKey, "web", "default")
	if err != nil {
		fmt.Fprintf(w, "data: Error: %s\n\n", err.Error())
		fmt.Fprintf(w, "data: [DONE]\n\n")
		return
	}

	chunks := splitIntoChunks(response, 10)
	for _, chunk := range chunks {
		fmt.Fprintf(w, "data: %s\n\n", chunk)
	}
	fmt.Fprintf(w, "data: [DONE]\n\n")
}

func (s *Server) handleStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	status := map[string]interface{}{
		"status":    "ok",
		"version":   "0.1.0",
		"gateway":   fmt.Sprintf("%s:%d", s.cfg.Gateway.Host, s.cfg.Gateway.Port),
		"web_ui":    fmt.Sprintf("http://%s:%d", s.cfg.Gateway.Host, s.cfg.Gateway.Port+1),
		"models":    len(s.cfg.ModelList),
		"heartbeat": s.cfg.Heartbeat.Enabled,
	}
	json.NewEncoder(w).Encode(status)
}

func (s *Server) saveConfig(w http.ResponseWriter) {
	homeDir, _ := os.UserHomeDir()
	configPath := path.Join(homeDir, ".picoclaw", "config.json")
	if err := config.SaveConfig(configPath, s.cfg); err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func getContentType(ext string) string {
	switch ext {
	case ".html":
		return "text/html"
	case ".css":
		return "text/css"
	case ".js":
		return "application/javascript"
	case ".json":
		return "application/json"
	case ".png":
		return "image/png"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".gif":
		return "image/gif"
	case ".svg":
		return "image/svg+xml"
	case ".ico":
		return "image/x-icon"
	default:
		return "application/octet-stream"
	}
}

func splitIntoChunks(s string, chunkSize int) []string {
	var chunks []string
	runes := []rune(s)
	for i := 0; i < len(runes); i += chunkSize {
		end := i + chunkSize
		if end > len(runes) {
			end = len(runes)
		}
		chunks = append(chunks, string(runes[i:end]))
	}
	return chunks
}
