package memorydb_test

import (
	"gorestapi/internal/database"
	"gorestapi/internal/database/memorydb"
	"testing"
)

func TestMemoryDbImplementsHerodb(t *testing.T) {
	var m interface{} = &memorydb.Memorydb{}
	if _, ok := m.(database.Herodb); !ok {
		t.Fatalf("memorydb does not implement herodb interface")
	}
	//var _ database.Herodb = &memorydb.Memorydb{}
}
