# GenericHold

[![Go](https://github.com/rlshukhov/generichold/actions/workflows/go.yml/badge.svg)]

GenericHold is [BadgerHold](https://github.com/timshannon/badgerhold) extension provided simplest, generics-powered, type-safe API.

GenericHold have rich tests coverage and pass all [BadgerHold](https://github.com/timshannon/badgerhold) unit tests.

## Usage example

```shell
go get github.com/rlshukhov/generichold
```

```go
package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/rlshukhov/generichold"
	"github.com/timshannon/badgerhold/v4"
)

type Item struct {
	ID       uint64 `badgerhold:"key"`
	Category string `badgerholdIndex:"Category"`
	Created  time.Time
}

func main() {
	// setup badger
	options := badgerhold.DefaultOptions
	options.InMemory = true
	options.Logger = nil

	// setup badgerhold
	bh, err := badgerhold.Open(options)
	defer bh.Close()
	if err != nil {
		log.Fatal(err)
	}

	store := generichold.Open[Item](bh)

	items := getItems()

	var ids []string
	for _, item := range items {
		err := store.Insert(badgerhold.NextSequence(), &item)
		if err != nil {
			log.Fatal(err)
		}

		// badgerhold.NextSequence() as key - will generate ID and set to original entity
		ids = append(ids, strconv.FormatUint(item.ID, 10))
	}
	fmt.Println(strings.Join(ids, ","))
	// Output: 0,1,2,3

	// Find all items in the blue category that have been created in the past hour
	result, err := store.Find(badgerhold.Where("Category").Eq("blue").And("Created").Ge(time.Now().Add(-1 * time.Hour)))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(len(result))
	// Output: 1

	fmt.Println(result[0].ID)
	// Output: 3
}

func getItems() []Item {
	return []Item{
		{
			Category: "blue",
			Created:  time.Now().Add(-4 * time.Hour),
		},
		{
			Category: "red",
			Created:  time.Now().Add(-3 * time.Hour),
		},
		{
			Category: "blue",
			Created:  time.Now().Add(-2 * time.Hour),
		},
		{
			Category: "blue",
			Created:  time.Now().Add(-20 * time.Minute),
		},
	}
}
```

## TODO

- Make `badgerhold.Criterion` generic version to avoid this limitation of BadgerHold:
> However if you have an existing slice of values to test against, you can't pass in that slice because it is not of type `[]interface{}`.
> ```go
> t := []string{"1", "2", "3", "4"}
> where := badgerhold.Where("Id").In(t...) // compile error
> ```
