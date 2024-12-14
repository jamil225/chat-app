package user_service

import (
	"log"
	"time"
)

// LogExecutionTime logs the time taken by a function to execute
func LogExecutionTime(start time.Time, functionName string) {
	elapsed := time.Since(start)
	log.Printf("Execution time for %s: %s\n", functionName, elapsed)
}
