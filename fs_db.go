package fsdb

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
)

// FsDb is a filesystem Database that automatically caches files for a duration of time.
type FsDb struct {
	Path        string                 `json:"-"`
	Collections map[string]*collection `json:"collections"`
}

// NewFsDb will return a new FsDb, or error if it can not be opened
func NewFsDb(path string) (*FsDb, error) {
	db := &FsDb{}
	db.Path = path
	db.Collections = make(map[string]*collection)

	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			// The path didn't exist, create it.
			if err := os.Mkdir(path, 0777); err != nil {
				return nil, err
			}
		} else {
			// Some other error, return it.
			return nil, err
		}
	}

	dbCollections, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	for _, f := range dbCollections {
		if f.IsDir() {
			c, err := db.newCollection(f.Name())
			if err != nil {
				return nil, err
			}
			c.Name = f.Name()
			db.Collections[f.Name()] = c
		}
	}

	return db, nil
}

// newCollection Returns a collection, populated with Items if the given path already exists.
func (db *FsDb) newCollection(name string) (*collection, error) {
	path := db.Path + string(os.PathSeparator) + name
	// If the path doesn't exist, create it. Return error if that fails.
	stat, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			// The path didn't exist, create it.
			if err := os.Mkdir(path, 0777); err != nil {
				return nil, err
			}
		} else {
			// Some other error, return it.
			return nil, err
		}
	}

	// Create the collection
	c := &collection{}
	c.Path = path
	c.Keys = make(map[string]*key)

	// If the directory did exist before, create keys.
	if stat != nil && stat.IsDir() {
		collectionKeys, err := ioutil.ReadDir(c.Path)
		if err != nil {
			return nil, err
		}

		for _, f := range collectionKeys {
			if !f.IsDir() {
				k, err := c.newKey(f.Name())
				if err != nil {
					return nil, err
				}
				k.Name = f.Name()
				c.Keys[k.Name] = k
			}
		}
	}
	return c, nil
}

// Write to key in collection
func (db *FsDb) Write(col, key string, v interface{}) error {
	if db.Collections[col] == nil {
		c, err := db.newCollection(col)
		if err != nil {
			return err
		}
		db.Collections[col] = c
	}
	return db.Collections[col].write(key, v)
}

// Read key from collection
func (db *FsDb) Read(col string, key string, v interface{}) error {
	if db.Collections[col] == nil {
		return errors.New("collection does not exist: " + col)
	}
	return db.Collections[col].read(key, v)
}

// Delete the key from collection
func (db *FsDb) Delete(col string, key string) error {
	if db.Collections[col] == nil {
		return errors.New("collection does not exist: " + col)
	}
	return db.Collections[col].delete(key)
}

// String will return the whole loaded database in json format (for debugging)
func (db *FsDb) String() string {
	b, err := json.Marshal(db)
	if err != nil {
		return ""
	}
	return "FSDB: " + string(b)
}
