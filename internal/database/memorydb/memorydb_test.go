package memorydb_test

import (
	"gorestapi/internal/database"
	"gorestapi/internal/database/memorydb"
	"gorestapi/internal/hero"
	"reflect"
	"testing"
)

func TestMemoryDbImplementsHerodb(t *testing.T) {
	var m interface{} = &memorydb.Memorydb{}
	if _, ok := m.(database.Herodb); !ok {
		t.Fatalf("memorydb does not implement herodb interface")
	}
	//var _ database.Herodb = &memorydb.Memorydb{}
}

// test that getting by a valid id works, as well as getting a non existing id
func TestGetById(t *testing.T) {
	// Initialize with test data
	db := memorydb.Memorydb{}
	db.Connect(nil)

	// Test for an existing hero
	existingID := 12
	hero, err := db.GetById(existingID)
	if err != nil {
		t.Errorf("GetById() with existing ID returned error: %v", err)
	}
	if hero == nil || hero.Id != existingID {
		t.Errorf("GetById() with existing ID returned incorrect hero: %+v", hero)
	}

	// Test for a non-existing hero
	nonExistingID := 999
	_, err = db.GetById(nonExistingID)
	if err == nil {
		t.Error("GetById() with non-existing ID did not return an error")
	}
}

// test that getting all items from database returns correct length
// as well as correct items
func TestGetAll(t *testing.T) {
	// initialize with test data
	db := memorydb.Memorydb{}
	db.Connect(nil)

	got, err := db.GetAll()
	if err != nil {
		t.Errorf("GetAll() error = %v, wantErr nil", err)
	}
	if len(got) != len(database.TestHeroes) {
		t.Errorf("GetAll() got %v items, want %v items", len(got), len(database.TestHeroes))
	}

	// assumes TestHeroes are sorted by ID (as Memorydb's Connect method sorts them)
	for i, hero := range got {
		if !reflect.DeepEqual(hero, database.TestHeroes[i]) {
			t.Errorf("GetAll() hero at index %d = %v, want %v", i, hero, database.TestHeroes[i])
		}
	}
}

// test that getting an existing hero works, as well as a non existing hero
func TestGetByName(t *testing.T) {
	// Initialize with test data
	db := memorydb.Memorydb{}
	db.Connect(nil)

	tests := []struct {
		name     string
		expected []hero.Hero
		wantErr  bool
	}{
		{"Nice", []hero.Hero{database.TestHeroes[0]}, false},
		{"nonexisting", []hero.Hero{}, false},
	}

	for _, tt := range tests {
		got, err := db.GetByName(tt.name)
		if (err != nil) != tt.wantErr {
			t.Errorf("GetByName(%s) error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if !reflect.DeepEqual(got, tt.expected) {
			t.Errorf("GetByName(%s) = %v, want %v", tt.name, got, tt.expected)
		}
	}
}

// test updating a hero that exists works, as well as updating a non-existing one
func TestUpdateHero(t *testing.T) {
	// Initialize with test data
	db := memorydb.Memorydb{}
	db.Connect(nil)

	// wantErr indicates that an error should occur in updating hero
	tests := []struct {
		hero    hero.Hero
		wantErr bool
	}{
		{hero.Hero{Id: 12, Name: "Updated Name", Power: "Updated Power"}, false},
		{hero.Hero{Id: 999, Name: "Nonexisting Hero"}, true},
	}

	for _, tt := range tests {
		err := db.UpdateHero(tt.hero)
		if (err != nil) != tt.wantErr {
			t.Errorf("UpdateHero() error = %v, wantErr %v", err, tt.wantErr)
		}

		// Verify update
		if !tt.wantErr {
			updatedHero, _ := db.GetById(tt.hero.Id)
			if !reflect.DeepEqual(updatedHero, &tt.hero) {
				t.Errorf("UpdateHero() failed to update hero, got %v, want %v", updatedHero, &tt.hero)
			}
		}
	}
}

// check that deleting a hero that exists works, as well as a non-existing one
func TestDeleteHero(t *testing.T) {
	// Initialize with test data
	db := memorydb.Memorydb{}
	db.Connect(nil)

	// wantErr indicates that an error should occur in deleting a hero
	tests := []struct {
		id      int
		wantErr bool
	}{
		{12, false},
		{999, true},
	}

	for _, tt := range tests {
		err := db.DeleteHero(tt.id)
		if (err != nil) != tt.wantErr {
			t.Errorf("DeleteHero() error = %v, wantErr %v", err, tt.wantErr)
		}

		if !tt.wantErr {
			// Verify deletion
			_, err := db.GetById(tt.id)
			if err == nil {
				t.Errorf("DeleteHero() failed to delete hero with ID %d", tt.id)
			}
		}
	}
}

// test that adding a new hero works, as well as a hero with an overlapping id
func TestAddHero(t *testing.T) {
	// Initialize with test data
	db := memorydb.Memorydb{}
	db.Connect(nil)

	// Test adding a new hero
	newHero := hero.Hero{Id: 21, Name: "New Hero", Power: "New Power"}
	err := db.AddHero(newHero)
	if err != nil {
		t.Errorf("AddHero() new hero error = %v, wantErr nil", err)
	}

	// Verify the new hero was added
	addedHero, err := db.GetById(21)
	if err != nil || !reflect.DeepEqual(addedHero, &newHero) {
		t.Errorf("AddHero() failed to add new hero, got %v, want %v", addedHero, newHero)
	}

	// Test adding a hero with an existing ID
	duplicateHero := hero.Hero{Id: 12, Name: "Duplicate Hero", Power: "Duplicate Power"}
	err = db.AddHero(duplicateHero)
	if err != nil {
		t.Errorf("AddHero() should return error on duplicate ID, got nil")
	}
}
