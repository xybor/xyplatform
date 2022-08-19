package xylog

// AddHandler adds a new handler to root logger.
func AddHandler(h *Handler) {
	rootLogger.AddHandler(h)
}

// RemoveHandler removes an existed handler from root logger.
func RemoveHandler(h *Handler) {
	rootLogger.RemoveHandler(h)
}

// AddFilter adds a specified filter to root logger.
func AddFilter(f Filter) {
	rootLogger.AddFilter(f)
}

// RemoveFilter removes an existed filter from root logger.
func RemoveFilter(f Filter) {
	rootLogger.RemoveFilter(f)
}

// SetLevel sets the new logging level for root logger.
func SetLevel(level int) {
	rootLogger.SetLevel(level)
}

// Log logs a message with a custom level by root logger.
func Log(level int, msg string, a ...any) {
	rootLogger.Log(level, msg, a...)
}

// Debug calls Log of root logger with DEBUG level.
func Debug(msg string, a ...any) {
	rootLogger.Debug(msg, a...)
}

// Info calls Log of root logger with INFO level.
func Info(msg string, a ...any) {
	rootLogger.Info(msg, a...)
}

// Warn calls Log of root logger with WARN level.
func Warn(msg string, a ...any) {
	rootLogger.Warn(msg, a...)
}

// Warning calls Log of root logger with WARNING level.
func Warning(msg string, a ...any) {
	rootLogger.Warning(msg, a...)
}

// Error calls Log of root logger with ERROR level.
func Error(msg string, a ...any) {
	rootLogger.Error(msg, a...)
}

// Fatal calls Log of root logger with FATAL level.
func Fatal(msg string, a ...any) {
	rootLogger.Fatal(msg, a...)
}

// Critical calls Log of root logger with CRITICAL level.
func Critical(msg string, a ...any) {
	rootLogger.Critical(msg, a...)
}
