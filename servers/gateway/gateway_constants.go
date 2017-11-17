package main

// Golang specifics
const (
    GOADDR = "GOADDR"
    GO_ENV = "GO_ENV"
    SIGNING_KEY = "thisissimplysigning"
    DEVELOPMENT = "development"
    DEFAULT_HTTPS_PORT = ":433"
    DEFAULT_HTTP_PORT = ":80"
)

// Certs
const (
    TLSKEY = "TLSKEY"
    TLSCERT = "TLSCERT"
)

// Microservices
const (
    MSGSVC_ADDRS = "MSGSVC_ADDRS"
    SUMMARYSVC_ADDRS = "SUMMARYSVC_ADDRS"
)

// Mongo
const (
    PROD_MONGO_SVR = "mongosvr:27017"
    DEV_MONGO_SVR = "localhost:27017"
    MONGO_DB_NAME = "info344"
    USERS_STORE = "users"
)

// Redis
const (
    DEV_REDIS_SVR = "localhost:6379"
    PROD_REDIS_SVR = "redissvr:6379"
)