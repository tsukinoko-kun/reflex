package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"example/out/frontend"
	"example/src/backend/config"
	"example/src/public"
)

func main() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	sock, err := net.Listen("tcp", config.Addr)
	if err != nil {
		log.Fatal(err)
	}

	defer fmt.Println("reflex exit")
	defer sock.Close()

	fmt.Printf("Listening on %s\n", sock.Addr())

	mux := http.NewServeMux()
	mux.HandleFunc("/", router)
	go http.Serve(sock, mux)

	<-sigs
}

func router(w http.ResponseWriter, r *http.Request) {
	if public.Handler(w, r) {
		return
	}
	if frontend.Handler(w, r) {
		return
	}

	fmt.Fprintf(w, "Hello, World!")
}
