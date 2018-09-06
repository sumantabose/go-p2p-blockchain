package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"bytes"
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

	///////
	fmt.Println("\nPOST I\n")

	sampleProduct := Product {
		Name: "XBOX",
		SerialNo: 25,
		CodeNo: 9123,
		Location: 2,
	}

	jsonValue, err := json.Marshal(sampleProduct)
	if err != nil {
		log.Fatal(err)
	}

	url := "http://localhost:8090/post"

    req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonValue))
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    response, err = client.Do(req)
    if err != nil {
        log.Fatal(err)
    }
    defer response.Body.Close()

    fmt.Println("response Status:", response.Status)
    fmt.Println("response Headers:", response.Header)
    body, err := ioutil.ReadAll(response.Body)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("response Body:", string(body))

    ///////

    fmt.Println("\nPOST II\n")

    sampleProduct2 := Product {
		Name: "YBOX",
		SerialNo: 35,
		CodeNo: 8123,
		Location: 3,
	}

	jsonValue2, err := json.Marshal(sampleProduct2)
	if err != nil {
		log.Fatal(err)
	}

	response, err = http.Post(url, "application/json", bytes.NewBuffer(jsonValue2))
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

    fmt.Println("response Status:", response.Status)
    fmt.Println("response Headers:", response.Header)
    body, err = ioutil.ReadAll(response.Body)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("response Body:", string(body))
}

