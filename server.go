package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	store "github.com/nplatte/DnDSpellAPI/stores"
)

var config_path string = "/home/nplatte/Desktop/DnDSpellAPI/stores/configs/my_config.json"

type SpellServer struct {
	db store.SpellStore
}

// loads an empty database with spell information
func InitializeServer() *SpellServer {
	server := StartServer()
	err := server.db.LoadStore()
	if err != nil {
		panic(err)
	}
	return server
}

// returns a Spell server to use
func StartServer() *SpellServer {
	server := SpellServer{}
	server.db = store.MakeInDatabaseStore(config_path)
	return &server
}

func (server *SpellServer) RegisterHandlers() error {
	e := echo.New()
	e.GET("/spells/:name", server.HandleGETSpellName)
	err := e.Start(":8080")
	if err != nil {
		panic(err)
	}
	return nil
}

func (server *SpellServer) HandleGETSpellName(ctx echo.Context) error {
	var spell store.Spell
	spell, err := server.db.GetSpell(ctx.Param("name"))
	if err != nil {
		return ctx.JSON(http.StatusNotFound, err)
	}
	return ctx.JSON(http.StatusOK, spell)
}

// I want a func that initializes the server
// This would start the webserver, connect to the DB, then load the DB with the json files
// I also want a func that starts the webserver
// This would do the same as above, but not load any new data into the database

func main() {
	var server *SpellServer
	// this is for starting an already existing server
	//server = StartServer()
	// this is for filling an existing database
	server = InitializeServer()
	server.RegisterHandlers()
}
