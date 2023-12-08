package utils

import (
	"fmt"
	"log"
	"os"
	"time"
)

// LogErrorToFile writes an error message to a file
func LogErrorToFile(operation, errorMessage string) error {
	// Open the file for appending
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current working directory:", err)
		return err
	}

	file, err := os.OpenFile(currentDir+"/logs/errors.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("error opening or creating log file: %v", err)
	}
	defer file.Close()
	logger := log.New(file, "", log.LstdFlags)
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	logger.Printf("[%s] [ERROR] [%s] %s\n", timestamp, operation, errorMessage)
	fmt.Printf("[%s] [ERROR] [%s] %s\n", timestamp, operation, errorMessage)

	return nil
}
