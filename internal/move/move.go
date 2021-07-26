package move

import (
	"errors"
	"reflect"
)

// Stringly typed :(
type Move struct {
	Images   string
	Name     string
	Input    string
	Damage   string
	Guard    string
	Startup  string
	Active   string
	Recovery string
	OnBlock  string
	OnHit    string
	Level    string
	CHType   string
	Hitboxes string
	Notes    string
	Type     string
	RISCGain string
	Prorate  string
	Invuln   string
	Cancel   string
}

func (m *Move) SetFieldByName(field string, value string) error {
	rv := reflect.ValueOf(m)

	if rv.Kind() != reflect.Ptr {
		return errors.New("move is not a pointer")
	}

	rv = rv.Elem()

	if rv.Kind() != reflect.Struct {
		return errors.New("move is not a struct")
	}

	f := rv.FieldByName(field)

	if !f.IsValid() || f.Kind() != reflect.String {
		return errors.New("field not valid")
	}

	f.SetString(value)

	return nil
}
