package tools

import (
    "regexp"
)

func Valid(param string,pattern string) bool{
    if ok,_ := regexp.Match(pattern,[]byte(param));!ok {
        return false
    }
    return true
}
