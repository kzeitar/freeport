// Package main provides a helper program for testing.
// It listens on a TCP port and blocks indefinitely.
package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"time"
)

var portFlag = flag.Int("port", 0, "Port to listen on (0 for random available port)")

func main() {
	flag.Parse()

	listener, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", *portFlag))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("listening on %s\n", listener.Addr())

	// Block indefinitely
	for {
		time.Sleep(time.Hour)
	}
}
