package main

import (
	"bufio"
	"fmt"
	"os"
    "strings"
    "slices"
    "strconv"
    "os/exec"
    "github.com/codecrafters-io/shell-starter-go/cmd/myshell/parsing"
)

var BUILTIN_LIST = []string {
    "exit",
    "echo",
    "type",
    "pwd",
    "cd",
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
        splt_input_word := parsing.Parse(input)

        splt_input := []string{}
        for i:=0; i < len(splt_input_word); i++ {
            splt_input = append(splt_input, splt_input_word[i].Val)
        }

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
                var args = []string{} 
                if len(splt_input) > 1 {
                    args = splt_input[1:]

                    for i:=len(args)-1; i>=0; i-- {
                        if strings.TrimSpace(args[i]) == "" {
                            args = slices.Delete(args, i, i+1)
                        }
                    }

                }
                out := strings.Join(args, " ")
                fmt.Fprintf(os.Stdout, "%s\n", out)
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
            } else if cmd == "cd" {
                chg_dir := string(splt_input[1])
                if chg_dir == "~" {
                    chg_dir = os.Getenv("HOME")
                }
                err := os.Chdir(chg_dir)
                if err != nil {
                    fmt.Fprintf(os.Stdout, "%s: No such file or directory\n", chg_dir)
                }

            }
        } else if cmd_path != "" {
            var args = []string{} 
            if len(splt_input) > 1 {
                args = splt_input[1:]

                for i:=len(args)-1; i>=0; i-- {
                    if strings.TrimSpace(args[i]) == "" {
                        args = slices.Delete(args, i, i+1)
                    }
                }
            }
            out, err := exec.Command(cmd_path, args...).Output()
            if err != nil {
                fmt.Print(len(args))
                fmt.Print(" ")
                fmt.Print(cmd_path)
                fmt.Print(" ")
                fmt.Print(args)
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
