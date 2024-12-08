package daemon

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net"
	"net/http"
	"time"

	"github.com/pushittoprod/bt-daemon/pkg/bluetooth"
	"github.com/pushittoprod/bt-daemon/pkg/utils"
)

type Peer struct {
	Addr        string
	DisplayName string
}

type Daemon struct {
	ServeAddr        string
	BluetoothManager bluetooth.BluetoothManager
	Peers            []Peer
}

func (d Daemon) setupMux() http.Handler {
	mux := http.NewServeMux()

	// /_self/ endpoints only get data about our own devices
	mux.HandleFunc("GET /_self/list", func(w http.ResponseWriter, r *http.Request) {
		devices, err := d.BluetoothManager.List()
		if err != nil {
			slog.Error("d.BluetoothManager.List", "err", err)
			http.Error(w, "error listing bluetooth devices", http.StatusInternalServerError)
			return
		}

		j, err := json.MarshalIndent(devices, "", "  ")
		if err != nil {
			slog.Error("json.MarshalIndent", "err", err)
			http.Error(w, "error formatting json", http.StatusInternalServerError)
			return
		}

		_, err = w.Write(j)
		if err != nil {
			slog.Error("w.Write", "err", err)
			return
		}
	})
	mux.HandleFunc("POST /_self/disconnect", func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			slog.Error("r.ParseForm", "err", err)
			http.Error(w, "could not parse request", http.StatusBadRequest)
			return
		}
		macAddr := r.FormValue("macAddr")
		if macAddr == "" {
			http.Error(w, "macAddr param missing or blank", http.StatusBadRequest)
			return
		}
		macAddr, ok := utils.NormalizeMac(macAddr)
		if !ok {
			http.Error(w, "invalid MAC address", http.StatusBadRequest)
			return
		}

		// confirm the device is known and connected
		_, err = d.BluetoothManager.Get(macAddr)
		if err != nil {
			// TODO: this probably means the device wasn't found, so we should
			// check the result and return a 404 instead of an internal server
			// error unless something actually went wrong
			http.Error(w, "failed to get device", http.StatusInternalServerError)
			return
		}

		// TODO: use a context with timeout here to avoid waiting forever if this hangs
		err = d.BluetoothManager.Disconnect(macAddr)
		if err != nil {
			http.Error(w, "failed to disconnect", http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "disconnected %q\n", macAddr)

	})

	// top-level endpoints get data about our own devices and all peers
	mux.HandleFunc("GET /list", func(w http.ResponseWriter, r *http.Request) {
		// TODO: list devices connected to this instance and all peers
		http.Error(w, "TODO", http.StatusNotImplemented)
	})

	return mux
}

func (d Daemon) RunServer(ctx context.Context) {
	mux := d.setupMux()

	server := &http.Server{
		Addr:    d.ServeAddr,
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
