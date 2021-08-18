package slashcommands

import (
	"errors"

	"github.com/bwmarrin/discordgo"
)

type Callback func(*discordgo.Session, *discordgo.InteractionCreate)

type SlashCommand struct {
	ApplicationCommand *discordgo.ApplicationCommand
	callback           Callback
}

type SubCommand struct {
	ApplicationCommandOption *discordgo.ApplicationCommandOption
}

type Choice struct {
	ApplicationCommandOptionChoice *discordgo.ApplicationCommandOptionChoice
}

func NewSlashCommand(name string, description string, callback Callback) *SlashCommand {
	return &SlashCommand{
		ApplicationCommand: &discordgo.ApplicationCommand{
			Name:        name,
			Description: description,
		},
		callback: callback,
	}
}

func NewSubCommand(commandType discordgo.ApplicationCommandOptionType, name string, description string, required bool) *SubCommand {
	return &SubCommand{&discordgo.ApplicationCommandOption{
		Type:        commandType,
		Name:        name,
		Description: description,
		Required:    required,
	}}
}
func (c *SubCommand) AddOption(commandType discordgo.ApplicationCommandOptionType, name string, description string, required bool) error {
	if len(name) == 0 || len(description) == 0 {
		return errors.New("Name or description empty")
	}
	subCommand := NewSubCommand(commandType, name, description, required)
	c.ApplicationCommandOption.Options = append(c.ApplicationCommandOption.Options, subCommand.ApplicationCommandOption)
	return nil
}

func (c *SlashCommand) AddOption(commandType discordgo.ApplicationCommandOptionType, name string, description string, required bool) error {
	if len(name) == 0 || len(description) == 0 {
		return errors.New("Name or description empty")
	}
	subCommand := NewSubCommand(commandType, name, description, required)
	c.ApplicationCommand.Options = append(c.ApplicationCommand.Options, subCommand.ApplicationCommandOption)
	return nil
}

func (c *SlashCommand) Run(s *discordgo.Session, i *discordgo.InteractionCreate) {
	c.callback(s, i)
}
