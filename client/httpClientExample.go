package main

import (
	"fmt"
	"github.com/norby7/shortening-service/client/http"
)

func main(){
	client := http.NewClient("http://localhost:3000")

	createReq := http.CreateRequest{
		Url:  "www.google.com",
	}

	url, err := client.Create(createReq)
	if err != nil{
		panic(err)
	}

	fmt.Println(url)

	counter, err := client.GetCounter(url.Id)
	if err != nil{
		panic(err)
	}

	fmt.Println(counter)

	err = client.Delete(url.Id)
	if err != nil{
		panic(err)
	}

	urlData, err := client.Get(url.Id)
	if err != nil{
		panic(err)
	}

	fmt.Println(urlData)
}