package memorydb_test

import (
	"gorestapi/internal/database"
	"gorestapi/internal/database/memorydb"
	"testing"
)

func TestMemoryDbImplementsHerodb(t *testing.T) {
	var _ database.Herodb = &memorydb.Memorydb{}
}
