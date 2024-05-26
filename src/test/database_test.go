package test_database

import (
	"team01/database"
	"testing"

	"github.com/google/uuid"
)

func TestDatabaseSet(t *testing.T) {
	db := database.New()

	key_uuid, _ := uuid.NewUUID()
	key := database.Key(key_uuid)
	value := database.Value{
		"key": "value",
	}

	if err := db.Set(key, value); err != nil {
		t.Error(err)
	}
}

func TestDatabaseLenAfterSet(t *testing.T) {
	db := database.New()

	key_uuid, _ := uuid.NewUUID()
	key := database.Key(key_uuid)
	value := database.Value{
		"key": "value",
	}

	db.Set(key, value)

	if len := db.Len(); len != 1 {
		t.Errorf("Expected 1, got %d", len)
	}
}

func TestDatabaseGet(t *testing.T) {
	db := database.New()

	key_uuid, _ := uuid.NewUUID()
	key := database.Key(key_uuid)
	value := database.Value{
		"key": "value",
	}

	db.Set(key, value)

	db_data, err := db.Get(key)
	if err != nil {
		t.Error(err)
	}

	if _, ok := db_data["key"].(string); !ok {
		t.Errorf("Expected 'value', got %s", db_data["key"])
	}

	if db_data["key"].(string) != "value" {
		t.Errorf("Expected 'value', got %s", db_data["key"])
	}
}

func TestDatabaseDelete(t *testing.T) {
	db := database.New()

	key_uuid, _ := uuid.NewUUID()
	key := database.Key(key_uuid)
	value := database.Value{
		"key": "value",
	}

	db.Set(key, value)
	if err := db.Delete(key); err != nil {
		t.Error(err)
	}

	if len := db.Len(); len != 0 {
		t.Errorf("Expected 0, got %d", len)
	}
}
