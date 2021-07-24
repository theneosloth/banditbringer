package main

import (
	"banditbringer/internal/commands"

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
		fmt.Fprintf(os.Stderr, "error %s", err)
		os.Exit(1)
	}

	err = dg.Open()
	if err != nil {
		fmt.Println("Error opening connection", err)
		os.Exit(1)
	}

	dg.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		fmt.Println("Started successfully")
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

	defer func() {
		sc := make(chan os.Signal, 1)
		signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
		<-sc
		// Cleanly close down the Discord session.
		fmt.Println("Shutting down")
		dg.Close()
	}()

}
