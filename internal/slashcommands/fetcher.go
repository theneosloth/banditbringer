package slashcommands

import (
	"banditbringer/internal/character"
	"banditbringer/internal/move"
	"banditbringer/internal/util"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

var Fetcher SlashCommand
var characterNames []string

func init() {
	characterNames = character.GetAllCharacters()
}

func generateMoveEmbed(character character.Character, m move.Move) *discordgo.MessageEmbed {

	description := ""
	if m.Name != "" {
		description = fmt.Sprintf("%s\n%s", m.Name, m.Input)
	} else {
		description = m.Input
	}

	image := "N/A"
	if m.Hitboxes != "" {
		image = m.Hitboxes
	} else if m.Images != "" {
		image = m.Images
	} else {
		image = character.Icon
	}

	return util.NewEmbed().
		SetTitle(character.GetReadableName()).
		SetDescription(description).
		SetThumbnail(image).
		SetColor(0xaf0016).
		AddField("Damage", m.Damage).
		AddField("Guard", m.Guard).
		AddField("Startup", m.Startup).
		AddField("On Block", m.OnBlock).
		AddField("On Hit", m.OnHit).
		InlineAllFields().
		MessageEmbed
}

func generateCharEmbed(c character.Character) *discordgo.MessageEmbed {
	return util.NewEmbed().
		SetTitle(c.GetReadableName()).
		SetThumbnail(c.Icon).
		SetColor(0xaf0016).
		SetURL(c.DustloopUrl).
		AddField("Defense", c.Defense).
		AddField("Guts", c.Guts).
		AddField("Prejump", c.Prejump).
		AddField("Backdash", c.Backdash).
		AddField("Weight", c.Weight).
		AddField("Unique Movement Options", c.UniqueMovementOptions).MessageEmbed
}

func init() {
	Fetcher = *NewSlashCommand("fetcher", "Frame data fetcher", func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		name := i.ApplicationCommandData().Options[0].StringValue()
		if name == "" {
			return
		}
		character := character.LoadChar(name)
		embed := generateCharEmbed(character)
		m := util.CreateInteractionEmbed(embed)
		s.InteractionRespond(i.Interaction, m)
	})

	nameOption := Fetcher.AddOption(discordgo.ApplicationCommandOptionString, "name", "The character name", true)
	for _, char := range characterNames {
		nameOption.AddChoice(char, char)
	}
}
