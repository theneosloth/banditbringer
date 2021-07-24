package commands

import (
	"banditbringer/internal/character"
	"banditbringer/internal/move"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
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

func generateEmbed(character character.Character, m move.Move) *discordgo.MessageEmbed {

	// Embed silently fails if given an empty string
	replaceEmptyString := func(s string) string {
		if len(s) == 0 {
			return "N/A"
		}
		return s
	}

	title := ""
	readableName := strings.Title(strings.ReplaceAll(m.Name, "_", " "))
	if m.Name != "" {
		title = fmt.Sprintf("%s\n%s", readableName, m.Input)
	} else {
		title = m.Input
	}

	return &discordgo.MessageEmbed{
		Author:      &discordgo.MessageEmbedAuthor{},
		Color:       0x00ff00, // Green
		Description: title,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL:      character.ImageUrl,
			ProxyURL: "https://www.guiltygear.com/ggst/assets/logo.jpg",
			Width:    50,
			Height:   50,
		},
		URL: character.DustloopUrl,

		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Damage",
				Value:  replaceEmptyString(m.Damage),
				Inline: true,
			},
			{
				Name:   "Guard",
				Value:  replaceEmptyString(m.Guard),
				Inline: true,
			},
			{
				Name:   "Startup",
				Value:  replaceEmptyString(m.Startup),
				Inline: true,
			},
			{
				Name:   "On Block",
				Value:  replaceEmptyString(m.OnBlock),
				Inline: true,
			},
			{
				Name:   "On Hit",
				Value:  replaceEmptyString(m.OnHit),
				Inline: true,
			},
		},
		Title: strings.Title(character.Name),
	}

}

func normalizeCommand(command string) string {
	pattern := regexp.MustCompile(`[\s.]+`)
	return pattern.ReplaceAllString(
		strings.TrimSpace(strings.ToLower(command)),
		"",
	)
}

func normalizeCompare(i string, j string) bool {
	return normalizeCommand(i) == normalizeCommand(j)
}

func init() {
	Fetcher = Command{
		Name:   "Fetcher",
		Prefix: "?",
		callback: func(s *discordgo.Session, m *discordgo.MessageCreate) {
			found := strings.SplitN(m.Content, " ", 2)

			if len(found) != 2 {
				return
			}
			char := found[0]

			normalizedName, charExists := character.IsValidCharName(char)

			if !charExists {
				s.ChannelMessageSend(m.ChannelID, "Character does not exist")
				return
			}

			character := loadChar(normalizedName)

			moves := character.Moves

			for _, move := range moves {
				if normalizeCompare(move.Input, found[1]) || normalizeCompare(move.Name, found[1]) {
					embed := generateEmbed(character, move)
					s.ChannelMessageSendEmbed(m.ChannelID, embed)
					return
				}
			}

			s.ChannelMessageSend(m.ChannelID, "Move not found")

		},
	}
}
