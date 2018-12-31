package main

import (
	"os"
	"fmt"
)


func main() {

	for _, username := range os.Args[1:] {
		fmt.Println("USERNAME:", username)
	}
}