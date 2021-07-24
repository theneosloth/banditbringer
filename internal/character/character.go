package character

import (
	"banditbringer/internal/move"
	"errors"
	"reflect"
	"strings"
)

type Character struct {
	Name                  string      `json:"name"`
	ImageUrl              string      `json:"image_url"`
	Defense               string      `json:"defense"`
	Guts                  string      `json:"guts"`
	Prejump               string      `json:"prejump"`
	Backdash              string      `json:"backdash"`
	Weight                string      `json:"weight"`
	UniqueMovementOptions string      `json:"unique_movement_options"`
	DustloopUrl           string      `json:"dustloop_url"`
	Moves                 []move.Move `json:"moves"`
	aliases               []string
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

func (c *Character) SetFieldByName(field string, value string) error {
	rv := reflect.ValueOf(c)

	if rv.Kind() != reflect.Ptr {
		return errors.New("character is not a pointer")
	}

	rv = rv.Elem()

	if rv.Kind() != reflect.Struct {
		return errors.New("character is not a struct")
	}

	f := rv.FieldByName(field)

	if !f.IsValid() || f.Kind() != reflect.String {
		return errors.New("trying to set a field that is not a string")
	}

	f.SetString(value)

	return nil
}
