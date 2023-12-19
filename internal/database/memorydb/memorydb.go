package memorydb

import (
	"errors"
	"fmt"
	"gorestapi/internal/database"
	"gorestapi/internal/hero"
	"sort"
	"strings"
	"sync"
)

type Memorydb struct {
	heroes  []hero.Hero
	rwMutex sync.RWMutex
}

func (m *Memorydb) Connect(p database.Params) error {
	m.heroes = []hero.Hero{
		{Id: 12, Name: "Dr. Nice", Power: "bein nice", AlterEgo: "nobody", PhotoURL: "http://localhost:8080/photo/minion.jpg"},
		{Id: 13, Name: "Bombasto", Power: "throwing stuf"},
		{Id: 14, Name: "Celeritas", Power: "celebrity", AlterEgo: "tom cruise"},
		{Id: 15, Name: "Magneta", Power: "not sure tbh"},
		{Id: 16, Name: "RubberMan", Power: "elastic arms", AlterEgo: "steve"},
		{Id: 17, Name: "Dynama", Power: "dynamite"},
		{Id: 18, Name: "Dr. IQ", Power: "talking", AlterEgo: "michael"},
		{Id: 19, Name: "Magma", Power: "making rocks"},
		{Id: 20, Name: "Tornado", Power: "spinning in an office chair", AlterEgo: "nick"},
	}

	m.SortHeroes()
	return nil
}

func (m *Memorydb) Disconnect() error {
	return nil
}

func (m *Memorydb) GetById(id int) (*hero.Hero, error) {
	for _, he := range m.heroes {
		if he.Id == id {
			return &he, nil
		}
	}
	return nil, nil
}

func (m *Memorydb) GetByName(name string) ([]hero.Hero, error) {
	hs := make([]hero.Hero, 0)
	for _, he := range m.heroes {
		if strings.Contains(strings.ToLower(he.Name), strings.ToLower(name)) {
			hs = append(hs, he)
		}
	}
	return hs, nil
}

func (m *Memorydb) GetAll() ([]hero.Hero, error) {
	return m.heroes, nil
}

func (m *Memorydb) UpdateHero(h hero.Hero) error {
	for i, he := range m.heroes {
		if he.Id == h.Id {
			m.heroes[i] = h
			return nil
		}
	}
	return errors.New("Hero with id " + fmt.Sprintf("%v", h.Id) + "not found")
}

func (m *Memorydb) DeleteHero(id int) error {
	del := -1
	for i, he := range m.heroes {
		if he.Id == id {
			del = i
			break
		}
	}
	if del == -1 {
		return errors.New("Could not delete hero with id " + fmt.Sprintf("%v", id))
	}
	m.heroes = append(m.heroes[:del], m.heroes[del+1:]...)
	m.SortHeroes()
	return nil
}

func (m *Memorydb) AddHero(h hero.Hero) error {
	// make sure id is not duplicate
	// works because heroes is sorted
	for _, he := range m.heroes {
		if he.Id == h.Id {
			h.Id++
		}
	}
	m.heroes = append(m.heroes, h)
	m.SortHeroes()
	return nil
}

func (m *Memorydb) SortHeroes() {
	// sort heroes
	sort.Slice(m.heroes, func(i int, j int) bool {
		return m.heroes[i].Id < m.heroes[j].Id
	})
}
