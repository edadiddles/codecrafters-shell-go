package parsing

import (
    "strings"
)

type Word struct {
    Val string
    IsQuoted bool
    IsOperator bool
}

const (
    ENCODING_SPACE = byte(32)
    ENCODING_SINGLE_QUOTE = byte(39)
    ENCODING_DOUBLE_QUOTE = byte(34)
    ENCODING_BACKSLASH = byte(92)
    ENCODING_DOLLAR_SIGN = byte(36)
    ENCODING_GREATER_THAN = byte(62)
)

const (
    QUOTE_TYPE_SINGLE = "single"
    QUOTE_TYPE_DOUBLE = "double"
)


func Parse(input string) []Word {
    words := make([]Word, 0)
    input = strings.TrimSpace(input)
    curr_word := []byte{}
    is_quoted := false
    var quote_type string
    for i := 0; i < len(input); i++ {
        if input[i] == ENCODING_SPACE && !is_quoted {
            if len(curr_word) > 0 {
                words = append(words, Word{Val: string(curr_word), IsQuoted: is_quoted, IsOperator: false})
                curr_word = []byte{}
                is_quoted = false
                quote_type = ""
            }
        } else if input[i] == ENCODING_SINGLE_QUOTE && (!is_quoted || quote_type == QUOTE_TYPE_SINGLE) {
            if is_quoted {
                //words = append(words, Word{Val: string(curr_word), IsQuoted: is_quoted, IsOperator: false})
                //curr_word = []byte{}
                is_quoted = false
                quote_type = ""
            } else {
                is_quoted = true
                quote_type = QUOTE_TYPE_SINGLE
            }
        } else if input[i] == ENCODING_DOUBLE_QUOTE && (!is_quoted || quote_type == QUOTE_TYPE_DOUBLE) {
            if is_quoted {
                //words = append(words, Word{Val: string(curr_word), IsQuoted: is_quoted, IsOperator: false})
                //curr_word = []byte{}
                is_quoted = false
                quote_type = ""
            } else {
                is_quoted = true
                quote_type = QUOTE_TYPE_DOUBLE 
            }
        } else if input[i] == ENCODING_BACKSLASH && (!is_quoted || quote_type == QUOTE_TYPE_DOUBLE) {
            letter := peek(input, i)
            if !is_quoted {
                curr_word = append(curr_word, letter)
                i++
            } else if quote_type == QUOTE_TYPE_DOUBLE {
                if letter == ENCODING_BACKSLASH || letter == ENCODING_DOUBLE_QUOTE || letter == ENCODING_DOLLAR_SIGN {
                    curr_word = append(curr_word, letter)
                    i++
                } else {
                    curr_word = append(curr_word, input[i])
                }
            }
        } else if input[i] == '1' || input[i] == '2' {
            
            letter := peek(input, i)
            if letter == ENCODING_GREATER_THAN {
                // save current word
                if len(curr_word) > 0 {
                    words = append(words, Word{Val: string(curr_word), IsQuoted: is_quoted, IsOperator: true})
                    curr_word = []byte{}
                    is_quoted = false
                    quote_type = "" 
                }
                // start tracking new word
                curr_word = append(curr_word, input[i])
                i++
                curr_word = append(curr_word, input[i])
                letter := peek(input, i)
                if letter == ENCODING_GREATER_THAN {
                    curr_word = append(curr_word, letter)
                    i++
                }
            } else {
                curr_word = append(curr_word, input[i])
            }
        } else if input[i] == ENCODING_GREATER_THAN {
            // save current word
            if len(curr_word) > 0 {
                words = append(words, Word{Val: string(curr_word), IsQuoted: is_quoted, IsOperator: true})
                curr_word = []byte{}
                is_quoted = false
                quote_type = ""
            }

            // start tracking new word
            curr_word = append(curr_word, input[i])
            letter := peek(input, i)
            if letter == ENCODING_GREATER_THAN {
                curr_word = append(curr_word, input[i+1])
                i++
            }
            
            words = append(words, Word{Val: string(curr_word), IsQuoted: is_quoted, IsOperator: true})
            curr_word = []byte{}
            is_quoted = false
            quote_type = ""
        } else {
            curr_word = append(curr_word, input[i])
        }
    }

    // save final word if not empty
    if len(curr_word) > 0 {
        words = append(words, Word{Val: string(curr_word), IsQuoted: false, IsOperator: false})
    }
    
    return words
}

func peek(input string, idx int) byte {
    if len(input) > idx + 1 {
        return input[idx+1]
    }

    return input[idx]
}
