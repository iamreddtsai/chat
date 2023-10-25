package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/iamreddtsai/chat/cmd/httpserver"
	"github.com/iamreddtsai/chat/cmd/websocket/app"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-quit
		cancel()
	}()

	a := app.New()
	srv := httpserver.New(httpserver.ServeMux(a.HttpHandler))
	srv.Execute(ctx)
}
