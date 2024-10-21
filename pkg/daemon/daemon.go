package daemon

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net"
	"net/http"
	"time"
)

func setupMux() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello ServeMux")
	})
	return mux
}

func RunServer(ctx context.Context) {
	mux := setupMux()

	serverAddr := ":8080"
	server := &http.Server{
		Addr:    serverAddr,
		Handler: mux,
	}

	slog.Info("starting server", "server", server)

	go func() {
		ln, err := net.Listen("tcp", server.Addr)
		if err != nil {
			slog.Error("net.Listen", "err", err)
			log.Panicf("net.Listen failed: %v", err)
			return
		}
		slog.Info("starting server", "addr", ln.Addr().String())
		log.Printf("starting server: http://%s", ln.Addr().String())
		if err := server.Serve(ln); err != nil {
			slog.Error("server.ListenAndServe", "err", err)
			return
		}
	}()

	<-ctx.Done()

	// TODO: I guess we need to return a shutdown function that takes a context
	// to use during shutdown. For now, we'll just use a context with a 5 second
	// timeout as a placeholder.
	shutdownCtx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	if err := server.Shutdown(shutdownCtx); err != nil {
		slog.Error("server.Close", "err", err)
	}
}
