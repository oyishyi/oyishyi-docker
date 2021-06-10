package main

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var initCommand = cli.Command{
	Name:  "init",
	Usage: "init a container",
	Action: func(context *cli.Context) error {
		logrus.Infof("Start initiating...")
		args := context.Args()
		containerCmd := args.Get(0)
		logrus.Infof("container command: %v", containerCmd)

		return nil
	},
}

var runCommand = cli.Command{
	Name:  "run",
	Usage: "Create a container",
	Flags: []cli.Flag{
		// integrate -i and -t for convenience
		&cli.BoolFlag{
			Name:  "it",
			Usage: "open an interactive tty(pseudo terminal)",
		},
	},
	Action: func(context *cli.Context) error {
		args := context.Args()
		if args.Len() == 0 {
			return errors.New("Run what?")
		}
		containerCmd := args.Get(0)        // command
		tty := context.Bool("it") // presudo terminal

		if err := Run(tty, containerCmd); err != nil {
			return err
		}

		return nil
	},
}

func Run(tty bool, containerCmd string) error {
	if !tty {
		return errors.New("tty fails to initiate, Did you type `-it`?")
	}
	fmt.Println(containerCmd)
	return nil
}
