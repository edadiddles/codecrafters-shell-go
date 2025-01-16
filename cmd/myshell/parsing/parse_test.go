package parsing

import (
    "testing"
)


func TestParseSimpleInput(t *testing.T) {
    words := Parse(" apple banana carrot ")

    if words[0].Val != "apple" || words[1].Val != "banana" || words[2].Val != "carrot" {
        t.Fatalf("not parsed correctly: expected [apple, bannana, carrot] -- actual [%s, %s, %s]", words[0].Val, words[1].Val, words[2].Val)
    }
}


func TestParseQuotedInput(t *testing.T) {
    words := Parse(" apple \"banana\" 'carrot' ")

    if words[0].Val != "apple" || words[1].Val != "banana" || words[2].Val != "carrot" {
        t.Fatalf("not parsed correctly: expected [apple, bannana, carrot] -- actual [%s, %s, %s]", words[0].Val, words[1].Val, words[2].Val)
    }
}

func TestParseQuotedInputWithWhiteSpace(t *testing.T) { 
    words := Parse("echo \"hello world\"")
    
    if words[0].Val != "echo" || words[1].Val != "hello world" {
        t.Fatalf("not parsed correctly: expected [echo, hello world] -- actual [%s, %s]", words[0].Val, words[1].Val)
    }
}

func TestParseQuotedInputWithEscapes(t *testing.T) {
    words := Parse("'$ \\\"string\\\"'")

    if words[0].Val != "$ \\\"string\\\"" { 
        t.Fatalf("not parsed correctly: expected [$ \\\"string\\\"] -- actual [%v]", words[0].Val)
    }
}

func TestParseCodeCraftersSingleQuotes(t *testing.T) {
    words := Parse("'shell hello' 'world     test'")

    if !(words[0].Val == "shell hello") || !(words[1].Val == "world     test") {
        t.Fatalf("not parsed correctly: expected [hello word, world     test] -- actual [%v, %v]", words[0].Val, words[1].Val)
    }
}

func TestParseCodeCraftersDoubleQuotes(t *testing.T) {
    words := Parse("\"quz  hello\"  \"bar\" \"bar\"  \"shell's\" \"foo\"")

    if !(words[0].Val == "quz  hello") || !(words[1].Val == "bar") || !(words[2].Val == "bar") || !(words[3].Val == "shell's") || !(words[4].Val == "foo") {
        t.Fatalf("not parsed correctly: expected [quz  hello, bar, bar, shell's, foo] -- actual [%v, %v, %v, %v, %v]", words[0].Val, words[1].Val, words[2].Val, words[3].Val, words[4].Val)
    }
}

func TestParseCodeCraftersEscapedUnquoted(t *testing.T) {
    words := Parse("before\\ \\ \\ \\ \\ \\ script")

    if !(words[0].Val == "before      script") {
        t.Fatalf("not parsed correctly: expected [before      script] -- actual [%v]", words[0].Val)
    }
}

func TestParseCodeCraftersEscapedSingleQuoted(t *testing.T) {
    words := Parse("'shell\\\\\\nscript'")

    if !(words[0].Val == "shell\\\\\\nscript") {
        t.Fatalf("not parsed correctly: expected [shell\\\\\\nscript] -- actual [%v]", words[0].Val)
    }
}

func TestParseCodeCraftersComplexSingleQuoted(t *testing.T) {
    words := Parse("'test   example' 'shell''hello'")

    if words[0].Val != "test   example shellhello" {
        t.Fatalf("not parsed correctly: expected [test   example shellhello] -- actual [%s]", words[0].Val)
    }
}

func TestParseCodeCraftersEscapedDoubleQuoted(t *testing.T) {
    words := Parse("'shell\\\\\\nscript'")

    if !(words[0].Val == "shell\\\\\\nscript") {
        t.Fatalf("not parsed correctly: expected [shell\\\\\\nscript] -- actual [%v]", words[0].Val)
    }
}

func TestParseCodeCraftersRedirectOutputDefault(t *testing.T) {
    words := Parse("ls /tmp/baz > /tmp/foo/baz.md")
    if words[0].Val != "ls" || words[1].Val != "/tmp/baz" || words[2].Val != ">" || words[3].Val != "/tmp/foo/baz.md" {
        t.Fatalf("not parsed correctly: expected [ls, /tmp/baz, >, /tmp/foo/baz.md] -- actual [%v, %v, %v, %v]", words[0].Val, words[1].Val, words[2].Val, words[3].Val)
    }
}

func TestParseCodeCraftersRedirectOutputStdin(t *testing.T) {
    words := Parse("echo 'Hello James' 1> /tmp/foo/baz.md")
    if words[0].Val != "echo" || words[1].Val != "Hello James" || words[2].Val != "1>" || words[3].Val != "/tmp/foo/baz.md" {
        t.Fatalf("not parsed correctly: expected [echo, Hello James, 1>, /tmp/foo/baz.md] -- actual [%v, %v, %v, %v]", words[0].Val, words[1].Val, words[2].Val, words[3].Val)
    }
}

func TestParseCodeCraftersRedirectOutputStderr(t *testing.T) {
    words := Parse("echo 'Maria file cannot be found' 2> /tmp/foo/baz.md")
    if words[0].Val != "echo" || words[1].Val != "Maria file cannot be found" || words[2].Val != "2>" || words[3].Val != "/tmp/foo/baz.md" {
        t.Fatalf("not parsed correctly: expected [echo, Maria file cannot be found, 2>, /tmp/foo/baz.md] -- actual [%v, %v, %v, %v]", words[0].Val, words[1].Val, words[2].Val, words[3].Val)
    }
}

