package main

import "github.com/tidwall/gjson"
import "fmt"
import "github.com/davecgh/go-spew/spew"
import "github.com/fatih/structs"
import "encoding/json"

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

    s := structs.New(*p)
    m = s.Map()              // Get a map[string]interface{}
    v := s.Values()           // Get a []interface{}
    f := s.Fields()           // Get a []*Field
    n := s.Names()            // // Get a []string
    nn := s.Name()

    fmt.Println("m = ", m)
    fmt.Printf("type of m = %T\n", m)
    fmt.Println("v = ", v)
    fmt.Println("f = ", f)
    fmt.Println("n = ", n)
    fmt.Println("nn = ", nn)

    bytes, _ := json.Marshal(*p)
    str := string(bytes)
    fmt.Println(str)
    fmt.Printf("type of str = %T\n", str)
    val := gjson.Get(str, "Age")
    println(val.String())

}