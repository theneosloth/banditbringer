package character

import (
	"banditbringer/internal/move"
	"errors"
	"reflect"
	"strings"
)

// https://dustloop.com/wiki/index.php?title=GGST/Sol_Badguy/Data
type Character struct {
	Name                            string      `json:"name"`
	Defense                         string      `json:"defense"`
	Guts                            string      `json:"guts"`
	RiscModifier                    string      `json:"risc_modifier"`
	Prejump                         string      `json:"prejump"`
	Backdash                        string      `json:"backdash"`
	ForwardDash                     string      `json:"forward_dash"`
	Weight                          string      `json:"weight"`
	JumpDuration                    string      `json:"jump_duration"`
	HighJumpDuration                string      `json:"high_jump_duration"`
	JumpHeight                      string      `json:"jump_height"`
	EarliestForwardInstantAirDash   string      `json:"earliest_forward_iad"`
	EarliestBackwardInstantAirDash  string      `json:"earliest_backward_iad"`
	ForwardAirDashDuration          string      `json:"forward_ad_duration"`
	BackwardAirDashDuration         string      `json:"backward_ad_duration"`
	ForwardAirDashAttackTransition  string      `json:"forward_ad_attack_transition"`
	BackwardAirDashAttackTransition string      `json:"backward_ad_attack_transition"`
	UniqueMovementOptions           string      `json:"unique_movement_options"`
	Portrait                        string      `json:"portrait"`
	Icon                            string      `json:"icon"`
	Voice                           string      `json:"voice"`
	Theme                           string      `json:"theme"`
	DustloopUrl                     string      `json:"dustloop_url"`
	Moves                           []move.Move `json:"moves"`
	aliases                         []string
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
	"ramlethal_valentine": {"ram"},
	"sol_badguy":          {"sol"},
	"zato-1":              {"zato", "eddie"},
	"goldlewis_dickinson": {"goldick", "golddick", "goldlewis"},
	"jack-o":              {"aria", "jacko"},
}

func IsValidName(name string) (normalizedName string, found bool) {

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

func (c *Character) GetReadableName() string {
	return strings.Title(strings.ReplaceAll(c.Name, "_", " "))
}

func (c *Character) SetName(name string) error {
	_, found := IsValidName(name)
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

	if !f.IsValid() {
		return errors.New("invalid field")
	}

	if f.Kind() != reflect.String {
		return errors.New("field is not a string")
	}

	f.SetString(value)

	return nil
}

func (c *Character) GetAllMoves() (moves []string) {
	moves = make([]string, len(c.Moves))
	for i, v := range c.Moves {
		moves[i] = v.Input
	}
	return moves
}
