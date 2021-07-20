package move

import (
	"testing"
)

func TestSetFieldByname(t *testing.T) {
	move := Move{}
	err := move.SetFieldByName("Input", "5P")

	if err != nil {
		t.Errorf("valid field set failed")
	}

	err = move.SetFieldByName("InvalidField", "5P")

	if err == nil {
		t.Errorf("invalid field set")
	}
}
