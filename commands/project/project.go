// Package project hold actions on the Bitbucket projects
package project

import (
	"github.com/daeMOn63/bitadmin/settings"
	"github.com/urfave/cli"
)

// Command define base struct for repository subcommands and actions
type Command struct {
	Settings *settings.BitAdminSettings
}

// GetCommand provide a ready to use cli.Command
func (command *Command) GetCommand() cli.Command {

	listProjectRepositories := &ListRepositoriesCommand{
		Settings: command.Settings,
		flags:    &ListRepositoriesFlags{},
	}

	return cli.Command{
		Name:  "project",
		Usage: "Project operations",
		Subcommands: []cli.Command{
			listProjectRepositories.GetCommand(),
		},
	}
}
