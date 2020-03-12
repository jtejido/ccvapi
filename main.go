package main

import (
	"context"
	"flag"
	"github.com/jtejido/ccvapi/config"
	"github.com/jtejido/ccvapi/logging"
	"github.com/jtejido/ccvapi/sanitation"
	"github.com/jtejido/ccvapi/server"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

var (
	conf                      *config.Config
	accessLogger, errorLogger logging.Logging
)

func main() {
	conf, err := config.LoadConfig("config.toml")
	if err != nil {
		panic(err)
	}

	// cmd flags
	flag.IntVar(&conf.Http.Host, "host", conf.Http.Host, "The local port to listen to.")
	flag.StringVar(&conf.Http.ErrorLog, "error-log", conf.Http.ErrorLog, "Location of the logfile.")
	flag.StringVar(&conf.Http.AccessLog, "access-log", conf.Http.AccessLog, "Location of the logfile.")
	flag.StringVar(&conf.Http.CardTypesPath, "card-path", conf.Http.CardTypesPath, "Location of the card types json file.")
	flag.Parse()

	err = sanitation.LoadCards(conf.Http.CardTypesPath)
	if err != nil {
		panic(err)
	}

	// log files
	accessLogger, _ = logging.NewLogger(conf.Http.AccessLog)
	errorLogger, _ = logging.NewLogger(conf.Http.ErrorLog)

	// just localhost for now
	frontend := ":" + strconv.Itoa(conf.Http.Host)

	s := server.NewServeMux(func(s *server.ServeMux) {
		s.Logger = accessLogger
	})

	srvr := &http.Server{Addr: frontend, Handler: s}
	accessLogger.Printf("starting http listening on " + frontend)

	go func() {
		if err := srvr.ListenAndServe(); err != http.ErrServerClosed {
			panic(err.Error())
		}
	}()

	select {}

	graceful(srvr, 5*time.Second)

}

func graceful(srvr *http.Server, timeout time.Duration) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, os.Kill, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := srvr.Shutdown(ctx); err != nil {
		errorLogger.Printf("server Error: %v\n", err)
	} else {
		accessLogger.Printf("http listening stopped. \n")
	}
}
