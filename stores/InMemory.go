package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

func MakeInMemoryStore() *InMemoryStore {
	store := &InMemoryStore{}
	data, _ := os.ReadFile("/home/nplatte/Desktop/DnDSpellAPI/stores/configs/test_config.json")
	json.Unmarshal(data, &store.Config)
	return store
}

type InMemoryStore struct {
	Spells []Spell
	Config ConfigData
}

func (store InMemoryStore) GetSpell(name string) (Spell, error) {
	for _, spell := range store.Spells {
		if spell.Name == name {
			return spell, nil
		}
	}
	return Spell{}, errors.New("Spell not in database")
}

// from the JSON config file, loads the list of strings into LoadedContent, then loads spells from each module provided
func (store *InMemoryStore) LoadStore() error {
	// this should take the config file and load each file specified
	var err error = nil
	for _, module := range store.Config.Modules {
		err = store.LoadModule(module)
	}
	return err
}

// When given a valid folder name, loads each file in folder
func (store *InMemoryStore) LoadModule(module string) error {
	var err error = nil
	for _, school := range store.Config.Files {
		err = store.LoadSchool(school, module)
	}
	return err
}

func (store *InMemoryStore) LoadSchool(school string, module string) error {
	schoolToLoad := fmt.Sprintf("/home/nplatte/Desktop/DnDSpellAPI/stores/spells/%s/%s.json", module, school)
	var spells []Spell
	rawData, err := os.ReadFile(schoolToLoad)
	if err != nil {
		return err
	}
	json.Unmarshal(rawData, &spells)
	for _, spell := range spells {
		spell.Source = module
		spell.School = school
		err = store.LoadSpell(spell)
		if err != nil {
			return err
		}
	}
	return nil
}

func (store *InMemoryStore) LoadSpell(spell Spell) error {
	store.Spells = append(store.Spells, spell)
	return nil
}
