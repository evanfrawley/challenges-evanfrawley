package users

import (
    "testing"
    "gopkg.in/mgo.v2/bson"
    "gopkg.in/mgo.v2"
    "log"
)

func TestMongoStore(t *testing.T) {

    // this code should probably be per ID, email, and username...
    testID := bson.NewObjectId()
    testFirstName := "fname"
    testLastName := "lname"
    testEmail := "test@test.com"
    testUsername := "testUsername"

    testUser := &User{
        FirstName: testFirstName,
        LastName: testLastName,
        Email: testEmail,
        UserName: testUsername,
        PhotoURL: "https://www.gravatar.com/avatar/b642b4217b34b1e8d3bd915fc65c4452",
    }

    testNewUser := &NewUser{
        FirstName: testFirstName,
        LastName: testLastName,
        Email: testEmail,
        UserName: testUsername,
        Password: "password",
        PasswordConf: "password",
    }

    mongoSession, err := mgo.Dial("localhost")
    if err != nil {
        log.Fatalf("error dialing mongo: %v", err)
    }

    store := NewMongoStore(mongoSession, "users", "users")

    if _, err := store.GetByID(testID); err != ErrUserNotFound {
        t.Errorf("incorrect error when getting user that was never stored: expected %v but got %v", ErrUserNotFound, err)
    }

    if _, err := store.GetByEmail(testEmail); err != ErrUserNotFound {
        t.Errorf("incorrect error when getting user that was never stored: expected %v but got %v", ErrUserNotFound, err)
    }

    if _, err := store.GetByUserName(testUsername); err != ErrUserNotFound {
        t.Errorf("incorrect error when getting user that was never stored: expected %v but got %v", ErrUserNotFound, err)
    }
    insertedUser, err := store.Insert(testNewUser)
    if  err != nil {
        t.Fatalf("error saving user: %v", err)
    }
    testID = insertedUser.ID
    testUser.ID = testID

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
