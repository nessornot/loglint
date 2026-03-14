package quicktest

import "log/slog"

func main() {
	slog.Info("Starting server on port 8080")	// exp fail: uppercase start
	slog.Error("запуск сервер")					// exp fail: not eng
	slog.Warn("server started! 🚀")				// exp fail: special chars
	slog.Debug("token=123")						// exp fail: sensitive data
	
	slog.Info("starting server on port :8080")	// ОК
	slog.Error("failed to connect to database")	// ОК
}