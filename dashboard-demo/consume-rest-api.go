package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)


type Product struct {
    Name string `json:"Name"`
    SerialNo int `json:"SerialNo"`
    CodeNo int `json:"CodeNo"`
    Location int `json:"Location"`
}

func main() {
	
	fmt.Println("\nGET\n")
	response, err := http.Get("http://localhost:8090/next")
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()
	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	var responseObject Product
	json.Unmarshal(responseData, &responseObject)

	fmt.Println(responseObject)
	fmt.Println(responseObject.Name)
	fmt.Println(responseObject.SerialNo)
	fmt.Println(responseObject.CodeNo)
	fmt.Println(responseObject.Location)

	///////
	
	fmt.Println("\nGET LOOP\n")
	response, err = http.Get("http://localhost:8090/next/5")
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()
	responseData, err = ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	var responseObject2 []Product
	json.Unmarshal(responseData, &responseObject2)

	fmt.Println(responseObject2)
}

