// logger/logger.go

package logger

import (
	"io"
	"log"
	"os"
	"sync"

	"github.com/fatih/color"
)

// Logger represents a prettified logger with color.
type Logger struct {
	*log.Logger
	mu sync.Mutex
}

// NewLogger creates a new Logger instance with color.
func NewLogger(logFilePath string) *Logger {
	file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}

	return &Logger{
		Logger: log.New(io.MultiWriter(os.Stdout, file), "[Duplicate Cleaner] ", log.LstdFlags),
	}
}

// SetOutput sets the output destination for the logger.
func (l *Logger) SetOutput(w io.Writer) {
	l.Logger.SetOutput(w)
}

// Printf prints a formatted message to the logger with color.
func (l *Logger) Printf(format string, v ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.Logger.Printf(format, v...)
}

// PastelPrintf prints a formatted message in pastel color.
func (l *Logger) PastelPrintf(colorCode color.Attribute, format string, v ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	color.New(color.FgHiWhite).Add(colorCode).Printf(format, v...)
}

// PrettyPrint prints a formatted message in pastel color with a specified prefix.
func (l *Logger) PrettyPrint(prefix string, colorCode color.Attribute, message string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	color.New(color.FgHiWhite).Add(colorCode).Printf("[%s] %s\n", prefix, message)
}

// PrettyError prints a formatted error message in pastel red color.
func (l *Logger) PrettyError(err error) {
	l.PrettyPrint("Error", color.FgHiRed, err.Error())
}

// PrettyInfo prints a formatted info message in pastel green color.
func (l *Logger) PrettyInfo(message string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	// Create a new color instance with pastel green color
	pastelGreen := color.New(color.FgHiGreen) // Change this line

	// Print the formatted info message with additional styling
	pastelGreen.Printf("[Info] %s\n", message)
}

// PrettyWarning prints a formatted warning message in pastel yellow color.
func (l *Logger) PrettyWarning(message string) {
	l.PrettyPrint("Warning", color.FgHiYellow, message)
}

// PrettySuccess prints a formatted success message in pastel blue color.
func (l *Logger) PrettySuccess(message string) {
	l.PrettyPrint("Success", color.FgHiBlue, message)
}

// PastelGreenPrint prints a message in pastel green color.
func (l *Logger) PastelGreenPrint(message string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	color.New(color.FgHiGreen).Print(message)
}

// PastelBluePrintf prints a formatted message in pastel blue color.
func (l *Logger) PastelBluePrintf(format string, v ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	color.New(color.FgHiBlue).Printf(format, v...)
}

// PastelGreenPrintf prints a formatted message in pastel green color.
func (l *Logger) PastelGreenPrintf(format string, v ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	color.New(color.FgHiGreen).Printf(format, v...)
}

// CyanPrint prints a message in cyan color.
func (l *Logger) CyanPrint(message string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	color.New(color.FgCyan).Print(message)
}

// YellowPrint prints a message in yellow color.
func (l *Logger) YellowPrint(message string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	color.New(color.FgHiYellow).Print(message)
}

// WhitePrint prints a message in white color.
func (l *Logger) WhitePrint(message string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	color.New(color.FgHiWhite).Print(message)
}

// CyanPrintf prints a formatted message in cyan color.
func (l *Logger) CyanPrintf(format string, v ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	color.New(color.FgCyan).Printf(format, v...)
}

// YellowPrintf prints a formatted message in yellow color.
func (l *Logger) YellowPrintf(format string, v ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	color.New(color.FgHiYellow).Printf(format, v...)
}

// WhitePrintf prints a formatted message in white color.
func (l *Logger) WhitePrintf(format string, v ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	color.New(color.FgHiWhite).Printf(format, v...)
}

// Close closes the underlying logger's output.
func (l *Logger) Close() {
	l.mu.Lock()
	defer l.mu.Unlock()
	// Implement any cleanup logic if needed
}
