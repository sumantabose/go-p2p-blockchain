package main

import (
    "fmt"
    "log"
    "io/ioutil"
    "strconv"


    "os"
    "github.com/kr/fs"
)

func main() {
    fmt.Println("hello world")
    searchFiles("data")

    //ExampleWalker("data")
}


func searchFiles(dir string) { // dir is the parent directory you what to search
    files, err := ioutil.ReadDir(dir)
    if err != nil {
        log.Fatal(err)
    }

    //var fileName string
    mostRecectFileNo := 0

    for _, file := range files {
        fmt.Println(file.Name()) 
        //fileName = file.Name()
        fileNo, _ := strconv.Atoi(file.Name()[len("product-data")+1:len(file.Name())-4])
        if fileNo > mostRecectFileNo { mostRecectFileNo = fileNo }
    }

    if mostRecectFileNo == 0 {
        fmt.Println("No recect file")
    } else {
        fmt.Println("Most Recect File No = ", mostRecectFileNo)
    }    
}


func ExampleWalker(dir string) {
    walker := fs.Walk(dir)
    for walker.Step() {
        if err := walker.Err(); err != nil {
            fmt.Fprintln(os.Stderr, err)
            continue
        }
        fmt.Println(walker.Path())
    }
}