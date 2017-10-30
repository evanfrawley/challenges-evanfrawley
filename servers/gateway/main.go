package main

import (
    "os"
    "net/http"
    "fmt"
    "log"
    "time"

    "github.com/go-redis/redis"
    "github.com/info344-a17/challenges-evanfrawley/servers/gateway/handlers"
    "github.com/info344-a17/challenges-evanfrawley/servers/gateway/sessions"
    "github.com/info344-a17/challenges-evanfrawley/servers/gateway/models/users"
    "gopkg.in/mgo.v2"
    "github.com/info344-a17/challenges-evanfrawley/servers/gateway/middleware"
)

//main is the main entry point for the server
func main() {
    localAddr := os.Getenv("GO_ADDR")
    if len(localAddr) == 0 {
        localAddr = "localhost:4000"
    }
    signingKey := "thisissimplysigning"

    //tlsKeyPath := os.Getenv("TLSKEY")
    //tlsCertPath := os.Getenv("TLSCERT")

    //tlsKeyPath := "/Users/evanfrawley/go/src/github.com/info344-a17/challenges-evanfrawley/servers/gateway/tls/privkey.pem"
    //tlsCertPath := "/Users/evanfrawley/go/src/github.com/info344-a17/challenges-evanfrawley/servers/gateway/tls/fullchain.pem"

    fmt.Printf("Go port: %s \n", localAddr)
    mux := http.NewServeMux()

    client := redis.NewClient(&redis.Options{
        Addr:     "localhost:6379",
        Password: "", // no password set
        DB:       0,  // use default DB
    })
    sessionStore := sessions.NewRedisStore(client, time.Hour)

    mongoSession, err := mgo.Dial("localhost")
    if err != nil {
        log.Fatalf("error dialing mongo: %v", err)
    }

    mongoStore := users.NewMongoStore(mongoSession, "users", "users")

    ctx := handlers.NewHandlerContext(mongoStore, sessionStore, signingKey)

    mux.HandleFunc("/v1/summary", handlers.SummaryHandler)
    mux.HandleFunc("/v1/users", ctx.UsersHandler)
    mux.HandleFunc("/v1/users/me", ctx.UsersMeHandler)
    mux.HandleFunc("/v1/sessions", ctx.SessionsHandler)
    mux.HandleFunc("/v1/sessions/mine", ctx.SessionsMineHandler)

    wrappedMux := middleware.NewCORS(mux)

    fmt.Printf("server is listening at http://%s \n", localAddr)
    log.Fatal(http.ListenAndServe(localAddr, wrappedMux))
    //log.Fatal(http.ListenAndServeTLS(localAddr, tlsCertPath, tlsKeyPath, mux))
}
