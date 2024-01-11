package database

import "gorestapi/internal/hero"

type Params interface{}

type Herodb interface {
	Connect(p Params) error
	Disconnect() error
	GetById(id int) (*hero.Hero, error)
	GetAll() ([]hero.Hero, error)
	GetByName(name string) ([]hero.Hero, error)
	UpdateHero(h hero.Hero) error
	DeleteHero(id int) error
	AddHero(h hero.Hero) error
}

var TestHeroes = []hero.Hero{
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
