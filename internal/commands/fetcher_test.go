package commands

import (
	"testing"
)

func TestNormalizeCommand(t *testing.T) {
	command := "f.S"

	normalizedCommand := normalizeCommand(command)
	expectedCommand := "fs"
	if normalizedCommand != "fs" {
		t.Errorf("normalized command should be equal to %s", expectedCommand)
	}

}

func testNormalizeCompare(t *testing.T) {
	command1 := "j.632146H"
	command2 := "J632136 H"

	if !normalizeCompare(command1, command2) {
		t.Errorf("commands should be the same %s %s", command1, command2)
	}
}
