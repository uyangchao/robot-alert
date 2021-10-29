package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"robot-alert/config"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-kit/log/level"
	"github.com/prometheus/common/promlog"
	"github.com/prometheus/common/version"
	"gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	os.Exit(run())
}

func run() int {
	var (
		sc = &config.SafeConfig{
			C: &config.Config{},
		}

		configFile    = kingpin.Flag("config.file", "Path to the configuration file.").Short('c').Default("config.yml").ExistingFile()
		listenAddress = kingpin.Flag("web.listen-address", "The address to listen on for web interface.").Default(":8060").String()
	)

	kingpin.HelpFlag.Short('h')
	kingpin.Parse()
	promlogConfig := &promlog.Config{}
	logger := promlog.New(promlogConfig)
	level.Info(logger).Log("msg", "Starting robot-alert", "version", version.Info())
	level.Info(logger).Log("msg", "Build context", version.BuildContext())

	if err := sc.ReloadConfig(*configFile, logger); err != nil {
		level.Error(logger).Log("msg", "Error loading config", "err", err)
		return 1
	}

	level.Info(logger).Log("msg", "Loaded config file")

	ctxWeb, cancelWeb := context.WithCancel(context.Background())
	defer cancelWeb()

	reloadCh := make(chan chan error)
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Post("/-/reload", func(w http.ResponseWriter, r *http.Request) {
		rc := make(chan error)
		reloadCh <- rc
		if err := <-rc; err != nil {
			http.Error(w, fmt.Sprintf("failed to reload config: %s", err), http.StatusInternalServerError)
		}
	})

	linstener, _ := net.Listen("tcp", *listenAddress)
	httpsrv := &http.Server{
		Handler: r,
	}
	srvCh := make(chan error)
	go func() {
		defer close(srvCh)
		level.Info(logger).Log("msg", "Start listening for connections", "address", listenAddress)
		if err := httpsrv.Serve(linstener); err != nil {
			level.Error(logger).Log("msg", "Error starting HTTP server", "err", err)
			srvCh <- err
		}

	}()

	var (
		hup  = make(chan os.Signal, 1)
		term = make(chan os.Signal, 1)
	)

	signal.Notify(hup, syscall.SIGHUP)
	signal.Notify(term, os.Interrupt, syscall.SIGTERM)

	go func() {
		for {
			select {
			case <-ctxWeb.Done():
				return
			case <-hup:
				if err := sc.ReloadConfig(*configFile, logger); err != nil {
					level.Error(logger).Log("msg", "Error reloading config", "err", err)
					continue
				}
				level.Info(logger).Log("msg", "Reloaded config file")
			case rc := <-reloadCh:
				if err := sc.ReloadConfig(*configFile, logger); err != nil {
					level.Error(logger).Log("msg", "Error reloading config", "err", err)
					rc <- err
				} else {
					level.Info(logger).Log("msg", "Reloaded config file")
					rc <- nil
				}
			}
		}
	}()

	for {
		select {
		case <-term:
			level.Info(logger).Log("msg", "Received SIGTERM, exiting gracefully...")
			httpsrv.Shutdown(ctxWeb)
			cancelWeb()
		case err := <-srvCh:
			if err != nil {
				return 1
			}

			return 0
		}
	}
}
