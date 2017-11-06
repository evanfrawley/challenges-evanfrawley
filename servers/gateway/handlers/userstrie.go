package handlers

import (
    "github.com/info344-a17/challenges-evanfrawley/servers/gateway/indexes"
    "github.com/info344-a17/challenges-evanfrawley/servers/gateway/models/users"
    "gopkg.in/mgo.v2/bson"
    "fmt"
)

func CreateTrieRoot() *indexes.TrieNode {
    root := &indexes.TrieNode{
        Data:           []rune(" ")[0],
        NextNodes:      map[rune]*indexes.TrieNode{},
        CompletedItems: []*indexes.CompletedItem{},
        Parent:         nil,
    }
    return root
}

func ConstructUsersTrie(usersStore *users.MongoStore) *indexes.TrieNode {
    root := CreateTrieRoot()
    userIterator := usersStore.GetUserIterator()
    tempUser := &users.User{}
    for userIterator.Next(tempUser) {
        root.InsertUser(tempUser)
    }
    return root
}

func (ctx *Context) DeleteUserFromTrie(userID bson.ObjectId) error {
    user, err := ctx.userMongoStore.GetByID(userID)
    if err != nil {
        return fmt.Errorf("error when getting user: %v", err)
    }
    if err := ctx.trieRoot.DeleteUser(user); err != nil {
        return fmt.Errorf("error removing user from trie: %v", err)
    }

    return nil
}
