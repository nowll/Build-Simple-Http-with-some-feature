package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	// "io/ioutil"
	"main/handler"
	"net/http"
)

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello World")
}

func postJson(w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)
	handler.ErrorHandler(err)

	var handleData map[int]string

	err = json.Unmarshal(data, &handleData)
	handler.ErrorHandler(err)

	fmt.Println("Data : ", handleData)

	fmt.Fprintln(w, "Accept")
}

func fileHandle(w http.ResponseWriter, request *http.Request) {
	fmt.Println("test")
	err := request.ParseMultipartForm(32 << 20) // maxMemory 32MB
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	file, h, err := request.FormFile("test")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	tmpfile, err := os.Create("./fileloc" + h.Filename)
	defer tmpfile.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = io.Copy(tmpfile, file)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(200)
	return
}

func main() {
	mux := http.NewServeMux()

	fmt.Println("Server is running")
	mux.HandleFunc("/", index)
	mux.HandleFunc("/postJson", postJson)
	mux.HandleFunc("/fileHandle", fileHandle)

	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	server.ListenAndServe()

}
