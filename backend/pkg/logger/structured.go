// Package logger provides structured logging for pipeline jobs.
// Per AGENT_TECH_SPEC.md: log source, duration_ms, status, artifact in JSON format.
package logger

import (
	"encoding/json"
	"log"
	"time"
)

// FetchLog represents a structured log entry for HTTP fetching.
type FetchLog struct {
	Timestamp  string `json:"timestamp"`   // ISO8601 timestamp
	Level      string `json:"level"`       // "info", "error", "warn"
	Source     string `json:"source"`      // Data source name (e.g., "TEPCO")
	DurationMs int64  `json:"duration_ms"` // Request duration in milliseconds
	Status     string `json:"status"`      // "success", "failure", "retry"
	Artifact   string `json:"artifact"`    // Output file path (if applicable)
	Message    string `json:"message"`     // Human-readable message
	Error      string `json:"error,omitempty"` // Error message (if any)
}

// Logger handles structured logging.
type Logger struct {
	jsonOutput bool
}

// New creates a new logger.
// jsonOutput: if true, logs in JSON format; otherwise, human-readable format.
func New(jsonOutput bool) *Logger {
	return &Logger{
		jsonOutput: jsonOutput,
	}
}

// LogFetch logs an HTTP fetch operation.
func (l *Logger) LogFetch(source, status, artifact, message string, duration time.Duration, err error) {
	entry := FetchLog{
		Timestamp:  time.Now().Format(time.RFC3339),
		Level:      l.deriveLevel(status, err),
		Source:     source,
		DurationMs: duration.Milliseconds(),
		Status:     status,
		Artifact:   artifact,
		Message:    message,
	}

	if err != nil {
		entry.Error = err.Error()
	}

	if l.jsonOutput {
		l.logJSON(entry)
	} else {
		l.logHuman(entry)
	}
}

// logJSON outputs the log entry in JSON format.
func (l *Logger) logJSON(entry FetchLog) {
	data, err := json.Marshal(entry)
	if err != nil {
		log.Printf("Failed to marshal log entry: %v", err)
		return
	}
	log.Println(string(data))
}

// logHuman outputs the log entry in human-readable format.
func (l *Logger) logHuman(entry FetchLog) {
	if entry.Error != "" {
		log.Printf("[%s] %s: %s (%dms) - %s", entry.Level, entry.Source, entry.Message, entry.DurationMs, entry.Error)
	} else {
		log.Printf("[%s] %s: %s (%dms)", entry.Level, entry.Source, entry.Message, entry.DurationMs)
	}
}

// deriveLevel determines the log level from status and error.
func (l *Logger) deriveLevel(status string, err error) string {
	if err != nil || status == "failure" {
		return "error"
	}
	if status == "retry" {
		return "warn"
	}
	return "info"
}

// Info logs an informational message.
func (l *Logger) Info(message string) {
	if l.jsonOutput {
		entry := map[string]string{
			"timestamp": time.Now().Format(time.RFC3339),
			"level":     "info",
			"message":   message,
		}
		data, _ := json.Marshal(entry)
		log.Println(string(data))
	} else {
		log.Printf("[info] %s", message)
	}
}

// Error logs an error message.
func (l *Logger) Error(message string, err error) {
	if l.jsonOutput {
		entry := map[string]string{
			"timestamp": time.Now().Format(time.RFC3339),
			"level":     "error",
			"message":   message,
			"error":     err.Error(),
		}
		data, _ := json.Marshal(entry)
		log.Println(string(data))
	} else {
		log.Printf("[error] %s: %v", message, err)
	}
}
