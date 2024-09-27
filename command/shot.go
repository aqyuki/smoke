package command

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func NewShotCommand() *Command {
	return &Command{
		Command: &discordgo.ApplicationCommand{
			Name:        "shot",
			Description: "指定したユーザーに飯テロを実行します",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "user",
					Description: "飯テロのターゲット",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionAttachment,
					Name:        "image",
					Description: "飯テロ画像",
					Required:    true,
				},
			},
		},
		Handler: func(s *discordgo.Session, ic *discordgo.InteractionCreate) {
			option := make(map[string]*discordgo.ApplicationCommandInteractionDataOption)
			fmt.Println("shot command")

			for _, v := range ic.ApplicationCommandData().Options {
				option[v.Name] = v
			}

			user := option["user"].UserValue(s)
			images := ic.ApplicationCommandData().Resolved.Attachments

			dm, err := s.UserChannelCreate(user.ID)
			if err != nil {
				fmt.Println("Failed to create DM channel")
				return
			}

			var content string
			for _, image := range images {
				content += image.URL + "\n"
			}

			if _, err := s.ChannelMessageSendComplex(dm.ID, &discordgo.MessageSend{
				Content: content,
			}); err != nil {
				_ = s.InteractionRespond(ic.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "飯テロに失敗しました",
					},
				})
			}

			_ = s.InteractionRespond(ic.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "飯テロに成功しました",
				},
			})
		},
	}
}
