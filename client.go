package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"main/handler"
	"mime/multipart"
	"net/http"
	"os"
	"time"
)

func main() {

	choice := 0

	for {
		fmt.Println("Welcome to the web")
		fmt.Println("1. Testing JSON")
		fmt.Println("2. Testing Multiple File Uploads")
		fmt.Println("3. Exit")
		fmt.Print(">> ")

		fmt.Scanf("%d\n", &choice)

		switch choice {
		case 1:
			uploadJson()
		case 2:
			multipartFile()
		case 3:
			return
		}
	}
}

func uploadJson() {

	transport := &http.Transport{
		DisableKeepAlives: true,
	}

	var client = &http.Client{
		Transport: transport,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel()
    req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost:8080/postJson", nil)
    handler.ErrorHandler(err)

    if _, err := client.Do(req); err != nil {
        log.Fatal(err)
    }

	
	data := map[int]string{
		1: "Asep",
		2: "Kuriniawan",
	}


	jsonData, err := json.Marshal(data)
	handler.ErrorHandler(err)

	jsonBuffer := bytes.NewBuffer(jsonData)

	request, err := http.NewRequest("POST", "http://localhost:8080/postJson", jsonBuffer)
	handler.ErrorHandler(err)
	request.Header.Set("Content-Type", "application/json")

	response, err := client.Do(request)
	handler.ErrorHandler(err)

	defer response.Body.Close()


	body, err := io.ReadAll(response.Body)
	handler.ErrorHandler(err)

	fmt.Println("response: ", string(body))

}

func multipartFile() {
	transport := &http.Transport{
		DisableKeepAlives: true,
	}

	client := &http.Client{
		Transport: transport,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel()
    req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost:8080/fileHandle", nil)
    handler.ErrorHandler(err)

    if _, err := client.Do(req); err != nil {
        log.Fatal(err)
    }
	// New multipart writer.
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Create a form field.
	fw, err := writer.CreateFormFile("test", "test.txt")
	handler.ErrorHandler(err)

	file, err := os.Open("test.txt")
	handler.ErrorHandler(err)


	_, err = io.Copy(fw, file)
	handler.ErrorHandler(err)


	writer.Close()
	req, err = http.NewRequest("POST", "http://localhost:8080/fileHandle", bytes.NewReader(body.Bytes()))
	handler.ErrorHandler(err)

	req.Header.Set("Content-Type", writer.FormDataContentType())
	rsp, _ := client.Do(req)
	if rsp.StatusCode != http.StatusOK {
		log.Printf("Request failed with response code: %d", rsp.StatusCode)
	}
}
