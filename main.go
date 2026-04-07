package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
)

func handleConn(conn net.Conn) {
	defer conn.Close()
	addr := conn.RemoteAddr().String()
	log.Printf("connection from %s", addr)
	if _, err := io.Copy(conn, conn); err != nil {
		log.Printf("connection %s closed: %v", addr, err)
		return
	}
	log.Printf("connection %s closed", addr)
}

func serve(ln net.Listener) {
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("accept error: %v", err)
			return
		}
		go handleConn(conn)
	}
}

func main() {
	port := flag.Int("port", 9095, "TCP port to listen on")
	flag.Parse()

	addr := fmt.Sprintf(":%d", *port)
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen on %s: %v", addr, err)
	}
	log.Printf("listening on %s", addr)
	serve(ln)
}
