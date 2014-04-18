package main

import (
    "fmt"
    "net/http"
    "log"
    "io"
    "os"
)

// test with curl -X POST localhost:1234/upload -F  file=@<FILE>
func HelloServer(w http.ResponseWriter, req *http.Request) {
    file_r, handler, err := req.FormFile("file")
    if err != nil {
        fmt.Println(err)
    }

    file_w, err := os.Create(handler.Filename)
    if err != nil { panic(err) }
    defer file_w.Close()

    length, err := io.Copy(file_w, file_r)
    fmt.Sprintf("copied %d bytes", length)
    if err != nil {
        panic(err)
    }

}

func main() {
    http.HandleFunc("/upload", HelloServer)
    err := http.ListenAndServe(":1234", nil)
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}
