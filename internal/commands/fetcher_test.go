package commands

import (
	"fmt"
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

func TestNormalizeCompare(t *testing.T) {
	command1 := "j.632146H"
	command2 := "J632146 H"

	if !normalizeCompare(command1, command2) {
		fmt.Println(normalizeCommand(command1), normalizeCommand(command2))
		t.Errorf("commands should be the same %s %s", command1, command2)
	}
}

func TestSameMove(t *testing.T) {
	command1 := "63214H"
	command2 := "624H"

	if !sameMove(command1, command2) {
		t.Errorf("hcf input should be the same %s %s", command1, command2)
	}
}
