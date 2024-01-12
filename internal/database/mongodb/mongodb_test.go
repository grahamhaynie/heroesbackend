package mongodb_test

import (
	"gorestapi/internal/database"
	"gorestapi/internal/database/mongodb"
	"testing"
)

func TestMongoDbImplementsHerodb(t *testing.T) {
	var m interface{} = &mongodb.Mongodb{}
	if _, ok := m.(database.Herodb); !ok {
		t.Fatalf("mongodb does not implement herodb interface")
	}
}

// TODO mock mongodb
// more of an integration test but mock or connect to real one and run unit tests?
