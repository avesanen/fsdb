package fsdb

import "testing"
import "os"

func TestOpen(t *testing.T) {
	dbPath := os.TempDir() + string(os.PathSeparator) + "db"
	if err := os.RemoveAll(dbPath); err != nil {
		t.Fatal(err.Error())
	}

	os.Mkdir(dbPath, 0777)

	os.Mkdir(dbPath+string(os.PathSeparator)+"col1", 0777)

	file, err := os.Create(dbPath + string(os.PathSeparator) + "col1" + string(os.PathSeparator) + "key1")

	if err != nil {
		t.Fatal(err.Error())
	} else {
		file.Write([]byte(`{"msg":"Hello world!"}`))
		file.Close()
	}

	fsdb, err := NewFsDb(dbPath)
	if err != nil {
		t.Fatal("should not fail if db does exist.")
	}

	if fsdb == nil {
		t.Fatal("fsdb not created, failing.")
	}

	// Database functions

	msg, err := fsdb.Read("col1", "key1")
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log(msg)

	msg, err = fsdb.Read("col1", "key2")
	if err == nil {
		t.Fatal("Key doesn't exist, should return error.")
	}

	msg = map[string]interface{}{
		"msg": "Hello world!",
	}

	if err := fsdb.Write("col1", "key1", msg); err != nil {
		t.Fatal(err.Error())
	}

	if err := fsdb.Write("col2", "key2", msg); err != nil {
		t.Fatal(err.Error())
	}

	if err := fsdb.Write("col3", "key3", msg); err != nil {
		t.Fatal(err.Error())
	}

	// Delete the key can test that reading it will return error.
	if err := fsdb.Delete("col3", "key3"); err != nil {
		t.Fatal(err.Error())
	}

	msg, err = fsdb.Read("col3", "key3")
	if err == nil {
		t.Fatal("key not deleted.")
	}

	t.Log(fsdb)
	if err := os.RemoveAll(dbPath); err != nil {
		t.Fatal(err.Error())
	}
}
