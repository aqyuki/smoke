package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/aqyuki/smoke/command"
	"github.com/bwmarrin/discordgo"
)

type exitCode int

const (
	ExitCodeOK exitCode = iota
	ExitCodeError
)

func main() {
	exit(realMain())
}

func realMain() exitCode {
	token, ok := os.LookupEnv("DISCORD_TOKEN")
	if !ok || token == "" {
		fmt.Println("DISCORD_TOKEN environment variable is not set")
		return ExitCodeError
	}

	// discordgo.New does not return an error.
	session, _ := discordgo.New("Bot " + token)

	if err := session.Open(); err != nil {
		fmt.Printf("Failed to open a connection to Discord Gateway API: %v\n", err)
		return ExitCodeError
	}

	router := command.NewCommandRouter()
	router.Add(command.NewShotCommand())
	session.AddHandler(router.OnInteractionCreate)
	if err := router.Register(session); err != nil {
		fmt.Printf("Failed to register commands: %v\n", err)
		_ = session.Close()
		return ExitCodeError
	}

	ctx, done := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer done()

	fmt.Println("Bot is running. SIGINT or SIGTERM to stop.")
	<-ctx.Done()

	if err := session.Close(); err != nil {
		fmt.Printf("Failed to close the connection to Discord Gateway API: %v\n", err)
		return ExitCodeError
	}
	return ExitCodeOK
}

func exit[T ~int](code T) {
	os.Exit(int(code))
}
