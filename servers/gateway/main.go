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
)

//main is the main entry point for the server
func main() {
    localAddr := os.Getenv("GO_ADDR")
    if len(localAddr) == 0 {
        localAddr = ":433"
    }
    env := os.Getenv("GO_ENV")

    signingKey := "thisissimplysigning"

    tlsKeyPath := os.Getenv("TLSKEY")
    tlsCertPath := os.Getenv("TLSCERT")

    fmt.Printf("Go port: %s \n", localAddr)
    mux := http.NewServeMux()


    var mongoSession *mgo.Session
    var redisAddr string
    var err error
    if env == "development" {
        mongoSession, err = mgo.Dial("localhost")
        redisAddr = "localhost:6379"
    } else {
        mongoSession, err = mgo.Dial("mongosvr")
        redisAddr = "redissvr:6379"
    }
    if err != nil {
        log.Fatalf("error dialing mongo: %v", err)
    }

    client := redis.NewClient(&redis.Options{
        Addr:     redisAddr,
        Password: "", // no password set
        DB:       0,  // use default DB
    })
    sessionStore := sessions.NewRedisStore(client, time.Hour)

    mongoStore := users.NewMongoStore(mongoSession, "users", "users")

    trieRoot := handlers.ConstructUsersTrie(mongoStore)

    ctx := handlers.NewHandlerContext(mongoStore, sessionStore, signingKey, trieRoot)

    mux.HandleFunc("/v1/summary", handlers.SummaryHandler)
    mux.HandleFunc("/v1/users", ctx.UsersHandler)
    mux.HandleFunc("/v1/users/me", ctx.UsersMeHandler)
    mux.HandleFunc("/v1/sessions", ctx.SessionsHandler)
    mux.HandleFunc("/v1/sessions/mine", ctx.SessionsMineHandler)

    wrappedMux := handlers.NewCORS(mux)

    fmt.Printf("server is listening at http://%s \n", localAddr)
    if env == "development" {
        log.Fatal(http.ListenAndServe(localAddr, wrappedMux))
    } else {
        log.Fatal(http.ListenAndServeTLS(localAddr, tlsCertPath, tlsKeyPath, wrappedMux))
    }
}
