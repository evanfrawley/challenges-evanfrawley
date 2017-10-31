package sessions

import (
	"time"

	"github.com/go-redis/redis"
    "fmt"
    "encoding/json"
    "reflect"
)

//RedisStore represents a session.Store backed by redis.
type RedisStore struct {
	//Redis client used to talk to redis server.
	Client *redis.Client
	//Used for key expiry time on redis.
	SessionDuration time.Duration
}

//NewRedisStore constructs a new RedisStore
func NewRedisStore(client *redis.Client, sessionDuration time.Duration) *RedisStore {
	//initialize and return a new RedisStore struct
	return &RedisStore {
		Client: client,
		SessionDuration: sessionDuration,
	}
}

//Store implementation

//Save saves the provided `sessionState` and associated SessionID to the store.
//The `sessionState` parameter is typically a pointer to a struct containing
//all the data you want to associated with the given SessionID.
func (rs *RedisStore) Save(sid SessionID, sessionState interface{}) error {
	//TODO: marshal the `sessionState` to JSON and save it in the redis database,
	//using `sid.getRedisKey()` for the key.
	//return any errors that occur along the way.
    sessionStateBytes, err := json.Marshal(sessionState)
    if err != nil {
        return fmt.Errorf("error when marshaling JSON: %v", err)
    }

    sessionStateJSONString := string(sessionStateBytes)

	_, err = rs.Client.Set(sid.getRedisKey(), sessionStateJSONString, rs.SessionDuration).Result()
	if err != nil {
        return fmt.Errorf("error saving to redis store: %v", err)
    }
	return nil
}

//Get populates `sessionState` with the data previously saved
//for the given SessionID
func (rs *RedisStore) Get(sid SessionID, sessionState interface{}) error {
	//TODO: get the previously-saved session state data from redis,
	//unmarshal it back into the `sessionState` parameter
	//and reset the expiry time, so that it doesn't get deleted until
	//the SessionDuration has elapsed.

	//for extra-credit using the Pipeline feature of the redis
	//package to do both the get and the reset of the expiry time
	//in just one network round trip!
    pipe := rs.Client.Pipeline()

	getCmd := pipe.Get(sid.getRedisKey())
    pipe.Expire(sid.getRedisKey(), rs.SessionDuration)

    _, err := pipe.Exec()
    if err != nil {
        return ErrStateNotFound
    }

    // unmarshal
    sessionStateUnmarshal := json.Unmarshal([]byte(getCmd.Val()), sessionState)
    sessionState = reflect.ValueOf(sessionStateUnmarshal)

	return nil
}

//Delete deletes all state data associated with the SessionID from the store.
func (rs *RedisStore) Delete(sid SessionID) error {
	//TODO: delete the data stored in redis for the provided SessionID
	_, err := rs.Client.Del(sid.getRedisKey()).Result()
	if err != nil {
	    return fmt.Errorf("error removing key '%s' from redis store: %v", sid, err)
    }
	return nil
}

//getRedisKey() returns the redis key to use for the SessionID
func (sid SessionID) getRedisKey() string {
	//convert the SessionID to a string and add the prefix "sid:" to keep
	//SessionID keys separate from other keys that might end up in this
	//redis instance
	return "sid:" + sid.String()
}
