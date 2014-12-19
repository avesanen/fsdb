# FsDb

FsDb is a small pet project that uses filesystem with json files as database.


## Installation

```bash
go get github.com/avesanen/fsdb
```


## Usage


```go
package main

import "github.com/avesanen/fsdb"
import "log"

func main() {
	fsdb := fsdb.NewFsDb('./db')

	item := map[string]interface{}{
		"msg": "Hello world!",
	}

	// Write the item 'hello' to collection 'messages'.
	if err := fsdb.Write('messages','hello') err != nil {
		panic(err)
	}

	// Read item 'hello' from collection 'messages'.
	msg, err := fsdb.Read('messages', 'hello')
	if err != nil {
		panic(err)
	}

	log.Println(msg)

	// Delete item from collection
	if err := fsdb.Delete('messages', 'hello'); err != nil {
		panic(err)
	}
}
```