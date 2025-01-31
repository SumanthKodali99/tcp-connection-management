package main

import (
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strings"
	"sync"
	"time"
)

const (
	port             = ":8080"
	simulatedDelay   = 20 * time.Second
	keepAliveTimeout = 120
)

var (
	mu          sync.Mutex
	activeConns int
)

func countActiveConnections() {
	cmd := exec.Command("sh", "-c", "netstat -tn | grep ':8080' | grep ESTABLISHED | wc -l")
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("[SERVER] Error executing netstat:", err)
		return
	}

	connectionsCount := strings.TrimSpace(string(output))
	fmt.Printf("[SERVER] Active TCP Connections: %s\n", connectionsCount)
}

func trackServerSockets() {
	for {
		mu.Lock()
		countActiveConnections()
		mu.Unlock()
		time.Sleep(5 * time.Second)
	}
}

func slowResponseHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	activeConns++
	fmt.Printf("[SERVER] New HTTP request received. Active connections: %d\n", activeConns)
	mu.Unlock()

	time.Sleep(simulatedDelay)

	select {
	case <-r.Context().Done():
		fmt.Println("[SERVER] Client closed connection before response. Aborting write.")
		return
	default:
	}

	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Keep-Alive", fmt.Sprintf("timeout=%d, max=100", keepAliveTimeout))

	mu.Lock()
	fmt.Println("[SERVER] Sending delayed response to client.")
	activeConns--
	mu.Unlock()

	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("Delayed response from server."))
	if err != nil {
		fmt.Println("[SERVER] Error writing response:", err)
	}
}

func main() {
	go trackServerSockets()

	server := &http.Server{
		Addr:    port,
		Handler: http.HandlerFunc(slowResponseHandler),
	}

	fmt.Println("[SERVER] HTTP Server started on", port)
	log.Fatal(server.ListenAndServe())
}
