package users

import (
    "gopkg.in/mgo.v2/bson"
    "gopkg.in/mgo.v2"
    "fmt"
)

type userFilter struct {
    Email    string        `bson:"email,omitempty"`
    ID       bson.ObjectId `bson:"_id,omitempty"`
    Username string        `bson:"username,omitempty"`
}

type userUpdateDoc struct {
    PassHash string `bson:"-,omitempty"`
    FirstName string `bson:"firstName,omitempty"`
    LastName string `bson:"lastName,omitempty"`
    Email string `bson:"email,omitempty"`
    Username string `bson:"username,omitempty"`
}

//MongoStore implements Store for MongoDB
type MongoStore struct {
    session *mgo.Session
    dbName  string
    colName string
}

//NewMongoStore constructs a new MongoStore
func NewMongoStore(sess *mgo.Session, dbName string, collectionName string) *MongoStore {
    if sess == nil {
        panic("nil pointer passed for session")
    }
    return &MongoStore{
        session: sess,
        dbName:  dbName,
        colName: collectionName,
    }
}

// Gets a single user by the given ID
func (s *MongoStore) GetByID(id bson.ObjectId) (*User, error) {
    user := &User{}
    userFilter := &userFilter{ID: id}
    return s.getWithFilter(user, userFilter)
}

// Gets a single user by the given email
func (s *MongoStore) GetByEmail(email string) (*User, error) {
    user := &User{}
    userFilter := &userFilter{Email: email}
    return s.getWithFilter(user, userFilter)
}

// Gets a single user by the given username
func (s *MongoStore) GetByUserName(username string) (*User, error) {
    user := &User{}
    userFilter := &userFilter{Username: username}
    return s.getWithFilter(user, userFilter)
}

// Gets a single user with a specific applied filter
func (s *MongoStore) getWithFilter(user *User, userFilter *userFilter) (*User, error) {
    col := s.session.DB(s.dbName).C(s.colName)
    if err := col.Find(userFilter).Limit(1).One(user); err != nil {
        return nil, ErrUserNotFound
    }
    return user, nil
}

// Inserts a new user into the DB
func (s *MongoStore) Insert(newUser *NewUser) (*User, error) {
    fmt.Printf("here at the beginning")
    user, err := newUser.ToUser()
    if err != nil {
        return nil, err
    }
    fmt.Printf("here at the middle")

    col := s.session.DB(s.dbName).C(s.colName)
    fmt.Printf("here at the after the db write")

    if err := col.Insert(user); err != nil {
        return nil, fmt.Errorf("error inerting user: %v", err)
    }
    fmt.Printf("here at the end")
    return user, nil
}

// Updates a specific user given the user ID
func (s *MongoStore) Update(userID bson.ObjectId, updates *Updates) error {
    upd := &userUpdateDoc {
        FirstName: updates.FirstName,
        LastName: updates.LastName,
    }

    change := mgo.Change{
        Update: bson.M{"$set": upd},
        ReturnNew: true,
    }

    user := &User{}
    col := s.session.DB(s.dbName).C(s.colName)
    if _, err := col.FindId(userID).Apply(change, user); err != nil {
        return fmt.Errorf("error updating task: %v", err)
    }

    return nil
}

// Deletes a user given the user ID
func (s *MongoStore) Delete(userID bson.ObjectId) error {
    return nil
}
