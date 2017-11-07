package handlers

import (
    "net/http"
    "fmt"
    "encoding/json"
    "github.com/info344-a17/challenges-evanfrawley/servers/gateway/models/users"
    "github.com/info344-a17/challenges-evanfrawley/servers/gateway/sessions"
    "time"
    "sort"
    "strings"
)

//var QueryParamNotFoundError = errors.New("query parameter was not found in request")

//UsersHandler handles requests for the /v1/user resource
func (ctx *Context) UsersHandler(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case "POST":
        nu := &users.NewUser{}
        if err := json.NewDecoder(r.Body).Decode(nu); err != nil {
            http.Error(w, fmt.Sprintf("error decoding json: %v", err), http.StatusBadRequest)
            return
        }

        if err := nu.Validate(); err != nil {
            http.Error(w, fmt.Sprintf("new user not valid: %v", err), http.StatusBadRequest)
            return
        }

        user, err := ctx.userMongoStore.Insert(nu)

        if err != nil {
            http.Error(w, fmt.Sprintf("error creating new user: %v", err), http.StatusInternalServerError)
            return
        }
        w.WriteHeader(http.StatusCreated)
        respond(w, user)
    case "GET":
        var usersSlice []*users.User
        queryParam := r.URL.Query().Get("q")
        items := ctx.trieRoot.GetUniqueUsersFromPrefix(queryParam)
        for _, item := range items {
            user, err := ctx.userMongoStore.GetByID(item.UserID)
            if err != nil {
                // idk if i should return here
                http.Error(w, fmt.Sprintf("error creating new user: %v", err), http.StatusInternalServerError)
                return
            }
            usersSlice = append(usersSlice, user)
        }
        sort.Slice(usersSlice[:], func(i, j int) bool {
            return strings.Compare(string(usersSlice[i].ID), string(usersSlice[j].ID)) == 1
        })
        if len(usersSlice) > 20 {
            usersSlice = usersSlice[:20]
        }
        w.WriteHeader(http.StatusCreated)
        respond(w, usersSlice)
    default:
        http.Error(w, "method must be POST", http.StatusMethodNotAllowed)
        return
    }
}

func (ctx *Context) UsersMeHandler(w http.ResponseWriter, r *http.Request) {
    user, err := ctx.GetAuthenticatedUser(r)
    if err != nil {
        http.Error(w, "please sign-in", http.StatusUnauthorized)
        return
    }
    switch r.Method {
    case "GET":
        respond(w, *user)
    case "PATCH":
        userUpdates := &users.Updates{}

        if err := json.NewDecoder(r.Body).Decode(userUpdates); err != nil {
            http.Error(w, fmt.Sprintf("error decoding json: %v", err), http.StatusBadRequest)
            return
        }
        if err := ctx.userMongoStore.Update(user.ID, userUpdates); err != nil {
            http.Error(w, fmt.Sprintf("error updating user: %v", err), http.StatusInternalServerError)
            return
        }
        user.ApplyUpdates(userUpdates)

        if err := ctx.UpdateAuthenticatedUserInSessionsStore(r, *user); err != nil {
            http.Error(w, "error saving updated user to sessionStore", http.StatusInternalServerError)
            return
        }

        respond(w, user)
    default:
        http.Error(w, "method must be GET or PATCH", http.StatusMethodNotAllowed)
        return
    }
}

func (ctx *Context) GetAuthenticatedUser(r *http.Request) (*users.User, error) {
    sessionState := SessionState{}
    _, err := sessions.GetSessionID(r, ctx.signingKey)

    if err != nil {
        return nil, fmt.Errorf("session id not valid: %v", err)
    }
    _, err = sessions.GetState(r, ctx.signingKey, ctx.sessionsStore, &sessionState)
    if err != nil {
        return nil, fmt.Errorf("authenticated user not found in session store: %v", err)
    }
    return &sessionState.User, nil
}

func (ctx *Context) UpdateAuthenticatedUserInSessionsStore(r *http.Request, user users.User) error {
    sessionState := SessionState{}
    sid, err := sessions.GetSessionID(r, ctx.signingKey)

    if err != nil {
        return fmt.Errorf("session id not valid: %v", err)
    }

    if err := ctx.sessionsStore.Get(sid, sessionState); err != nil {
        return fmt.Errorf("authenticated user not found in session store: %v", err)
    }
    sessionState.User = user

    if err := ctx.sessionsStore.Save(sid, sessionState); err != nil {
        return fmt.Errorf("authenticated user not found in session store: %v", err)
    }

    return nil
}

// Sessions
func (ctx *Context) SessionsHandler(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case "POST":
        creds := &users.Credentials{}
        if err := json.NewDecoder(r.Body).Decode(creds); err != nil {
            http.Error(w, fmt.Sprintf("error decoding json: %v", err), http.StatusBadRequest)
            return
        }
        user, err := ctx.userMongoStore.GetByEmail(creds.Email)
        if err != nil {
            http.Error(w, "user not found", http.StatusUnauthorized)
            return
        }
        if err := user.Authenticate(creds.Password); err != nil {
            http.Error(w, "password not valid", http.StatusUnauthorized)
            return
        }
        sessionState := &SessionState{
            Created: time.Now(),
            User: *user,
        }
        _, err = sessions.BeginSession(ctx.signingKey, ctx.sessionsStore, sessionState, w)
        if err != nil {
            http.Error(w, "password not valid", http.StatusInternalServerError)
            return
        }
        respond(w, "session created")
    default:
        http.Error(w, "only POST is allowed at this endpoint", http.StatusMethodNotAllowed)
        return
    }
}

func (ctx *Context) SessionsMineHandler(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case "DELETE":
        _, err := sessions.EndSession(r, ctx.signingKey, ctx.sessionsStore)
        if err != nil {
            http.Error(w, "error signing out", http.StatusUnauthorized)
            return
        }
        respond(w, "signed out")
    default:
        http.Error(w, "only DELETE is allowed at this endpoint", http.StatusMethodNotAllowed)
        return
    }
}
