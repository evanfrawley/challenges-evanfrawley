package users

import (
    "testing"
    "golang.org/x/crypto/bcrypt"
)

//TODO: add tests for the various functions in user.go, as described in the assignment.
//use `go test -cover` to ensure that you are covering all or nearly all of your code paths.

func TestValidateNewUser(t *testing.T) {
    cases := []struct {
        name        string
        hint        string
        newUser     *NewUser
        expectError bool
    }{
        {
            name: "Valid NewUser",
            hint: "",
            newUser: &NewUser{
                Email:        "test@test.edu",
                Password:     "password",
                PasswordConf: "password",
                UserName:     "test",
                FirstName:    "test_fname",
                LastName:     "test_lname",
            },
            expectError: false,
        },
        {
            name: "Invalid NewUser - Email",
            hint: "",
            newUser: &NewUser{
                Email:        "bademailaddr",
                Password:     "password",
                PasswordConf: "password",
                UserName:     "test",
                FirstName:    "test_fname",
                LastName:     "test_lname",
            },
            expectError: true,
        },
        {
            name: "Invalid NewUser - Password short",
            hint: "",
            newUser: &NewUser{
                Email:        "test@test.com",
                Password:     "pass",
                PasswordConf: "pass",
                UserName:     "test",
                FirstName:    "test_fname",
                LastName:     "test_lname",
            },
            expectError: true,
        },
        {
            name: "Invalid NewUser - Password confirmation doesn't match",
            hint: "Check to see if the passwords are matching",
            newUser: &NewUser{
                Email:        "test@test.com",
                Password:     "password",
                PasswordConf: "password1",
                UserName:     "test",
                FirstName:    "test_fname",
                LastName:     "test_lname",
            },
            expectError: true,
        },
        {
            name: "Invalid NewUser - Username empty",
            hint: "Test to see if there is anything in the username",
            newUser: &NewUser{
                Email:        "test@test.com",
                Password:     "password",
                PasswordConf: "password",
                UserName:     "",
                FirstName:    "test_fname",
                LastName:     "test_lname",
            },
            expectError: true,
        },
        {
            name: "Invalid NewUser - Username whitespace tab",
            hint: "\t",
            newUser: &NewUser{
                Email:        "test@test.com",
                Password:     "password",
                PasswordConf: "password",
                UserName:     "\t",
                FirstName:    "test_fname",
                LastName:     "test_lname",
            },
            expectError: true,
        },
        {
            name: "Invalid NewUser - Username whitespace newline",
            hint: "\n",
            newUser: &NewUser{
                Email:        "test@test.com",
                Password:     "password",
                PasswordConf: "password",
                UserName:     "\n",
                FirstName:    "test_fname",
                LastName:     "test_lname",
            },
            expectError: true,
        },
        {
            name: "Invalid NewUser - Username whitespace space",
            hint: " ",
            newUser: &NewUser{
                Email:        "test@test.com",
                Password:     "password",
                PasswordConf: "password",
                UserName:     " ",
                FirstName:    "test_fname",
                LastName:     "test_lname",
            },
            expectError: true,
        },
    }
    for _, c := range cases {
        err := c.newUser.Validate()
        if err != nil && !c.expectError {
            t.Errorf("case %s: unexpected error: %v\nHINT: %s", c.name, err, c.hint)
        }
        if c.expectError && err == nil {
            t.Errorf("case %s: expected error but didn't get one\nHINT: %s", c.name, c.hint)
        }
    }
}

func TestNewUserToUser(t *testing.T) {
    cases := []struct {
        name        string
        hint        string
        newUser     *NewUser
        expectError bool
    }{
        {
            name: "NewUser to User",
            hint: "",
            newUser: &NewUser{
                Email:        "test@test.edu",
                Password:     "password",
                PasswordConf: "password",
                UserName:     "test",
                FirstName:    "test_fname",
                LastName:     "test_lname",
            },
            expectError: false,
        },
    }

    for _, c := range cases {
        user, err := c.newUser.ToUser()
        if err != nil && !c.expectError {
            t.Errorf("case %s: unexpected error: %v\nHINT: %s", c.name, err, c.hint)
        }
        if c.expectError && err == nil {
            t.Errorf("case %s: expected error but didn't get one\nHINT: %s", c.name, c.hint)
        }
        if !c.expectError && user == nil {
            t.Errorf("case %s: nil user returned\nHINT: %s", c.name, c.hint)
        }
    }
}

func TestUserFullName(t *testing.T) {
    cases := []struct {
        name         string
        hint         string
        user         *User
        expectedName string
    }{
        {
            name:         "NewUser to User",
            hint:         "Make sure that the password is being hashed properly",
            user:         &User{
                FirstName: "fname",
                LastName: "lname",
            },
            expectedName: "fname lname",
        },
        {
            name:         "NewUser to User",
            hint:         "Check setting a no length password",
            user:        &User{
                FirstName: "fname",
                LastName: "",
            },
            expectedName: "fname",
        },
        {
            name:         "NewUser to User",
            hint:         "Check setting a no length password",
            user:        &User{
                FirstName: "",
                LastName: "lname",
            },
            expectedName: "lname",
        },
    }

    for _, c := range cases {
        fullName := c.user.FullName()
        if fullName != c.expectedName {
            t.Errorf("case %s: full name didnt equal expected name\nGot: %s | expected: %s\nHINT: %s", c.name, fullName, c.expectedName, c.hint)
        }
    }
}

func TestUserSetPassword(t *testing.T) {
    cases := []struct {
        name        string
        hint        string
        user        *User
        password    string
        expectError bool
    }{
        {
            name:        "NewUser to User",
            hint:        "Make sure that the password is being hashed properly",
            user:        &User{},
            password:    "password",
            expectError: false,
        },
        {
            name:        "NewUser to User",
            hint:        "Check setting a no length password",
            user:        &User{},
            password:    "",
            expectError: true,
        },
    }

    for _, c := range cases {
        err := c.user.SetPassword(c.password)
        if err != nil && !c.expectError {
            t.Errorf("case %s: unexpected error: %v\nHINT: %s", c.name, err, c.hint)
        }
        if c.expectError && err == nil {
            t.Errorf("case %s: expected error but didn't get one\nHINT: %s", c.name, c.hint)
        }
        err = bcrypt.CompareHashAndPassword(c.user.PassHash, []byte(c.password))
        if err != nil && !c.expectError {
            t.Errorf("case %s: password hash comparison and password had err: V%v\nHINT: %s", c.name, err, c.hint)
        }
    }
}

func TestUserAuthenticate(t *testing.T) {
    cases := []struct {
        name        string
        hint        string
        user        *User
        password    string
        expectError bool
    }{
        {
            name:        "NewUser to User",
            hint:        "Make sure that the password is being hashed properly",
            user:        &User{
                PassHash: []byte("$2a$13$gbfbvZNI1TdPiMf7Zo3gXOmjeZWR3nfrijQAoh2sWx9CLRfRDbKsO"),
            },
            password:    "password",
            expectError: false,
        },
        {
            name:        "NewUser to User",
            hint:        "Check setting a no length password",
            user:        &User{
                PassHash: []byte("invalid_hash"),
            },
            password:    "password",
            expectError: true,
        },
    }

    for _, c := range cases {
        err := c.user.Authenticate(c.password)
        if err != nil && !c.expectError {
            t.Errorf("case %s: unexpected error: %v\nHINT: %s", c.name, err, c.hint)
        }
        if c.expectError && err == nil {
            t.Errorf("case %s: expected error but didn't get one\nHINT: %s", c.name, c.hint)
        }
    }
}

func TestUserApplyUpdates(t *testing.T) {
    cases := []struct {
        name        string
        hint        string
        user        *User
        updates    *Updates
        expectError bool
    }{
        {
            name:        "NewUser to User",
            hint:        "Make sure that the password is being hashed properly",
            user:        &User{
                FirstName: "fname",
                LastName: "lname",
            },
            updates:    &Updates{
                FirstName: "new_fname",
                LastName: "new_lname",
            },
            expectError: false,
        },
        {
            name:        "ApplyUpdates fail due to no last name",
            hint:        "Make sure there is a no len check",
            user:        &User{
                FirstName: "fname",
                LastName: "lname",
            },
            updates:    &Updates{
                FirstName: "new_fname",
                LastName: "",
            },
            expectError: true,
        },
        {
            name:        "ApplyUpdates fail due to no last name",
            hint:        "Make sure there is a no len check",
            user:        &User{
                FirstName: "fname",
                LastName: "lname",
            },
            updates:    &Updates{
                FirstName: "",
                LastName: "new_lname",
            },
            expectError: true,
        },
    }

    for _, c := range cases {
        err := c.user.ApplyUpdates(c.updates)
        if err != nil && !c.expectError {
            t.Errorf("case %s: unexpected error: %v\nHINT: %s", c.name, err, c.hint)
        }
        if c.expectError && err == nil {
            t.Errorf("case %s: expected error but didn't get one\nHINT: %s", c.name, c.hint)
        }
    }
}