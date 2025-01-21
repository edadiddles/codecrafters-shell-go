package main

import (
	"fmt"
	"os"
    "strings"
    "slices"
    "strconv"
    "os/exec"
    "github.com/codecrafters-io/shell-starter-go/cmd/myshell/parsing"
    "syscall"
    "unsafe"
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

        
        fd := int(os.Stdin.Fd())

        // Save the original terminal state
        originalTermios, err := getTermios(fd)
        if err != nil {
            fmt.Println("Error getting terminal attributes:", err)
            return
        }
        defer setTermios(fd, originalTermios) // Restore on exit

        // Set terminal to raw mode
        raw := *originalTermios
        raw.Lflag &^= syscall.ECHO | syscall.ICANON // Disable echo and canonical mode
        raw.Cc[syscall.VMIN] = 1                    // Minimum number of characters to read
        raw.Cc[syscall.VTIME] = 0                   // No timeout
        if err := setTermios(fd, &raw); err != nil {
            fmt.Println("Error setting terminal to raw mode:", err)
            return
        }


        buffer := make([]byte, 1)
        input_buffer := make([]byte, 0)
        splt_input := []string{}
        tab_pressed := false
        for {
            _, err := os.Stdin.Read(buffer)
            if err != nil {
                fmt.Println("Error reading input:", err)
                break
            }

            if buffer[0] == '\t' {
                if len(input_buffer) == 0 {
                    continue
                }
                splt_input = split_input(string(input_buffer))
                autocomplete_list := chk_cmd_autocomplete(splt_input[0])
                if len(autocomplete_list) == 0 {
                    tab_pressed = false
                    fmt.Print("\a")
                } else if len(autocomplete_list) == 1 {
                    tab_pressed = false
                    input_buffer = []byte(autocomplete_list[0] + " ")
                    fmt.Print("\r\x1b[K")
                    fmt.Printf("$ %s", input_buffer) 
                } else if len(autocomplete_list) > 1 {
                    if !tab_pressed {
                        p_cmd := find_common_name(autocomplete_list)
                        if p_cmd != string(input_buffer) { 
                            input_buffer = []byte(p_cmd)
                            fmt.Print("\r\x1b[K")
                            fmt.Printf("$ %s", input_buffer) 
                        } else {
                            tab_pressed = true
                            fmt.Print("\a")
                        }
                    } else if tab_pressed {
                        tab_pressed = false
                        fmt.Print("\n")
                        cmd_list_str:= ""
                        for _, autocomplete := range autocomplete_list {
                            cmd_list_str += fmt.Sprintf("%s  ", autocomplete)
                        }
                        fmt.Printf("%s\n", strings.TrimSpace(cmd_list_str))
                        fmt.Printf("$ %s", input_buffer)
                    }
                }
            } else if buffer[0] == '\n' {
                fmt.Print("\n")
                splt_input = split_input(string(input_buffer))
                break
            } else if buffer[0] == '\x7f' {
                tab_pressed = false
                if len(input_buffer) > 0 {
                    fmt.Print("\r\x1b[K")
                    input_buffer = input_buffer[0:len(input_buffer)-1]
                    fmt.Printf("$ %s", input_buffer) 
                }
            } else {
                tab_pressed = false
                input_buffer = append(input_buffer, buffer[0])
                fmt.Printf("%s", string(buffer[0]))
            }
        }
        status = execute(splt_input)

        if status != -1 {
            break
        }
    }

    if status != 0 {
        panic("program ended with non-zero status")
    }
}

func separate_args(args []string) ([]string, []string) {
    global_args := make([]string, 0)
    input_args := make([]string, 0)
    for _, arg := range args {
        if len(arg) == 0 {
            continue
        } else if arg[0] == '-' {
            global_args = append(global_args, arg)
        } else {
            input_args = append(input_args, arg)
        }
    }

    return global_args, input_args
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
            r.write_stdout("")
        } else if inpts[i] == "2>" {
            r.Stderr = inpts[i+1]
            inpts = slices.Delete(inpts, i, i+2)
            r.write_stderr("")
        } else if inpts[i] == ">>" || inpts[i] == "1>>" {
            r.ShallAppend = true
            r.Stdout = inpts[i+1]
            inpts = slices.Delete(inpts, i, i+2)
            r.write_stdout("")
        } else if inpts[i] == "2>>" {
            r.ShallAppend = true
            r.Stderr = inpts[i+1]
            inpts = slices.Delete(inpts, i, i+2)
            r.write_stderr("")
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

func chk_cmd_autocomplete(p_cmd string) []string {
    autocompletes := make([]string, 0)
    t_cmd := strings.TrimSpace(p_cmd)
    
    //need to check builtins
    for _, builtin := range BUILTIN_LIST {
        if strings.Index(builtin, t_cmd) == 0 {
            if chk_uniqueness(autocompletes, builtin) {
                autocompletes = append(autocompletes, builtin)
            }
        }
    }

    //need to check PATH
    paths := strings.Split(os.Getenv("PATH"), ":")
    for _, path := range paths {
        files, _ := os.ReadDir(path)
        for _, file := range files {
            if strings.Index(file.Name(), t_cmd) == 0 {
                if chk_uniqueness(autocompletes, file.Name()) {
                    autocompletes = append(autocompletes, file.Name())
                }
            }
        }
    }

    //need to check relative path
    files, _ := os.ReadDir(t_cmd)
    for _, file := range files {
        if strings.Index(file.Name(), t_cmd) == 0 {
            if chk_uniqueness(autocompletes, file.Name()) {
                autocompletes = append(autocompletes, file.Name())
            }
        }
    }
    slices.Sort(autocompletes)
    return autocompletes
}

func chk_uniqueness(cmds []string, p_cmd string) bool {
    is_unique := true
    for _, cmd := range cmds {
       if p_cmd == cmd {
           is_unique = false
           break
       }
   }
   return is_unique
}

// Termios structure for terminal attributes
type termios struct {
	Iflag  uint32
	Oflag  uint32
	Cflag  uint32
	Lflag  uint32
	Cc     [20]byte
	Ispeed uint32
	Ospeed uint32
}

// Get the current terminal attributes
func getTermios(fd int) (*termios, error) {
	var t termios
	_, _, errno := syscall.Syscall6(syscall.SYS_IOCTL, uintptr(fd), uintptr(syscall.TCGETS), uintptr(unsafe.Pointer(&t)), 0, 0, 0)
	if errno != 0 {
		return nil, errno
	}
	return &t, nil
}

// Set terminal attributes
func setTermios(fd int, t *termios) error {
	_, _, errno := syscall.Syscall6(syscall.SYS_IOCTL, uintptr(fd), uintptr(syscall.TCSETS), uintptr(unsafe.Pointer(t)), 0, 0, 0)
	if errno != 0 {
		return errno
	}
	return nil
}
func split_input(input string) []string {
    input = strings.TrimSpace(input)
    splt_input_word := parsing.Parse(input)

    splt_input := []string{}
    for i:=0; i < len(splt_input_word); i++ {
        splt_input = append(splt_input, splt_input_word[i].Val)
    }

    return splt_input
}

func execute(splt_input []string) int {
    status := -1

    cmd := string(splt_input[0])
    is_builtin := is_valid_builtin(cmd)
    redirect := Redirect{ Stdout: "", Stderr: "", ShallAppend: false }
    redirect.setup_redirect(splt_input)
    cmd_path := chk_path(cmd)

    // Wait for user input
    if len(splt_input) == 0 {
       return  status
    } else if is_builtin {
        if cmd == "exit" {
            status, _ = strconv.Atoi(splt_input[1])
            return status
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
        cmd_parts := strings.Split(cmd_path, "/")
        cmd := cmd_parts[len(cmd_parts)-1]
        if len(args) == 0 {
            output, err := exec.Command(cmd).Output()
            if err != nil {
                out := fmt.Sprintf("%s: %s\n", cmd, "Error")
                redirect.write_stderr(out)
            } else {
                out := fmt.Sprintf("%s", output)
                redirect.write_stdout(out)
            }
        }
        g_args, i_args := separate_args(args)
        for _, arg := range i_args {
            c_args := append(g_args, arg)
            output, err := exec.Command(cmd, c_args...).Output()
            if err != nil {
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

    return status
}

func find_common_name(cmd_list []string) string {
    common_cmd := ""

    is_common := true
    for i := 0; is_common; i++ {
        var curr_letter byte
        for j:=0; j < len(cmd_list); j++ {
            if i >= len(cmd_list[j]) {
                is_common = false
                break
            } else if j == 0 {
                curr_letter = cmd_list[j][i]
            } else if curr_letter != cmd_list[j][i] {
                is_common = false
                break
            }
        }

        if is_common {
            common_cmd += string(curr_letter)
        }
    }

    return common_cmd
}
            
