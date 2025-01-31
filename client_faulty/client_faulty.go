package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"math/rand"
	"net/http"
	"os/exec"
	"strings"
	"time"
)

const requestURL = "http://localhost:8080"

func NewRestClient() *http.Client {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client.Transport = tr
	return client
}

func countOpenSockets() {
	cmd := exec.Command("sh", "-c", "netstat -tn | grep '8080' || true")
	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		fmt.Println("Failed to execute netstat:", err)
		return
	}

	rawOutput := strings.TrimSpace(out.String())
	lines := strings.Split(rawOutput, "\n")

	if len(lines) == 1 && lines[0] == "" {
		fmt.Println("No open TCP connections on port 8080.")
		return
	}

	established := 0
	for _, line := range lines {
		if strings.Contains(line, "ESTABLISHED") {
			established++
		}
	}

	fmt.Printf("[CLIENT] Total TCP Connections: %d | ESTABLISHED: %d\n", len(lines), established)
}

func makeRequest(label string) {
	client := NewRestClient()

	fmt.Printf("[CLIENT] (%s) Sending request to: %s\n", label, requestURL)

	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		fmt.Printf("[CLIENT] (%s) Request creation failed: %v\n", label, err)
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("[CLIENT] (%s) Request failed: %v\n", label, err)
		return
	}
	defer resp.Body.Close()

	fmt.Printf("[CLIENT] (%s) Connection established, monitoring TCP connections...\n", label)

	buffer := make([]byte, 1024)

	for {
		_, err := resp.Body.Read(buffer)
		if err != nil {
			fmt.Printf("[CLIENT] (%s) Server closed connection or read error: %v\n", label, err)
			return
		}

		countOpenSockets()
		time.Sleep(5 * time.Second)
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	for {
		countOpenSockets()
		go makeRequest("Main")

		randomDelay := time.Duration(rand.Intn(20)) * time.Second
		if rand.Intn(2) == 0 {
			time.Sleep(randomDelay)
			fmt.Printf("[CLIENT] Sending an additional request after %v seconds\n", randomDelay)
			go makeRequest("Extra")
		}

		time.Sleep(20 * time.Second)
	}
}
