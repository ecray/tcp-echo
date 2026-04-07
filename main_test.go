package main

import (
	"net"
	"testing"
	"time"
)

func startTestServer(t *testing.T) net.Listener {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to start test server: %v", err)
	}
	go serve(ln)
	return ln
}

func TestEcho(t *testing.T) {
	ln := startTestServer(t)
	defer ln.Close()

	conn, err := net.DialTimeout("tcp", ln.Addr().String(), 2*time.Second)
	if err != nil {
		t.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()

	msg := []byte("hello, echo!")
	if _, err := conn.Write(msg); err != nil {
		t.Fatalf("write error: %v", err)
	}

	buf := make([]byte, len(msg))
	if _, err := conn.Read(buf); err != nil {
		t.Fatalf("read error: %v", err)
	}

	if string(buf) != string(msg) {
		t.Errorf("expected %q, got %q", msg, buf)
	}
}

func TestMultipleMessages(t *testing.T) {
	ln := startTestServer(t)
	defer ln.Close()

	conn, err := net.DialTimeout("tcp", ln.Addr().String(), 2*time.Second)
	if err != nil {
		t.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()

	messages := []string{"ping", "pong", "foo", "bar"}
	for _, msg := range messages {
		if _, err := conn.Write([]byte(msg)); err != nil {
			t.Fatalf("write error: %v", err)
		}
		buf := make([]byte, len(msg))
		if _, err := conn.Read(buf); err != nil {
			t.Fatalf("read error: %v", err)
		}
		if string(buf) != msg {
			t.Errorf("expected %q, got %q", msg, buf)
		}
	}
}

func TestMultipleClients(t *testing.T) {
	ln := startTestServer(t)
	defer ln.Close()

	done := make(chan struct{}, 3)
	for i := 0; i < 3; i++ {
		go func(id int) {
			defer func() { done <- struct{}{} }()
			conn, err := net.DialTimeout("tcp", ln.Addr().String(), 2*time.Second)
			if err != nil {
				t.Errorf("client %d: failed to connect: %v", id, err)
				return
			}
			defer conn.Close()

			msg := []byte("client message")
			conn.Write(msg)
			buf := make([]byte, len(msg))
			conn.Read(buf)
			if string(buf) != string(msg) {
				t.Errorf("client %d: expected %q, got %q", id, msg, buf)
			}
		}(i)
	}
	for i := 0; i < 3; i++ {
		<-done
	}
}
