package indexes

import (
    "testing"
    "github.com/info344-a17/challenges-evanfrawley/servers/gateway/models/users"
    "gopkg.in/mgo.v2/bson"
)

func TestTrieNodeInsert(t *testing.T) {
    userID1 := bson.NewObjectId()
    userID2 := bson.NewObjectId()

    insertCases := []struct {
        name string
        user users.User
        emailSize int
        prefixSize int
        prefix string
        expectError bool
    }{
        {
            name: "Test user",
            user: users.User{
                ID: userID1,
                FirstName: "test",
                LastName: "user",
                Email: "test@test.com",
                UserName: "testuser",
            },
            emailSize: 1,
            prefixSize: 3,
            prefix: "test",
            expectError: false,
        },
        {
            name: "Evan user",
            user: users.User{
                ID: userID2,
                FirstName: "evan",
                LastName: "frawley",
                Email: "frawley@uw.edu",
                UserName: "frawley",
            },
            emailSize: 1,
            prefixSize: 3,
            prefix: "frawley",
            expectError: true,
        },
    }

    root := &TrieNode{
        Data:           []rune(" ")[0],
        NextNodes:      map[rune]*TrieNode{},
        CompletedItems: []*CompletedItem{},
        Parent:         nil,
    }

    for _, c := range insertCases {
        root.InsertUser(&c.user)
        items := root.FindCompletedItemsWithPrefix(c.user.Email)
        if len(items) != c.emailSize {
            t.Errorf("%s:\nemail size did not match.\nGOT: %v\nEXPECTED: %v\n", c.name, len(items), c.emailSize)
        }
        items = root.FindCompletedItemsWithPrefix(c.prefix)
        if len(items) != c.prefixSize {
            t.Errorf("%s:\nprefix size did not match.\nGOT: %v\nEXPECTED: %v\n", c.name, len(items), c.prefixSize)
        }
    }
}

func TestTrieNodeDeleteKey(t *testing.T) {
    userID1 := bson.NewObjectId()
    userID2 := bson.NewObjectId()

    cases := []struct {
        name string
        user users.User
        prefixSize int
        keyToDelete string
        expectError bool
    }{
        {
            name: "Test user",
            user: users.User{
                ID: userID1,
                FirstName: "test",
                LastName: "user",
                Email: "test@test.com",
                UserName: "testuser",
            },
            prefixSize: 0,
            keyToDelete: "test@test.com",
            expectError: false,
        },
        {
            name: "Evan user",
            user: users.User{
                ID: userID2,
                FirstName: "evan",
                LastName: "frawley",
                Email: "frawley@uw.edu",
                UserName: "frawley",
            },
            prefixSize: 3,
            keyToDelete: "frawley@uw.edu",
            expectError: true,
        },
    }

    root := &TrieNode{
        Data:           []rune(" ")[0],
        NextNodes:      map[rune]*TrieNode{},
        CompletedItems: []*CompletedItem{},
        Parent:         nil,
    }

    for _, c := range cases {
        root.InsertUser(&c.user)
        root.DeleteKey(c.keyToDelete, c.user.ID)
        items := root.FindCompletedItemsWithPrefix(c.keyToDelete)
        if len(items) != 0 {
            t.Errorf("%s:\nexpected items to be empty, but got size: %v", c.name, len(items))
        }
        var ci []CompletedItem
        root.dfs(&ci)
    }
}

func TestTrieNodeDeleteUser(t *testing.T) {
    userID1 := bson.NewObjectId()
    userID2 := bson.NewObjectId()
    userID3 := bson.NewObjectId()

    casesInsert := []struct {
        name string
        user users.User
        expectError bool
    }{
        {
            name: "Test user",
            user: users.User{
                ID: userID1,
                FirstName: "test",
                LastName: "user",
                Email: "test@test.com",
                UserName: "testuser",
            },
            expectError: false,
        },
        {
            name: "Evan user",
            user: users.User{
                ID: userID2,
                FirstName: "evan",
                LastName: "frawley",
                Email: "frawley@uw.edu",
                UserName: "frawley",
            },

            expectError: false,
        },
    }

    casesDelete := []struct {
        name string
        user users.User
        expectError bool
    }{
        {
            name: "Does not exist",
            user: users.User{
                ID: bson.NewObjectId(),
                FirstName: "doesnot",
                LastName: "exist",
                Email: "dne@dne.com",
                UserName: "doesnotexist",
            },
            expectError: true,
        },
    }

    casesDelete = append(casesInsert, casesDelete...)

    root := &TrieNode{
        Data:           []rune(" ")[0],
        NextNodes:      map[rune]*TrieNode{},
        CompletedItems: []*CompletedItem{},
        Parent:         nil,
    }

    user3 := &users.User{
        ID: userID3,
        FirstName: "tes",
        LastName: "user",
        Email: "test@test.com1",
        UserName: "monkeys",
    }

    root.InsertUser(user3)

    for _, c := range casesInsert {
        root.InsertUser(&c.user)
    }

    for _, c := range casesDelete {
        err := root.DeleteUser(&c.user)
        items := root.GetUniqueUsersFromPrefix(c.user.UserName)
        if !c.expectError && err != nil && len(items) != 0 {
            t.Errorf("%s:\nexpected items to be empty, but got size: %v", c.name, len(items))
        }

        if c.expectError && err == nil {
            t.Errorf("%s:\nexpected and error but did not get one", c.name)
        }
    }

    leftOverUser := root.GetUniqueUsersFromPrefix(user3.UserName)
    if len(leftOverUser) != 1 {
        t.Errorf("expected to get 1 user back for the test user username, but got: %v", len(leftOverUser))
    }

    leftOverUser = root.GetUniqueUsersFromPrefix(user3.Email)
    if len(leftOverUser) != 1 {
        t.Errorf("expected to get 1 user back for the test user email, but got: %v", len(leftOverUser))
    }

    leftOverUser = root.GetUniqueUsersFromPrefix(user3.LastName)
    if len(leftOverUser) != 1 {
        t.Errorf("expected to get 1 user back for the test user last name, but got: %v", len(leftOverUser))
    }

    leftOverUser = root.GetUniqueUsersFromPrefix(user3.FirstName)
    if len(leftOverUser) != 1 {
        t.Errorf("expected to get 1 user back for the test user first name, but got: %v", len(leftOverUser))
    }
}