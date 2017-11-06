package handlers

import "net/http"

/* TODO: implement a CORS middleware handler, as described
in https://drstearns.github.io/tutorials/cors/ that responds
with the following headers to all requests:

  Access-Control-Allow-Origin: *
  Access-Control-Allow-Methods: GET, PUT, POST, PATCH, DELETE
  Access-Control-Allow-Headers: Content-Type, Authorization
  Access-Control-Expose-Headers: Authorization
  Access-Control-Max-Age: 600
*/
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
    w.Header().Add(AccessControlAllowOriginKey, AccessControlAllowOriginVal)
    w.Header().Add(AccessControlAllowMethodsKey, AccessControlAllowMethodsVal)
    w.Header().Add(AccessControlAllowHeadersKey, AccessControlAllowHeadersVal)
    w.Header().Add(AccessControlExposeHeadersKey, AccessControlExposeHeadersVal)
    w.Header().Add(AccessControlMaxAgeKey, AccessControlMaxAgeVal)
    w.Header().Add(ContentTypeKey, ContentTypeJSONUTF8Val)
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
