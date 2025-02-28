package store

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

var config_path string = "/home/nplatte/Desktop/DnDSpellAPI/stores/configs/test_config.json"

var firebolt Spell = Spell{
	Name:          "Fire Bolt",
	Range:         "120 feet",
	Level:         "cantrip",
	CastTime:      "1 action",
	Description:   "You hurl a mote of fire at a creature or object within range. Make a ranged spell attack against the target. On a hit, the target takes 1d10 fire damage. A flammable object hit by this spell ignites if it isn't being worn or carried.",
	Duration:      "Instantaneous",
	Concentration: "False",
	Ritual:        "False",
	ClassList:     []string{"Artificer", "Sorcerer", "Wizard"},
	Source:        "PlayersHB",
	School:        "Evocation",
	Components:    []string{"V", "S"},
	HigherLevels:  "This spell’s damage increases by 1d10 when you reach 5th level (2d10), 11th level (3d10), and 17th level (4d10).",
}

func TestMakeInMemoryStore(t *testing.T) {
	// checks to make sure returns correct type
	t.Run("Assert Store is correct Types", func(t *testing.T) {
		assert.IsType(t, &InMemoryStore{}, MakeInMemoryStore())
		assert.IsType(t, &InDatabaseStore{}, MakeInDatabaseStore(config_path))
	})
}

func TestDatabaseGetConn(t *testing.T) {
	store := &InDatabaseStore{}
	t.Run("test Gen Conn", func(t *testing.T) {
		connected, err := store.GetDBConn(config_path)
		assert.Equal(t, nil, err)
		assert.True(t, connected)
	})
}

func assertStoreGETSpell(t *testing.T, store SpellStore, want Spell) {
	t.Helper()
	t.Run("Can Get Spell", func(t *testing.T) {
		got, err := store.GetSpell(want.Name)
		assert.Equal(t, nil, err)
		assert.Equal(t, want, got)
	})
}

func assertSpellNotFound(t *testing.T, store SpellStore) {
	t.Helper()
	t.Run("Assert Spell Not Found", func(t *testing.T) {
		want := errors.New("Spell not in database")
		spell, got := store.GetSpell("Ray Of Frost")
		assert.Equal(t, spell, Spell{})
		assert.Equal(t, want, got)
	})
}

func TestMemoryGetSpell(t *testing.T) {
	store := MakeInMemoryStore()
	store.Spells = []Spell{firebolt}
	assertStoreGETSpell(t, store, firebolt)
	assertSpellNotFound(t, store)
}

func TestDatabaseGetSpell(t *testing.T) {
	store := MakeInDatabaseStore(config_path)
	assertSpellNotFound(t, store)
	err := loadSpellIntoSQLDB(firebolt, store)
	if err != nil {
		panic(err)
	}
	assertStoreGETSpell(t, store, firebolt)

}

func assertStoreLoads(t *testing.T, store SpellStore, name string) {
	t.Run(name, func(t *testing.T) {
		err := store.LoadStore()
		assert.Equal(t, nil, err)
		spell, err := store.GetSpell(firebolt.Name)
		// assert Spell in database
		assert.Equal(t, nil, err)
		assert.Equal(t, firebolt, spell)
	})
}

func TestLoadStore(t *testing.T) {
	assertStoreLoads(t, MakeInMemoryStore(), "In Memory Loads Store")
	assertStoreLoads(t, MakeInDatabaseStore(config_path), "In Database Loads Store")
}

func assertLoadModule(t *testing.T, store SpellStore, name string) {
	folder := "PlayersHB"
	t.Run(name, func(t *testing.T) {
		err := store.LoadModule(folder)
		assert.Equal(t, nil, err)
		spell, err := store.GetSpell(firebolt.Name)
		// assert Spell in database
		assert.Equal(t, nil, err)
		assert.Equal(t, firebolt, spell)
	})
}

func TestMemoryLoadModule(t *testing.T) {
	assertLoadModule(t, MakeInMemoryStore(), "In Memory Load Module")
	assertLoadModule(t, MakeInDatabaseStore(config_path), "In Database Load Module")
}

func assertLoadSchool(t *testing.T, store SpellStore, name string) {
	t.Run(name, func(t *testing.T) {
		err := store.LoadSchool("Evocation", "PlayersHB")
		assert.Equal(t, nil, err)
		spell, err := store.GetSpell(firebolt.Name)
		// assert Spell in database
		assert.Equal(t, nil, err)
		assert.Equal(t, firebolt, spell)
	})
}

func TestLoadSchool(t *testing.T) {
	assertLoadSchool(t, MakeInMemoryStore(), "Memory Load School")
	assertLoadSchool(t, MakeInDatabaseStore(config_path), "Database Load School")
}

func TestMemoryLoadSpell(t *testing.T) {
	store := MakeInMemoryStore()
	t.Run("Test Load Spell", func(t *testing.T) {
		assert.ElementsMatch(t, []Spell{}, store.Spells)
		err := store.LoadSpell(firebolt)
		assert.Equal(t, nil, err)
		spell, err := store.GetSpell(firebolt.Name)
		// assert Spell in database
		assert.Equal(t, nil, err)
		assert.Equal(t, firebolt, spell)
		assert.ElementsMatch(t, []Spell{firebolt}, store.Spells)
	})
}

func TestDatabaseLoadSpell(t *testing.T) {
	store := MakeInDatabaseStore(config_path)
	t.Run("Test Load Spell", func(t *testing.T) {
		// assert the DB is empty
		var spell Spell
		var rawClassesJSON []byte
		var rawComponentsJSON []byte
		query := "SELECT * FROM spells"
		row := store.db.QueryRow(query)
		err := row.Scan(&spell.Name, &spell.Range, &spell.Level, &spell.CastTime, &spell.Description, &spell.Duration, &spell.Concentration, &spell.Ritual, &rawClassesJSON, &spell.Source, &spell.School, &rawComponentsJSON, &spell.HigherLevels)
		assert.Equal(t, err.Error(), "sql: no rows in result set")
		// load the spell into the db
		err = store.LoadSpell(firebolt)
		if err != nil {
			panic(err)
		}
		// get the spell
		spell, err = store.GetSpell(firebolt.Name)
		if err != nil {
			panic(err)
		}
		assert.Equal(t, spell, firebolt)
	})
}

func loadSpellIntoSQLDB(spell Spell, store *InDatabaseStore) error {
	classJSON, _ := json.Marshal(spell.ClassList)
	componentsJSON, _ := json.Marshal(spell.Components)
	insert := "INSERT INTO spells (Name, SpellRange, Level, CastTime, Description, Duration, Concentration, Ritual, ClassList, Source, School, Components, HigherLevels) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);"
	_, err := store.db.Exec(insert, spell.Name, spell.Range, spell.Level, spell.CastTime, spell.Description, spell.Duration, spell.Concentration, spell.Ritual, classJSON, spell.Source, spell.School, componentsJSON, spell.HigherLevels)
	return err
}
