package handlers

import (
    "net/http"
    "fmt"
    "encoding/json"
    "github.com/info344-a17/challenges-evanfrawley/servers/gateway/models/users"
    "gopkg.in/mgo.v2/bson"
    "path"
    "errors"
    "net/url"
)

var QueryParamNotFoundError = errors.New("query parameter was not found in request")

//UsersHandler handles requests for the /v1/user resource
func (ctx *Context) UsersHandler(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case "GET":
        user, err := ctx.getUserFromQueryParam(r.URL)
        if err != nil {
            http.Error(w, fmt.Sprintf("error getting user: %v", err), http.StatusInternalServerError)
            return
        }
        respond(w, user)
    case "POST":
        nu := &users.NewUser{}
        if err := json.NewDecoder(r.Body).Decode(nu); err != nil {
            http.Error(w, fmt.Sprintf("error decoding json: %v", err), http.StatusBadRequest)
            return
        }
        user, err := ctx.userMongoStore.Insert(nu)

        if err != nil {
            http.Error(w, fmt.Sprintf("error creating new user: %v", err), http.StatusInternalServerError)
            return
        }
        respond(w, user)
    default:
        http.Error(w, "method must be GET or POST", http.StatusMethodNotAllowed)
        return
    }
}

// UsersSpecificHandler handles requests for the /v1/user/{userID}
func (ctx *Context) UsersSpecificHandler(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case "PATCH": {
        // TODO do stuff of getting the current user to update
        uid := bson.ObjectIdHex(path.Base(r.URL.Path))

        uu := &users.Updates{}
        if err := json.NewDecoder(r.Body).Decode(uu); err != nil {
            http.Error(w, fmt.Sprintf("error decoding json: %v", err), http.StatusBadRequest)
            return
        }

        if err := ctx.userMongoStore.Update(uid, uu); err != nil {
            http.Error(w, fmt.Sprintf("error udpating user with id: %v with error: %v", uid, err), http.StatusBadRequest)
            return
        }

    }
    default:
        http.Error(w, "method must be a PATCH", http.StatusMethodNotAllowed)
        return
    }
}

func (ctx *Context) getUserFromQueryParam(url *url.URL) (*users.User, error) {
    if uid := url.Query().Get(userIDQueryParam); uid != "" {
        user, err := ctx.userMongoStore.GetByID(bson.ObjectIdHex(uid))
        if err != nil {
            return nil, fmt.Errorf("error getting user by id: %v", err)
        }
        return user, nil
    } else if email := url.Query().Get(emailQueryParam); email != "" {
        user, err :=  ctx.userMongoStore.GetByEmail(email)
        if err != nil {
            return nil, fmt.Errorf("error getting user by email: %v", err)
        }
        return user, nil
    } else if username := url.Query().Get(usernameQueryParam); username != "" {
        user, err :=  ctx.userMongoStore.GetByUserName(username)
        if err != nil {
            return nil, fmt.Errorf("error getting user by username: %v", err)
        }
        return user, nil
    } else {
        return nil, QueryParamNotFoundError
    }
}