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
