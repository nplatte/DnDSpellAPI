package store

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

func MakeInDatabaseStore(config_path string) *InDatabaseStore {
	store := &InDatabaseStore{}
	store.GetDBConn(config_path)
	err := store.CreateTables()
	if err != nil {
		panic(err)
	}
	return store
}

type InDatabaseStore struct {
	db     *sql.DB
	Config ConfigData
}

func (store *InDatabaseStore) GetSpell(name string) (Spell, error) {
	var spell Spell
	query := fmt.Sprintf("SELECT * FROM spells WHERE Name='%s'", name)
	row := store.db.QueryRow(query)
	var rawClassesJSON []byte
	var rawComponentsJSON []byte
	err := row.Scan(&spell.Name, &spell.Range, &spell.Level, &spell.CastTime, &spell.Description, &spell.Duration, &spell.Concentration, &spell.Ritual, &rawClassesJSON, &spell.Source, &spell.School, &rawComponentsJSON, &spell.HigherLevels)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return Spell{}, errors.New("Spell not in database")
		}
	}
	json.Unmarshal(rawClassesJSON, &spell.ClassList)
	json.Unmarshal(rawComponentsJSON, &spell.Components)
	return spell, err
}

func (store *InDatabaseStore) GetDBConn(config_path string) (bool, error) {
	raw, err := os.ReadFile(config_path)
	if err != nil {
		panic(err)
	}
	json.Unmarshal(raw, &store.Config)
	store.db, err = sql.Open("mysql", store.Config.DBinfo.FormatDSN())
	return err == nil, err
}

func (store *InDatabaseStore) CreateTables() error {
	var err error
	drop := "DROP TABLE IF EXISTS spells;"
	_, err = store.db.Exec(drop)
	if err != nil {
		return err
	}
	create := "CREATE TABLE spells (Name VARCHAR(50), SpellRange VARCHAR(50), Level VARCHAR(30), CastTime VARCHAR(50), Description VARCHAR(2000), Duration VARCHAR(50), Concentration VARCHAR(10), Ritual VARCHAR(10), ClassList JSON, Source VARCHAR(50), School VARCHAR(50), Components JSON, HigherLevels VARCHAR(500));"
	_, err = store.db.Exec(create)
	return err
}

// Goes through each module in the modules and loads it
func (store *InDatabaseStore) LoadStore() error {
	for _, module := range store.Config.Modules {
		err := store.LoadModule(module)
		if err != nil {
			return err
		}
	}
	return nil
}

func (store *InDatabaseStore) LoadModule(module string) error {
	var err error = nil
	for _, school := range store.Config.Files {
		err = store.LoadSchool(school, module)
	}
	return err
}

func (store *InDatabaseStore) LoadSchool(school string, module string) error {
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

func (store *InDatabaseStore) LoadSpell(spell Spell) error {
	classJSON, _ := json.Marshal(spell.ClassList)
	componentsJSON, _ := json.Marshal(spell.Components)
	insert := "INSERT INTO spells (Name, SpellRange, Level, CastTime, Description, Duration, Concentration, Ritual, ClassList, Source, School, Components, HigherLevels) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);"
	_, err := store.db.Exec(insert, spell.Name, spell.Range, spell.Level, spell.CastTime, spell.Description, spell.Duration, spell.Concentration, spell.Ritual, classJSON, spell.Source, spell.School, componentsJSON, spell.HigherLevels)
	return err
}
