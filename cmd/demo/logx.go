package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/garaekz/tfx/logx"
)

func runLogxDemo() {
	fmt.Println("📝 TFX Logging Demonstrations")
	fmt.Println("=============================")

	// 1. Basic Logging Levels
	fmt.Println("\n1. BASIC LOGGING LEVELS:")
	logx.Trace("This is a trace message - very detailed debugging")
	logx.Debug("This is a debug message - development info")
	logx.Info("This is an info message - general information")
	logx.Success("This is a success message - operation completed!")
	logx.Warn("This is a warning message - something to watch")
	logx.Error("This is an error message - something went wrong")
	// logx.Fatal("This would be a fatal message - program exits")

	time.Sleep(500 * time.Millisecond)

	// 2. If Variants - Conditional Logging
	fmt.Println("\n2. CONDITIONAL LOGGING (If variants):")
	
	// Simulate some operations with potential errors
	err1 := errors.New("connection timeout")
	err2 := errors.New("invalid input")
	var err3 error // nil error

	// These return true if they logged, false if no error
	if logx.ErrorIf(err1, "Failed to connect to database") {
		fmt.Println("   → Database error was logged")
	}

	if logx.WarnIf(err2, "Input validation issue") {
		fmt.Println("   → Warning was logged")
	}

	if !logx.InfoIf(err3, "This won't be logged because err3 is nil") {
		fmt.Println("   → No error, so nothing was logged")
	}

	time.Sleep(500 * time.Millisecond)

	// 3. Modern Badge System
	fmt.Println("\n3. MODERN BADGE SYSTEM:")
	logx.SuccessBadge("API", "Successfully connected to external service")
	logx.ErrorBadge("DB", "Connection pool exhausted")
	logx.WarnBadge("CACHE", "Cache hit ratio below threshold")
	logx.InfoBadge("SYS", "System startup completed")
	logx.DebugBadge("AUTH", "JWT token validation process")

	time.Sleep(500 * time.Millisecond)

	// 4. Themed Badges - Domain Specific
	fmt.Println("\n4. THEMED DOMAIN BADGES:")
	logx.DatabaseBadge("Connected to PostgreSQL", true)
	logx.DatabaseBadge("Failed to migrate schema", false)
	logx.APIBadge("External API call succeeded", true)
	logx.APIBadge("Rate limit exceeded", false)
	logx.AuthBadge("User authenticated successfully", true)
	logx.AuthBadge("Invalid credentials provided", false)
	logx.CacheBadge("Redis cache warmed up", true)
	logx.CacheBadge("Cache miss for user data", false)
	logx.SystemBadge("Application ready to serve requests")
	logx.ConfigBadge("Configuration loaded from environment")

	time.Sleep(500 * time.Millisecond)

	// 5. Fluent/DSL Interface
	fmt.Println("\n5. FLUENT/DSL INTERFACE:")
	
	testErr := errors.New("network unreachable")
	
	// Fluent conditional logging
	logx.If(testErr).AsError().Msg("Network operation failed")
	logx.If(testErr).AsWarn().WithField("retry_count", 3).Msg("Will retry operation")
	
	// More complex fluent usage
	logx.If(testErr).
		AsError().
		WithField("component", "network").
		WithField("operation", "fetch_data").
		WithField("timeout", "30s").
		Msg("Data fetch operation failed")

	time.Sleep(500 * time.Millisecond)

	// 6. Context-based Logging
	fmt.Println("\n6. CONTEXT-BASED LOGGING:")
	
	// Create a context with fields
	ctx := logx.WithFields(map[string]interface{}{
		"request_id": "req-12345",
		"user_id":    "user-67890",
		"session":    "sess-abcdef",
	})

	ctx.Info("User logged into the system")
	ctx.Success("Profile update completed successfully")
	
	// Add more context
	ctxWithMore := ctx.WithField("operation", "file_upload")
	ctxWithMore.Info("Starting file upload process")
	
	// Conditional logging with context
	uploadErr := errors.New("file too large")
	if ctxWithMore.ErrorIf(uploadErr, "File upload failed") {
		ctxWithMore.Warn("Will notify user about file size limit")
	}

	time.Sleep(500 * time.Millisecond)

	// 7. Multipath API Examples
	fmt.Println("\n7. MULTIPATH API EXAMPLES:")
	
	// EXPRESS: Quick defaults
	expressLogger := logx.Log()
	expressLogger.Info("This is from the express logger")

	// CONFIG: Explicit configuration
	configLogger := logx.NewWithConfig(logx.LogOptions{
		Level:     logx.LevelDebug,
		Timestamp: true,
	})
	configLogger.Debug("This is from the config logger with timestamp")

	// FLUENT: Functional options
	fluentLogger := logx.NewWith(
		logx.WithLevel(logx.LevelInfo),
		logx.WithTimestamp(false),
		logx.WithDevelopment(),
	)
	fluentLogger.Success("This is from the fluent logger")

	time.Sleep(500 * time.Millisecond)

	// 8. Preset Loggers
	fmt.Println("\n8. PRESET LOGGER EXAMPLES:")
	
	devLogger := logx.DevLogger()
	devLogger.Debug("Development logger - shows debug info")

	consoleLogger := logx.ConsoleLogger()
	consoleLogger.Info("Console logger - colorful output")

	structuredLogger := logx.StructuredLogger()
	structuredLogger.Info("Structured logger - JSON format")

	time.Sleep(500 * time.Millisecond)

	// 9. Custom Modern Badges
	fmt.Println("\n9. CUSTOM MODERN BADGES:")
	
	logx.ModernBadge("CUSTOM", "This is a custom badge",
		logx.WithBadgeColor(logx.White),
		logx.WithBadgeBackground(logx.Purple),
		logx.WithBold(),
	)

	logx.ModernBadge("HTTP", "Request processed successfully",
		logx.WithBadgeColor(logx.Black),
		logx.WithBadgeBackground(logx.Green),
		logx.WithBadgeLevel(logx.LevelSuccess),
	)

	logx.ModernBadge("WARN", "Memory usage is high",
		logx.WithBadgeColor(logx.Black),
		logx.WithBadgeBackground(logx.Yellow),
		logx.WithBadgeLevel(logx.LevelWarn),
		logx.WithBold(),
	)

	fmt.Println("\n✅ Logging demonstration completed!")
}

// Example of how to use FatalIf (commented out since it would exit)
func exampleFatalIf() {
	criticalErr := errors.New("database connection lost")
	
	// This would log the error and exit the program with status 1
	// logx.FatalIf(criticalErr, "Cannot continue without database")
	
	// For demo purposes, just show what it would do:
	fmt.Println("   → FatalIf would log this error and call os.Exit(1)")
	fmt.Printf("   → Error: %v\n", criticalErr)
}