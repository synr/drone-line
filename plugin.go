package main

import (
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/line/line-bot-sdk-go/linebot"
)

type (
	// Repo information.
	Repo struct {
		Owner string
		Name  string
	}

	// Build information.
	Build struct {
		Event  string
		Number int
		Commit string
		Branch string
		Author string
		Status string
		Link   string
	}

	// Config for the plugin.
	Config struct {
		ChannelID     string
		ChannelSecret string
		MID           string
		To            []string
		Message       string
	}

	// Plugin values.
	Plugin struct {
		Repo   Repo
		Build  Build
		Config Config
	}
)

// Exec executes the plugin.
func (p Plugin) Exec() error {

	if len(p.Config.ChannelID) == 0 || len(p.Config.ChannelSecret) == 0 || len(p.Config.MID) == 0 {
		log.Println("missing line bot config")

		return errors.New("missing line bot config")
	}

	ChannelID, err := strconv.ParseInt(p.Config.ChannelID, 10, 64)
	if err != nil {
		log.Println("wrong channel id")

		return err
	}

	bot, _ := linebot.NewClient(ChannelID, p.Config.ChannelSecret, p.Config.MID)

	if len(p.Config.To) == 0 {
		log.Println("missing line user config")

		return errors.New("missing line user config")
	}

	var message string
	if p.Config.Message != "" {
		message = p.Config.Message
	} else {
		message = p.Message(p.Repo, p.Build)
	}

	_, err = bot.SendText(p.Config.To, message)

	if err != nil {
		log.Println(err.Error())

		return err
	}

	return nil
}

// Message is line default message.
func (p Plugin) Message(repo Repo, build Build) string {
	return fmt.Sprintf("[%s] <%s|%s/%s#%s> (%s) by %s",
		build.Status,
		build.Link,
		repo.Owner,
		repo.Name,
		build.Commit[:8],
		build.Branch,
		build.Author,
	)
}
