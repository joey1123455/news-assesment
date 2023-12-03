package middleware

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

// Logger defines the interface for logging operations
type Logger interface {
	Log(message string)
}

// FileLogger implements Logger for writing logs to a file
type FileLogger struct {
	file *os.File
}

// Log writes a log message to the file
func (f *FileLogger) Log(message string) {
	log.New(f.file, "", log.LstdFlags).Println(message)
}

// NewFileLogger creates a new FileLogger instance
func NewFileLogger(filename string) (*FileLogger, error) {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		// Create the file if it doesn't exist
		file, createErr := os.Create(filename)
		if createErr != nil {
			return nil, fmt.Errorf("error creating file: %v", createErr)
		}
		return &FileLogger{file: file}, nil
	} else if err != nil {
		// Handle other errors
		return nil, fmt.Errorf("error checking file existence: %v", err)
	}

	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	return &FileLogger{file: file}, nil
}

// RequestLogger logs the HTTP request details to a file
func RequestLogger(logger Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()

		c.Next()

		latency := time.Since(t)

		// Log the request details using the provided logger
		logger.Log(fmt.Sprintf("%s %s %s %s",
			c.Request.Method,
			c.Request.RequestURI,
			c.Request.Proto,
			latency,
		))
	}
}

// ResponseLogger logs the HTTP response details to a file
func ResponseLogger(logger Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("X-Content-Type-Options", "nosniff")

		c.Next()

		// Log the response details using the provided logger
		logger.Log(fmt.Sprintf("%d %s %s",
			c.Writer.Status(),
			c.Request.Method,
			c.Request.RequestURI,
		))
	}
}
