package main

import (
	"bremlin/parser"
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/k0kubun/pp"
)

func main() {
	fmt.Println("Enter Bremlin command:")
	scanner := bufio.NewScanner(os.Stdin)
	var lines []string
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			break
		}
		lines = append(lines, line)
	}
	cmd := strings.Join(lines, "\n")

	chain, err := parser.ParseCommand(cmd)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	} else {
		pp.Println(chain)
	}
}
