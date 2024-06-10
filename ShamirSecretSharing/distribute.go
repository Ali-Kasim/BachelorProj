package main

import (
    "bytes"
    "fmt"
    "io"
    "io/ioutil"
    "log"
    "mime/multipart"
    "net/http"
    "os"
    "os/exec"
    "path/filepath"
)

func uploadFile(url, filePath string) error {
    file, err := os.Open(filePath)
    if err != nil {
        return err
    }
    defer file.Close()

    body := &bytes.Buffer{}
    writer := multipart.NewWriter(body)
    part, err := writer.CreateFormFile("file", filepath.Base(file.Name()))
    if err != nil {
        return err
    }

    io.Copy(part, file)

    writer.Close()

    req, _ := http.NewRequest("POST", url+"/"+filepath.Base(filePath), body)
    req.Header.Set("Content-Type", writer.FormDataContentType())

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    responseBody, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return err
    }

    fmt.Printf("Upload to %s response: %s\n", url, responseBody)
    return nil
}

func setupServers(baseDir string, numServers int, startPort int) ([]string, error) {
    serverGoPath := "C:/Users/Ali-h/Downloads/BachProj/TaskE/server.go"
    serverUrls := make([]string, numServers)

    for i := 1; i <= numServers; i++ {
        serverDir := filepath.Join(baseDir, fmt.Sprintf("server%d", i))
        if _, err := os.Stat(serverDir); os.IsNotExist(err) {
            err := os.MkdirAll(serverDir, 0755)
            if err != nil {
                return nil, fmt.Errorf("error creating server directory: %v", err)
            }
        }

        destServerGoPath := filepath.Join(serverDir, "server.go")
        err := copyFile(serverGoPath, destServerGoPath)
        if err != nil {
            return nil, fmt.Errorf("error copying server.go file: %v", err)
        }

        // Build the server.go file
        serverExecutablePath := filepath.Join(serverDir, "server.exe")
        cmd := exec.Command("go", "build", "-o", serverExecutablePath, destServerGoPath)
        cmd.Dir = serverDir
        buildOutput, err := cmd.CombinedOutput()
        if err != nil {
            return nil, fmt.Errorf("error building server.go: %v, output: %s", err, string(buildOutput))
        }

        // Check if the server executable exists
        if _, err := os.Stat(serverExecutablePath); os.IsNotExist(err) {
            return nil, fmt.Errorf("server executable not found: %v", serverExecutablePath)
        }

        // Start the server
        port := startPort + i - 1
        serverUrl := fmt.Sprintf("http://localhost:%d", port)
        serverUrls[i-1] = serverUrl

        cmd = exec.Command(serverExecutablePath, fmt.Sprintf("-port=%d", port))
        cmd.Dir = serverDir
        cmd.Stdout = os.Stdout
        cmd.Stderr = os.Stderr
        err = cmd.Start()
        if err != nil {
            return nil, fmt.Errorf("error starting server: %v", err)
        }
    }

    return serverUrls, nil
}


// Copy a file from src to dst
func copyFile(src, dst string) error {
    sourceFile, err := os.Open(src)
    if err != nil {
        return err
    }
    defer sourceFile.Close()

    destinationFile, err := os.Create(dst)
    if err != nil {
        return err
    }
    defer destinationFile.Close()

    _, err = io.Copy(destinationFile, sourceFile)
    return err
}

func distributeKeys(sharesDir string, keyServers []string) {
    for i, server := range keyServers {
        filePath := fmt.Sprintf("%s/share%d.txt", sharesDir, i+1)
        err := uploadFile(server, filePath)
        if err != nil {
            log.Fatalf("Failed to upload %s: %v", filePath, err)
        }
    }
}

func distributeDataFiles(encryptedFile, ivFile string, dataServers []string) {
    for _, server := range dataServers {
        err := uploadFile(server, encryptedFile)
        if err != nil {
            log.Fatalf("Failed to upload %s: %v", encryptedFile, err)
        }

        err = uploadFile(server, ivFile)
        if err != nil {
            log.Fatalf("Failed to upload %s: %v", ivFile, err)
        }
    }
}

func main() {
    sharesDir := "C:/Users/Ali-h/Downloads/BachProj/TaskE/shares" 

    // Get the number of shares
    shareFiles, err := ioutil.ReadDir(sharesDir)
    if err != nil {
        log.Fatalf("Failed to read shares directory: %v", err)
    }
    numShares := len(shareFiles)

    // Setup key servers
    keyServers, err := setupServers("C:/Users/Ali-h/Downloads/BachProj/TaskE/key_servers", numShares, 8001)
    if err != nil {
        log.Fatalf("Failed to setup key servers: %v", err)
    }

    // Setup data servers
    dataServers, err := setupServers("C:/Users/Ali-h/Downloads/BachProj/TaskE/data_servers", 3, 9001)
    if err != nil {
        log.Fatalf("Failed to setup data servers: %v", err)
    }

    // Distribute the shares and encrypted file
    distributeKeys(sharesDir, keyServers)
    distributeDataFiles("C:/Users/Ali-h/Downloads/BachProj/TaskE/encrypted.dat", "C:/Users/Ali-h/Downloads/BachProj/TaskE/iv.dat", dataServers)

    fmt.Println("Distribution complete.")
}
