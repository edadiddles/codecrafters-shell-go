package main

import (
    "testing"
    "fmt"
)

func TestChkCmdAutocomplete(t *testing.T) {
    autocompletes := chk_cmd_autocomplete("echo")
    fmt.Println("--------")
    fmt.Println(autocompletes)
}
