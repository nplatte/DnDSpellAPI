package store

import "github.com/go-sql-driver/mysql"

type SpellStore interface {
	GetSpell(string) (Spell, error)
	LoadStore() error
	LoadModule(string) error
	LoadSchool(string, string) error
}

type Spell struct {
	Name          string
	Range         string
	Level         string
	CastTime      string
	Description   string
	Duration      string
	Concentration string
	Ritual        string
	ClassList     []string
	Source        string
	School        string
	Components    []string
	HigherLevels  string
}

type ConfigData struct {
	Modules []string
	Files   []string
	DBinfo  mysql.Config
}
