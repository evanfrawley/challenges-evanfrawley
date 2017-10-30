package handlers

//TODO: define a handler context struct that
//will be a receiver on any of your HTTP
//handler functions that need access to
//globals, such as the key used for signing
//and verifying SessionIDs, the session store
//and the user store
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
