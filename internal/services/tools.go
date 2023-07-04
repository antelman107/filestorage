package services

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

const gracefulShutdownDeadline = 15 * time.Second

func runEchoWithGracefulShutdown(ctx context.Context, e *echo.Echo, listenPort string, zapLogger *zap.Logger) error {
	go func() {
		if err := e.Start(listenPort); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal(err)
		}
	}()

	<-ctx.Done()
	zapLogger.Info("got shutdown signal, stopping")

	shutDownCtx, cancel := context.WithTimeout(context.Background(), gracefulShutdownDeadline)
	defer cancel()

	if err := e.Shutdown(shutDownCtx); err != nil {
		return fmt.Errorf("failed to shutdown echo: %w", err)
	}

	return nil
}

func isDirectoryEmpty(name string) bool {
	f, err := os.Open(name)
	if err != nil {
		return false
	}
	defer f.Close()

	_, err = f.Readdirnames(1) // Or f.Readdir(1)
	if err == io.EOF {
		return true
	}
	return false
}

func getPathWithDirectories(fileName string) string {
	result := ""
	for i := 0; i < numSubdirectories; i++ {
		result += fileName[i*2 : (i+1)*2]
		if i < numSubdirectories-1 {
			result += string(os.PathSeparator)
		}
	}

	result += string(os.PathSeparator) + fileName

	return result
}
