// Copyright 2019 Tim Shannon. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package generichold_test

import (
	"github.com/rlshukhov/generichold"
	"testing"
	"time"

	"github.com/dgraph-io/badger/v4"
	"github.com/timshannon/badgerhold/v4"
)

func TestGet(t *testing.T) {
	testWrap(t, func(bh *badgerhold.Store, t *testing.T) {
		store := generichold.Open[ItemTest](bh)

		key := "testKey"
		data := &ItemTest{
			Name:    "Test Name",
			Created: time.Now(),
		}
		err := store.Insert(key, data)
		if err != nil {
			t.Fatalf("Error creating data for get test: %s", err)
		}

		result, err := store.Get(key)
		if err != nil {
			t.Fatalf("Error getting data from badgerhold: %s", err)
		}

		if !data.equal(&result) {
			t.Fatalf("Got %v wanted %v.", result, data)
		}
	})
}

func TestIssue36(t *testing.T) {
	testWrap(t, func(bh *badgerhold.Store, t *testing.T) {
		type Tag1 struct {
			ID uint64 `badgerholdKey`
		}

		type Tag2 struct {
			ID uint64 `badgerholdKey:"Key"`
		}

		type Tag3 struct {
			ID uint64 `badgerhold:"key"`
		}

		type Tag4 struct {
			ID uint64 `badgerholdKey:""`
		}

		store1 := generichold.Open[Tag1](bh)
		data1 := []Tag1{{}, {}, {}}
		for i := range data1 {
			ok(t, store1.Insert(badgerhold.NextSequence(), &data1[i]))
			equals(t, uint64(i), data1[i].ID)
		}

		store2 := generichold.Open[Tag2](bh)
		data2 := []Tag2{{}, {}, {}}
		for i := range data2 {
			ok(t, store2.Insert(badgerhold.NextSequence(), &data2[i]))
			equals(t, uint64(i), data2[i].ID)
		}

		store3 := generichold.Open[Tag3](bh)
		data3 := []Tag3{{}, {}, {}}
		for i := range data3 {
			ok(t, store3.Insert(badgerhold.NextSequence(), &data3[i]))
			equals(t, uint64(i), data3[i].ID)
		}

		store4 := generichold.Open[Tag4](bh)
		data4 := []Tag4{{}, {}, {}}
		for i := range data4 {
			ok(t, store4.Insert(badgerhold.NextSequence(), &data4[i]))
			equals(t, uint64(i), data4[i].ID)
		}

		// Get
		for i := range data1 {
			get1, err := store1.Get(data1[i].ID)
			ok(t, err)
			equals(t, data1[i], get1)
		}

		for i := range data2 {
			get2, err := store2.Get(data2[i].ID)
			ok(t, err)
			equals(t, data2[i], get2)
		}

		for i := range data3 {
			get3, err := store3.Get(data3[i].ID)
			ok(t, err)
			equals(t, data3[i], get3)
		}

		for i := range data4 {
			get4, err := store4.Get(data4[i].ID)
			ok(t, err)
			equals(t, data4[i], get4)
		}

		// Find

		for i := range data1 {
			find1, err := store1.Find(badgerhold.Where(badgerhold.Key).Eq(data1[i].ID))
			ok(t, err)
			assert(t, len(find1) == 1, "incorrect rows returned")
			equals(t, find1[0], data1[i])
		}

		for i := range data2 {
			find2, err := store2.Find(badgerhold.Where(badgerhold.Key).Eq(data2[i].ID))
			ok(t, err)
			assert(t, len(find2) == 1, "incorrect rows returned")
			equals(t, find2[0], data2[i])
		}

		for i := range data3 {
			find3, err := store3.Find(badgerhold.Where(badgerhold.Key).Eq(data3[i].ID))
			ok(t, err)
			assert(t, len(find3) == 1, "incorrect rows returned")
			equals(t, find3[0], data3[i])
		}

		for i := range data4 {
			find4, err := store4.Find(badgerhold.Where(badgerhold.Key).Eq(data4[i].ID))
			ok(t, err)
			assert(t, len(find4) == 1, "incorrect rows returned")
			equals(t, find4[0], data4[i])
		}
	})
}

func TestTxGetBadgerError(t *testing.T) {
	testWrap(t, func(bh *badgerhold.Store, t *testing.T) {
		store := generichold.Open[ItemTest](bh)

		key := "testKey"
		data := &ItemTest{
			Name:    "Test Name",
			Created: time.Now(),
		}
		err := store.Insert(key, data)
		if err != nil {
			t.Fatalf("Error creating data for TxGet test: %s", err)
		}

		txn := store.Badger().NewTransaction(false)
		txn.Discard()

		_, err = store.TxGet(txn, key)
		if err != badger.ErrDiscardedTxn {
			t.Fatalf("TxGet didn't fail! Expected %s got %s", badger.ErrDiscardedTxn, err)
		}
	})
}
