package main

import (
	"banditbringer/internal/commands"
	"banditbringer/internal/slashcommands"
	"log"

	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

func main() {

	Token := os.Getenv("DISCORD_TOKEN")
	dg, err := discordgo.New("Bot " + Token)

	if err != nil {
		log.Panicf("error %s", err)
	}

	err = dg.Open()
	if err != nil {
		log.Panicf("Error opening connection %s", err)
	}

	dg.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Println("Started successfully")
	})

	activeCommands := []commands.Command{
		commands.Fetcher,
	}

	dg.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == s.State.User.ID {
			return
		}

		for _, command := range activeCommands {
			command.Run(s, m)
		}

	})

	command := slashcommands.Fetcher

	_, err = dg.ApplicationCommandCreate(dg.State.User.ID, "404881906690293760", command.ApplicationCommand)
	if err != nil {
		log.Panicf("Cannot create '%v' command: %v", command.ApplicationCommand.Name, err)
	}
	dg.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		command.Run(s, i)
	})

	defer func() {
		sc := make(chan os.Signal, 1)
		signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
		<-sc
		// Cleanly close down the Discord session.
		fmt.Println("Shutting down")
		dg.Close()
	}()

}
