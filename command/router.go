package command

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type Command struct {
	Command *discordgo.ApplicationCommand
	Handler  func(*discordgo.Session, *discordgo.InteractionCreate)
}

type CommandRouter struct {
	commands   []*discordgo.ApplicationCommand
	registered []*discordgo.ApplicationCommand
	handlers   map[string]func(*discordgo.Session, *discordgo.InteractionCreate)
}

func NewCommandRouter() *CommandRouter {
	return &CommandRouter{
		commands: make([]*discordgo.ApplicationCommand, 0),
		handlers: make(map[string]func(*discordgo.Session, *discordgo.InteractionCreate)),
	}
}

func (h *CommandRouter) Add(command *Command) {
	h.commands = append(h.commands, command.Command)
	h.handlers[command.Command.Name] = command.Handler
}

func (h *CommandRouter) Register(session *discordgo.Session) error {
	for _, command := range h.commands {
		cmd, err := session.ApplicationCommandCreate(session.State.User.ID, "", command)
		if err != nil {
			return fmt.Errorf("failed to register command %s: %w", command.Name, err)
		}
		fmt.Printf("success to register command : %s\n", cmd.Name)
		h.registered = append(h.registered, cmd)
	}
	return nil
}

func (h *CommandRouter) Unregister(session *discordgo.Session) error {
	for _, command := range h.registered {
		if err := session.ApplicationCommandDelete(session.State.User.ID, "", command.ID); err != nil {
			return fmt.Errorf("failed to unregister command %s: %w", command.Name, err)
		}
		fmt.Printf("success to unregister command : %s\n", command.Name)
	}
	return nil
}

func (h *CommandRouter) Commands() []*discordgo.ApplicationCommand {
	return h.commands
}

func (h *CommandRouter) OnInteractionCreate(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	if interaction.Type != discordgo.InteractionApplicationCommand {
		return
	}

	if handler, ok := h.handlers[interaction.ApplicationCommandData().Name]; ok {
		handler(session, interaction)
	}
}
