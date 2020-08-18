package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/lupinthe14th/todo/intenal/todo"
)

const (
	ExitOK    = 0
	ExitError = 1
)

func main() {
	var addr string
	flag.StringVar(&addr, "addr", ":8080", "host:port")
	flag.Parse()
	os.Exit(run(addr))
}

func run(addr string) int {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	errCh := make(chan error)

	db := todo.NewMemoryDB()
	srv := todo.NewServer(addr, db)

	go func() {
		errCh <- srv.Start()
	}()

	select {
	case err := <-errCh:
		if err != nil {
			log.Println(err)
			return ExitError
		}
	case <-sigCh:
		if err := srv.Stop(context.Background()); err != nil {
			log.Println(err)
			return ExitError
		}
	}
	return ExitOK
}
