// Copyright 2019 Tim Shannon. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package generichold_test

import (
	"encoding/json"
	"fmt"
	"github.com/rlshukhov/generichold"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"

	"github.com/timshannon/badgerhold/v4"
)

func TestOpen(t *testing.T) {
	opt := testOptions()
	bh, err := badgerhold.Open(opt)
	if err != nil {
		t.Fatalf("Error opening %s: %s", opt.Dir, err)
	}

	store := generichold.Open[any](bh)

	if store == nil {
		t.Fatalf("store is null!")
	}

	err = store.Close()
	if err != nil {
		t.Fatal(err)
	}
	err = os.RemoveAll(opt.Dir)
	if err != nil {
		t.Fatal(err)
	}
}

func TestBadger(t *testing.T) {
	testWrap(t, func(bh *badgerhold.Store, t *testing.T) {
		store := generichold.Open[any](bh)

		b := store.Badger()
		if b == nil {
			t.Fatalf("Badger is null in generichold")
		}
	})
}

func TestAlternateEncoding(t *testing.T) {
	opt := testOptions()
	opt.Encoder = json.Marshal
	opt.Decoder = json.Unmarshal
	bh, err := badgerhold.Open(opt)
	if err != nil {
		t.Fatalf("Error opening %s: %s", opt.Dir, err)
	}
	store := generichold.Open[ItemTest](bh)

	defer os.RemoveAll(opt.Dir)
	defer store.Close()

	insertTestData(t, store)

	tData := testData[3]

	result, _ := store.Find(badgerhold.Where(badgerhold.Key).Eq(tData.Key))

	if len(result) != 1 {
		if testing.Verbose() {
			t.Fatalf("Find result count is %d wanted %d.  Results: %v", len(result), 1, result)
		}
		t.Fatalf("Find result count is %d wanted %d.", len(result), 1)
	}

	if !result[0].equal(&tData) {
		t.Fatalf("Results not equal! Wanted %v, got %v", tData, result[0])
	}

}

func TestGetUnknownType(t *testing.T) {
	opt := testOptions()
	bh, err := badgerhold.Open(opt)
	if err != nil {
		t.Fatalf("Error opening %s: %s", opt.Dir, err)
	}

	type test struct {
		Test string
	}
	store := generichold.Open[test](bh)

	defer os.RemoveAll(opt.Dir)
	defer store.Close()

	_, err = store.Get("unknownKey")
	if err != badgerhold.ErrNotFound {
		t.Errorf("Expected error of type ErrNotFound, not %T", err)
	}
}

type ItemWithStorer struct{ Name string }

func (i *ItemWithStorer) Type() string { return "Item" }
func (i *ItemWithStorer) Indexes() map[string]badgerhold.Index {
	return map[string]badgerhold.Index{
		"Name": {
			IndexFunc: func(_ string, value interface{}) ([]byte, error) {
				// If the upsert wants to delete an existing value first,
				// value could be a **Item instead of *Item
				// panic: interface conversion: interface {} is **Item, not *Item
				v := value.(*ItemWithStorer).Name
				return badgerhold.DefaultEncode(v)
			},
			Unique: false,
		},
	}
}

func TestIssue115(t *testing.T) {
	testWrap(t, func(bh *badgerhold.Store, t *testing.T) {
		store := generichold.Open[ItemWithStorer](bh)

		item := ItemWithStorer{"Name"}
		for i := 0; i < 2; i++ {
			err := store.Upsert("key", &item)
			if err != nil {
				t.Fatal(err)
			}
		}
	})
}

func TestIssue70TypePrefixCollision(t *testing.T) {
	testWrap(t, func(bh *badgerhold.Store, t *testing.T) {
		type TestStruct struct {
			Value int
		}

		type TestStructCollision struct {
			Value int
		}

		store := generichold.Open[TestStruct](bh)
		storeCollision := generichold.Open[TestStructCollision](bh)

		for i := 0; i < 5; i++ {
			ok(t, store.Insert(i, &TestStruct{Value: i}))
			ok(t, storeCollision.Insert(i, &TestStructCollision{Value: i}))
		}

		query := badgerhold.Where(badgerhold.Key).In(0, 1, 2, 3, 4)
		results, err := store.Find(query)
		ok(t, err)

		equals(t, 5, len(results))
	})
}

func TestIssue71IndexByCustomName(t *testing.T) {
	testWrap(t, func(bh *badgerhold.Store, t *testing.T) {
		type Person struct {
			Name     string
			Division string `badgerholdIndex:"IdxDivision"`
		}

		store := generichold.Open[Person](bh)

		record := Person{Name: "test", Division: "testDivision"}

		ok(t, store.Insert(1, &record))

		_, err := store.Find(badgerhold.Where("Division").Eq(record.Division).Index("IdxDivision"))
		ok(t, err)
	})
}

// utilities

func testWrap(t *testing.T, tests func(store *badgerhold.Store, t *testing.T)) {
	testWrapWithOpt(t, testOptions(), tests)
}

func testWrapWithOpt(t *testing.T, opt badgerhold.Options, tests func(store *badgerhold.Store, t *testing.T)) {
	var err error
	store, err := badgerhold.Open(opt)
	if err != nil {
		t.Fatalf("Error opening %s: %s", opt.Dir, err)
	}

	if store == nil {
		t.Fatalf("store is null!")
	}

	tests(store, t)
	store.Close()
	os.RemoveAll(opt.Dir)
}

type emptyLogger struct{}

func (e emptyLogger) Errorf(msg string, data ...interface{})   {}
func (e emptyLogger) Infof(msg string, data ...interface{})    {}
func (e emptyLogger) Warningf(msg string, data ...interface{}) {}
func (e emptyLogger) Debugf(msg string, data ...interface{})   {}

func testOptions() badgerhold.Options {
	opt := badgerhold.DefaultOptions
	opt.InMemory = true
	//opt.Dir = tempdir()
	//opt.ValueDir = opt.Dir
	opt.Logger = emptyLogger{}
	// opt.ValueLogLoadingMode = options.FileIO // slower but less memory usage
	// opt.TableLoadingMode = options.FileIO
	// opt.NumMemtables = 1
	// opt.NumLevelZeroTables = 1
	// opt.NumLevelZeroTablesStall = 2
	// opt.NumCompactors = 1
	return opt
}

// tempdir returns a temporary dir path.
func tempdir() string {
	name, err := ioutil.TempDir("", "generichold-")
	if err != nil {
		panic(err)
	}
	return name
}

// assert fails the test if the condition is false.
func assert(tb testing.TB, condition bool, msg string, v ...interface{}) {
	if !condition {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: "+msg+"\033[39m\n\n", append([]interface{}{filepath.Base(file), line}, v...)...)
		tb.FailNow()
	}
}

// ok fails the test if an err is not nil.
func ok(tb testing.TB, err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: unexpected error: %s\033[39m\n\n", filepath.Base(file), line, err.Error())
		tb.FailNow()
	}
}

// equals fails the test if exp is not equal to act.
func equals(tb testing.TB, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\033[39m\n\n", filepath.Base(file), line, exp, act)
		tb.FailNow()
	}
}
