// === builder.go ===
package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	key := generateKey()
	fmt.Println("Generated Key:", key)
	templateFiles := []string{"gs-client.go", "gs-listener.go"}
	for _, file := range templateFiles {
		input, err := os.ReadFile(file)
		if err != nil {
			fmt.Println("Failed to read", file)
			continue
		}
		updated := strings.ReplaceAll(string(input), "<SHARED_KEY_PLACEHOLDER>", key)
		os.WriteFile(file, []byte(updated), 0644)
		fmt.Println("Injected key into", file)
	}
}
