package commands

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
	"github.com/hbollon/go-edlib"
)

var Fetcher Command

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
func replaceEmptyString(s string) string {
	if len(s) == 0 {
		return "N/A"
	}
	// TODO: Extend to escape all markdown
	return strings.ReplaceAll(s, "*", "\\*")
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
	// Embed silently fails if given an empty string
	replaceEmptyString := func(s string) string {
		if len(s) == 0 {
			return "N/A"
		}
		// TODO: Extend to escape all markdown
		return strings.ReplaceAll(s, "*", "\\*")
	}

	return &discordgo.MessageEmbed{
		Author:      &discordgo.MessageEmbedAuthor{},
		Color:       0x00ff00,
		Description: "",
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL:    c.Icon,
			Width:  50,
			Height: 50,
		},
		URL: c.DustloopUrl,

		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Defense",
				Value:  replaceEmptyString(c.Defense),
				Inline: true,
			},
			{
				Name:   "Guts",
				Value:  replaceEmptyString(c.Guts),
				Inline: true,
			},
			{
				Name:   "Prejump",
				Value:  replaceEmptyString(c.Prejump),
				Inline: true,
			},
			{
				Name:   "Backdash",
				Value:  replaceEmptyString(c.Backdash),
				Inline: true,
			},
			{
				Name:   "Weight",
				Value:  replaceEmptyString(c.Weight),
				Inline: true,
			},
			{
				Name:  "Unique Movement Options",
				Value: replaceEmptyString(c.UniqueMovementOptions),
			},
		},
		Title: c.GetReadableName(),
	}

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
	Fetcher = Command{
		Name:   "Fetcher",
		Prefix: "?",
		callback: func(s *discordgo.Session, m *discordgo.MessageCreate) {
			found := strings.SplitN(m.Content, " ", 2)

			if len(found[0]) == 0 {
				return
			}

			char := found[0]

			normalizedName, charExists := character.IsValidName(char)

			if !charExists {
				s.ChannelMessageSend(m.ChannelID, "Character does not exist")
				return
			}

			character := loadChar(normalizedName)

			if len(found) == 1 {
				embed := generateCharEmbed(character)
				s.ChannelMessageSendEmbed(m.ChannelID, embed)
				return
			}

			if len(found) != 2 {
				return
			}

			moves := character.Moves

			for _, move := range moves {
				if sameMove(move.Input, found[1]) || normalizeCompare(move.Name, found[1]) {
					embed := generateMoveEmbed(character, move)
					s.ChannelMessageSendEmbed(m.ChannelID, embed)
					return
				}
			}

			res, err := edlib.FuzzySearchSet(strings.ToUpper(found[1]), character.GetAllMoves(), 4, edlib.Levenshtein)
			if err != nil {
				fmt.Println(err)
			} else if len(res) > 0 {
				msg := fmt.Sprintf("Move not found. Did you mean one of the following: %s?", strings.Join(res, ", "))
				s.ChannelMessageSend(m.ChannelID, msg)
				return
			}
			s.ChannelMessageSend(m.ChannelID, "Move not found")

		},
	}
}
