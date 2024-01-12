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
	m.rwMutex.Lock()
	defer m.rwMutex.Unlock()
	m.heroes = database.TestHeroes
	m.SortHeroes()
	return nil
}

func (m *Memorydb) Disconnect() error {
	return nil
}

func (m *Memorydb) GetById(id int) (*hero.Hero, error) {
	m.rwMutex.Lock()
	defer m.rwMutex.Unlock()

	for _, he := range m.heroes {
		if he.Id == id {
			return &he, nil
		}
	}
	return nil, errors.New("hero with id not found")
}

func (m *Memorydb) GetByName(name string) ([]hero.Hero, error) {
	m.rwMutex.Lock()
	defer m.rwMutex.Unlock()

	hs := make([]hero.Hero, 0)
	for _, he := range m.heroes {
		if strings.Contains(strings.ToLower(he.Name), strings.ToLower(name)) {
			hs = append(hs, he)
		}
	}
	return hs, nil
}

func (m *Memorydb) GetAll() ([]hero.Hero, error) {
	m.rwMutex.Lock()
	defer m.rwMutex.Unlock()
	return m.heroes, nil
}

func (m *Memorydb) UpdateHero(h hero.Hero) error {
	m.rwMutex.Lock()
	defer m.rwMutex.Unlock()
	for i, he := range m.heroes {
		if he.Id == h.Id {
			m.heroes[i] = h
			return nil
		}
	}
	return errors.New("Hero with id " + fmt.Sprintf("%v", h.Id) + "not found")
}

func (m *Memorydb) DeleteHero(id int) error {
	m.rwMutex.Lock()
	defer m.rwMutex.Unlock()
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
	m.rwMutex.Lock()
	defer m.rwMutex.Unlock()

	for _, he := range m.heroes {
		if he.Id == h.Id {
			return errors.New("hero with id already exists")
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
