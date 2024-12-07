package main

import (
	"bufio"
	"fmt"
	"os"
    "strings"
    "slices"
    "strconv"
    "os/exec"
)

var BUILTIN_LIST = []string {
    "exit",
    "echo",
    "type",
    "pwd",
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

        cmd := string(splt_input[0])
        is_builtin := is_valid_builtin(cmd)
        cmd_path := chk_path(cmd)
        if len(splt_input) == 0 {
           return 
        } else if is_builtin {
            if cmd == "exit" {
                status, err = strconv.Atoi(splt_input[1])
                break
            } else if cmd == "echo" {
                fmt.Fprintf(os.Stdout, "%s\n", splt_input[1])
            } else if cmd == "type" {
                chk := string(splt_input[1])

                is_builtin = is_valid_builtin(chk)
                cmd_path = chk_path(chk)
                if is_builtin {
                    fmt.Fprintf(os.Stdout, "%s is a shell builtin\n", chk)
                } else if cmd_path != "" {
                    fmt.Fprintf(os.Stdout, "%s is %s\n", chk, cmd_path)
                } else {
                    fmt.Fprintf(os.Stdout, "%s: not found\n", chk) 
                }
            } else if cmd == "pwd" {
                dir, err := os.Getwd()
                if err != nil {
                    panic(err)
                }
                fmt.Fprintf(os.Stdout, "%s\n", dir)
            }
        } else if cmd_path != "" {
            var args = []string{} 
            if len(splt_input) > 1 {
                args = strings.Split(string(splt_input[1]), " ")
            }
            out, err := exec.Command(cmd_path, args...).Output()
            if err != nil {
                panic(err)
            }
            fmt.Fprintf(os.Stdout, "%s", out)
        } else {
            fmt.Fprintf(os.Stdout, "%s: command not found\n", cmd)
        }

    }

    if status != 0 {
        panic("program ended with non-zero status")
    }
}

func is_valid_builtin(cmd string) bool {
    return slices.Contains(BUILTIN_LIST, cmd)
}

func chk_path(cmd string) string {
    paths := strings.Split(os.Getenv("PATH"), ":")

    for _, path := range paths {
        files, err := os.ReadDir(path)
        if err != nil {
            continue
        }

        for _, file := range files {
            if cmd == file.Name() {
                return path + "/" + file.Name()
            }
        }
    }

    return ""
}

