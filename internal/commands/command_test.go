package commands

import (
	"strings"
	"testing"

	"github.com/bwmarrin/discordgo"
)

var (
	commandName  = "TestCommand"
	prefix       = "!"
	validCommand = "!echo"
)

func TestCheckCommand(t *testing.T) {
	testCommand := Command{
		Name:     commandName,
		Prefix:   prefix,
		callback: nil,
	}

	testValid := testCommand.checkCommand(validCommand)

	if !testValid {
		t.Errorf("%s should be a valid command", validCommand)
	}

	invalidCommand := "?echo"
	testInvalid := testCommand.checkCommand(invalidCommand)
	if testInvalid {
		t.Errorf("%s should be an invalid command", invalidCommand)
	}

	noWord := "!!!!!?"
	testWord := testCommand.checkCommand(noWord)
	if testWord {
		t.Errorf("%s should be an invalid command", noWord)
	}

}

func TestRun(t *testing.T) {

	mockMessage := discordgo.Message{
		Content: validCommand,
	}

	mockMessageCreate := discordgo.MessageCreate{
		Message: &mockMessage,
	}

	mockSession := discordgo.Session{}

	testCommand := Command{
		Name:   commandName,
		Prefix: prefix,
		callback: func(s *discordgo.Session, m *discordgo.MessageCreate) {
			if len(m.Content) == 0 {
				t.Error("passed message empty")
			}
			if strings.HasPrefix(m.Content, prefix) {
				t.Error("message still has prefix")
			}
		},
	}

	testCommand.Run(&mockSession, &mockMessageCreate)
}
