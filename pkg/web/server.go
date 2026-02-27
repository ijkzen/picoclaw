package web

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/sipeed/picoclaw/pkg/agent"
	"github.com/sipeed/picoclaw/pkg/bus"
	"github.com/sipeed/picoclaw/pkg/config"
	"github.com/sipeed/picoclaw/pkg/logger"
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
	mux.HandleFunc("/api/status", s.handleStatus)
	mux.HandleFunc("/api/gateway/restart", s.handleGatewayRestart)

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
	logWebOperation(r, "config")

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method == http.MethodGet {
		json.NewEncoder(w).Encode(s.cfg)
		return
	}

	if r.Method == http.MethodPost {
		var newConfig config.Config
		if err := json.NewDecoder(r.Body).Decode(&newConfig); err != nil {
			logWebError("config_decode_failed", err)
			http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusBadRequest)
			return
		}
		homeDir, _ := os.UserHomeDir()
		configPath := path.Join(homeDir, ".picoclaw", "config.json")
		if err := config.SaveConfig(configPath, &newConfig); err != nil {
			logWebError("config_save_failed", err)
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
	logWebOperation(r, "models")

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method == http.MethodGet {
		json.NewEncoder(w).Encode(s.cfg.ModelList)
		return
	}
	http.Error(w, `{"error": "Method not allowed"}`, http.StatusMethodNotAllowed)
}

func (s *Server) handleDefaultModel(w http.ResponseWriter, r *http.Request) {
	logWebOperation(r, "default_model")

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
		logWebError("default_model_decode_failed", err)
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusBadRequest)
		return
	}

	s.cfg.Agents.Defaults.ModelName = req.ModelName
	s.saveConfig(w)
}

func (s *Server) handleChat(w http.ResponseWriter, r *http.Request) {
	logWebOperation(r, "chat")

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
		logWebError("chat_decode_failed", err)
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusBadRequest)
		return
	}

	if req.SessionKey == "" {
		req.SessionKey = "web:default"
	}

	ctx := context.Background()
	response, err := s.agentLoop.ProcessDirectWithChannel(ctx, req.Content, req.SessionKey, "web", "default")
	if err != nil {
		logWebError("chat_process_failed", err)
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"response": response})
}

func (s *Server) handleStatus(w http.ResponseWriter, r *http.Request) {
	logWebOperation(r, "status")

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

func (s *Server) handleGatewayRestart(w http.ResponseWriter, r *http.Request) {
	logWebOperation(r, "gateway_restart")

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method != http.MethodPost {
		http.Error(w, `{"error": "Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	exe, err := os.Executable()
	if err != nil {
		logWebError("gateway_restart_resolve_exe_failed", err)
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	// Delay restart slightly so this HTTP response can be returned before
	// the running gateway process is stopped.
	cmd := exec.Command(exe, "gateway", "restart", "--delay", "1")
	cmd.Stdin = nil
	cmd.Stdout = nil
	cmd.Stderr = nil
	if err := cmd.Start(); err != nil {
		logWebError("gateway_restart_start_failed", err)
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusInternalServerError)
		return
	}
	_ = cmd.Process.Release()

	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{"status": "restarting"})
}

func (s *Server) saveConfig(w http.ResponseWriter) {
	homeDir, _ := os.UserHomeDir()
	configPath := path.Join(homeDir, ".picoclaw", "config.json")
	if err := config.SaveConfig(configPath, s.cfg); err != nil {
		logWebError("save_config_failed", err)
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
func logWebOperation(r *http.Request, operation string) {
	if r == nil {
		return
	}

	logger.InfoCF("web", "Web operation",
		map[string]any{
			"operation": operation,
			"method":    r.Method,
			"path":      r.URL.Path,
		})
}

func logWebError(operation string, err error) {
	if err == nil {
		return
	}

	logger.ErrorCF("web", "Web operation failed",
		map[string]any{
			"operation": operation,
			"error":     err.Error(),
		})
}
