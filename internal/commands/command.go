package commands

import (
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type Callback func(*discordgo.Session, *discordgo.MessageCreate)

type Runnable interface {
	checkCommand(string)
	Run(func(*discordgo.Session, discordgo.MessageCreate))
}

type Command struct {
	Name     string
	Prefix   string
	callback Callback
}

func (c *Command) checkCommand(m string) bool {
	hasLetter := regexp.MustCompile(".*[A-z].*")

	isWord := hasLetter.MatchString(m)

	if !isWord {
		return false
	}

	return strings.HasPrefix(m, c.Prefix)
}

func (c *Command) Run(s *discordgo.Session, m *discordgo.MessageCreate) {
	msg := strings.Trim(strings.ToLower(m.Content), " ")
	if c.checkCommand(msg) {
		m.Content = strings.TrimPrefix(m.Content, c.Prefix)
		c.callback(s, m)
	}
}
