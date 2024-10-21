package daemon

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net"
	"net/http"
)

func RunServer(ctx context.Context) {
	serverAddr := ":8080"
	server := &http.Server{Addr: serverAddr}

	slog.Info("starting server", "server", server)

	go func() {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "Hello world")
		})

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
	if err := server.Close(); err != nil {
		slog.Error("server.Close", "err", err)
	}
}
