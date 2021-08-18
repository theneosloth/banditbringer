package util

import (
	"github.com/bwmarrin/discordgo"
)

// Constants for message embed character limits
const (
	limit = 2048
)

func CreateInteractionMessage(text string) *discordgo.InteractionResponse {
	if len(text) > limit {
		return nil
	}
	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: text,
		},
	}
}

func CreateInteractionEmbed(embed *discordgo.MessageEmbed) *discordgo.InteractionResponse {
	embeds := []*discordgo.MessageEmbed{
		embed,
	}
	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: embeds,
		},
	}
}
