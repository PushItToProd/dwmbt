package daemon

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net"
	"net/http"
	"time"

	"github.com/pushittoprod/bt-daemon/pkg/bluetooth"
)

const (
	DefaultRequestTimeout  = 5 * time.Second
	DefaultShutdownTimeout = 5 * time.Second
)

type Peer struct {
	Addr        string
	DisplayName string
}

type Daemon struct {
	ServeAddr        string
	BluetoothManager bluetooth.BluetoothManager
	Peers            []Peer
	RequestTimeout   time.Duration
	ShutdownTimeout  time.Duration
}

func InitDaemon(d *Daemon) {
	if d.BluetoothManager == nil {
		d.BluetoothManager = bluetooth.NewBluetoothManager()
	}
	if d.RequestTimeout == 0 {
		d.RequestTimeout = DefaultRequestTimeout
	}
	if d.ShutdownTimeout == 0 {
		d.ShutdownTimeout = DefaultShutdownTimeout
	}
}

func (d Daemon) setupMux() http.Handler {
	mux := http.NewServeMux()

	// /_self/ endpoints only get data about our own devices

	// GET /_self/list lists Bluetooth devices connected to this host
	mux.HandleFunc("GET /_self/list", func(w http.ResponseWriter, r *http.Request) {
		devices, err := d.BluetoothManager.List(r.Context())
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

	// POST /_self/disconnect takes a form parameter `macAddr` and disconnects the device with that MAC address if
	// possible.
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

		// confirm the device is known and connected
		_, err = d.BluetoothManager.Get(r.Context(), macAddr)
		if errors.Is(err, bluetooth.ErrInvalidMac) {
			http.Error(w, "invalid MAC address", http.StatusBadRequest)
			return
		}
		if err != nil {
			// TODO: this probably means the device wasn't found, so we should
			// check the result and return a 404 instead of an internal server
			// error unless something actually went wrong
			slog.Error("d.BluetoothManager.Get", "err", err)
			http.Error(w, "failed to get device", http.StatusInternalServerError)
			return
		}

		// TODO: the TimeoutHandler used to wrap the mux should prevent this
		// from hanging forever, but it would be better to support a more async
		// approach. For example, if host A tells hosts B, C, and D to
		// disconnect device X, instead of waiting for the requests to complete
		// on each host, B, C, and D should reply immediately with either "that
		// device isn't connected" or "okay, disconnecting" and, in the latter
		// case, start disconnection on their end. Then, when disconnection
		// finishes, they should send another message to A saying "device X has
		// been disconnected".
		err = d.BluetoothManager.Disconnect(r.Context(), macAddr)
		if err != nil {
			http.Error(w, "failed to disconnect", http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "disconnected %q\n", macAddr) // TODO: return JSON
	})

	// top-level endpoints get data about our own devices and all peers

	// GET /list returns a list of all devices connected to this instance and its active peers.
	mux.HandleFunc("GET /list", func(w http.ResponseWriter, r *http.Request) {
		// TODO: list devices connected to this instance and all peers
		http.Error(w, "TODO", http.StatusNotImplemented)
	})

	return mux
}

func (d *Daemon) RunServer(ctx context.Context) {
	if d == nil {
		log.Panic("daemon is nil")
	}
	InitDaemon(d)

	mux := d.setupMux()

	// Wrap the mux in a timeout handler so requests won't hang forever.
	h := http.TimeoutHandler(mux, d.RequestTimeout, "timeout")
	server := &http.Server{
		Addr:    d.ServeAddr,
		Handler: h,
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
		if err := server.Serve(ln); err != nil {
			slog.Error("server.ListenAndServe", "err", err)
			return
		}
	}()

	// Wait for the server to stop.
	<-ctx.Done()

	// TODO: I guess we need to return a shutdown function that takes a context
	// to use during shutdown. For now, we'll just use a context with a 5 second
	// timeout as a placeholder.
	shutdownCtx, cancel := context.WithTimeout(context.Background(), d.ShutdownTimeout)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		slog.Error("server.Close", "err", err)
	}
}
