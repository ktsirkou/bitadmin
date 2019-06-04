package project

import (
	"fmt"
	"github.com/daeMOn63/bitclient"
	"github.com/daeMOn63/bitadmin/settings"
	"github.com/urfave/cli"
	"github.com/daeMOn63/bitadmin/helper"
)

// ListRepositoriesCommand define base struct for ListRepositories action
type ListRepositoriesCommand struct {
	Settings *settings.BitAdminSettings
	flags    *ListRepositoriesFlags
}

// ListRepositoriesFlags define flags required by the ListRepositories action
type ListRepositoriesFlags struct {
	project    string
}

// GetCommand provide a ready to use cli.Command
func (command *ListRepositoriesCommand) GetCommand() cli.Command {
	return cli.Command{
		Name:   "list-repositories",
		Usage:  "Show repositories on project",
		Action: command.ListRepositoriesAction,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:        "project",
				Usage:       "The `<project>` to list work on",
				Destination: &command.flags.project,
			},
		},
		BashComplete: func(c *cli.Context) {
			helper.AutoComplete(c, command.Settings.GetFileCache())
		},
	}
}

type OutputRow struct {
	Slug       string
}

type OutputRows []OutputRow

// ListRepositoriesAction display the current user / group permissions on given repository
func (command *ListRepositoriesCommand) ListRepositoriesAction(context *cli.Context) error {

	client, err := command.Settings.GetAPIClient()
	if err != nil {
		return err
	}

	repositoriesResponse, err := client.GetRepositories(
		command.flags.project,
		bitclient.PagedRequest{Limit:250},
	)
	if err != nil {
		return err
	}
	var rows OutputRows
	for _, repositoryItem := range repositoriesResponse.Values {
		row := OutputRow{Slug: repositoryItem.Slug}

		rows = rows.appendIfNew(row)
	}


	fmt.Printf("%s\n", rows)

	return nil
}

func (rows OutputRows) appendIfNew(row OutputRow) OutputRows {
	for _, r := range rows {
		if r.Slug == row.Slug {
			return rows
		}
	}
	rows = append(rows, row)
	return rows
}

func (rows OutputRows) String() string {
	out := "Slug\n"
	for _, row := range rows {
		out += fmt.Sprintf("%s", row)
	}

	return out
}

func (row OutputRow) String() string {
	return fmt.Sprintf("%s\n",
		row.Slug,
	)
}
