package indexes

import (
    "fmt"
    "gopkg.in/mgo.v2/bson"

    "github.com/info344-a17/challenges-evanfrawley/servers/gateway/models/users"
    "strings"
)

const (
    FullName  = "fullName"
    FirstName = "firstName"
    LastName  = "lastName"
    UserName  = "username"
    Email     = "email"
)

type TrieNode struct {
    Data           rune
    CompletedItems []*CompletedItem
    NextNodes      map[rune]*TrieNode
    Parent         *TrieNode
}

type CompletedItem struct {
    //Val    string        `json:"val,omitempty"`
    Type   string        `json:"type,omitempty"`
    UserID bson.ObjectId `json:"userId,omitempty"`
}

func (n *TrieNode) InsertUser(user *users.User) {
    n.InsertItem(strings.ToLower(user.FirstName), FirstName, user)
    n.InsertItem(strings.ToLower(user.LastName), LastName, user)
    n.InsertItem(strings.ToLower(user.Email), Email, user)
    n.InsertItem(strings.ToLower(user.UserName), UserName, user)
}

//func (n *TrieNode) InsertItem(baseString, currentString, stringType string, u *users.User) {
func (n *TrieNode) InsertItem(currentString, stringType string, u *users.User) {
    if len(currentString) == 0 {
        return
    }

    // pluck first rune from the string
    nextRune := []rune(currentString)[0]
    // if no node exists for plucked rune, then created one
    nextNode, found := n.NextNodes[nextRune]
    if !found {
        newNode := &TrieNode{
            Data:           nextRune,
            NextNodes:      map[rune]*TrieNode{},
            CompletedItems: []*CompletedItem{},
            Parent:         n,
        }
        n.NextNodes[nextRune] = newNode
        nextNode = newNode
    }
    if len(currentString) == 1 {
        // Case when we've reached the end of the string and just added the last rune to the children
        completedItem := &CompletedItem{
            //Val:    baseString,
            Type:   stringType,
            UserID: u.ID,
        }
        nextNode.CompletedItems = append(nextNode.CompletedItems, completedItem)
    } else {
        // Otherwise, pass down the string without the first element down the chain
        restOfString := currentString[1:]
        nextNode.InsertItem(restOfString, stringType, u)
        //nextNode.InsertItem(baseString, restOfString, stringType, u)
    }
}

func (n *TrieNode) findCompletedItemsWithPrefix(prefix string) []CompletedItem {
    var completedItems []CompletedItem
    // not really sure if this is the right logic here
    targetNode, err := n.getTargetNode(prefix)

    if err != nil {
        // if not found, return an empty array
        fmt.Printf("getting an error: %v\n", err)
        return completedItems
    }

    targetNode.dfs(&completedItems)

    return completedItems
}

func (n *TrieNode) GetUniqueUsersFromPrefix(prefix string) []CompletedItem {
    prefix = strings.ToLower(prefix)
    items := n.findCompletedItemsWithPrefix(prefix)
    var ci []CompletedItem
    idToCompletedItemMap := map[bson.ObjectId]CompletedItem{}
    for _, item := range items {
        idToCompletedItemMap[item.UserID] = item
    }
    for _, val := range idToCompletedItemMap {
        // ensures that the items returned are unique
        ci = append(ci, val)
    }
    return ci
}

func (n *TrieNode) getTargetNode(key string) (*TrieNode, error) {
    var targetNode = n
    prefixLen := len(key)
    // get to the node at the base of the prefix, at most iterating through the prefix length
    for i := 0; i < prefixLen; i++ {
        runeKey := []rune(key)[0]
        key = key[1:]
        nextNode, found := targetNode.NextNodes[runeKey]
        // if the prefix does not exist in the trie, return the empty items slice
        if !found {
            return nil, fmt.Errorf("key: \"%v\" not found in trie", string(runeKey))
        }
        targetNode = nextNode
    }
    return targetNode, nil
}

func (n *TrieNode) dfs(completedItems *[]CompletedItem) {
    // do a dfs through the children nodes
    for _, node := range n.NextNodes {
        node.dfs(completedItems)
    }
    // any completed items that are found, add them to the completed items
    for _, completedItem := range n.CompletedItems {
        *completedItems = append(*completedItems, *completedItem)
    }
}

func (n *TrieNode) DeleteUser(user *users.User) error {
    // four cases
    // 1. leaf with only 1 thing
    // 2. has children part of a greater string
    // 3. popping up and there is a child of someehere in the string
    // 4. not in the trie

    // first, check to see if the user is indeed within the trie
    items := n.findCompletedItemsWithPrefix(user.Email)
    if len(items) == 0 {
        return fmt.Errorf("error: user was not found in the trie")
    }

    // now we know that the user is in the trie
    n.DeleteKey(user.FirstName, user.ID)
    n.DeleteKey(user.LastName, user.ID)
    n.DeleteKey(user.Email, user.ID)
    n.DeleteKey(user.UserName, user.ID)

    return nil
}

func (n *TrieNode) DeleteKey(key string, userIDToDelete bson.ObjectId) error {
    targetNode, err := n.getTargetNode(key)
    if err != nil {
        return err
    }

    // evaluate target node
    // first remove the completed item
    for i, item := range targetNode.CompletedItems {
        if item.UserID == userIDToDelete {
            if len(targetNode.CompletedItems) > 1 {
                targetNode.CompletedItems = append(targetNode.CompletedItems[:i], targetNode.CompletedItems[i+1:]...)
            } else if len(targetNode.CompletedItems) == 1 && targetNode.CompletedItems[0].UserID == userIDToDelete {
                targetNode.CompletedItems = []*CompletedItem{}
            } else {
                return fmt.Errorf("tried to delete a user that was not part of the trie")
            }
        }
    }

    // remove nodes equal to the number of elements in the key
    for i := 0; i < len(key); i++ {
        if len(targetNode.NextNodes) == 0 {
            if len(targetNode.CompletedItems) == 0 {
                runeValueToRemove := targetNode.Data
                targetNode = targetNode.Parent
                delete(targetNode.NextNodes, runeValueToRemove)
            } else {
                break
            }
        } else {
            break
        }
    }

    return nil
}
