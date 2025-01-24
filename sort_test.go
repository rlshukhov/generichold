// Copyright 2019 Tim Shannon. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package generichold_test

import (
	"github.com/rlshukhov/generichold"
	"testing"

	"github.com/timshannon/badgerhold/v4"
)

var sortTests = []test{
	{
		name:   "Sort By Name",
		query:  badgerhold.Where("Category").Eq("animal").SortBy("Name"),
		result: []int{9, 5, 14, 8, 13, 2, 16},
	},
	{
		name:   "Sort By Name Reversed",
		query:  badgerhold.Where("Category").Eq("animal").SortBy("Name").Reverse(),
		result: []int{16, 2, 13, 8, 14, 5, 9},
	},
	{
		name:   "Sort By Multiple Fields",
		query:  badgerhold.Where("ID").In(8, 3, 13).SortBy("Category", "Name"),
		result: []int{13, 15, 4, 3},
	},
	{
		name:   "Sort By Multiple Fields Reversed",
		query:  badgerhold.Where("ID").In(8, 3, 13).SortBy("Category", "Name").Reverse(),
		result: []int{3, 4, 15, 13},
	},
	{
		name:   "Sort By Duplicate Field Names",
		query:  badgerhold.Where("ID").In(8, 3, 13).SortBy("Category", "Name", "Category"),
		result: []int{13, 15, 4, 3},
	},
	{
		name:   "Sort By Name with limit",
		query:  badgerhold.Where("Category").Eq("animal").SortBy("Name").Limit(3),
		result: []int{9, 5, 14},
	},
	{
		name:   "Sort By Name with skip",
		query:  badgerhold.Where("Category").Eq("animal").SortBy("Name").Skip(3),
		result: []int{8, 13, 2, 16},
	},
	{
		name:   "Sort By Name with skip and limit",
		query:  badgerhold.Where("Category").Eq("animal").SortBy("Name").Skip(2).Limit(3),
		result: []int{14, 8, 13},
	},
	{
		name:   "Sort By Name Reversed with limit",
		query:  badgerhold.Where("Category").Eq("animal").SortBy("Name").Skip(2).Limit(3),
		result: []int{14, 8, 13},
	},
	{
		name:   "Sort By Name Reversed with skip",
		query:  badgerhold.Where("Category").Eq("animal").SortBy("Name").Skip(4),
		result: []int{13, 2, 16},
	},
	{
		name:   "Sort By Name Reversed with skip and limit",
		query:  badgerhold.Where("Category").Eq("animal").SortBy("Name").Skip(2).Limit(3),
		result: []int{14, 8, 13},
	},
	{
		name:   "Sort By Name with skip greater than length",
		query:  badgerhold.Where("Category").Eq("animal").SortBy("Name").Skip(10),
		result: []int{},
	},
}

func TestSortedFind(t *testing.T) {
	testWrap(t, func(bh *badgerhold.Store, t *testing.T) {
		store := generichold.Open[ItemTest](bh)

		insertTestData(t, store)

		for _, tst := range sortTests {
			t.Run(tst.name, func(t *testing.T) {
				result, err := store.Find(tst.query)
				if err != nil {
					t.Fatalf("Error finding sort data from badgerhold: %s", err)
				}
				if len(result) != len(tst.result) {
					if testing.Verbose() {
						t.Fatalf("Sorted Find result count is %d wanted %d.  Results: %v", len(result),
							len(tst.result), result)
					}
					t.Fatalf("Sorted Find result count is %d wanted %d.", len(result), len(tst.result))
				}

				for i := range result {
					if !result[i].equal(&testData[tst.result[i]]) {
						if testing.Verbose() {
							t.Fatalf("Expected index %d to be %v, Got %v Results: %v", i, &testData[tst.result[i]],
								result[i], result)
						}
						t.Fatalf("Expected index %d to be %v, Got %v", i, &testData[tst.result[i]], result[i])
					}
				}
			})
		}
	})
}

func TestSortedUpdateMatching(t *testing.T) {
	for _, tst := range sortTests {
		t.Run(tst.name, func(t *testing.T) {
			testWrap(t, func(bh *badgerhold.Store, t *testing.T) {
				store := generichold.Open[ItemTest](bh)

				insertTestData(t, store)

				err := store.UpdateMatching(tst.query, func(record *ItemTest) error {
					update := record

					update.UpdateField = "updated"
					update.UpdateIndex = "updated index"

					return nil
				})

				if err != nil {
					t.Fatalf("Error updating data from badgerhold: %s", err)
				}

				result, err := store.Find(badgerhold.Where("UpdateIndex").Eq("updated index").And("UpdateField").Eq("updated"))
				if err != nil {
					t.Fatalf("Error finding result after update from badgerhold: %s", err)
				}

				if len(result) != len(tst.result) {
					if testing.Verbose() {
						t.Fatalf("Find result count after update is %d wanted %d.  Results: %v",
							len(result), len(tst.result), result)
					}
					t.Fatalf("Find result count after update is %d wanted %d.", len(result),
						len(tst.result))
				}

				for i := range result {
					found := false
					for k := range tst.result {
						if result[i].Key == testData[tst.result[k]].Key &&
							result[i].UpdateField == "updated" &&
							result[i].UpdateIndex == "updated index" {
							found = true
							break
						}
					}

					if !found {
						if testing.Verbose() {
							t.Fatalf("Could not find %v in the update result set! Full results: %v",
								result[i], result)
						}
						t.Fatalf("Could not find %v in the updated result set!", result[i])
					}
				}

			})

		})
	}
}

func TestSortedDeleteMatching(t *testing.T) {
	for _, tst := range sortTests {
		t.Run(tst.name, func(t *testing.T) {
			testWrap(t, func(bh *badgerhold.Store, t *testing.T) {
				store := generichold.Open[ItemTest](bh)

				insertTestData(t, store)

				err := store.DeleteMatching(tst.query)
				if err != nil {
					t.Fatalf("Error deleting data from badgerhold: %s", err)
				}

				result, err := store.Find(nil)
				if err != nil {
					t.Fatalf("Error finding result after delete from badgerhold: %s", err)
				}

				if len(result) != (len(testData) - len(tst.result)) {
					if testing.Verbose() {
						t.Fatalf("Delete result count is %d wanted %d.  Results: %v", len(result),
							(len(testData) - len(tst.result)), result)
					}
					t.Fatalf("Delete result count is %d wanted %d.", len(result),
						(len(testData) - len(tst.result)))

				}

				for i := range result {
					found := false
					for k := range tst.result {
						if result[i].equal(&testData[tst.result[k]]) {
							found = true
							break
						}
					}

					if found {
						if testing.Verbose() {
							t.Fatalf("Found %v in the result set when it should've been deleted! Full results: %v", result[i], result)
						}
						t.Fatalf("Found %v in the result set when it should've been deleted!", result[i])
					}
				}

			})

		})
	}
}

func TestSortOnKey(t *testing.T) {
	testWrap(t, func(bh *badgerhold.Store, t *testing.T) {
		store := generichold.Open[ItemTest](bh)

		defer func() {
			if r := recover(); r == nil {
				t.Fatalf("Running Sort on Key field did not panic!")
			}
		}()

		_, _ = store.Find(badgerhold.Where("Name").Eq("blah").SortBy(badgerhold.Key))
	})
}

func TestSortedFindOnInvalidFieldName(t *testing.T) {
	testWrap(t, func(bh *badgerhold.Store, t *testing.T) {
		store := generichold.Open[ItemTest](bh)

		insertTestData(t, store)

		_, err := store.Find(badgerhold.Where("BadFieldName").Eq("test").SortBy("BadFieldName"))
		if err == nil {
			t.Fatalf("Sorted find query against a bad field name didn't return an error!")
		}
	})
}

// NOTE: unnecessary due to compile-time type assertion
//func TestSortedFindWithNonSlicePtr(t *testing.T) {
//	testWrap(t, func(bh *badgerhold.Store, t *testing.T) {
//		store := generichold.Open[ItemTest](bh)
//
//		defer func() {
//			if r := recover(); r == nil {
//				t.Fatalf("Running Find with non-slice pointer did not panic!")
//			}
//		}()
//		_, _ = store.Find(badgerhold.Where("Name").Eq("blah").SortBy("Name"))
//	})
//}
