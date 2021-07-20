package character

import (
	"banditbringer/internal/move"
	"errors"
	"strings"
)

type Character struct {
	Name    string
	Moves   []move.Move
	aliases []string
}

var validCharacters = map[string][]string{
	"ky_kiske":            {"ky", "kyle"},
	"anji_mito":           {"anji"},
	"axl_low":             {"axl"},
	"chipp_zanuff":        {"chipp", "chip"},
	"faust":               {"doctor", "faust"},
	"giovanna":            {"gio"},
	"i-no":                {"ino"},
	"leo_whitefang":       {"leo"},
	"may":                 {""},
	"millia_rage":         {"millia"},
	"nagoriyuki":          {"nago"},
	"potemkin":            {"pot"},
	"ramlethal_valentine": {"ram", "valentine"},
	"sol_badguy":          {"sol"},
	"zato-1":              {"zato", "eddie"},
}

func IsValidCharName(name string) (normalizedName string, found bool) {

	name = strings.ToLower(strings.Replace(name, " ", "_", -1))
	_, exists := validCharacters[name]

	if exists {
		return name, true
	}

	for k, v := range validCharacters {
		for _, v0 := range v {
			if v0 == name {
				return k, true
			}
		}
	}
	return "", false
}
func (c *Character) SetName(name string) error {
	_, found := IsValidCharName(name)
	if found {
		c.Name = name
		return nil
	}
	return errors.New("Not a valid character")
}
