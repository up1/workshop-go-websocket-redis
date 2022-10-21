package main

import (
	"context"
	"demo/process"
	"demo/routes"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const port = "8080"

func main() {
	http.HandleFunc("/healthz", func(rw http.ResponseWriter, req *http.Request) {
		rw.Write([]byte("OK"))
	})
	http.Handle("/ws/", http.HandlerFunc(routes.WebSocketHandler))
	server := http.Server{Addr: ":" + port, Handler: nil}
	go func() {
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatal("failed to start server", err)
		}
	}()

	exit := make(chan os.Signal)
	signal.Notify(exit, syscall.SIGTERM, syscall.SIGINT)
	<-exit

	log.Println("exit signalled")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	process.Cleanup()
	server.Shutdown(ctx)

	log.Println("chat app exited")
}