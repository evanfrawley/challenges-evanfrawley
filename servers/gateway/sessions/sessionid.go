package sessions

import (
	"crypto/sha256"
	"errors"
	"encoding/base64"
	"fmt"
	"crypto/rand"
	"crypto/hmac"
    "crypto/subtle"
)

//InvalidSessionID represents an empty, invalid session ID
const InvalidSessionID SessionID = ""

//idLength is the length of the ID portion
const idLength = 32

//signedLength is the full length of the signed session ID
//(ID portion plus signature)
const signedLength = idLength + sha256.Size

//SessionID represents a valid, digitally-signed session ID.
//This is a base64 URL encoded string created from a byte slice
//where the first `idLength` bytes are crytographically random
//bytes representing the unique session ID, and the remaining bytes
//are an HMAC hash of those ID bytes (i.e., a digital signature).
//The byte slice layout is like so:
//+-----------------------------------------------------+
//|...32 crypto random bytes...|HMAC hash of those bytes|
//+-----------------------------------------------------+
type SessionID string

//ErrInvalidID is returned when an invalid session id is passed to ValidateID()
var ErrInvalidID = errors.New("invalid Session ID")
var ErrEmptySigningKey = errors.New("signing key is empty")

//NewSessionID creates and returns a new digitally-signed session ID,
//using `signingKey` as the HMAC signing key. An error is returned only
//if there was an error generating random bytes for the session ID
func NewSessionID(signingKey string) (SessionID, error) {
	if len(signingKey) == 0 {
		return InvalidSessionID, ErrEmptySigningKey
	}

	key := []byte(signingKey)

	signature := make([]byte, signedLength)
	_, err := rand.Read(signature[:idLength])
	if err != nil {
		return InvalidSessionID, fmt.Errorf("error creating random bytes with err: %v", err)
	}
	h := hmac.New(sha256.New, key)
	h.Write(signature[:idLength])
	hmacRandomBytes := h.Sum(nil)
	copy(signature[idLength:], hmacRandomBytes)
	sessionID := SessionID(base64.URLEncoding.EncodeToString(signature))

	return sessionID, nil
}

//ValidateID validates the string in the `id` parameter
//using the `signingKey` as the HMAC signing key
//and returns an error if invalid, or a SessionID if valid
func ValidateID(id string, signingKey string) (SessionID, error) {
    if len(signingKey) == 0 {
        return InvalidSessionID, ErrEmptySigningKey
    }

    key := []byte(signingKey)
    decodedIdBytes, err := base64.URLEncoding.DecodeString(id)
    if err != nil {
        return InvalidSessionID, fmt.Errorf("error decoding sessionID: %v", err)
    }


	if len(decodedIdBytes) != signedLength {
		return InvalidSessionID, fmt.Errorf("length is not same as signed length: %v", err)
	}

	h := hmac.New(sha256.New, key)
    _, err = h.Write(decodedIdBytes[:idLength])
    if err != nil {
    	// TODO
    	return InvalidSessionID, fmt.Errorf("ran into err: %v", err)
    }
    newSigBytes := h.Sum(nil)
    if subtle.ConstantTimeCompare(newSigBytes, decodedIdBytes[idLength:]) == 1 {
        return SessionID(id), nil
    } else {
        return InvalidSessionID, ErrInvalidID
    }
}

//String returns a string representation of the sessionID
func (sid SessionID) String() string {
	return string(sid)
}
