package main

import (
    "fmt"
    "strings"
)

func main() {
    data := "cat"
    newData := ""

    result := strings.SplitAfter(data, ".")

    for index, item := range result {
        fmt.Println(index, item)
        //l := len(result) ; 
        if index > 0 {
        	newData = newData + "#."
        }
        newData = newData + item
    }
    
    fmt.Println(newData)
    fmt.Println(len(result))
}