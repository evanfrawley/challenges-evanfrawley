package sessions

import (
    "errors"
    "net/http"
    "fmt"
    "strings"
)

const headerAuthorization = "Authorization"
const paramAuthorization = "auth"
const schemeBearer = "Bearer "

//ErrNoSessionID is used when no session ID was found in the Authorization header
var ErrNoSessionID = errors.New("no session ID found in " + headerAuthorization + " header")

//ErrInvalidScheme is used when the authorization scheme is not supported
var ErrInvalidScheme = errors.New("authorization scheme not supported")

var ErrNoSigningKey = errors.New("no signing key was provided")

//BeginSession creates a new SessionID, saves the `sessionState` to the store, adds an
//Authorization header to the response with the SessionID, and returns the new SessionID
func BeginSession(signingKey string, store Store, sessionState interface{}, w http.ResponseWriter) (SessionID, error) {
    if len(signingKey) == 0 {
        return InvalidSessionID, ErrNoSigningKey
    }
    sid, err := NewSessionID(signingKey)
    sid = SessionID(sid)
    if err != nil {
        return InvalidSessionID, nil
    }
    store.Save(sid, sessionState)

    w.Header().Add(headerAuthorization, fmt.Sprintf("%s%s", schemeBearer, sid))
    w.WriteHeader(http.StatusCreated)

    return sid, nil
}

//GetSessionID extracts and validates the SessionID from the request headers
func GetSessionID(r *http.Request, signingKey string) (SessionID, error) {
    sidWithBearer := r.Header.Get(headerAuthorization)
    sid := ""
    if len(sidWithBearer) != 0 {
        if strings.HasPrefix(sidWithBearer, schemeBearer) {
            sid = strings.Replace(sidWithBearer, schemeBearer, "", 1)
        } else {
            return InvalidSessionID, ErrInvalidScheme
        }
    } else {
        authQueryParam := r.URL.Query().Get("auth")
        if len(authQueryParam) == 0 {
            return InvalidSessionID, ErrInvalidScheme
        }
        if !strings.HasPrefix(authQueryParam, schemeBearer) {
            return InvalidSessionID, ErrInvalidScheme
        }

        sid = strings.Replace(authQueryParam, schemeBearer, "", 1)
    }
    _, err := ValidateID(sid, signingKey)
    if err != nil {
        return InvalidSessionID, ErrNoSessionID
    }

    return SessionID(sid), nil
}

//GetState extracts the SessionID from the request,
//gets the associated state from the provided store into
//the `sessionState` parameter, and returns the SessionID
func GetState(r *http.Request, signingKey string, store Store, sessionState interface{}) (SessionID, error) {
    sid, err := GetSessionID(r, signingKey)
    if err != nil {
        return InvalidSessionID, ErrNoSessionID
    }

    err = store.Get(sid, sessionState)
    if err != nil {
        return InvalidSessionID, ErrStateNotFound
    }

    return sid, nil
}

//EndSession extracts the SessionID from the request,
//and deletes the associated data in the provided store, returning
//the extracted SessionID.
func EndSession(r *http.Request, signingKey string, store Store) (SessionID, error) {
    sid, err := GetSessionID(r, signingKey)
    if err != nil {
        return InvalidSessionID, ErrInvalidScheme
    }

    if err := store.Delete(SessionID(sid)); err != nil {
        return InvalidSessionID, err
    }

    return SessionID(sid), nil
}
