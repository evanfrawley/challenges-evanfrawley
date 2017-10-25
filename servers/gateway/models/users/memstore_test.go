package users

import (
    "encoding/json"
    "testing"
    "gopkg.in/mgo.v2/bson"
    "time"
    "reflect"
)

func TestMemStore(t *testing.T) {
     testID := bson.NewObjectId()
     user := &User{
         ID: testID,
         FirstName: "fname",
         LastName: "lname",
         Email: "test@test.com",
     }

     userRet := &User{}

    store := NewMemStore(time.Hour, time.Minute)

    if err := store.Get(testID, user); err != ErrUserNotFound {
        t.Errorf("incorrect error when getting user that was never stored: expected %v but got %v", ErrUserNotFound, err)
    }

    if err := store.Save(testID, user); err != nil {
        t.Fatalf("error saving user: %v", err)
    }

    if err := store.Get(testID, userRet); err != nil {
        t.Fatalf("error getting user: %v", err)
    }

    if !reflect.DeepEqual(user, userRet) {
        jexp, _ := json.MarshalIndent(user, "", "  ")
        jact, _ := json.MarshalIndent(userRet, "", "  ")
        t.Errorf("incorrect state retrieved:\nEXPECTED\n%s\nACTUAL\n%s", string(jexp), string(jact))
    }

    if err := store.Delete(testID); err != nil {
        t.Errorf("error deleting user: %v", err)
    }

    if err := store.Get(testID, userRet); err != ErrUserNotFound {
        t.Fatalf("incorrect error when getting user that was deleted: expected %v but got %v", ErrUserNotFound, err)
    }
}

func TestMemStoreSaveUnmarshalble(t *testing.T) {
    //verify that saving an umarshalalbe session state
    //generates an error
    user := func() {} //function values can't be marshaled into JSON

    userID := bson.NewObjectId()
    store := NewMemStore(time.Hour, time.Minute)
    if err := store.Save(userID, user); err == nil {
        t.Error("expected error when attempting to save a session state with an unmarshalable field")
    }
}

// todo update testing