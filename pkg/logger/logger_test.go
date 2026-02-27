package logger

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestLogLevelFiltering(t *testing.T) {
	initialLevel := GetLevel()
	defer SetLevel(initialLevel)

	SetLevel(WARN)

	tests := []struct {
		name      string
		level     LogLevel
		shouldLog bool
	}{
		{"DEBUG message", DEBUG, false},
		{"INFO message", INFO, false},
		{"WARN message", WARN, true},
		{"ERROR message", ERROR, true},
		{"FATAL message", FATAL, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.level {
			case DEBUG:
				Debug(tt.name)
			case INFO:
				Info(tt.name)
			case WARN:
				Warn(tt.name)
			case ERROR:
				Error(tt.name)
			case FATAL:
				if tt.shouldLog {
					t.Logf("FATAL test skipped to prevent program exit")
				}
			}
		})
	}

	SetLevel(INFO)
}

func TestLoggerWithComponent(t *testing.T) {
	initialLevel := GetLevel()
	defer SetLevel(initialLevel)

	SetLevel(DEBUG)

	tests := []struct {
		name      string
		component string
		message   string
		fields    map[string]any
	}{
		{"Simple message", "test", "Hello, world!", nil},
		{"Message with component", "discord", "Discord message", nil},
		{"Message with fields", "telegram", "Telegram message", map[string]any{
			"user_id": "12345",
			"count":   42,
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch {
			case tt.fields == nil && tt.component != "":
				InfoC(tt.component, tt.message)
			case tt.fields != nil:
				InfoF(tt.message, tt.fields)
			default:
				Info(tt.message)
			}
		})
	}

	SetLevel(INFO)
}

func TestLogLevels(t *testing.T) {
	tests := []struct {
		name  string
		level LogLevel
		want  string
	}{
		{"DEBUG level", DEBUG, "DEBUG"},
		{"INFO level", INFO, "INFO"},
		{"WARN level", WARN, "WARN"},
		{"ERROR level", ERROR, "ERROR"},
		{"FATAL level", FATAL, "FATAL"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if logLevelNames[tt.level] != tt.want {
				t.Errorf("logLevelNames[%d] = %s, want %s", tt.level, logLevelNames[tt.level], tt.want)
			}
		})
	}
}

func TestSetGetLevel(t *testing.T) {
	initialLevel := GetLevel()
	defer SetLevel(initialLevel)

	tests := []LogLevel{DEBUG, INFO, WARN, ERROR, FATAL}

	for _, level := range tests {
		SetLevel(level)
		if GetLevel() != level {
			t.Errorf("SetLevel(%v) -> GetLevel() = %v, want %v", level, GetLevel(), level)
		}
	}
}

func TestLoggerHelperFunctions(t *testing.T) {
	initialLevel := GetLevel()
	defer SetLevel(initialLevel)

	SetLevel(INFO)

	Debug("This should not log")
	Info("This should log")
	Warn("This should log")
	Error("This should log")

	InfoC("test", "Component message")
	InfoF("Fields message", map[string]any{"key": "value"})

	WarnC("test", "Warning with component")
	ErrorF("Error with fields", map[string]any{"error": "test"})

	SetLevel(DEBUG)
	DebugC("test", "Debug with component")
	WarnF("Warning with fields", map[string]any{"key": "value"})
}

func TestBuildLogDetailRedactsSensitiveFields(t *testing.T) {
	detail := buildLogDetail("agent", "Tool call requested", map[string]any{
		"tool":    "web_search",
		"content": "secret prompt",
		"args":    map[string]any{"query": "hello"},
		"chat_id": "default",
	})

	if strings.Contains(detail, "secret prompt") {
		t.Fatalf("detail should not contain sensitive content: %s", detail)
	}
	if strings.Contains(detail, "query") {
		t.Fatalf("detail should not contain sensitive args: %s", detail)
	}
	if !strings.Contains(detail, "tool=web_search") {
		t.Fatalf("detail should keep non-sensitive fields: %s", detail)
	}
}

func TestCleanupOldFilesLockedKeepsRecentSevenDays(t *testing.T) {
	tmpDir := t.TempDir()
	now := time.Date(2026, time.February, 27, 10, 0, 0, 0, time.Local)

	// Create 10 daily log files.
	for i := 0; i < 10; i++ {
		day := now.AddDate(0, 0, -i).Format(dateLayout)
		p := filepath.Join(tmpDir, day+".log")
		if err := os.WriteFile(p, []byte("x"), 0o644); err != nil {
			t.Fatalf("failed to create test file: %v", err)
		}
	}

	logger.logDir = tmpDir
	logger.retentionDays = 7

	if err := cleanupOldFilesLocked(now); err != nil {
		t.Fatalf("cleanupOldFilesLocked failed: %v", err)
	}

	files, err := os.ReadDir(tmpDir)
	if err != nil {
		t.Fatalf("failed to list log directory: %v", err)
	}

	if len(files) != 7 {
		t.Fatalf("expected 7 retained files, got %d", len(files))
	}
}

func TestWriteFileLogLineFormat(t *testing.T) {
	tmpDir := t.TempDir()
	now := time.Date(2026, time.February, 27, 10, 11, 12, 0, time.Local)

	logger.fileMu.Lock()
	oldDir := logger.logDir
	oldRetention := logger.retentionDays
	oldDate := logger.currentDate
	oldFile := logger.file

	logger.logDir = tmpDir
	logger.retentionDays = 7
	logger.currentDate = ""
	logger.file = nil
	logger.fileMu.Unlock()

	writeFileLogLine(INFO, "agent", "Conversation started", map[string]any{
		"session_key": "web:default",
		"content":     "hello world",
	}, now)

	logger.fileMu.Lock()
	if logger.file != nil {
		_ = logger.file.Close()
	}
	logger.logDir = oldDir
	logger.retentionDays = oldRetention
	logger.currentDate = oldDate
	logger.file = oldFile
	logger.fileMu.Unlock()

	p := filepath.Join(tmpDir, "2026-02-27.log")
	data, err := os.ReadFile(p)
	if err != nil {
		t.Fatalf("failed to read log file: %v", err)
	}

	line := strings.TrimSpace(string(data))
	parts := strings.SplitN(line, " | ", 3)
	if len(parts) != 3 {
		t.Fatalf("expected 3-part log line, got %d: %q", len(parts), line)
	}
	if parts[0] != "INFO" {
		t.Fatalf("unexpected level part: %q", parts[0])
	}
	if !strings.Contains(parts[2], "event=Conversation started") {
		t.Fatalf("unexpected detail part: %q", parts[2])
	}
	if strings.Contains(parts[2], "hello world") {
		t.Fatalf("detail should not contain sensitive content: %q", parts[2])
	}
}
