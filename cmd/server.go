package cmd

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/philiplambok/task-api/internal/transport"
	"github.com/spf13/cobra"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the HTTP server",
	Long:  `Start the HTTP server on port 8000 to handle incoming requests and serve API endpoints.`,
	Run: func(cmd *cobra.Command, args []string) {
		config, err := newConfig()
		if err != nil {
			slog.Error("failed to load config", slog.String("err", err.Error()))
			return
		}

		db, err := newDatabase(config)
		if err != nil {
			slog.Error("failed to connect to database", slog.String("err", err.Error()))
			return
		}

		httpServer := transport.NewHTTPServer(config, db)

		// Create a channel to listen for OS signals
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

		// Create a channel to capture server errors
		serverErrors := make(chan error, 1)

		// Start the server in a goroutine
		go func() {
			slog.Info("Starting HTTP server", slog.Int("port", config.HTTPServer.Port))
			if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				serverErrors <- err
			}
		}()

		// Block until we receive a signal or an error
		select {
		case sig := <-sigChan:
			slog.Info("Received shutdown signal", slog.String("signal", sig.String()))
		case err := <-serverErrors:
			slog.Error("Server error", slog.String("err", err.Error()))
			return
		}

		// Create a context with timeout for graceful shutdown
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		slog.Info("Shutting down server gracefully...")
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			slog.Error("Failed to shutdown server gracefully", slog.String("err", err.Error()))
			return
		}

		slog.Info("Server stopped successfully")
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serverCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serverCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
