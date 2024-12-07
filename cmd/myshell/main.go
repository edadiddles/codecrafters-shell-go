package main

import (
	"bufio"
	"fmt"
	"os"
    "strings"
    "slices"
    "strconv"
)

var CMD_LIST = []string {
    "exit",
    "echo",
}


func main() {
    status := 0
    for {
        fmt.Fprint(os.Stdout, "$ ")

        // Wait for user input
        input, err := bufio.NewReader(os.Stdin).ReadString('\n')
        if err != nil {
            panic(err)
        }

        input = strings.TrimSpace(input)
        splt_input := strings.SplitN(input, " ", 2)

        if len(splt_input) == 0 {
           return 
        }
        if !is_valid_cmd(splt_input[0]) {
            fmt.Fprintf(os.Stdout, "%s: command not found\n", splt_input[0])
        }

        cmd := string(splt_input[0])
        if cmd == "exit" {
            status, err = strconv.Atoi(splt_input[1])
            break
        }
        if cmd == "echo" {
            fmt.Fprintf(os.Stdout, "%s\n", splt_input[1])
        }
    }

    if status != 0 {
        panic("program ended with non-zero status")
    }
}

func is_valid_cmd(cmd string) bool {
    return slices.Contains(CMD_LIST, cmd)
}
