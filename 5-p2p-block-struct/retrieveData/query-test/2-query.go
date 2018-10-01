package main

import "github.com/tidwall/gjson"
import "fmt"
import "github.com/davecgh/go-spew/spew"
import "github.com/fatih/structs"

const jsondata = `{"name":{"first":"Janet","last":"Prichard"},"age":47}`

type person struct {
    Name Name
    Age  int
}

type Name struct {
	Firstname string
	Lastname string
}


func main() {
	value := gjson.Get(jsondata, "name.last")
	println(value.String())

	var p *person

	p = &person{
        Age: 25,
        Name: Name{
            Firstname: "Sumanta",
            Lastname:    "Bose",
        },
    }

    fmt.Println(p)
    spew.Dump(*p)

    fmt.Printf("type of jsondata = %T\n", jsondata)
    fmt.Printf("type of *p = %T\n", *p)

    m := structs.Map(*p)
    fmt.Printf("type of m = %T\n", m)
    fmt.Println("m = ", m)


}