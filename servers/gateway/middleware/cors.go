package middleware

import (
    "net/http"
    "github.com/info344-a17/challenges-evanfrawley/servers/gateway/handlers"
)

//CORS is a middleware handler that adds CORS support
type CORS struct {
    handler http.Handler
}

//ServeHTTP handles the request by adding the CORS headers
//and calling the real handler if the method is something other
//then OPTIONS (used for pre-flight requests)
func (c *CORS) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    //set the various CORS response headers depending on
    //what you want your server to allow
    w.Header().Add(handlers.AccessControlAllowOriginKey, handlers.AccessControlAllowOriginVal)
    //...more CORS response headers...

    //if this is preflight request, the method will
    //be OPTIONS, so call the real handler only if
    //the method is something else
    if r.Method != "OPTIONS" {
        c.handler.ServeHTTP(w, r)
    }
}

//NewCORS constructs a new CORS middleware handler
func NewCORS(handlerToWrap http.Handler) *CORS {
    return &CORS{handlerToWrap}
}
