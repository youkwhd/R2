package commands

import (
	"fmt"
	"strings"
	"time"

	db "R2/internal/db/mem"

	"github.com/bwmarrin/discordgo"
)

type CommandHandlerFn func(botSession *discordgo.Session, i *discordgo.InteractionCreate)

type Command struct {
	Information *discordgo.ApplicationCommand
	Handler CommandHandlerFn
}

var /* const */ COMMANDS = [...]Command{
	{
		Information: &discordgo.ApplicationCommand{
			Name: "ping",
			Description: "<test> returns back pong",
		},
		Handler: func(botSession *discordgo.Session, i *discordgo.InteractionCreate) {
			botSession.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Pong!",
				},
			})
		},
	},
	{
		Information: &discordgo.ApplicationCommand{
			Name: "react",
			Description: "Add a new reaction role",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "message_link",
					Description: "Message link",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionRole,
					Name:        "role",
					Description: "Role",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "emoji",
					Description: "Emoji",
					Required:    true,
				},
			},
		},
		Handler: func(botSession *discordgo.Session, i *discordgo.InteractionCreate) {
			_args := i.ApplicationCommandData().Options

			args := make(map[string]*discordgo.ApplicationCommandInteractionDataOption)
			for _, arg := range _args {
				args[arg.Name] = arg
			}

			messageLink := args["message_link"].Value.(string)
			// TODO: update go to ver 1.20, use CutPrefix
			_, messageLink, _ = strings.Cut(messageLink, "https://discord.com/channels/")
			_, messageLink, _ = strings.Cut(messageLink, "/")
			channelID, messageLink, _ := strings.Cut(messageLink, "/")

			messageID := messageLink
			role := args["role"].Value.(string)
			emoji := args["emoji"].Value.(string)

			// ?? what is this type rule, golang??
			db.Messages[db.MessageID(messageID)] = db.NewRoleReactionMessage(channelID)
			db.Messages[db.MessageID(messageID)].Reactions[db.Emoji(emoji)] = db.Role(role)

			botSession.MessageReactionAdd(channelID, messageID, emoji)

			roles, _ := botSession.GuildRoles(i.GuildID)

			// TODO: Golang really does not have a .Find() function
			var R *discordgo.Role
			for _, r := range roles {
				if r.ID == role {
					R = r
				}
			}

			botSession.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("Registered (**@%s**)[%s]", R.Name, emoji),
				},
			})

			time.Sleep(time.Second * 2)
			botSession.InteractionResponseDelete(i.Interaction)
		},
	},
}