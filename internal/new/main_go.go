package new

import (
	"bytes"
	"text/template"
)

type mainGoData struct {
	Title      string
	BackendDir string
	PublicDir  string
	OutDir     string
}

func mainGo(data mainGoData) string {
	tmpl := `package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"{{.Title}}/{{.OutDir}}/frontend"
	"{{.Title}}/{{.BackendDir}}/config"
	"{{.Title}}/{{.PublicDir}}"
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
`
	t := template.Must(template.New("main").Parse(tmpl))
	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return ""
	}
	return buf.String()
}
