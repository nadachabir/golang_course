// Command api is the entry point for the TaskFlow HTTP server.
//
// Phase 0 ships a placeholder so the module compiles and runs. In Phase 1 you
// will replace main() with a real net/http server that wires dependencies
// (config, logger, router) and listens for requests.
package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/wchabir/taskflow/internal/handler"
)

func main() {
	r := gin.Default()
	r.GET("/health", handler.Health)
	r.StaticFile("/", "./static/index.html")
	srv := &http.Server{Addr: ":8080", Handler: r}
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server failed: %v", err)
		}
	}()
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	srv.Shutdown(shutdownCtx)
}
