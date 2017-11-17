package handlers

const (
    headerContentType="Content-Type"
    contentTypeJSON="application/json"
)

const (
    userIDQueryParam="id"
    usernameQueryParam="username"
    emailQueryParam="email"
)

const (
    AccessControlAllowOriginKey   = "Access-Control-Allow-Origin"
    AccessControlAllowOriginVal   = "*"
    ContentTypeKey                = "Content-Type"
    ContentTypeJSONUTF8Val        = "application/json; charset=utf-8"
    AccessControlAllowMethodsKey  = "Access-Control-Allow-Methods"
    AccessControlAllowMethodsVal  = "GET, PUT, POST, PATCH, DELETE"
    AccessControlAllowHeadersKey  = "Access-Control-Allow-Headers"
    AccessControlAllowHeadersVal  = "Content-Type, Authorization"
    AccessControlExposeHeadersKey = "Access-Control-Expose-Headers"
    AccessControlExposeHeadersVal = "Authorization"
    AccessControlMaxAgeKey        = "Access-Control-Max-Age"
    AccessControlMaxAgeVal        = "600"
)