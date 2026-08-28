package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
)

func getWSLIP() (string, error) {
	out, err := exec.Command("wsl", "-d", "podman-machine-default", "ip", "addr", "show", "eth0").Output()
	if err != nil {
		return "", err
	}
	lines := strings.Split(string(out), "\n")
	for _, l := range lines {
		l = strings.TrimSpace(l)
		if strings.HasPrefix(l, "inet ") {
			parts := strings.Fields(l)
			if len(parts) >= 2 {
				ip := strings.Split(parts[1], "/")[0]
				return ip, nil
			}
		}
	}
	return "", fmt.Errorf("ip not found")
}

func forward(localPort int, remoteIP string, remotePort int) {
	listener, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", localPort))
	if err != nil {
		// Already listening or forwarded
		return
	}
	defer listener.Close()

	for {
		client, err := listener.Accept()
		if err != nil {
			time.Sleep(100 * time.Millisecond)
			continue
		}

		go func(c net.Conn) {
			defer c.Close()
			target, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", remoteIP, remotePort), 2*time.Second)
			if err != nil {
				return
			}
			defer target.Close()

			var wg sync.WaitGroup
			wg.Add(2)
			go func() {
				defer wg.Done()
				io.Copy(target, c)
			}()
			go func() {
				defer wg.Done()
				io.Copy(c, target)
			}()
			wg.Wait()
		}(client)
	}
}

func main() {
	ip, err := getWSLIP()
	if err != nil {
		fmt.Printf("Error getting WSL IP: %v\n", err)
		os.Exit(1)
	}

	ports := []int{8080, 9090, 3000, 5432}
	for _, p := range ports {
		go forward(p, ip, p)
	}

	// Keep alive
	select {}
}
