package users

import (
    "gopkg.in/mgo.v2/bson"
)

//MemStore represents an in-process memory user store.
//This should be used only for testing and prototyping.
//Production systems should use a shared server store like redis
type MemStore struct {
    users []User
}

//NewMemStore constructs and returns a new MemStore
func NewMemStore() *MemStore {
    return &MemStore{
        users: []User{},
    }
}

func (ms *MemStore) GetByID(uid bson.ObjectId) (*User, error) {
    for _, user := range ms.users {
        if user.ID == uid {
            return &user, nil
        }
    }

    return nil, ErrUserNotFound
}

func (ms *MemStore) GetByEmail(email string) (*User, error) {
    for _, user := range ms.users {
        if user.Email == email {
            return &user, nil
        }
    }

    return nil, ErrUserNotFound
}

func (ms *MemStore) GetByUserName(username string) (*User, error) {
    for _, user := range ms.users {
        if user.UserName == username {
            return &user, nil
        }
    }

    return nil, ErrUserNotFound
}

//Delete deletes all user data associated with the user ID from the store.
func (ms *MemStore) Delete(uid bson.ObjectId) error {
    for index, user := range ms.users {
        if user.ID == uid {
            ms.users = append(ms.users[:index], ms.users[index + 1:]...)
            return nil
        }
    }

    return ErrUserNotFound
}

//Inserts the user data associated with the user ID from the store.
func (ms *MemStore) Insert(userToInsert *User) (*User, error) {
    ms.users = append(ms.users, *userToInsert)
    return userToInsert, nil
}

func (ms *MemStore) Update(uid bson.ObjectId, updates *Updates) error {
    for index, user := range ms.users {
        if user.ID == uid {
            if updates.FirstName != "" {
                user.FirstName = updates.FirstName
            }
            if updates.LastName != "" {
                user.LastName = updates.LastName
            }
            ms.users[index] = user
            return nil
        }
    }

    return ErrUserNotFound
}