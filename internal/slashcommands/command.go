package slashcommands

import (
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
	if required {
		return &SubCommand{&discordgo.ApplicationCommandOption{
			Type:        commandType,
			Name:        name,
			Description: description,
		}}

	}
	return &SubCommand{&discordgo.ApplicationCommandOption{
		Type:        commandType,
		Name:        name,
		Description: description,
		Required:    required,
	}}
}

func NewChoice(name string, value interface{}) *Choice {
	return &Choice{
		&discordgo.ApplicationCommandOptionChoice{
			Name:  name,
			Value: value,
		}}
}

func (c *SubCommand) AddOption(commandType discordgo.ApplicationCommandOptionType, name string, description string, required bool) *SubCommand {
	if len(name) == 0 || len(description) == 0 {
		return nil
	}
	subCommand := NewSubCommand(commandType, name, description, required)
	c.ApplicationCommandOption.Options = append(c.ApplicationCommandOption.Options, subCommand.ApplicationCommandOption)
	return c
}

func (c *SubCommand) AddChoice(name string, value interface{}) *SubCommand {
	if len(name) == 0 {
		return nil
	}
	choice := NewChoice(name, value)
	c.ApplicationCommandOption.Choices = append(c.ApplicationCommandOption.Choices, choice.ApplicationCommandOptionChoice)
	return c
}

func (c *SlashCommand) AddOption(commandType discordgo.ApplicationCommandOptionType, name string, description string, required bool) *SubCommand {
	if len(name) == 0 || len(description) == 0 {
		return nil
	}
	subCommand := NewSubCommand(commandType, name, description, required)
	c.ApplicationCommand.Options = append(c.ApplicationCommand.Options, subCommand.ApplicationCommandOption)
	return subCommand
}

func (c *SlashCommand) Run(s *discordgo.Session, i *discordgo.InteractionCreate) {
	c.callback(s, i)
}
