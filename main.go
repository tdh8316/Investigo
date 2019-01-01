package main

import (
	"os"
	"fmt"
	"./web"
)


func main() {
	for _, username := range os.Args[1:] {
		fmt.Printf("Searching username %s on:\n", username)
		web.Sniffer(username)
	}
}
