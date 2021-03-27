package tools

import (
    "regexp"
    "reflect"
)

func Valid(param string,pattern string) bool{
    if ok,_ := regexp.Match(pattern,[]byte(param));!ok {
        return false
    }
    return true
}

func IsZeroValue(i interface{}) bool{
    defer func() {
        recover()
    }()
    vi := reflect.ValueOf(i)
    return !vi.IsValid()
}
