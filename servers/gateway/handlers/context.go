package handlers

import (
    "github.com/info344-a17/challenges-evanfrawley/servers/gateway/models/users"
)

//Context holds context values
//used by multiple handler functions.
type Context struct {
    userMongoStore users.Store
}

func NewHandlerContext(userStore users.Store) *Context {
    return &Context {
        userMongoStore: userStore,
    }
}