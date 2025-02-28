package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testCase struct {
	name      string
	spellName string
	json      string
}

func assertGETSpell(t *testing.T, test testCase) {
	e := echo.New()
	s := InitializeServer()
	response := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/", nil)
	ctx := e.NewContext(request, response)
	ctx.SetPath("/spells/{name}")
	ctx.SetParamNames("name")
	ctx.SetParamValues(test.spellName)
	t.Run(test.name, func(t *testing.T) {
		s.HandleGETSpellName(ctx)
		assert.Equal(t, http.StatusOK, response.Code)
		require.JSONEq(t, test.json, response.Body.String())
	})
}

// Checks that the server will respond with firebolt
func TestGETFirebolt(t *testing.T) {
	test := testCase{
		name:      "Get Firebolt returns Firebolt",
		spellName: "Fire Bolt",
		json: `{
			"Name" : "Fire Bolt",
			"Range" : "120 feet",
			"Level" : "cantrip",
			"CastTime" : "1 action",
			"Description" : "You hurl a mote of fire at a creature or object within range. Make a ranged spell attack against the target. On a hit, the target takes 1d10 fire damage. A flammable object hit by this spell ignites if it isn't being worn or carried.",
			"Duration" : "Instantaneous",
			"Concentration" : "False",
			"Ritual" : "False",
			"ClassList" : ["Artificer", "Sorcerer", "Wizard"],
			"Source" : "PlayersHB",
			"School" : "Evocation",
			"Components" : ["V", "S"],
			"HigherLevels" : "This spell’s damage increases by 1d10 when you reach 5th level (2d10), 11th level (3d10), and 17th level (4d10)."
		}`,
	}
	assertGETSpell(t, test)
}

func TestGETTollTheDead(t *testing.T) {
	test := testCase{
		name:      "Get Toll the Dead returns Toll the Dead",
		spellName: "Spare The Dying",
		json: `{
			"Name" : "Spare the Dying",
			"Range" : "Touch",
			"Level" : "cantrip",
			"CastTime" : "1 action",
			"Description" : "You touch a living creature that has 0 hit points. The creature becomes stable. This spell has no effect on undead or constructs.",
			"Duration" : "Instantaneous",
			"Concentration" : "False",
			"Ritual" : "False",
			"ClassList" : ["Artificer", "Cleric"],
			"Source" : "PlayersHB",
			"School" : "Necromancy",
			"Components" : ["V", "S"],
			"HigherLevels" : ""
		}`,
	}
	assertGETSpell(t, test)
}
