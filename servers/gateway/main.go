package main

import (
    "os"
    "net/http"
    "fmt"
    "log"
    "github.com/info344-a17/challenges-evanfrawley/servers/gateway/handlers"
    "github.com/info344-a17/challenges-evanfrawley/servers/gateway/sessions"
    //"github.com/go-redis/redis"
    //"time"
    "github.com/info344-a17/challenges-evanfrawley/servers/gateway/models/users"
    "gopkg.in/mgo.v2"
)

//main is the main entry point for the server
func main() {
    localAddr := os.Getenv("GO_ADDR")
    if len(localAddr) == 0 {
        localAddr = ":443"
    }

    //tlsKeyPath := os.Getenv("TLSKEY")
    //tlsCertPath := os.Getenv("TLSCERT")

    //tlsKeyPath := "/Users/evanfrawley/go/src/github.com/info344-a17/challenges-evanfrawley/servers/gateway/tls/privkey.pem"
    //tlsCertPath := "/Users/evanfrawley/go/src/github.com/info344-a17/challenges-evanfrawley/servers/gateway/tls/fullchain.pem"

    fmt.Printf("Go port: %s \n", localAddr)
    mux := http.NewServeMux()
    sessionID, _ := sessions.NewSessionID("nice")
    sessions.ValidateID(string(sessionID), "nice")

    //client := redis.NewClient(&redis.Options{
    //    Addr:     "localhost:6379",
    //    Password: "", // no password set
    //    DB:       0,  // use default DB
    //})
    //store := sessions.NewRedisStore(client, time.Hour)
    //
    //pong, err := client.Ping().Result()

    //sessions.BeginSession("nice", store, )
    mongoSession, err := mgo.Dial("localhost")
    if err != nil {
        log.Fatalf("error dialing mongo: %v", err)
    }

    mongoStore := users.NewMongoStore(mongoSession, "users", "users")
    ctx := handlers.NewHandlerContext(mongoStore)

    mux.HandleFunc("/v1/summary", handlers.SummaryHandler)
    mux.HandleFunc("/v1/users", ctx.UsersHandler)
    mux.HandleFunc("/v1/user", ctx.UsersSpecificHandler)

    fmt.Printf("server is listening at http://%s \n", localAddr)
    log.Fatal(http.ListenAndServe(localAddr, mux))
    //log.Fatal(http.ListenAndServeTLS(localAddr, tlsCertPath, tlsKeyPath, mux))
}
