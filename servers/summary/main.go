package main

import (
    "os"
    "net/http"
    "fmt"
    "log"
)

//main is the main entry point for the server
func main() {
    summaryAddr := os.Getenv("SUMMARY_ADDR")
    if len(summaryAddr) == 0 {
        summaryAddr = ":80"
    }
    fmt.Printf("Go port: %s \n", summaryAddr)
    mux := http.NewServeMux()

    mux.HandleFunc("/v1/summary", LinkSummaryHandler)

    fmt.Printf("server is listening at http://%s \n", summaryAddr)
    log.Fatal(http.ListenAndServe(summaryAddr, mux))
}