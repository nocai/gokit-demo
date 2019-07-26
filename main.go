package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd/consul"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"nocai/gokit-demo/book"
	"nocai/gokit-demo/infra"
	"nocai/gokit-demo/infra/client/auth"
	"nocai/gokit-demo/infra/constants"
	"os"
	"os/signal"
	"syscall"
)

var (
	port       int
	consulAddr string

	configPath string

	ctx = context.Background()
)

func init() {
	flag.IntVar(&port, "http.port", 8888, "http port")
	flag.StringVar(&consulAddr, "consul.addr", "127.0.0.1:8500", "Consul address")
	flag.StringVar(&configPath, "consul.config.path", constants.AppName, "Config path at Consul")
	flag.Parse()
}

func main() {
	var l log.Logger
	{
		l = log.NewLogfmtLogger(os.Stdout)
		l = log.With(l, "ts", log.DefaultTimestampUTC)
		l = log.With(l, "caller", log.DefaultCaller)
	}

	var consulClient consul.Client
	{
		consulApi := infra.ConsulApi(l, consulAddr)
		infra.ConsulKv(l, consulApi, consulAddr, configPath)

		consulClient = consul.NewClient(consulApi)
		infra.ConsulResgister(l, consulClient, port)
	}

	{
		http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
			l.Log("msg", "ping")
			w.Write([]byte("pong"))
		})
		http.Handle("/metrics", promhttp.Handler())
	}

	authClient := auth.New(l, consulClient)

	var bs book.Service
	{
		bs = book.NewService(l, authClient)
	}




	http.Handle("/books/", accessControl(book.MakeHandler(bs, l)))

	errs := make(chan error)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		l.Log("transport", "HTTP", "port", port)
		errs <- http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	}()

	l.Log("exit", <-errs)
}

func accessControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")

		if r.Method == "OPTIONS" {
			return
		}

		h.ServeHTTP(w, r)
	})
}
