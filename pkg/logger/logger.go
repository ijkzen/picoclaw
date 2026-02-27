package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
	FATAL
)

var (
	logLevelNames = map[LogLevel]string{
		DEBUG: "DEBUG",
		INFO:  "INFO",
		WARN:  "WARN",
		ERROR: "ERROR",
		FATAL: "FATAL",
	}

	currentLevel = INFO
	logger       *Logger
	once         sync.Once
	mu           sync.RWMutex
)

type Logger struct {
	file          *os.File
	logDir        string
	retentionDays int
	currentDate   string
	fileMu        sync.Mutex
}

type LogEntry struct {
	Level     string         `json:"level"`
	Timestamp string         `json:"timestamp"`
	Component string         `json:"component,omitempty"`
	Message   string         `json:"message"`
	Fields    map[string]any `json:"fields,omitempty"`
	Caller    string         `json:"caller,omitempty"`
}

const (
	dateLayout              = "2006-01-02"
	defaultRetentionDays    = 7
	maxFieldValueLogRunes   = 120
	maxDetailLengthLogRunes = 800
)

var sensitiveFieldKeyTokens = []string{
	"content",
	"text",
	"message",
	"prompt",
	"response",
	"arguments",
	"argument",
	"args",
	"input",
	"output",
	"preview",
	"history",
	"messages_json",
	"tools_json",
	"query",
	"error",
}

func init() {
	once.Do(func() {
		logger = &Logger{}
	})
}

func SetLevel(level LogLevel) {
	mu.Lock()
	defer mu.Unlock()
	currentLevel = level
}

func GetLevel() LogLevel {
	mu.RLock()
	defer mu.RUnlock()
	return currentLevel
}

func EnableFileLogging(filePath string) error {
	logDir := filepath.Dir(filePath)
	if err := EnableDailyFileLogging(logDir, defaultRetentionDays); err != nil {
		return err
	}

	log.Println("File logging enabled:", filePath)
	return nil
}

func EnableDailyFileLogging(logDir string, retentionDays int) error {
	if strings.TrimSpace(logDir) == "" {
		return fmt.Errorf("log directory is required")
	}
	if retentionDays <= 0 {
		retentionDays = defaultRetentionDays
	}
	if err := os.MkdirAll(logDir, 0o755); err != nil {
		return fmt.Errorf("failed to create log directory: %w", err)
	}

	mu.Lock()
	defer mu.Unlock()

	logger.fileMu.Lock()
	defer logger.fileMu.Unlock()

	logger.logDir = logDir
	logger.retentionDays = retentionDays

	if err := rotateDailyFileLocked(time.Now()); err != nil {
		return err
	}

	if err := cleanupOldFilesLocked(time.Now()); err != nil {
		log.Printf("failed to cleanup old log files: %v", err)
	}

	log.Printf("Daily file logging enabled: dir=%s retention_days=%d", logDir, retentionDays)
	return nil
}

func DisableFileLogging() {
	mu.Lock()
	defer mu.Unlock()

	logger.fileMu.Lock()
	defer logger.fileMu.Unlock()

	if logger.file != nil {
		logger.file.Close()
		logger.file = nil
		log.Println("File logging disabled")
	}

	logger.logDir = ""
	logger.currentDate = ""
	logger.retentionDays = 0
}

func logMessage(level LogLevel, component string, message string, fields map[string]any) {
	if level < currentLevel {
		return
	}

	entry := LogEntry{
		Level:     logLevelNames[level],
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Component: component,
		Message:   message,
		Fields:    fields,
	}

	if pc, file, line, ok := runtime.Caller(2); ok {
		fn := runtime.FuncForPC(pc)
		if fn != nil {
			entry.Caller = fmt.Sprintf("%s:%d (%s)", file, line, fn.Name())
		}
	}

	writeFileLogLine(level, component, message, fields, time.Now())

	var fieldStr string
	if len(fields) > 0 {
		fieldStr = " " + formatFields(fields)
	} else {
		fieldStr = ""
	}

	logLine := fmt.Sprintf("[%s] [%s]%s %s%s",
		entry.Timestamp,
		logLevelNames[level],
		formatComponent(component),
		message,
		fieldStr,
	)

	log.Println(logLine)

	if level == FATAL {
		os.Exit(1)
	}
}

func writeFileLogLine(level LogLevel, component, message string, fields map[string]any, now time.Time) {
	logger.fileMu.Lock()
	defer logger.fileMu.Unlock()

	if logger.logDir == "" {
		return
	}

	if err := rotateDailyFileLocked(now); err != nil {
		log.Printf("failed to rotate log file: %v", err)
		return
	}

	if logger.file == nil {
		return
	}

	detail := buildLogDetail(component, message, fields)
	line := fmt.Sprintf("%s | %s | %s\n", logLevelNames[level], now.Format(time.RFC3339), detail)
	_, _ = logger.file.WriteString(line)
}

func rotateDailyFileLocked(now time.Time) error {
	if logger.logDir == "" {
		return nil
	}

	currentDate := now.Format(dateLayout)
	if logger.file != nil && logger.currentDate == currentDate {
		return nil
	}

	if logger.file != nil {
		_ = logger.file.Close()
		logger.file = nil
	}

	logPath := filepath.Join(logger.logDir, currentDate+".log")
	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return fmt.Errorf("failed to open daily log file: %w", err)
	}

	logger.file = file
	logger.currentDate = currentDate

	if err := cleanupOldFilesLocked(now); err != nil {
		log.Printf("failed to cleanup old log files: %v", err)
	}

	return nil
}

func cleanupOldFilesLocked(now time.Time) error {
	if logger.logDir == "" {
		return nil
	}

	retentionDays := logger.retentionDays
	if retentionDays <= 0 {
		retentionDays = defaultRetentionDays
	}

	entries, err := os.ReadDir(logger.logDir)
	if err != nil {
		return fmt.Errorf("failed to read log directory: %w", err)
	}

	cutoff := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).
		AddDate(0, 0, -(retentionDays - 1))

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		if !strings.HasSuffix(name, ".log") {
			continue
		}

		datePart := strings.TrimSuffix(name, ".log")
		fileDate, err := time.ParseInLocation(dateLayout, datePart, now.Location())
		if err != nil {
			continue
		}

		if fileDate.Before(cutoff) {
			_ = os.Remove(filepath.Join(logger.logDir, name))
		}
	}

	return nil
}

func buildLogDetail(component, message string, fields map[string]any) string {
	parts := make([]string, 0, len(fields)+1)

	if component != "" {
		parts = append(parts, "component="+sanitizeDetailValue(component))
	}

	baseMessage := sanitizeDetailValue(message)
	if baseMessage != "" {
		parts = append(parts, "event="+baseMessage)
	}

	for key, value := range fields {
		if isSensitiveFieldKey(key) {
			continue
		}
		formatted := formatFieldValue(value)
		if formatted == "" {
			continue
		}
		parts = append(parts, sanitizeFieldKey(key)+"="+sanitizeDetailValue(formatted))
	}

	if len(parts) == 0 {
		return "event=log"
	}

	detail := strings.Join(parts, " ")
	return truncateRunes(detail, maxDetailLengthLogRunes)
}

func sanitizeFieldKey(key string) string {
	k := strings.TrimSpace(strings.ToLower(key))
	k = strings.ReplaceAll(k, " ", "_")
	if k == "" {
		return "field"
	}
	return k
}

func sanitizeDetailValue(v string) string {
	s := strings.TrimSpace(v)
	if s == "" {
		return ""
	}
	s = strings.ReplaceAll(s, "\n", " ")
	s = strings.ReplaceAll(s, "\r", " ")
	s = strings.ReplaceAll(s, "|", "/")
	return truncateRunes(s, maxFieldValueLogRunes)
}

func isSensitiveFieldKey(key string) bool {
	lower := strings.ToLower(strings.TrimSpace(key))
	for _, token := range sensitiveFieldKeyTokens {
		if strings.Contains(lower, token) {
			return true
		}
	}
	return false
}

func formatFieldValue(v any) string {
	switch value := v.(type) {
	case nil:
		return ""
	case string:
		return value
	case fmt.Stringer:
		return value.String()
	case error:
		return value.Error()
	case bool:
		if value {
			return "true"
		}
		return "false"
	case int:
		return strconv.Itoa(value)
	case int8, int16, int32, int64:
		return fmt.Sprintf("%d", value)
	case uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", value)
	case float32, float64:
		return fmt.Sprintf("%v", value)
	default:
		return fmt.Sprintf("%v", value)
	}
}

func truncateRunes(s string, max int) string {
	if max <= 0 {
		return s
	}
	r := []rune(s)
	if len(r) <= max {
		return s
	}
	return string(r[:max]) + "..."
}

func formatComponent(component string) string {
	if component == "" {
		return ""
	}
	return fmt.Sprintf(" %s:", component)
}

func formatFields(fields map[string]any) string {
	parts := make([]string, 0, len(fields))
	for k, v := range fields {
		parts = append(parts, fmt.Sprintf("%s=%v", k, v))
	}
	return fmt.Sprintf("{%s}", strings.Join(parts, ", "))
}

func Debug(message string) {
	logMessage(DEBUG, "", message, nil)
}

func DebugC(component string, message string) {
	logMessage(DEBUG, component, message, nil)
}

func DebugF(message string, fields map[string]any) {
	logMessage(DEBUG, "", message, fields)
}

func DebugCF(component string, message string, fields map[string]any) {
	logMessage(DEBUG, component, message, fields)
}

func Info(message string) {
	logMessage(INFO, "", message, nil)
}

func InfoC(component string, message string) {
	logMessage(INFO, component, message, nil)
}

func InfoF(message string, fields map[string]any) {
	logMessage(INFO, "", message, fields)
}

func InfoCF(component string, message string, fields map[string]any) {
	logMessage(INFO, component, message, fields)
}

func Warn(message string) {
	logMessage(WARN, "", message, nil)
}

func WarnC(component string, message string) {
	logMessage(WARN, component, message, nil)
}

func WarnF(message string, fields map[string]any) {
	logMessage(WARN, "", message, fields)
}

func WarnCF(component string, message string, fields map[string]any) {
	logMessage(WARN, component, message, fields)
}

func Error(message string) {
	logMessage(ERROR, "", message, nil)
}

func ErrorC(component string, message string) {
	logMessage(ERROR, component, message, nil)
}

func ErrorF(message string, fields map[string]any) {
	logMessage(ERROR, "", message, fields)
}

func ErrorCF(component string, message string, fields map[string]any) {
	logMessage(ERROR, component, message, fields)
}

func Fatal(message string) {
	logMessage(FATAL, "", message, nil)
}

func FatalC(component string, message string) {
	logMessage(FATAL, component, message, nil)
}

func FatalF(message string, fields map[string]any) {
	logMessage(FATAL, "", message, fields)
}

func FatalCF(component string, message string, fields map[string]any) {
	logMessage(FATAL, component, message, fields)
}
