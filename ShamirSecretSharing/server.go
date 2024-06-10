package main

import (
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "os"
    "io"
    "flag"
)

func uploadHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        return
    }

    file, _, err := r.FormFile("file")
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    defer file.Close()

    out, err := os.Create("." + r.URL.Path)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer out.Close()

    _, err = io.Copy(out, file)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Success"))
}

func fileHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        return
    }

    fileName := "." + r.URL.Path
    content, err := ioutil.ReadFile(fileName)
    if err != nil {
        http.Error(w, "File not found", http.StatusNotFound)
        return
    }

    w.Header().Set("Content-Type", "text/plain") // Adjust content type if needed
    w.Write(content)
}

func main() {
    port := flag.Int("port", 8000, "port to serve on")
    flag.Parse()

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        if r.Method == http.MethodGet {
            fileHandler(w, r)
        } else {
            uploadHandler(w, r)
        }
    })

    log.Printf("Serving on port %d", *port)
    log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
}
