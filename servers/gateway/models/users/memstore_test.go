package users

import (
    "encoding/json"
    "testing"
    "gopkg.in/mgo.v2/bson"
    "reflect"
    "fmt"
)

func TestMemStore(t *testing.T) {
    // this code should probably be per ID, email, and username...
    testID := bson.NewObjectId()
    testFirstName := "fname"
    testLastName := "lname"
    testEmail := "test@test.com"
    testUsername := "testUsername"

    testUser := &User{
        ID: testID,
        FirstName: testFirstName,
        LastName: testLastName,
        Email: testEmail,
        UserName: testUsername,
    }

    store := NewMemStore()

    if u, err := store.GetByID(testID); u == nil && err != ErrUserNotFound {
        t.Errorf("incorrect error when getting user that was never stored: expected %v but got %v", ErrUserNotFound, err)
    }

    if u, err := store.GetByEmail(testEmail); u == nil && err != ErrUserNotFound {
        t.Errorf("incorrect error when getting user that was never stored: expected %v but got %v", ErrUserNotFound, err)
    }

    if u, err := store.GetByUserName(testUsername); u == nil && err != ErrUserNotFound {
        t.Errorf("incorrect error when getting user that was never stored: expected %v but got %v", ErrUserNotFound, err)
    }

    if _, err := store.Insert(testUser); err != nil {
        t.Fatalf("error saving user: %v", err)
    }

    user, err := store.GetByID(testID)

    if err != nil {
        t.Fatalf("error getting user: %v", err)
    }

    if err = checkEqualToTestValue(user, testUser); err != nil {
        t.Error(err)
    }

    user, err = store.GetByEmail(testEmail)

    if err != nil {
        t.Fatalf("error getting user: %v", err)
    }

    if err = checkEqualToTestValue(user, testUser); err != nil {
        t.Error(err)
    }

    user, err = store.GetByUserName(testUsername)

    if err != nil {
        t.Fatalf("error getting user: %v", err)
    }

    if err = checkEqualToTestValue(user, testUser); err != nil {
        t.Error(err)
    }

    if err := store.Delete(testID); err != nil {
        t.Errorf("error deleting user: %v", err)
    }

    if u, err := store.GetByID(testID); u == nil && err != ErrUserNotFound {
        t.Fatalf("incorrect error when getting user that was deleted: expected %v but got %v", ErrUserNotFound, err)
    }
}

func checkEqualToTestValue(user, userExpected *User) error {
    if !reflect.DeepEqual(user, userExpected) {
        jexp, _ := json.MarshalIndent(user, "", "  ")
        jact, _ := json.MarshalIndent(userExpected, "", "  ")
        return fmt.Errorf("incorrect state retrieved:\nEXPECTED\n%s\nACTUAL\n%s", string(jexp), string(jact))
    }
    return nil
}