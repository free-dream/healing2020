package tools

import (
    "regexp"
    "reflect"
    //"fmt"
)

func Valid(param string,pattern string) bool{
    if ok,_ := regexp.Match(pattern,[]byte(param));!ok {
        //fmt.Println(err)
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
