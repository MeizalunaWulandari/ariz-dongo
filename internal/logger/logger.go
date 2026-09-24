package logger

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

const (
	maxLines = 200
	logDir   = "logs"
	logFile  = "logs/app.log"
)

var (
	mu sync.Mutex
)

func Init() error {
	mu.Lock()
	defer mu.Unlock()

	if err := os.MkdirAll(logDir, 0755); err != nil {
		return fmt.Errorf("create log directory: %w", err)
	}

	// Pastikan file log sudah ada.
	file, err := os.OpenFile(
		logFile,
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0644,
	)

	if err != nil {
		return fmt.Errorf("open log file: %w", err)
	}

	file.Close()

	// Rapikan jika file sudah lebih dari 200 baris.
	return trim()
}

func Close() {
	// Tidak ada file handle permanen.
	// Setiap write membuka dan menutup file.
}

func Info(format string, args ...any) {
	write("INFO", format, args...)
}

func Error(format string, args ...any) {
	write("ERROR", format, args...)
}

func Warn(format string, args ...any) {
	write("WARN", format, args...)
}

func Debug(format string, args ...any) {
	write("DEBUG", format, args...)
}

func write(level string, format string, args ...any) {
	mu.Lock()
	defer mu.Unlock()

	message := fmt.Sprintf(format, args...)

	line := fmt.Sprintf(
		"%s [%s] %s",
		time.Now().Format("2006-01-02 15:04:05"),
		level,
		message,
	)

	// ==========================================
	// TERMINAL
	// ==========================================

	fmt.Println(line)

	// ==========================================
	// FILE
	// ==========================================

	file, err := os.OpenFile(
		logFile,
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0644,
	)

	if err != nil {
		fmt.Printf(
			"%s [ERROR] Cannot open log file: %v\n",
			time.Now().Format("2006-01-02 15:04:05"),
			err,
		)

		return
	}

	_, err = file.WriteString(line + "\n")

	file.Close()

	if err != nil {
		fmt.Printf(
			"%s [ERROR] Cannot write log file: %v\n",
			time.Now().Format("2006-01-02 15:04:05"),
			err,
		)

		return
	}

	// ==========================================
	// MAX 200 LINES
	// ==========================================

	if err := trim(); err != nil {
		fmt.Printf(
			"%s [ERROR] Cannot trim log file: %v\n",
			time.Now().Format("2006-01-02 15:04:05"),
			err,
		)
	}
}

func trim() error {
	data, err := os.ReadFile(logFile)

	if err != nil {
		return err
	}

	content := string(data)

	if content == "" {
		return nil
	}

	lines := strings.Split(
		strings.TrimRight(content, "\r\n"),
		"\n",
	)

	if len(lines) <= maxLines {
		return nil
	}

	// Ambil 200 baris terakhir.
	lines = lines[len(lines)-maxLines:]

	newContent := strings.Join(lines, "\n") + "\n"

	return os.WriteFile(
		logFile,
		[]byte(newContent),
		0644,
	)
}
