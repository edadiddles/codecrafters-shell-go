package main

import (
	"bufio"
	"fmt"
	"os"
    "strings"
)

func main() {
    for {
        repl()
    }
}

func repl() {
    fmt.Fprint(os.Stdout, "$ ")

	// Wait for user input
    input, err := bufio.NewReader(os.Stdin).ReadString('\n')
    if err != nil {
        panic(err)
    }

    input = strings.TrimSpace(input)
    if !is_valid_cmd(input) {
        fmt.Fprintf(os.Stdout, "%s: command not found\n", input)
    }
}

func is_valid_cmd(cmd string) bool {
    return false
}
