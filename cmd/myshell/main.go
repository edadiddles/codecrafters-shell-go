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
        splt_input := parse_inputs(input)
       // splt_input := strings.SplitN(input, " ", 2)

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

func parse_inputs(input string) []string {
    sp_enc := byte(32) // encoding for {white space}
    bs_enc := byte(92) // encoding for \
    sq_enc := byte(39) // encoding for '
    dq_enc := byte(34) // encoding for "
    ds_enc := byte(36) // encoding for $

    var output_args = []string{}
    var curr_arg = []byte{}
    is_escaped := false
    is_single_quoted := false
    is_double_quoted := false
    for i:=0; i < len(input); i++ {
        if !is_escaped && input[i] == bs_enc {
            if !is_double_quoted && !is_single_quoted {
                is_escaped = true
                continue
            } else if is_double_quoted && len(input) > i+1{
                pk := input[i+1]
                if pk == bs_enc || pk == dq_enc || pk == ds_enc {
                    is_escaped = true
                    continue
                }
            }
        } else if input[i] == sq_enc {
            if !is_double_quoted {
                if !is_escaped {
                    is_single_quoted = !is_single_quoted
                    continue
                }
            }

        } else if input[i] == dq_enc {
            if !is_single_quoted {
                if !is_escaped {
                    is_double_quoted = !is_double_quoted
                    continue
                }
            }
        } else if !is_escaped && !is_double_quoted && !is_single_quoted && input[i] == sp_enc {
            output_args = slices.Insert(output_args, len(output_args), string(curr_arg))
            curr_arg = []byte{}
            continue
        }
        
        curr_arg = slices.Insert(curr_arg, len(curr_arg), input[i])

        if is_escaped {
            is_escaped = false
        }
    }
    output_args = slices.Insert(output_args, len(output_args), string(curr_arg))
    return output_args
}
