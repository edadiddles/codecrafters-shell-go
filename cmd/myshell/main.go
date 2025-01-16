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

type Redirect struct {
    Stdout string
    Stderr string
    ShallAppend bool
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
        redirect := Redirect{ Stdout: "", Stderr: "", ShallAppend: false }
        redirect.setup_redirect(splt_input)
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
                output := strings.Join(args, " ")
                out := fmt.Sprintf("%s\n", output)
                redirect.write_stdout(out)
            } else if cmd == "type" {
                chk := string(splt_input[1])

                is_builtin = is_valid_builtin(chk)
                cmd_path = chk_path(chk)
                if is_builtin {
                    out := fmt.Sprintf("%s is a shell builtin\n", chk)
                    redirect.write_stdout(out)
                } else if cmd_path != "" {
                    out := fmt.Sprintf("%s is %s\n", chk, cmd_path)
                    redirect.write_stdout(out)
                } else {
                    out := fmt.Sprintf("%s: not found\n", chk) 
                    redirect.write_stdout(out)
                }
            } else if cmd == "pwd" {
                dir, err := os.Getwd()
                if err != nil {
                    panic(err)
                }
                out := fmt.Sprintf("%s\n", dir)
                redirect.write_stdout(out)
            } else if cmd == "cd" {
                chg_dir := string(splt_input[1])
                if chg_dir == "~" {
                    chg_dir = os.Getenv("HOME")
                }
                err := os.Chdir(chg_dir)
                if err != nil {
                    out := fmt.Sprintf("%s: No such file or directory\n", chg_dir)
                    redirect.write_stdout(out)
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
            if len(args) == 0 {
                cmd_parts := strings.Split(cmd_path, "/")
                cmd := cmd_parts[len(cmd_parts)-1]
                output, err := exec.Command(cmd).Output()
                if err != nil {
                    out := fmt.Sprintf("%s: %s\n", cmd, "Error")
                    redirect.write_stderr(out)
                } else {
                    out := fmt.Sprintf("%s", output)
                    redirect.write_stdout(out)
                }
            }
            for _, arg := range args {
                output, err := exec.Command(cmd, arg).Output()
                if err != nil {
                    cmd := strings.Split(cmd_path, "/")
                    out := fmt.Sprintf("%s: %s: %s\n", cmd, arg, "No such file or directory")
                    redirect.write_stderr(out)
                } else {
                    out := fmt.Sprintf("%s", output)
                    redirect.write_stdout(out)
                }
            }
        } else {
            out := fmt.Sprintf("%s: command not found\n", cmd)
            redirect.write_stdout(out)
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

func (r *Redirect) setup_redirect(inpts []string) {
    r.Stdout = "stdout"
    r.Stderr = "stderr"
    r.ShallAppend = false
    for i:=len(inpts)-1; i > 0; i-- {
        if inpts[i] == ">" || inpts[i] == "1>" {
            r.Stdout = inpts[i+1]
            inpts = slices.Delete(inpts, i, i+2)
        } else if inpts[i] == "2>" {
            r.Stderr = inpts[i+1]
            inpts = slices.Delete(inpts, i, i+2)
        } else if inpts[i] == ">>" || inpts[i] == "1>>" {
            r.ShallAppend = true
            r.Stdout = inpts[i+1]
            inpts = slices.Delete(inpts, i, i+2)
        } else if inpts[i] == "2>>" {
            r.ShallAppend = true
            r.Stderr = inpts[i+1]
            inpts = slices.Delete(inpts, i, i+2)
        }
    }
}

func (r *Redirect) write_stdout(str string) {
    if r.Stdout == "stdout" {
        fmt.Fprint(os.Stdout, str)
    } else { 
        if r.ShallAppend {
            f, err := os.OpenFile(r.Stdout, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
            if err != nil {
                panic(err)
            }
            _, err = f.Write([]byte(str))
            if err != nil {
                f.Close()
                panic(err)
            }
            err = f.Close()
            if err != nil {
                panic(err)
            }
        } else {
            f, err := os.OpenFile(r.Stdout, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
            if err != nil {
                panic(err)
            }
            _, err = f.Write([]byte(str))
            if err != nil {
                f.Close()
                panic(err)
            }
            err = f.Close()
            if err != nil {
                panic(err)
            }
        }
    }
}

func (r *Redirect) write_stderr(str string) {
    if r.Stderr == "stderr" {
        fmt.Fprint(os.Stderr, str)
    } else { 
        if r.ShallAppend {
            f, err := os.OpenFile(r.Stderr, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
            if err != nil {
                panic(err)
            }
            _, err = f.Write([]byte(str))
            if err != nil {
                f.Close()
                panic(err)
            }
            err = f.Close()
            if err != nil {
                panic(err)
            }
        } else {
            f, err := os.OpenFile(r.Stderr, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0777)
            if err != nil {
                panic(err)
            }
            _, err = f.Write([]byte(str))
            if err != nil {
                f.Close()
                panic(err)
            }
            err = f.Close()
            if err != nil {
                panic(err)
            }
        }
    }
}
