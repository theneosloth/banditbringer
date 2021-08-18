package slashcommands

import (
	"banditbringer/internal/character"
	"banditbringer/internal/move"
	"banditbringer/internal/util"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var Fetcher SlashCommand

func loadChar(name string) (character character.Character) {
	name = strings.Replace(name, " ", "_", -1)

	fname := fmt.Sprintf("%s.json", name)
	fpath, err := filepath.Abs(filepath.Join("json", fname))

	file, err := ioutil.ReadFile(fpath)

	if err != nil {
		panic(err)
	}

	json.Unmarshal(file, &character)

	return character
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

func normalizeCommand(command string) string {
	pattern := regexp.MustCompile(`[\s.,]+`)
	return pattern.ReplaceAllString(
		strings.TrimSpace(strings.ToLower(command)),
		"",
	)
}

func removeDiagonals(command string) string {
	diagonals := regexp.MustCompile(`[1379]+`)
	return diagonals.ReplaceAllString(command, "")
}

func normalizeCompare(i string, j string) bool {
	return normalizeCommand(i) == normalizeCommand(j)
}

func sameMove(command1 string, command2 string) bool {
	normalizedEqual := normalizeCompare(command1, command2)
	// Try again with diagonals removed for HCF motions
	if !normalizedEqual && len(command1) > 3 && len(command2) > 3 {
		normalizedEqual = normalizeCompare(removeDiagonals(command1), removeDiagonals(command2))
	}
	return normalizedEqual
}

func init() {
	Fetcher = *NewSlashCommand("fetcher", "Frame data fetcher", func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		character := loadChar("May")
		embed := generateCharEmbed(character)
		m := util.CreateInteractionEmbed(embed)
		s.InteractionRespond(i.Interaction, m)
	})
}
