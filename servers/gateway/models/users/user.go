package users

import (
    "gopkg.in/mgo.v2/bson"
    "net/mail"
    "fmt"
    "strings"
    "golang.org/x/crypto/bcrypt"
    "github.com/info344-a17/challenges-evanfrawley/servers/gateway/helpers"
)

const gravatarBasePhotoURL = "https://www.gravatar.com/avatar/"

var bcryptCost = 13

//User represents a user account in the database
type User struct {
    ID        bson.ObjectId `json:"id" bson:"_id"`
    Email     string        `json:"email" bson:"email"`
    PassHash  []byte        `json:"-"` //stored, but not encoded to clients
    UserName  string        `json:"userName" bson:"userName"`
    FirstName string        `json:"firstName" bson:"firstName"`
    LastName  string        `json:"lastName" bson:"lastName"`
    PhotoURL  string        `json:"photoURL" bson:"photoURL"`
}

//Credentials represents user sign-in credentials
type Credentials struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

//NewUser represents a new user signing up for an account
type NewUser struct {
    Email        string `json:"email"`
    Password     string `json:"password"`
    PasswordConf string `json:"passwordConf"`
    UserName     string `json:"userName"`
    FirstName    string `json:"firstName"`
    LastName     string `json:"lastName"`
}

//Updates represents allowed updates to a user profile
type Updates struct {
    FirstName string `json:"firstName" bson:"firstName"`
    LastName  string `json:"lastName" bson:"lastName"`
}

//Validate validates the new user and returns an error if
//any of the validation rules fail, or nil if its valid
func (nu *NewUser) Validate() error {
    // email
    _, err := mail.ParseAddress(nu.Email)
    if err != nil {
        return fmt.Errorf("user provided an invalid email address")
    }

    // pw len
    if len(nu.Password) < 6 {
        return fmt.Errorf("user provided too short of a password")
    }

    // pw match
    if nu.Password != nu.PasswordConf {
        return fmt.Errorf("password and password confirmation do not match")
    }

    // username != 0
    if len(strings.TrimSpace(nu.UserName)) == 0 {
        return fmt.Errorf("username must not be empty")
    }

    return nil
}

//ToUser converts the NewUser to a User, setting the
//PhotoURL and PassHash fields appropriately
func (nu *NewUser) ToUser() (*User, error) {
    emailHash := helpers.GetMD5Hash(nu.Email)
    gravatarURL := fmt.Sprintf("%s%s", gravatarBasePhotoURL, emailHash)

    bsonID := bson.NewObjectId()

    user := &User{
        ID:        bsonID,
        Email:     nu.Email,
        UserName:  nu.UserName,
        FirstName: nu.FirstName,
        LastName:  nu.LastName,
        PhotoURL:  gravatarURL,
    }

    if err := user.SetPassword(nu.Password); err != nil {
        return nil, fmt.Errorf("error while setting password: %v", err)
    }

    return user, nil
}

//FullName returns the user's full name, in the form:
// "<FirstName> <LastName>"
//If either first or last name is an empty string, no
//space is put betweeen the names
func (u *User) FullName() string {
    nameString := fmt.Sprintf("%s %s", u.FirstName, u.LastName)
    return strings.TrimSpace(nameString)
}

//SetPassword hashes the password and stores it in the PassHash field
func (u *User) SetPassword(password string) error {
    if len(password) == 0 {
        return fmt.Errorf("attempting to set a zero lendth password")
    }
    pwBytes := []byte(password)
    pwHashBytes, err := bcrypt.GenerateFromPassword(pwBytes, bcryptCost)
    if err != nil {
        return fmt.Errorf("error generating a hash from password: %v", err)
    }

    u.PassHash = pwHashBytes
    return nil
}

//Authenticate compares the plaintext password against the stored hash
//and returns an error if they don't match, or nil if they do
func (u *User) Authenticate(password string) error {
    err := bcrypt.CompareHashAndPassword(u.PassHash, []byte(password))
    if err != nil {
        return fmt.Errorf("error comparing password to hash: %v", err)
    }
    return nil
}

//ApplyUpdates applies the updates to the user. An error
//is returned if the updates are invalid
func (u *User) ApplyUpdates(updates *Updates) error {
    updateFailString := "tried updating %s with an empty string"
    if len(updates.FirstName) == 0 {
        return fmt.Errorf(updateFailString, "first name")
    }

    if len(updates.LastName) == 0 {
        return fmt.Errorf(updateFailString, "last name")
    }

    u.FirstName = updates.FirstName
    u.LastName = updates.LastName
    return nil
}
