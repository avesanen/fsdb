package fsdb

import (
	"errors"
	"os"
	"sync"
)

// collection is a subdirectory of fsdb path
type collection struct {
	Name string          `json:"name"`
	Path string          `json:"-"`
	Keys map[string]*key `json:"keys"`
}

// read will read key to v
func (c *collection) read(key string, v interface{}) error {
	if c.Keys[key] == nil {
		return errors.New("key does not exist: " + key)
	}
	return c.Keys[key].read(v)
}

// write will write v to key
func (c *collection) write(key string, v interface{}) error {
	if c.Keys[key] == nil {
		k, err := c.newKey(key)
		if err != nil {
			return err
		}
		c.Keys[key] = k
	}
	return c.Keys[key].write(v)
}

// delete will delete the key
func (c *collection) delete(key string) error {
	if c.Keys[key] == nil {
		return errors.New("key does not exist: " + key)
	}
	err := c.Keys[key].delete()
	if err != nil {
		return err
	}
	delete(c.Keys, key)
	return nil
}

// newKey returns a new key
func (c *collection) newKey(name string) (*key, error) {
	path := c.Path + string(os.PathSeparator) + name
	k := &key{}
	k.Name = name
	k.Path = path
	k.RWMutex = &sync.RWMutex{}
	return k, nil
}
