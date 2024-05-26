package database

import (
	"errors"

	"github.com/google/uuid"
)

type Value map[string]any
type Key uuid.UUID

type IDataBase interface {
	Len() int
	Set(key Key, value Value) error
	Get(key Key) (Value, error)
	Delete(key Key) error
}

type Storage map[Key]Value
type database struct {
	storage Storage
}

func New() IDataBase {
	return &database{
		storage: make(Storage),
	}
}

func (d *database) Len() int {
	return len(d.storage)
}

func (d *database) Set(key Key, value Value) error {
	d.storage[key] = value
	return nil
}

func (d *database) Get(key Key) (Value, error) {
	if data, ok := d.storage[key]; ok {
		return data, nil
	}
	return nil, errors.New("key not found")
}

func (d *database) Delete(key Key) error {
	if _, ok := d.storage[key]; ok {
		delete(d.storage, key)
		return nil
	}

	return errors.New("key not found")
}
