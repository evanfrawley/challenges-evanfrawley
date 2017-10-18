package main

import (
    "os"
    "net/http"
    "fmt"
    "log"
    "github.com/info344-a17/challenges-evanfrawley/servers/gateway/handlers"
)

//main is the main entry point for the server
func main() {
    localAddr := os.Getenv("GO_ADDR")
    if len(localAddr) == 0 {
        localAddr = ":443"
    }

    tlsKeyPath := os.Getenv("TLSKEY")
    tlsCertPath := os.Getenv("TLSCERT")

    fmt.Printf("Go port: %s \n", localAddr)
    mux := http.NewServeMux()

    mux.HandleFunc("/v1/summary", handlers.SummaryHandler)

    fmt.Printf("server is listening at http://%s \n", localAddr)
    log.Fatal(http.ListenAndServeTLS(localAddr, tlsCertPath, tlsKeyPath, mux))
}
