package users

import (
    "encoding/json"
    "time"

    "github.com/patrickmn/go-cache"
    "gopkg.in/mgo.v2/bson"
)

//MemStore represents an in-process memory user store.
//This should be used only for testing and prototyping.
//Production systems should use a shared server store like redis
type MemStore struct {
    entries *cache.Cache
}

//NewMemStore constructs and returns a new MemStore
func NewMemStore(userDuration time.Duration, purgeInterval time.Duration) *MemStore {
    return &MemStore{
        entries: cache.New(userDuration, purgeInterval),
    }
}

//Save saves the provided `user` and associated user ID to the store.
//The `user` parameter is typically a pointer to a struct containing
//all the data you want to associated with the given user ID.
func (ms *MemStore) Save(uid bson.ObjectId, user interface{}) error {
    j, err := json.Marshal(user)
    if nil != err {
        return err
    }
    ms.entries.Set(uid.String(), j, cache.DefaultExpiration)
    return nil
}

//Get populates `user` with the data previously saved
//for the given user id
func (ms *MemStore) Get(uid bson.ObjectId, user interface{}) error {
    j, found := ms.entries.Get(uid.String())
    if !found {
        return ErrUserNotFound
    }
    //reset TTL
    ms.entries.Set(uid.String(), j, 0)
    return json.Unmarshal(j.([]byte), user)
}

//Delete deletes all user data associated with the user id from the store.
func (ms *MemStore) Delete(uid bson.ObjectId) error {
    ms.entries.Delete(uid.String())
    return nil
}
