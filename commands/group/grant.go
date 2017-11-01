package group

import (
	"fmt"
	"github.com/daeMOn63/bitadmin/helper"
	"github.com/daeMOn63/bitadmin/settings"
	"github.com/daeMOn63/bitclient"
	"github.com/urfave/cli"
)

type GroupGrantCommand struct {
	Settings *settings.BitAdminSettings
	flags    *GroupGrantCommandFlags
}

type GroupGrantCommandFlags struct {
	repositories cli.StringSlice
	names        cli.StringSlice
	permission   string
}

func (command *GroupGrantCommand) GetCommand() cli.Command {
	return cli.Command{
		Name:   "grant",
		Usage:  "Grant groups permission on repositories",
		Action: command.GrantAction,
		Flags: []cli.Flag{
			cli.StringSliceFlag{
				Name:  "repository",
				Usage: "The `<repository_slug>` the user will be added on",
				Value: &command.flags.repositories,
			},
			cli.StringSliceFlag{
				Name:  "name",
				Usage: "The `<name>` of the group to be added on the repository. Can be repeated multiple times",
				Value: &command.flags.names,
			},
			cli.StringFlag{
				Name:        "permission",
				Usage:       "The `<permission>` level the user will have (one of REPO_READ, REPO_WRITE, REPO_ADMIN)",
				Destination: &command.flags.permission,
			},
		},
		BashComplete: func(c *cli.Context) {
			helper.AutoComplete(c, command.Settings.GetFileCache())
		},
	}
}

func (command *GroupGrantCommand) GrantAction(context *cli.Context) error {

	if len(command.flags.repositories) == 0 {
		return fmt.Errorf("flag --repository is required.")
	}

	if len(command.flags.names) == 0 {
		return fmt.Errorf("At least one --name is required.")
	}

	if len(command.flags.permission) == 0 {
		return fmt.Errorf("flag --permission is required.")
	}

	fileCache := command.Settings.GetFileCache()

	client, err := command.Settings.GetApiClient()
	if err != nil {
		return err
	}

	for _, repositorySlug := range command.flags.repositories {

		repo, err := fileCache.SearchRepositorySlug(repositorySlug)

		if err != nil {
			return err
		}

		for _, name := range command.flags.names {
			params := bitclient.SetRepositoryGroupPermissionRequest{
				Name:       name,
				Permission: command.flags.permission,
			}

			err := client.SetRepositoryGroupPermission(repo.Project.Key, repositorySlug, params)

			if err != nil {
				fmt.Printf("[KO] rep%s - %s\n", name, err)
				return fmt.Errorf("repo %s, group %s, permission %s - reason: %s\n", repositorySlug, name, command.flags.permission, err)
			} else {
				fmt.Printf("[OK] repo %s, group %s, permission %s\n", repositorySlug, name, command.flags.permission)
			}
		}
	}

	fmt.Printf("Done granting group permissions\n")

	return nil
}