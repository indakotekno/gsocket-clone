// === gserver.go ===
package main

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"net"
	"strings"
	"sync"
)

var (
	clientConn   net.Conn
	listenerConn net.Conn
	sharedKey    = generateKey()
	mu           sync.Mutex
)

func generateKey() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func encryptAES(key, plaintext []byte) []byte {
	block, _ := aes.NewCipher(key)
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	rand.Read(iv)
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)
	return ciphertext
}

func decryptAES(key, ciphertext []byte) []byte {
	block, _ := aes.NewCipher(key)
	iv := ciphertext[:aes.BlockSize]
	plaintext := make([]byte, len(ciphertext[aes.BlockSize:]))
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(plaintext, ciphertext[aes.BlockSize:])
	return plaintext
}

func main() {
	ln, err := net.Listen("tcp", ":443")
	if err != nil {
		panic(err)
	}
	fmt.Println("[+] Relay server listening on :443")
	fmt.Println("[+] Shared key:", sharedKey)

	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}
		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	roleLine, _ := reader.ReadString('\n')
	keyLine, _ := reader.ReadString('\n')

	if !strings.HasPrefix(roleLine, "ROLE:") || !strings.HasPrefix(keyLine, "KEY:") {
		return
	}

	role := strings.TrimSpace(strings.TrimPrefix(roleLine, "ROLE:"))
	key := strings.TrimSpace(strings.TrimPrefix(keyLine, "KEY:"))

	if key != sharedKey {
		fmt.Println("[!] Invalid key from", conn.RemoteAddr())
		return
	}

	mu.Lock()
	defer mu.Unlock()

	switch role {
	case "client":
		clientConn = conn
		fmt.Println("[*] Client connected:", conn.RemoteAddr())
		waitForListener()
	case "listener":
		listenerConn = conn
		fmt.Println("[*] Listener connected:", conn.RemoteAddr())
		waitForClient()
	}
}

func waitForListener() {
	if listenerConn != nil {
		go bridge(clientConn, listenerConn)
	}
}

func waitForClient() {
	if clientConn != nil {
		go bridge(listenerConn, clientConn)
	}
}

func bridge(a, b net.Conn) {
	defer a.Close()
	defer b.Close()
	go io.Copy(a, b)
	io.Copy(b, a)
}
