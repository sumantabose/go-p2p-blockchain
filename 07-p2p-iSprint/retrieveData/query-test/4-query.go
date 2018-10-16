package main

import (
    "github.com/tidwall/gjson"
    "fmt"
    _ "github.com/davecgh/go-spew/spew"
    _ "github.com/fatih/structs"
    "encoding/json"
)




type Person struct {
    Name Name
    Age  int
}

type Name struct {
	Firstname string
	Lastname string
}

var Persons []Person

func main() {

	p1 := &Person{
        Age: 28,
        Name: Name{
            Firstname: "Sumanta",
            Lastname:    "Bose",
        },
    }

    p2 := &Person{
        Age: 60,
        Name: Name{
            Firstname: "PM",
            Lastname:    "Lee",
        },
    }

    Persons = append(Persons, *p1, *p2)

    bytes, _ := json.Marshal(*p1)
    str := string(bytes)
    fmt.Println(str)
    fmt.Printf("type of str = %T\n", str)
    val := gjson.Get(str, "Name.Firstname")
    fmt.Printf("%T, %s\n", val, val)

    Persons = append(Persons, *p1, *p2)
    bytes, _ = json.Marshal(Persons)
    str = string(bytes)
    gjson.ForEachLine(str, func(line gjson.Result) bool{
    println(line.String())
    return true
})

}