package config

import (
	"flag"
	"os"
)

var (
	Addr = ":80"
)

func init() {
	// look for environment variables or command line flags to override default port
	if addr, ok := os.LookupEnv("ADDR"); ok {
		Addr = addr
	}

	// look for cli flags
	flag.StringVar(&Addr, "addr", Addr, "address to listen on")
	flag.Parse()
}
