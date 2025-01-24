// Copyright 2025 Lane Shukhov. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package generichold

import (
	"github.com/dgraph-io/badger/v4"
	"github.com/timshannon/badgerhold/v4"
)

type store[T any] struct {
	store *badgerhold.Store
}

type Store[T any] interface {
	FindAggregate(query *badgerhold.Query, groupBy ...string) ([]*badgerhold.AggregateResult, error)
	TxFindAggregate(tx *badger.Txn, query *badgerhold.Query, groupBy ...string) ([]*badgerhold.AggregateResult, error)
	Delete(key any) error
	DeleteMatching(query *badgerhold.Query) error
	TxDelete(tx *badger.Txn, key any) error
	TxDeleteMatching(tx *badger.Txn, query *badgerhold.Query) error
	Count(query *badgerhold.Query) (uint64, error)
	Find(query *badgerhold.Query) ([]T, error)
	FindOne(query *badgerhold.Query) (T, error)
	ForEach(query *badgerhold.Query, fn any) error
	Get(key any) (T, error)
	TxCount(tx *badger.Txn, query *badgerhold.Query) (uint64, error)
	TxFind(tx *badger.Txn, query *badgerhold.Query) ([]T, error)
	TxFindOne(tx *badger.Txn, query *badgerhold.Query) (T, error)
	TxForEach(tx *badger.Txn, query *badgerhold.Query, fn any) error
	TxGet(tx *badger.Txn, key any) (T, error)
	Insert(key any, data *T) error
	TxInsert(tx *badger.Txn, key any, data *T) error
	TxUpdate(tx *badger.Txn, key any, data *T) error
	TxUpdateMatching(tx *badger.Txn, query *badgerhold.Query, update func(record *T) error) error
	TxUpsert(tx *badger.Txn, key any, data *T) error
	Update(key any, data *T) error
	UpdateMatching(query *badgerhold.Query, update func(record *T) error) error
	Upsert(key any, data *T) error
	Badger() *badger.DB
	Close() error
}

func Open[T any](s *badgerhold.Store) Store[T] {
	return &store[T]{store: s}
}

func (s *store[T]) FindAggregate(query *badgerhold.Query, groupBy ...string) ([]*badgerhold.AggregateResult, error) {
	return s.store.FindAggregate(s.zeroValue(), query, groupBy...)
}

func (s *store[T]) TxFindAggregate(tx *badger.Txn, query *badgerhold.Query, groupBy ...string) ([]*badgerhold.AggregateResult, error) {
	return s.store.TxFindAggregate(tx, s.zeroValue(), query, groupBy...)
}

func (s *store[T]) Delete(key any) error {
	return s.store.Delete(key, s.zeroValue())
}

func (s *store[T]) DeleteMatching(query *badgerhold.Query) error {
	return s.store.DeleteMatching(s.zeroValue(), query)
}

func (s *store[T]) TxDelete(tx *badger.Txn, key any) error {
	return s.store.TxDelete(tx, key, s.zeroValue())
}

func (s *store[T]) TxDeleteMatching(tx *badger.Txn, query *badgerhold.Query) error {
	return s.store.TxDeleteMatching(tx, s.zeroValue(), query)
}

func (s *store[T]) Count(query *badgerhold.Query) (uint64, error) {
	return s.store.Count(s.zeroValue(), query)
}

func (s *store[T]) Find(query *badgerhold.Query) ([]T, error) {
	var result []T
	// pointer to slice needs to badger limitation: panic: result argument must be a slice address
	err := s.store.Find(&result, query)
	return result, err
}

func (s *store[T]) FindOne(query *badgerhold.Query) (T, error) {
	var result T
	err := s.store.FindOne(&result, query)
	return result, err
}

func (s *store[T]) ForEach(query *badgerhold.Query, fn any) error {
	return s.store.ForEach(query, fn)
}

func (s *store[T]) Get(key any) (T, error) {
	var result T
	err := s.store.Get(key, &result)
	return result, err
}

func (s *store[T]) TxCount(tx *badger.Txn, query *badgerhold.Query) (uint64, error) {
	return s.store.TxCount(tx, s.zeroValue(), query)
}

func (s *store[T]) TxFind(tx *badger.Txn, query *badgerhold.Query) ([]T, error) {
	var result []T
	// pointer to slice needs to badger limitation: panic: result argument must be a slice address
	err := s.store.TxFind(tx, &result, query)
	return result, err
}

func (s *store[T]) TxFindOne(tx *badger.Txn, query *badgerhold.Query) (T, error) {
	var result T
	err := s.store.TxFindOne(tx, &result, query)
	return result, err
}

func (s *store[T]) TxForEach(tx *badger.Txn, query *badgerhold.Query, fn any) error {
	return s.store.TxForEach(tx, query, fn)
}

func (s *store[T]) TxGet(tx *badger.Txn, key any) (T, error) {
	var result T
	err := s.store.TxGet(tx, key, &result)
	return result, err
}

func (s *store[T]) Insert(key any, data *T) error {
	return s.store.Insert(key, data)
}

func (s *store[T]) TxInsert(tx *badger.Txn, key any, data *T) error {
	return s.store.TxInsert(tx, key, data)
}

func (s *store[T]) TxUpdate(tx *badger.Txn, key any, data *T) error {
	return s.store.TxUpdate(tx, key, data)
}

func (s *store[T]) TxUpdateMatching(tx *badger.Txn, query *badgerhold.Query, update func(record *T) error) error {
	anyUpdate := func(record any) error {
		return update(record.(*T))
	}
	return s.store.TxUpdateMatching(tx, s.zeroValue(), query, anyUpdate)
}

func (s *store[T]) TxUpsert(tx *badger.Txn, key any, data *T) error {
	return s.store.TxUpsert(tx, key, data)
}

func (s *store[T]) Update(key any, data *T) error {
	return s.store.Update(key, data)
}

func (s *store[T]) UpdateMatching(query *badgerhold.Query, update func(record *T) error) error {
	anyUpdate := func(record any) error {
		return update(record.(*T))
	}
	return s.store.UpdateMatching(s.zeroValue(), query, anyUpdate)
}

func (s *store[T]) Upsert(key any, data *T) error {
	return s.store.Upsert(key, data)
}

func (s *store[T]) Badger() *badger.DB {
	return s.store.Badger()
}

func (s *store[T]) Close() error {
	return s.store.Close()
}

func (s *store[T]) zeroValue() T {
	var zero T
	return zero
}
