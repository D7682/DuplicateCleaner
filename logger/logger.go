// logger/logger.go

package logger

import (
	"github.com/fatih/color"
	"log"
	"os"
)

// Logger represents a prettified logger with color.
type Logger struct {
	*log.Logger
	Color *color.Color
}

// NewLogger creates a new Logger instance with color.
func NewLogger() *Logger {
	return &Logger{
		Logger: log.New(os.Stdout, "[Duplicate Cleaner] ", log.LstdFlags),
		Color:  color.New(color.FgWhite), // Set your desired color
	}
}
