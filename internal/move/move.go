package move

import (
	"errors"
	"reflect"
)

// Stringly typed :(
type Move struct {
	Input    string
	Name     string
	Damage   string
	Guard    string
	Startup  string
	Active   string
	Recovery string
	OnBlock  string
	OnHit    string
	RiscGain string
	Level    string
	Invuln   string
	Prorate  string
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
