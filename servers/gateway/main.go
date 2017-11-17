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
    "net/http/httputil"
    "sync"
    "encoding/json"
    "strings"
)

func GetCurrentUser(r *http.Request, ctx *handlers.Context) *users.User {
    //does some magic with our sessions packages
    user, err := ctx.GetAuthenticatedUser(r)
    if err != nil {
        fmt.Printf("error when getitng auth user with err: %v\n", err)
        return &users.User{}
    } else {
        return user
    }
}

func NewServiceProxy(addrs []string, ctx *handlers.Context) *httputil.ReverseProxy {
    nextIndex := 0
    mx := sync.Mutex{}
    return &httputil.ReverseProxy{
        Director: func(r *http.Request) {
            //modify the request to indicate
            //remote host
            user := GetCurrentUser(r, ctx)
            userJSON, err := json.Marshal(user)
            if err != nil {
                log.Printf("error marshaling user: %v", err)
            }
            r.Header.Add("X-User", string(userJSON))

            mx.Lock()
            r.URL.Host = addrs[nextIndex%len(addrs)]
            nextIndex++
            fmt.Printf("making request to %s\n", r.URL.Host)
            mx.Unlock()
            r.URL.Scheme = "http"
        },
    }
}

//main is the main entry point for the server
func main() {
    env := os.Getenv(GO_ENV)
    localAddr := os.Getenv(GOADDR)
    if len(localAddr) == 0 {
        if env == DEVELOPMENT {
            localAddr = DEFAULT_HTTP_PORT
        } else {
            localAddr = DEFAULT_HTTPS_PORT
        }
    }

    signingKey := SIGNING_KEY

    tlsKeyPath := os.Getenv(TLSKEY)
    tlsCertPath := os.Getenv(TLSCERT)

    fmt.Printf("Go port: %s \n", localAddr)
    mux := http.NewServeMux()

    // Microservice Stuff
    // Messaging
    msgSvcAddrs := os.Getenv(MSGSVC_ADDRS)
    splitMsgSvcAddrs := strings.Split(msgSvcAddrs, ",")
    // Summary
    summarySvcAddrs := os.Getenv(SUMMARYSVC_ADDRS)
    splitSummarySvcAddrs := strings.Split(summarySvcAddrs, ",")

    // MONGO SET UP
    mgoAddr := os.Getenv("MONGO_ADDR")
    if len(mgoAddr) == 0 {
        if env == DEVELOPMENT {
            mgoAddr = DEV_MONGO_SVR
        } else {
            mgoAddr = PROD_MONGO_SVR
        }
    }

    mongoSession, err := mgo.Dial(mgoAddr)
    if err != nil {
        fmt.Printf("mgo addr: %v\n", mgoAddr)
        log.Fatalf("error dialing mongo: %v", err)
    }

    mongoStore := users.NewMongoStore(mongoSession, MONGO_DB_NAME, USERS_STORE)

    // REDIS SET UP
    redisAddr := os.Getenv("REDIS_ADDR")
    if len(redisAddr) == 0 {
        if env == DEVELOPMENT {
            redisAddr = DEV_REDIS_SVR
        } else {
            redisAddr = PROD_REDIS_SVR
        }
    }

    fmt.Printf("redis addr: %v\n", redisAddr)
    client := redis.NewClient(&redis.Options{
        Addr:     redisAddr,
        Password: "", // no password set
        DB:       0,  // use default DB
    })

    sessionStore := sessions.NewRedisStore(client, time.Hour)

    // USER TRIE SET UP
    trieRoot := handlers.ConstructUsersTrie(mongoStore)

    // CREATE CONTEXT
    ctx := handlers.NewHandlerContext(mongoStore, sessionStore, signingKey, trieRoot)

    // RESOURCE HANDLERS
    mux.Handle("/v1/summary", NewServiceProxy(splitSummarySvcAddrs, ctx))
    mux.HandleFunc("/v1/users", ctx.UsersHandler)
    mux.HandleFunc("/v1/users/me", ctx.UsersMeHandler)
    mux.HandleFunc("/v1/sessions", ctx.SessionsHandler)
    mux.HandleFunc("/v1/sessions/mine", ctx.SessionsMineHandler)
    mux.Handle("/v1/channels", NewServiceProxy(splitMsgSvcAddrs, ctx))
    mux.Handle("/v1/channels/", NewServiceProxy(splitMsgSvcAddrs, ctx))
    mux.Handle("/v1/messages", NewServiceProxy(splitMsgSvcAddrs, ctx))
    mux.Handle("/v1/messages/", NewServiceProxy(splitMsgSvcAddrs, ctx))

    // MIDDLEWARE
    wrappedMux := handlers.NewCORS(mux)

    // LISTEN AND SERVE
    fmt.Printf("server is listening at http://%s \n", localAddr)
    if env == DEVELOPMENT {
        log.Fatal(http.ListenAndServe(localAddr, wrappedMux))
    } else {
        log.Fatal(http.ListenAndServeTLS(localAddr, tlsCertPath, tlsKeyPath, wrappedMux))
    }
}
