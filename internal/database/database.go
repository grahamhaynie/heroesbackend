package database

import (
	"context"
	"gorestapi/internal/hero"
)

type Params interface{}

type Herodb interface {
	Connect(ctx context.Context, p Params) error
	Disconnect(ctx context.Context) error
	GetById(ctx context.Context, id int) (*hero.Hero, error)
	GetAll(ctx context.Context) ([]hero.Hero, error)
	GetByName(ctx context.Context, name string) ([]hero.Hero, error)
	UpdateHero(ctx context.Context, h hero.Hero) error
	DeleteHero(ctx context.Context, id int) error
	AddHero(ctx context.Context, h hero.Hero) error
}

// sample heroes to populate database with
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
